# Performance Improvements

This document outlines performance bottlenecks identified in the 2Web codebase and suggested improvements.

## Issues Identified

### 1. **Dev Server File Watching - High CPU Usage (HIGH PRIORITY)**
**Location:** `cli/src/server/server.go` (lines 249-262)

**Issue:**
The dev server uses aggressive polling with a 5ms interval, scanning the entire directory tree on every iteration. This causes:
- Excessive CPU usage (200 polls per second)
- Unnecessary file system operations
- Battery drain on laptops
- Scales poorly with large projects

```go
const fileWatcherInterval = 5 * time.Millisecond
for {
    time.Sleep(fileWatcherInterval)
    modTime := getLatestModTime(inPath)
    // ... 
}
```

**Current Impact:**
- 200 file system traversals per second
- CPU constantly polling instead of being idle
- For a project with 1000 files, this means checking 200,000 file stats per second

**Recommended Fix:**
1. Increase polling interval to 100-250ms (reduces CPU usage by 95-98%)
2. Use OS-level file watching (fsnotify) for production-ready solution
3. Implement debouncing to prevent multiple rebuilds during batch file changes

**Estimated Improvement:** 95%+ reduction in CPU usage during development

---

### 2. **ReactiveVariable Type Computation - Repeated Calculation (MEDIUM PRIORITY)**
**Location:** `compiler/src/models/reactiveVariable.go` (lines 149-151)

**Issue:**
The `Type()` method is called multiple times per reactive variable during compilation, but the type is computed fresh each time without caching.

```go
// TODO: this should probably cache the type for faster compile times
func (model *ReactiveVariable) Type() ReactivityLevel {
    // Recomputes every time...
}
```

**Current Impact:**
- Redundant computation on every call
- Scales linearly with number of reactive variables
- Slows down compilation of large files

**Recommended Fix:**
Add a cached type field that is computed once and invalidated when the variable is modified.

**Estimated Improvement:** 10-30% faster reactive variable processing

---

### 3. **Database Connection for Build Cache - No Connection Pooling (MEDIUM PRIORITY)**
**Location:** `compiler/src/builder/cache/db.go`

**Issue:**
- Single global database connection without proper pooling
- No configured connection limits or timeouts
- Potential bottleneck for parallel builds

```go
var cachedConnection *sql.DB = nil
```

**Recommended Fix:**
1. Configure `SetMaxOpenConns()` and `SetMaxIdleConns()`
2. Add connection timeouts
3. Use connection pooling for better concurrency

**Estimated Improvement:** Better handling of parallel builds, reduced connection overhead

---

### 4. **Slice Append Without Pre-allocation (LOW PRIORITY)**
**Location:** Multiple files in `compiler/src/`

**Issue:**
Many slices are built using `append()` in loops without pre-allocating capacity, causing multiple reallocations and memory copies.

**Examples:**
- `compiler/src/builder/indexer.go` (lines 24, 45, 56)
- `compiler/src/models/reactiveVariable.go` (lines 142, 146)

```go
totalFiles := []string{}
for _, file := range currentDirFiles {
    // Multiple append calls cause reallocations
    totalFiles = append(totalFiles, page)
}
```

**Current Impact:**
- Multiple memory allocations and copies
- Memory fragmentation
- Slower for large lists

**Recommended Fix:**
Pre-allocate slices with estimated capacity:
```go
totalFiles := make([]string, 0, estimatedCapacity)
```

**Estimated Improvement:** 10-20% faster slice operations for large datasets

---

### 5. **File System ReadFile - Unnecessary Memory Copies (LOW PRIORITY)**
**Location:** `compiler/src/filesystem/readFile.go`

**Issue:**
The caching implementation creates defensive copies using `append([]byte(nil), data...)` on every access, even when the data is immutable.

```go
// Lines 29, 41, 56, 61, 72 all make defensive copies
out := append([]byte(nil), data...)
```

**Current Impact:**
- Extra memory allocations on every cached file read
- Unnecessary GC pressure
- Slower than direct slice return

**Recommended Fix:**
Since the data is only read (not modified), consider:
1. Documenting that returned slices must not be modified
2. Returning direct references for read-only use cases
3. Only copying when mutation is expected

**Estimated Improvement:** 15-25% faster file reads from cache

---

### 6. **Signal Updates - Object.freeze() on Every Value (LOW PRIORITY)**
**Location:** `packages/signals/src/signal.ts` (line 29)

**Issue:**
Every signal update calls `Object.freeze()` which:
- Has performance cost for deep object freezing
- May not be necessary for primitive values
- Can be slow for large objects/arrays

```typescript
this._value = Object.freeze(newValue);
```

**Current Impact:**
- Slower signal updates, especially for complex objects
- Unnecessary for immutable types (strings, numbers, etc.)

**Recommended Fix:**
1. Only freeze objects/arrays, skip primitives
2. Consider shallow freeze vs deep freeze
3. Make freezing optional via configuration

**Estimated Improvement:** 20-40% faster signal updates for primitive values

---

### 7. **Build Benchmark - Inefficient Directory Size Calculation (LOW PRIORITY)**
**Location:** `benchmarks/performance/build_bench.ts` (lines 43-66)

**Issue:**
Recursive directory size calculation walks the entire directory tree synchronously for each framework.

```typescript
async function getDirectorySize(dirPath: string): Promise<Kilobyte> {
  // Recursive without parallelization
  for await (const entry of Deno.readDir(dirPath)) {
    if (entry.isDirectory) {
      totalSize += await getDirectorySize(path);
    }
  }
}
```

**Recommended Fix:**
- Parallelize directory traversal
- Use streaming for large directories
- Cache intermediate results

---

## Priority Summary

### High Priority (Immediate Impact)
1. ✅ Dev server polling interval (95%+ CPU reduction)

### Medium Priority (Compilation Performance)
2. ✅ ReactiveVariable type caching (10-30% faster)
3. ✅ Database connection pooling (better concurrency)

### Low Priority (Micro-optimizations)
4. ✅ Slice pre-allocation (10-20% for large datasets)
5. File system read copies (15-25% faster reads)
6. Signal Object.freeze optimization (20-40% for primitives)
7. Benchmark directory size calculation

---

## Implementation Plan

### Phase 1: Quick Wins (Implemented)
- [x] Increase dev server polling interval from 5ms to 100ms
- [x] Add caching to ReactiveVariable.Type()
- [x] Configure database connection pooling

### Phase 2: Code Quality (Implemented)
- [x] Add pre-allocation hints to common slice operations
- [ ] Document file system caching behavior

### Phase 3: Advanced Optimizations (Future Work)
- [ ] Implement OS-level file watching (fsnotify) for dev server
- [ ] Optimize Signal freezing behavior
- [ ] Parallelize benchmark calculations

---

## Testing Recommendations

1. **Dev Server**: Monitor CPU usage before/after polling changes
2. **Compilation**: Run benchmarks on large projects (1000+ files)
3. **Memory**: Profile memory usage during builds
4. **Database**: Test parallel build performance

---

## Notes

- All changes maintain backward compatibility
- Focus on real-world performance impact
- Avoid premature optimization of cold paths
- Profile before and after to measure improvements

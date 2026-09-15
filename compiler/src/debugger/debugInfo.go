package debugger

import (
	"sync"
)

// The compiler runs page compilation on a thread pool, so every buffer in this
// package is guarded by a mutex.
var debugMutex sync.Mutex

var bufferedDebugInfo []DebugInfo = []DebugInfo{}
var bufferedPageDebugInfo []PageDebugInfo = []PageDebugInfo{}

func AddDebugInfo(info DebugInfo) {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	bufferedDebugInfo = append(bufferedDebugInfo, info)
}

// GetDebugInfo returns a copy of the buffered compiler messages so that
// callers can't race on the underlying slice.
func GetDebugInfo() []DebugInfo {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	messages := make([]DebugInfo, len(bufferedDebugInfo))
	copy(messages, bufferedDebugInfo)
	return messages
}

// AddPageDebugInfo records the reactive graph (variables, properties, and
// events) that was compiled for a single page.
func AddPageDebugInfo(pageInfo PageDebugInfo) {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	bufferedPageDebugInfo = append(bufferedPageDebugInfo, pageInfo)
}

// GetPageDebugInfo returns a copy of the buffered reactive graphs.
func GetPageDebugInfo() []PageDebugInfo {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	pages := make([]PageDebugInfo, len(bufferedPageDebugInfo))
	copy(pages, bufferedPageDebugInfo)
	return pages
}

// currentBuildPages holds the input file paths of the build that is currently
// running. The debug file writer uses it to prune stale reactive graphs that
// were carried over from previous builds (e.g. pages that were deleted from
// the site, or pages that belong to a different input directory).
var currentBuildPages = map[string]bool{}

// SetCurrentBuildPages records the input file paths that the current build
// compiles. It must be called before the pages are compiled.
func SetCurrentBuildPages(pages []string) {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	currentBuildPages = make(map[string]bool, len(pages))
	for _, pagePath := range pages {
		currentBuildPages[pagePath] = true
	}
}

// HasCurrentBuildPage returns whether the input file path belongs to the
// current build.
func HasCurrentBuildPage(pagePath string) bool {
	debugMutex.Lock()
	defer debugMutex.Unlock()

	return currentBuildPages[pagePath]
}

type DebugInfo struct {
	Message string
}

// VariableInfo describes a reactive variable declaration
// (e.g. `$count = 0;` inside a <script compiled> block).
//
// ReactivityClass mirrors the compiler's reactivity levels:
// "unused", "static", "static-property", "assignment", "reactive".
type VariableInfo struct {
	Name            string   `json:"name"`
	InitialValue    string   `json:"initialValue"`
	ReactivityClass string   `json:"reactivity"`
	DependsOn       []string `json:"dependsOn"`
}

// PropertyInfo describes a reactive property binding
// (e.g. *hidden="$isOpen" or {{ $count }}).
//
// Dependencies lists the reactive variables that the property's reducer
// references.
type PropertyInfo struct {
	PropName     string   `json:"prop"`
	Reducer      string   `json:"reducer"`
	Dependencies []string `json:"dependencies"`
}

// EventInfo describes a reactive event binding
// (e.g. @click="$count = $count + 1").
//
// Sink is the variable that the event assigns to (empty when the event doesn't
// assign to a declared reactive variable) and Dependencies lists the reactive
// variables that the event reducer reads.
type EventInfo struct {
	EventName    string   `json:"event"`
	Reducer      string   `json:"reducer"`
	Sink         string   `json:"sink"`
	Dependencies []string `json:"dependencies"`
}

// PageDebugInfo is the reactive graph for a single compiled page.
type PageDebugInfo struct {
	Page       string         `json:"page"`
	Variables  []VariableInfo `json:"variables"`
	Properties []PropertyInfo `json:"properties"`
	Events     []EventInfo    `json:"events"`
}

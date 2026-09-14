package nodes

import (
	"strings"
	"testing"

	lexer "hudson-newey/2web/src/compiler/2-lexer"
	"hudson-newey/2web/src/compiler/io/reader"
)

func extractTestFunctions(t *testing.T, blockSource string) []compiledFunction {
	t.Helper()

	page := "<html><body><script compiled>\n" + blockSource + "\n</script></body></html>"

	lexInstance := lexer.NewLexer(reader.NewReader("test.html", page))
	structure := lexInstance.Execute()

	return extractCompiledFunctions(structure)
}

func TestExtractCompiledFunctions(t *testing.T) {
	functions := extractTestFunctions(t, "$count = 1;\n\nfunction bump() {\n  $count = $count + 1;\n}")

	if len(functions) != 1 {
		t.Fatalf("expected 1 function, got %d", len(functions))
	}

	fn := functions[0]

	if fn.name != "bump" {
		t.Errorf("expected function name 'bump', got %q", fn.name)
	}

	if !strings.Contains(fn.source, "function bump()") || !strings.HasSuffix(fn.source, "}") {
		t.Errorf("unexpected function source: %q", fn.source)
	}

	// The function's span must cover the whole declaration so that reactive
	// variable statements inside the body can be told apart from the block's
	// top level declarations.
	if fn.start.row >= fn.end.row {
		t.Errorf("expected the function span to cover multiple rows, got (%d,%d)-(%d,%d)", fn.start.row, fn.start.col, fn.end.row, fn.end.col)
	}
}

func TestExtractNestedFunctionsStayInside(t *testing.T) {
	functions := extractTestFunctions(t, "function outer() {\n  function inner() {\n    $count++;\n  }\n  $count--;\n}")

	if len(functions) != 1 {
		t.Fatalf("expected only the top level function to be extracted, got %d", len(functions))
	}

	if functions[0].name != "outer" {
		t.Errorf("expected 'outer', got %q", functions[0].name)
	}

	// The nested function's body (and the outer body) still count as
	// assignment sinks.
	sinks := map[string]bool{}
	for _, reference := range functions[0].references {
		if reference.kind != refRead {
			sinks[reference.selector] = true
		}
	}

	if !sinks["$count"] {
		t.Errorf("expected $count to be an assignment sink")
	}
}

func TestScanReactiveReferences(t *testing.T) {
	source := "$count = $count + 1; if ($count > 2) { $name++; $label += 'x'; $other = getValue($name); $count--; }"
	locate := func(int) sourcePosition { return sourcePosition{row: 1, col: 1} }

	references := scanReactiveReferences(source, locate)

	// Every non read reference, in source order. $count appears twice: as a
	// plain assignment and as the decrement shorthand.
	expectedAssignments := []struct {
		selector string
		kind     reactiveReferenceKind
	}{
		{"$count", refAssign},
		{"$name", refIncrement},
		{"$label", refCompoundAssign},
		{"$other", refAssign},
		{"$count", refIncrement},
	}

	assignments := []reactiveReference{}
	for _, reference := range references {
		if reference.kind != refRead {
			assignments = append(assignments, reference)
		}
	}

	if len(assignments) != len(expectedAssignments) {
		t.Fatalf("expected %d assignments, got %d (%v)", len(expectedAssignments), len(assignments), assignments)
	}

	for i, expected := range expectedAssignments {
		if assignments[i].selector != expected.selector || assignments[i].kind != expected.kind {
			t.Errorf("assignment %d: expected %s (%d), got %s (%d)", i, expected.selector, expected.kind, assignments[i].selector, assignments[i].kind)
		}
	}

	// Reads that are not part of an assignment statement. (Reads inside an
	// assignment's value expression are consumed by the assignment reference
	// and resolved when the assignment is rewritten.)
	readCount := 0
	for _, reference := range references {
		if reference.kind == refRead {
			readCount++
		}
	}

	if readCount != 1 {
		t.Errorf("expected 1 read (the $count comparison), got %d", readCount)
	}
}

func TestRewriteCompiledFunction(t *testing.T) {
	functions := extractTestFunctions(t, "$count = 1;\n\nfunction bump(step) {\n  $count += step;\n  $count++;\n  $message = 'count: ' + $count;\n}")

	if len(functions) != 1 {
		t.Fatalf("expected 1 function, got %d", len(functions))
	}

	// Build the index state by hand: $count has a runtime representation and
	// an update function; $message only exists as a function assignment.
	count := &reactiveVariableNode{variableName: "count", initialValue: "1"}
	message := &reactiveVariableNode{variableName: "message", initialValue: `""`}

	index := &ReactiveIndex{
		Variables:         []*reactiveVariableNode{count, message},
		dependencies:      map[string][]string{},
		runtimeVariables:  map[string]string{},
		rpcFunctions:      map[string]bool{},
		compiledFunctions: map[string]bool{},
		functionSinks:     map[string]bool{},
		handlerNames:      map[string]string{},
	}

	index.RegisterRuntimeVariable("$count", "count")
	index.RegisterHandlerName("$count", "updateCount")
	index.RegisterRuntimeVariable("$message", "message")
	index.RegisterHandlerName("$message", "updateMessage")

	rewritten := rewriteCompiledFunction(functions[0], index)

	expectedStatements := []string{
		"count += step; updateCount(count);",
		"count++; updateCount(count);",
		"message = 'count: ' + count; updateMessage(message);",
	}

	for _, statement := range expectedStatements {
		if !strings.Contains(rewritten, statement) {
			t.Errorf("expected the rewritten function to contain %q:\n%s", statement, rewritten)
		}
	}

	if strings.Contains(rewritten, "$") {
		t.Errorf("expected every reactive reference to be rewritten:\n%s", rewritten)
	}
}

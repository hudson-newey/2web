package cfIf

import "slices"

// TODO: if conditions do not currently support reactive variables of other
// javascript code like iife expressions
func IfConditionContent(condition string, conditionalContent string) string {
	// We follow the "C like" if statements, where there are clearly defined
	// "falsy" values (0 in C), and any other value is counted as a "truthy"
	// value.
	//
	// TODO: This has some typing issues where in an ideal world would be caught
	// by the compiler if I ever got some time to add typing information to this
	// compiler.
	falsyConditions := []string{"false", "0", "undefined", "null"}

	if slices.Contains(falsyConditions, condition) {
		// return nothing, meaning that the if condition compiler selector will be
		// replaced with an empty string (removed)
		return ""
	}

	return conditionalContent
}

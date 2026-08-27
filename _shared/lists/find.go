package lists

import "errors"

// TODO: I'm supposed to have access to slices.Find in my go.mod version, but my
// lsp can't seem to find it.
func Find[T any](items []T, condition func(x T) bool) (T, error) {
	for i := range items {
		if condition(items[i]) {
			return items[i], nil
		}
	}
	return items[0], errors.New("Could not find match in Find condition")
}

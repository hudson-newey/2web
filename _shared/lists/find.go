package lists

import "errors"

// TODO: I'm supposed to have access to slices.Find in my go.mod version, but my
// lsp can't seem to find it.
func Find[T any](items []T, condition func(x T) bool) (T, error) {
	var noMatchSentinel T
	if len(items) == 0 {
		return noMatchSentinel, errors.New("Could not find match in Find condition")
	}

	for i := range items {
		if condition(items[i]) {
			return items[i], nil
		}
	}

	return noMatchSentinel, errors.New("Could not find match in Find condition")
}

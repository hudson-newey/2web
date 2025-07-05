package helpers

// Creates a discriminated union between T | nil
func Optional[T any](value ...T) *T {
	if len(value) == 0 {
		var void *T
		return void
	}

	return &value[0]
}

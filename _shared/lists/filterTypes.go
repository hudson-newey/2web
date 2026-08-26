package lists

func FilterTypes[T any, I any](elements []I) []T {
	filtered := make([]T, 0)

	// Use index here to prevent creating a copy
	for _, item := range elements {
		if val, ok := any(item).(T); ok {
			filtered = append(filtered, val)
		}
	}

	return filtered
}

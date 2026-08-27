package lists

func Filter[T any](items []T, condition func(x T) bool) []T {
	result := make([]T, 0)
	for i := range items {
		if condition(items[i]) {
			result = append(result, items[i])
		}
	}
	return result
}

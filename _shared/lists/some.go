package lists

func Some[T any](items []T, condition func(x T) bool) bool {
	_, err := Find(items, condition)
	return err == nil
}

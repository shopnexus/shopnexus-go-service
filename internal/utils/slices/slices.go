package slices

func NonEmptySlice[T any](arr []T) []T {
	if arr == nil {
		return []T{}
	}
	return arr
}

func Last[T any](slice []T) *T {
	if len(slice) == 0 {
		return nil
	}
	return &slice[len(slice)-1]
}

// UnsafeConvertSlice converts a slice of one type to another type
func UnsafeConvertSlice[From any, To any](items []From) []To {
	result := make([]To, len(items))
	for i, v := range items {
		result[i] = any(v).(To)
	}
	return result
}

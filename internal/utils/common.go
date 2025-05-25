package util

import "golang.org/x/exp/constraints"

func Diff[T comparable](a, b []T) (added, removed []T) {
	aMap := make(map[T]struct{}, len(a))

	for _, item := range a {
		aMap[item] = struct{}{}
	}

	// Track items that are in b but not in a
	for _, item := range b {
		if _, exists := aMap[item]; !exists {
			added = append(added, item)
		} else {
			// Remove from aMap to avoid unnecessary iteration later
			delete(aMap, item)
		}
	}

	// Remaining items in aMap are the removed ones
	for item := range aMap {
		removed = append(removed, item)
	}

	return added, removed
}

func Min[T constraints.Ordered](x T, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Ordered](x T, y T) T {
	if x > y {
		return x
	}
	return y
}

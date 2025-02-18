package util

import (
	"runtime"
	"time"
)

func ToPtr[T any](v T) *T {
	return &v
}

func DerefOrNil[T any](v *T) any {
	if v == nil {
		return nil
	}

	return *v
}

func DerefDefault[T any](v *T, d T) T {
	if v == nil {
		return d
	}

	return *v
}

func UnsafeDerefOrNil[T any](v *T) any {
	if v == nil {
		return nil
	}

	return *v
}

// PtrMilisToTime converts a millisecond timestamp to a time.Time pointer
func PtrMilisToTime(v *int64) *time.Time {
	if v == nil {
		return nil
	}
	return ToPtr(time.UnixMilli(*v))
}

func BrandedToStringPtr[T ~string](b *T) *string {
	if b == nil {
		return nil
	}
	str := string(*b)
	return &str
}

func Trace(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc)
	return fn.Name()
}

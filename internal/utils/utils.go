package utils

func Ptr[T any](value T) *T {
	return &value
}

func NewPtr[T any](value T) *T {
	p := new(T)
	*p = value
	return p
}

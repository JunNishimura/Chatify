package utils

type Info[T any] struct {
	Value      T
	HasChanged bool
}

func SetInfo[T any](i *Info[T], v T) {
	i.Value = v
	i.HasChanged = true
}

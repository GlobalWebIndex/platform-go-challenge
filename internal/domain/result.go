package domain

type Result[T any] struct {
	Val T
	Err error
}

func (r Result[T]) Ok() bool {
	return r.Err == nil
}

func (r Result[T]) ValueOr(v T) T {
	if r.Ok() {
		return r.Val
	}
	return v
}
func (r Result[T]) ValueOrPanic() T {
	if r.Ok() {
		return r.Val
	}
	panic(r.Val)
}

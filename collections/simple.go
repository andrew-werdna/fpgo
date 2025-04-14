package collections

func Map[T, U any](collection []T, fn func(T) U) []U {
	out := make([]U, 0)
	for _, v := range collection {
		out = append(out, fn(v))
	}
	return out
}

func Filter[T any](collection []T, predicate func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range collection {
		if predicate(v) {
			out = append(out, v)
		}
	}
	return out
}

func Reduce[T, U any](collection []T, fn func(U, T, int) U, init U) U {
	acc := init
	for i := range collection {
		current := collection[i]
		acc = fn(acc, current, i)
	}
	return acc
}

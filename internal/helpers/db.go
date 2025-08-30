package helpers

func ValueOr[T any](valid bool, val T, def T) T {
	if valid {
		return val
	}
	return def
}

func PtrFromNull[T any](valid bool, val T) *T {
	if !valid {
		return nil
	}
	return &val
}

func ToNull[T any, N any](p *T, wrap func(T) N) N {
	var zero N
	if p == nil {
		return zero
	}
	return wrap(*p)
}

package ptr

func Value[T any](in *T, fallback T) T {
	if in == nil {
		return fallback
	}
	return *in
}

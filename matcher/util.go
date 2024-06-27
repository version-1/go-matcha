package matcher

func typeMatch[T any](v any) bool {
	switch v.(type) {
	case T:
		return true
	default:
		return false
	}
}

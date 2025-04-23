package utils

type Number interface {
	~int | ~uint
}

func DefaultString(src, dst string) string {
	if src == "" {
		return dst
	}
	return src
}

func DefalutInt[T Number](src, dst T) T {
	if src == 0 {
		return dst
	}
	return src
}

func DefaultBool(src, dst bool) bool {
	if src {
		return src
	}
	return dst
}

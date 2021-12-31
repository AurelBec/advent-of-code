package utils

type integer interface {
	int | uint | uint16 | uint32 | uint64 | byte | rune
}

type number interface {
	integer | float32 | float64
}

// Max returns the biggest number between A & B
func Max[K number](a, b K) K {
	if a > b {
		return a
	}
	return b
}

// Min returns the smallest number between A & B
func Min[K number](a, b K) K {
	if a < b {
		return a
	}
	return b
}

// Abs returns the absolute value of A
func Abs[K number](a K) K {
	if a < 0 {
		return -a
	}
	return a
}

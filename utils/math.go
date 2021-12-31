package utils

import (
	"math"
)

type integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type number interface {
	integer | float32 | float64
}

// Max returns the biggest number between A & B
func Max[K number](x, y K) K {
	return max(x, y)
}

// Min returns the smallest number between A & B
func Min[K number](x, y K) K {
	return min(x, y)
}

// Abs returns the absolute value of A
func Abs[K number](a K) K {
	if a < 0 {
		return -a
	}
	return a
}

// Sign returns a unit number with same sign as a
func Sign[K number](a K) (s K) {
	if a < 0 {
		return s - 1
	} else if a > 0 {
		return s + 1
	} else {
		return s
	}
}

// Mod returns the modulo value i%n, always positive
func Mod[K integer](i, n K) K {
	return ((i % n) + n) % n
}

// GCD returns the Greatest Common Divisor via Euclidean algorithm
func GCD[K integer](a, b K) K {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM return the Least Common Multiple via GCD
func LCM[K integer](integers ...K) K {
	if l := len(integers); l == 0 {
		return 0
	} else if l == 1 {
		return integers[0]
	}

	result := integers[0] * integers[1] / GCD(integers[0], integers[1])
	for i := 2; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}
	return result
}

// Interpolation012 returns the quadratic interpolation for an order 2 polynomial defined for x=0,1,2
// from Lagrange:
//
// y0(x-x1)(x-x2)/(x0-x1)/(x0-x2)+y1(x-x0)(x-x2)/(x1-x0)/(x1-x2)+y2(x-x0)(x-x1)/(x2-x0)/(x1-x0)
//
// using x0=0, x1=1, x2=2
// y0(x-1)(x-2)/-1/-2+y1(x)(x-2)/1/-1+y2(x)(x-1)/21/1
// (y0(x-1)(x-2)-2y1(x)(x-2)+y2(x)(x-1))/2
// (y0(x²-3x+2)-y1(2x²-4x)+y2(x²-x))/2
// (x²(y0-2y1+y2) + x(-3y0+4y1-y2) + 2y0) / 2
func Interpolation012[K number](x float64, y0, y1, y2 K) K {
	a := float64(y0 - 2*y1 + y2)
	b := float64(4*y1 - 3*y0 - y2)
	c := float64(2 * y0)
	return K((a*x*x + b*x + c) / 2)
}

// Interpolation returns the quadratic interpolation for an order 2 polynomial defined by 3 F(Xi)=Yi pairs
func Interpolation[K number](x, x0, x1, x2, y0, y1, y2 K) K {
	return y0*(x-x1)/(x0-x1)*(x-x2)/(x0-x2) + y1*(x-x0)/(x1-x0)*(x-x2)/(x1-x2) + y2*(x-x0)/(x2-x0)*(x-x1)/(x1-x0)
}

// Sum sums all elements in input range
func Sum[S ~[]E, E number | string](inputs S) (sum E) {
	for _, input := range inputs {
		sum += input
	}
	return
}

// SumFunc sums all callback results for all elements in input range
func SumFunc[S ~[]E, E any, T number | string](inputs S, callback func(E, ...int) T) (sum T) {
	for i, input := range inputs {
		sum += callback(input, i)
	}
	return
}

// Multiply multiplies all elements in input range
func Multiply[S ~[]E, E number](inputs S) (mul E) {
	mul = 1
	for _, input := range inputs {
		mul *= input
	}
	return
}

// MultiplyFunc multiplies all callback results for all elements in input range
func MultiplyFunc[S ~[]E, E any, T number](inputs S, callback func(E, ...int) T) (mul T) {
	mul = 1
	for i, input := range inputs {
		mul *= callback(input, i)
	}
	return
}

// GaussianEliminationSolver solves linear equation system using gaussian elimination
func GaussianEliminationSolver[E number](A [][]E, B []E) []E {
	n := len(A)

	// generate the extended matrix M=(A|B)
	M := make([][]float64, n)
	for i := range M {
		M[i] = make([]float64, n+1)
		for j, v := range A[i] {
			M[i][j] = float64(v)
		}
		M[i][n] = float64(B[i])
	}

	// set to 0 the lower triangle
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if M[i][j] == 0 {
				continue
			}
			f := M[i][j] / M[j][j]
			for x := range M[i] {
				M[i][x] = M[i][x] - M[j][x]*f
			}
		}
	}

	// set to 0 the upper triangle
	for i := n - 2; i >= 0; i-- {
		for j := i + 1; j < n; j++ {
			if M[i][j] == 0 {
				continue
			}
			f := M[i][j] / M[j][j]
			for x := range M[i] {
				M[i][x] = M[i][x] - M[j][x]*f
			}
		}
	}

	// retrieve solution
	R := make([]E, n)
	for i := range R {
		R[i] = E(math.Round(M[i][n] / M[i][i]))
	}

	return R
}

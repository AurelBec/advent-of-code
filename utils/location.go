package utils

// Location2D represents a 2D X/Y coordinate
type Location2D[K number] struct {
	X, Y K
}

// NewLocation is a quick way to get a Location2D without worrying about specification and value assignation
func NewLocation2D[K number](x, y K) Location2D[K] {
	return Location2D[K]{X: x, Y: y}
}

// ManhattanDist returns the manhattan distance between two Location2D
func (lhs Location2D[K]) ManhattanDist(rhs Location2D[K]) K {
	return Abs(lhs.X-rhs.X) + Abs(lhs.Y-rhs.Y)
}

func (loc Location2D[K]) MovedBy(dx, dy K) Location2D[K] {
	return NewLocation2D(loc.X+dx, loc.Y+dy)
}

// Location3D represents a 3D X/Y/Z coordinate
type Location3D[K number] struct {
	X, Y, Z K
}

// NewLocation is a quick way to get a Location3D without worrying about specification and value assignation
func NewLocation3D[K number](x, y, z K) Location3D[K] {
	return Location3D[K]{X: x, Y: y, Z: z}
}

package hexcoord

import (
	"math"
)

// HexFractional is fractional hex coordinates in
// cubic coordinate system.
type HexFractional struct {
	Q float64
	R float64
}

// S is the implicit additional coordinate when using cubic coordinate system.
func (h HexFractional) S() float64 {
	return -1 * (h.Q + h.R)
}

// ToHex takes in fractional hex coordinates in
// cubic coordinates and rounds them to the nearest
// actual hex coordinate. This is all in normal coordinate
// space, not screen space.
func (h HexFractional) ToHex() Hex {
	q := round(h.Q)
	r := round(h.R)
	s := round(h.S())

	qd := math.Abs(float64(q) - h.Q)
	rd := math.Abs(float64(r) - h.R)
	sd := math.Abs(float64(s) - h.S())

	if qd > rd && qd > sd {
		q = -r - s
	} else if rd > sd {
		r = -q - s
	}

	return Hex{
		Q: q,
		R: r,
	}
}

func round(f float64) int {
	if f > 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
}

func closeEnough(a, b float64) bool {
	if a == b {
		return true
	}
	return math.Abs(a-b) < 1e-10
}

// AlmostEquals returns true when h and x are equal or close
// enough to equal for practical matters.
func (h HexFractional) AlmostEquals(x HexFractional) bool {
	return closeEnough(h.Q, x.Q) && closeEnough(h.R, x.R)
}

// Add combines two hexes.
func (h HexFractional) Add(x HexFractional) HexFractional {
	o := HexFractional{
		Q: x.Q + h.Q,
		R: x.R + h.R,
	}
	return o
}

// Subtract combines two hexes.
func (h HexFractional) Subtract(x HexFractional) HexFractional {
	o := HexFractional{
		Q: h.Q - x.Q,
		R: h.R - x.R,
	}
	return o
}

// Multiply scales a hex by a scalar value.
func (h HexFractional) Multiply(k float64) HexFractional {
	o := HexFractional{
		Q: h.Q * k,
		R: h.R * k,
	}
	return o
}

// LerpHexFractional finds a point between a and b weighted by t.
// See https://en.wikipedia.org/wiki/Linear_interpolation
func LerpHexFractional(a HexFractional, b HexFractional, t float64) HexFractional {
	return HexFractional{
		lerpFloat(a.Q, b.Q, t),
		lerpFloat(a.R, b.R, t),
	}
}

// Length gets the length of the hex to the grid origin.
// This is the Euclidean Distance.
func (h HexFractional) Length() float64 {
	return h.DistanceTo(HexFractional{0, 0})
}

// DistanceTo returns the distance between two hexes.
// This is the Euclidean Distance.
func (h HexFractional) DistanceTo(x HexFractional) float64 {
	d := h.Subtract(x)
	return math.Sqrt(d.Q*d.Q + d.R*d.R + d.Q*d.R)
}

// Normalize returns a vector that points in the same direction
// but has a length of 1.
func (h HexFractional) Normalize() HexFractional {
	return h.Multiply(1.0 / h.Length())
}

// DotProduct returns the dot product.
// See https://en.wikipedia.org/wiki/Dot_product
func (h HexFractional) DotProduct(x HexFractional) float64 {
	return h.Q*x.Q + h.R*x.R
}

// ProjectOn projects h onto x.
// It returns a vector parallel to x.
func (h HexFractional) ProjectOn(x HexFractional) HexFractional {
	return h.Multiply(h.DotProduct(x) / x.DotProduct(x))
}

// Rotate should move a hex about a center point counterclockwise
// by some number of radians.
func (h HexFractional) Rotate(center HexFractional, radians float64) HexFractional {

	cartX, cartY := h.Subtract(center).ToCartesian()

	rotation := complex(math.Cos(radians), math.Sin(radians))
	rotated := complex(cartX, cartY) * rotation

	return HexFractionalFromCartesian(real(rotated), imag(rotated)).Add(center)
}

var sqrt3 float64

func init() {
	sqrt3 = math.Sqrt(3.0)
}

// ToCartesian returns the hex in Cartesian Coordinates.
func (h HexFractional) ToCartesian() (x, y float64) {
	x = sqrt3*h.Q + sqrt3*h.R/2.0
	y = 1.5 * h.R
	return
}

// HexFractionalFromCartesian returns the hex in Cartesian Coordinates.
func HexFractionalFromCartesian(x, y float64) HexFractional {
	// rotate y by 30 degrees to get R
	return HexFractional{
		Q: x*sqrt3/3.0 - y*1.0/3.0,
		R: 2.0 / 3.0 * y,
	}
}
package hexcoord

// CurveSegmenter is a continuous curve segment.
// It can be used to draw a curve.
type CurveSegmenter interface {
	// Sample returns a point on the curve.
	// t is valid for 0 to 1, inclusive.
	Sample(t float64) (position, tangent, curvature HexFractional)

	// Length returns the length of the curve.
	Length() float64
}

// LineSegment is a CurveSegmenter.
type LineSegment struct {
	i      HexFractional
	e      HexFractional
	length float64
	slope  HexFractional
}

// Sample returns a point on the curve.
// t is valid for 0 to 1, inclusive.
func (ls LineSegment) Sample(t float64) (position, tangent, curvature HexFractional) {
	position = LerpHexFractional(ls.i, ls.e, t)
	tangent = ls.slope
	curvature = HexFractional{0.0, 0.0}
	return
}

// Length returns the length of the curve.
func (ls LineSegment) Length() float64 {
	return ls.i.DistanceTo(ls.e)
}

// NewLineSegment creates a line segment curve.
// Inputs are start and end points.
func NewLineSegment(i, e HexFractional) LineSegment {

	return LineSegment{
		i:      i,
		e:      e,
		length: i.DistanceTo(e),
		slope:  i.Subtract(e).Normalize(),
	}
}

// ArcSegment is a CurveSegmenter.
type ArcSegment struct {
	sample func(t float64) (position, tangent, curvature HexFractional)
	length float64
}

// Sample returns a point on the curve.
// t is valid for 0 to 1, inclusive.
func (ac ArcSegment) Sample(t float64) (position, tangent, curvature HexFractional) {
	return ac.sample(t)
}

// Length returns the length of the curve.
func (ac ArcSegment) Length() float64 {
	return ac.length
}

// NewArcSegment creates a circular arc segment curve.
func NewArcSegment(pi, tiu, pe HexFractional) ArcSegment {
	// TODO
	panic("not implemented yet")
}

// combinationSegment is a CurveSegmenter.
type combinationSegment struct {
	segments []CurveSegmenter
	length   float64
}

// Sample returns a point on the curve.
// t is valid for 0 to 1, inclusive.
func (cs combinationSegment) Sample(t float64) (position, tangent, curvature HexFractional) {
	lenT := t * cs.length
	// determine which sub-segment t lands us in.
	prevLength := float64(0.0)
	runningLength := float64(0.0)
	for _, a := range cs.segments {
		runningLength += a.Length()
		if lenT <= runningLength {
			// we are in the current segment
			// now we need to remap t for it

			remappedT := (lenT - prevLength) / runningLength
			return a.Sample(remappedT)
		}
		prevLength += a.Length()
	}
	panic("t is out of scope")
}

// Length returns the length of the curve.
func (cs combinationSegment) Length() float64 {
	return cs.length
}

// JoinSegments creates a multipart curve.
func JoinSegments(arcs []CurveSegmenter) CurveSegmenter {

	// Don't wrap a single element.
	if len(arcs) == 1 {
		return arcs[0]
	}

	cs := combinationSegment{
		segments: arcs,
		length:   float64(0.0),
	}

	for _, a := range arcs {
		cs.length += a.Length()
	}

	return cs
}

// NewCurveSegmenter converts a circular arc into a sample-able curve.
func NewCurveSegmenter(ca circularArc) CurveSegmenter {
	// This is split into 3 cases in an attempt to work
	// around inaccuracies introduced by using floating points.
	v := ca.e.Subtract(ca.i)
	vtDot := v.Normalize().DotProduct(ca.tiu)
	if closeEnough(vtDot, 1.0) {
		// This is the line segment case, where ca.i + ca.tiu is collinear with ca.e.
		return NewLineSegment(ca.i, ca.e)
	}
	// This is the circular arc case.
	return NewArcSegment(ca.i, ca.tiu, ca.e)
}

// this should be dropped and generalized
/*func semiCircleSegment(pi, tiu, pe HexFractional) CurveSegmenter {
	diameter := pi.DistanceTo(pe)
	center := pe.Subtract(pi).Multiply(0.5).Add(pi)
	arcLength := math.Pi * diameter / 2.0
	scalarCurvature := 2.0 / diameter
	centralAngle := math.Pi

	// t = 0 is arcLength = 0 = arclength to pi
	// t = 1 is arcLength = max arcLength = arclength to pe

	return curveSegmenterImpl{
		sample: func(t float64) (position, tangent, curvature HexFractional) {
			// sweep by some ratio of the maximal central angle to get position.
			position = pi.Rotate(center, t*centralAngle)

			tangent = getTangent(center, position)

			// curvature points toward the center of the circle
			curvature = position.Subtract(center).Normalize().Multiply(scalarCurvature * (-1.0))
			return
		},
		length: arcLength,
	}
}
*/
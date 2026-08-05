package drawing

import (
	"math"
	"slices"
)

const (
	// CurveRecursionLimit represents the maximum recursion that is really necessary to subsivide a curve into straight lines
	CurveRecursionLimit = 32
	// maxArcSegments bounds the segments an arc may be tessellated into
	maxArcSegments = 1 << 14
)

// hasNonFinite reports whether any value is NaN or infinite.
func hasNonFinite(vals []float64) bool {
	return slices.ContainsFunc(vals, func(v float64) bool {
		return math.IsNaN(v) || math.IsInf(v, 0)
	})
}

// Cubic
//	x1, y1, cpx1, cpy1, cpx2, cpy2, x2, y2 float64

// SubdivideCubic a Bezier cubic curve in 2 equivalents Bezier cubic curves.
// c1 and c2 parameters are the resulting curves
func SubdivideCubic(c, c1, c2 []float64) {
	// First point of c is the first point of c1
	c1[0], c1[1] = c[0], c[1]
	// Last point of c is the last point of c2
	c2[6], c2[7] = c[6], c[7]

	// Subdivide segment using midpoints
	c1[2] = (c[0] + c[2]) / 2
	c1[3] = (c[1] + c[3]) / 2

	midX := (c[2] + c[4]) / 2
	midY := (c[3] + c[5]) / 2

	c2[4] = (c[4] + c[6]) / 2
	c2[5] = (c[5] + c[7]) / 2

	c1[4] = (c1[2] + midX) / 2
	c1[5] = (c1[3] + midY) / 2

	c2[2] = (midX + c2[4]) / 2
	c2[3] = (midY + c2[5]) / 2

	c1[6] = (c1[4] + c2[2]) / 2
	c1[7] = (c1[5] + c2[3]) / 2

	// Last Point of c1 is equal to the first point of c2
	c2[0], c2[1] = c1[6], c1[7]
}

// TraceCubic subdivides the cubic curve into line segments emitted through t.
// flatteningThreshold sets the flatness tolerance; non-finite control points produce no output.
func TraceCubic(t Liner, cubic []float64, flatteningThreshold float64) {
	// skip non-finite control points to avoid unbounded subdivision
	if hasNonFinite(cubic[:8]) {
		return
	}
	const lastIteration = CurveRecursionLimit - 1

	// Allocation curves
	var curves [CurveRecursionLimit * 8]float64
	copy(curves[0:8], cubic[0:8])

	// current curve
	for i := 0; i >= 0; {
		c := curves[i*8:]
		dx := c[6] - c[0]
		dy := c[7] - c[1]

		d2 := math.Abs((c[2]-c[6])*dy - (c[3]-c[7])*dx)
		d3 := math.Abs((c[4]-c[6])*dy - (c[5]-c[7])*dx)

		// degenerate point can't flatten; emit and pop
		isPoint := dx == 0 && dy == 0 && c[2] == c[0] && c[3] == c[1] && c[4] == c[0] && c[5] == c[1]

		// if it's flat then trace a line
		if (d2+d3)*(d2+d3) < flatteningThreshold*(dx*dx+dy*dy) || isPoint || i == lastIteration {
			t.LineTo(c[6], c[7])
			i--
		} else {
			// second half of bezier go lower onto the stack
			SubdivideCubic(c, curves[(i+1)*8:], curves[i*8:])
			i++
		}
	}
}

// Quad
// x1, y1, cpx1, cpy2, x2, y2 float64

// SubdivideQuad a Bezier quad curve in 2 equivalents Bezier quad curves.
// c1 and c2 parameters are the resulting curves
func SubdivideQuad(c, c1, c2 []float64) {
	// First point of c is the first point of c1
	c1[0], c1[1] = c[0], c[1]
	// Last point of c is the last point of c2
	c2[4], c2[5] = c[4], c[5]

	// Subdivide segment using midpoints
	c1[2] = (c[0] + c[2]) / 2
	c1[3] = (c[1] + c[3]) / 2
	c2[2] = (c[2] + c[4]) / 2
	c2[3] = (c[3] + c[5]) / 2
	c1[4] = (c1[2] + c2[2]) / 2
	c1[5] = (c1[3] + c2[3]) / 2
	c2[0], c2[1] = c1[4], c1[5]
}

func traceWindowIndices(i int) (startAt, endAt int) {
	startAt = i * 6
	endAt = startAt + 6
	return
}

func traceCalcDeltas(c []float64) (dx, dy, d float64) {
	dx = c[4] - c[0]
	dy = c[5] - c[1]
	d = math.Abs((c[2]-c[4])*dy - (c[3]-c[5])*dx)
	return
}

func traceIsFlat(dx, dy, d, threshold float64) bool {
	return (d * d) < threshold*(dx*dx+dy*dy)
}

func traceGetWindow(curves []float64, i int) []float64 {
	startAt, endAt := traceWindowIndices(i)
	return curves[startAt:endAt]
}

// TraceQuad subdivides the quadratic curve into line segments emitted through t.
// flatteningThreshold sets the flatness tolerance; non-finite control points produce no output.
func TraceQuad(t Liner, quad []float64, flatteningThreshold float64) {
	// skip non-finite control points to avoid unbounded subdivision
	if hasNonFinite(quad[:6]) {
		return
	}
	const curveLen = CurveRecursionLimit * 6
	const lastIteration = CurveRecursionLimit - 1

	// Allocates curves stack
	curves := make([]float64, curveLen)

	// copy 6 elements from the quad path to the stack
	copy(curves[0:6], quad[0:6])

	var i int
	var c []float64
	var dx, dy, d float64

	for i >= 0 {
		c = traceGetWindow(curves, i)
		dx, dy, d = traceCalcDeltas(c)

		// degenerate point can't flatten; emit and pop
		isPoint := dx == 0 && dy == 0 && c[2] == c[0] && c[3] == c[1]

		// if it's flat then trace a line
		if traceIsFlat(dx, dy, d, flatteningThreshold) || isPoint || i == lastIteration {
			t.LineTo(c[4], c[5])
			i--
		} else {
			SubdivideQuad(c, traceGetWindow(curves, i+1), traceGetWindow(curves, i))
			i++
		}
	}
}

// TraceArc traces an elliptical arc into line segments emitted through t, returning the arc end
// point. Non-finite parameters produce no segments, the end point is still returned.
func TraceArc(t Liner, x, y, rx, ry, start, angle, scale float64) (lastX, lastY float64) {
	end := start + angle
	endX, endY := x+math.Cos(end)*rx, y+math.Sin(end)*ry
	params := [7]float64{x, y, rx, ry, start, angle, scale}
	if hasNonFinite(params[:]) {
		return endX, endY
	}
	clockWise := angle >= 0
	ra := (math.Abs(rx) + math.Abs(ry)) / 2
	da := math.Acos(ra/(ra+0.125/scale)) * 2
	// floor the step so an underflowed or NaN da still sweeps the arc within the segment bound
	if minDA := math.Abs(angle) / maxArcSegments; !(da > minDA) {
		da = minDA
	}
	//normalize
	if !clockWise {
		da = -da
	}
	angle = start + da
	var curX, curY float64
	// bounded as a backstop, a floored da always exits on the condition below first
	for i := 0; i < maxArcSegments; i++ {
		if (angle < end-da/4) != clockWise {
			break
		}
		curX = x + math.Cos(angle)*rx
		curY = y + math.Sin(angle)*ry

		prevAngle := angle
		angle += da
		t.LineTo(curX, curY)
		if angle == prevAngle {
			break // da below the ULP of angle, the step can never advance again
		}
	}
	return endX, endY
}

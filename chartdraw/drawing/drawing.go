package drawing

// FillRule defines the type for fill rules
type FillRule int

const (
	// FillRuleEvenOdd determines the "insideness" of a point in the shape
	// by drawing a ray from that point to infinity in any direction
	// and counting the number of path segments from the given shape that the ray crosses.
	// If this number is odd, the point is inside; if even, the point is outside.
	FillRuleEvenOdd FillRule = iota
	// FillRuleWinding determines the "insideness" of a point in the shape
	// by drawing a ray from that point to infinity in any direction
	// and then examining the places where a segment of the shape crosses the ray.
	// Starting with a count of zero, add one each time a path segment crosses
	// the ray from left to right and subtract one each time
	// a path segment crosses the ray from right to left. After counting the crossings,
	// if the result is zero then the point is outside the path. Otherwise, it is inside.
	FillRuleWinding
)

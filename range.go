package charts

import (
	"math"
)

const defaultAxisDivideCount = 6

type axisRange struct {
	p           *Painter
	divideCount int
	min         float64
	max         float64
	size        int
	boundary    bool
}

type AxisRangeOption struct {
	Painter *Painter
	// The min value of axis
	Min float64
	// The max value of axis
	Max float64
	// The size of axis
	Size int
	// Boundary gap
	Boundary bool
	// The count of divide
	DivideCount int
}

// NewRange returns a axis range
func NewRange(opt AxisRangeOption) axisRange {
	max := opt.Max
	min := opt.Min

	max += math.Abs(max * 0.1)
	min -= math.Abs(min * 0.1)
	divideCount := opt.DivideCount
	r := math.Abs(max - min)

	// minimum unit calculation
	unit := 1
	if r > 5 {
		unit = 2
	}
	if r > 10 {
		unit = 4
	}
	if r > 30 {
		unit = 5
	}
	if r > 100 {
		unit = 10
	}
	if r > 200 {
		unit = 20
	}
	unit = int((r/float64(divideCount))/float64(unit))*unit + unit

	if min != 0 {
		isLessThanZero := min < 0
		min = float64(int(min/float64(unit)) * unit)
		// if less than zero, int is rounded up, so adjust
		if min < 0 ||
			(isLessThanZero && min == 0) {
			min -= float64(unit)
		}
	}
	max = min + float64(unit*divideCount)
	expectMax := opt.Max * 2
	if max > expectMax {
		max = float64(ceilFloatToInt(expectMax))
	}
	return axisRange{
		p:           opt.Painter,
		divideCount: divideCount,
		min:         min,
		max:         max,
		size:        opt.Size,
		boundary:    opt.Boundary,
	}
}

// Values returns values of range
func (r axisRange) Values() []string {
	offset := (r.max - r.min) / float64(r.divideCount)
	values := make([]string, 0)
	formatter := commafWithDigits
	if r.p != nil && r.p.valueFormatter != nil {
		formatter = r.p.valueFormatter
	}
	for i := 0; i <= r.divideCount; i++ {
		v := r.min + float64(i)*offset
		value := formatter(v)
		values = append(values, value)
	}
	return values
}

func (r *axisRange) getHeight(value float64) int {
	if r.max <= r.min {
		return 0
	}
	v := (value - r.min) / (r.max - r.min)
	return int(v * float64(r.size))
}

func (r *axisRange) getRestHeight(value float64) int {
	return r.size - r.getHeight(value)
}

// GetRange returns a range of index
func (r *axisRange) GetRange(index int) (float64, float64) {
	unit := float64(r.size) / float64(r.divideCount)
	return unit * float64(index), unit * float64(index+1)
}

// AutoDivide divides the axis
func (r *axisRange) AutoDivide() []int {
	return autoDivide(r.size, r.divideCount)
}

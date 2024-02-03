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
}

const rangeMinPaddingPercent = 2.0
const rangeMaxPaddingPercent = 10.0 // max padding percent per-side

// NewRange returns a range of data for an axis, this range will have padding to better present the data.
func NewRange(painter *Painter, size, divideCount int, min, max float64, addPadding bool) axisRange {
	if addPadding {
		min, max = padRange(min, max)
	}
	return axisRange{
		p:           painter,
		divideCount: divideCount,
		min:         min,
		max:         max,
		size:        size,
	}
}

func padRange(min, max float64) (float64, float64) {
	span := max - min
	spanIncrement := span * 0.01 // must be 1% of the span
	var spanIncrementMultiplier float64
	targetFound := false
targetLoop:
	for expo := -1.0; expo < 6; expo++ {
		// use 10^expo so that we prefer 0, 10, 100, etc numbers
		targetVal := math.Floor(math.Pow(10, expo)) // Math.Floor so -1 expo will convert 0.1 to 0
		if min < 0 {
			targetVal *= -1
		}
		// set default and then check if we can even look for this target within our padding range
		spanIncrementMultiplier = rangeMinPaddingPercent
		if targetVal > min-(spanIncrement*rangeMinPaddingPercent) {
			continue targetLoop // we need to get further from zero to get into our minimum padding
		} else if targetVal < min-(spanIncrement*rangeMaxPaddingPercent) {
			break targetLoop // no match possible, use the min padding as a default
		}

		for ; spanIncrementMultiplier < rangeMaxPaddingPercent; spanIncrementMultiplier++ {
			adjustedMin := min - (spanIncrement * spanIncrementMultiplier)
			if adjustedMin <= targetVal { // we found our target value between the min and adjustedMin
				// set the min to the target value and the multiplier so the max can be set
				targetFound = true
				spanIncrementMultiplier = (min - targetVal) / spanIncrement
				min = targetVal
				break targetLoop
			}
		}
	}
	if !targetFound {
		min, spanIncrementMultiplier = roundLimit(min, spanIncrement, rangeMinPaddingPercent, false)
	}
	// update max and match based off the ideal padding
	max, _ = roundLimit(max, spanIncrement, spanIncrementMultiplier, true)
	return min, max
}

func roundLimit(val, increment, multiplier float64, add bool) (float64, float64) {
	for orderOfMagnitude := math.Floor(math.Log10(val)); orderOfMagnitude > 0; orderOfMagnitude-- {
		roundValue := math.Pow(10, orderOfMagnitude)
		var proposedVal float64
		var proposedMultiplier float64
		for roundAdjust := 0.0; roundAdjust < 9.0; roundAdjust++ {
			if add {
				proposedVal = (math.Ceil(val/roundValue) * roundValue) + (roundValue * roundAdjust)
				proposedMultiplier = (proposedVal - val) / increment
			} else {
				proposedVal = (math.Floor(val/roundValue) * roundValue) - (roundValue * roundAdjust)
				proposedMultiplier = (val - proposedVal) / increment
			}

			if proposedMultiplier > rangeMaxPaddingPercent {
				break // shortcut inner loop as multiplier will only go up
			} else if proposedMultiplier >= rangeMinPaddingPercent {
				return proposedVal, proposedMultiplier
			}
		}
		if proposedMultiplier < rangeMinPaddingPercent {
			break // shortcut outer loop if multiplier is below the min after adjust check, as this will only get smaller
		}
	}
	// no rounder alternative found, just adjust based off initial multiplier
	if add {
		return val + (increment * multiplier), multiplier
	} else {
		return val - (increment * multiplier), multiplier
	}
}

// Values returns values of range
func (r axisRange) Values() []string {
	offset := (r.max - r.min) / float64(r.divideCount)
	formatter := commafWithDigits
	if r.p != nil && r.p.valueFormatter != nil {
		formatter = r.p.valueFormatter
	}
	values := make([]string, r.divideCount+1)
	for i := 0; i <= r.divideCount; i++ {
		v := r.min + float64(i)*offset
		values[i] = formatter(v)
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

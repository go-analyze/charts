package charts

import (
	"math"
)

const rangeMinPaddingPercentMin = 0.0 // increasing could result in forced negative y-axis minimum
const rangeMinPaddingPercentMax = 10.0
const rangeMaxPaddingPercentMin = 5.0  // set minimum spacing at the top of the graph
const rangeMaxPaddingPercentMax = 20.0 // larger to allow better chances of finding a good interval
const rangeDefaultPaddingPercent = 5.0

type axisRange struct {
	p           *Painter
	divideCount int
	min         float64
	max         float64
	size        int
}

// NewRange returns a range of data for an axis, this range will have padding to better present the data.
func NewRange(painter *Painter, size, labelCount int, min, max float64, addPadding bool) axisRange {
	if addPadding {
		min, max = padRange(labelCount, min, max)
	}
	return axisRange{
		p:           painter,
		divideCount: labelCount,
		min:         min,
		max:         max,
		size:        size,
	}
}

func padRange(labelCount int, min, max float64) (float64, float64) {
	minResult := min
	maxResult := max
	span := max - min
	spanIncrement := span * 0.01 // must be 1% of the span
	var spanIncrementMultiplier float64
	// find a min value to start our range from
	// we prefer (in order, negative if necessary), 0, 1, 10, 100, ..., 2, 20, ..., 5, 50, ...
rootLoop:
	for _, multiple := range []float64{1.0, 2.0, 5.0} {
	expoLoop:
		for expo := -1.0; expo < 6; expo++ {
			// use 10^expo so that we prefer 0, 10, 100, etc numbers
			targetVal := math.Floor(math.Pow(10, expo)) * multiple // Math.Floor so -1 expo will convert 0.1 to 0
			if targetVal == 0 && multiple != 1 {
				continue expoLoop // we already tested this value
			}
			if min < 0 {
				targetVal *= -1
			}
			// set default and then check if we can even look for this target within our padding range
			if targetVal > min-(spanIncrement*rangeMinPaddingPercentMin) {
				continue expoLoop // we need to get further from zero to get into our minimum padding
			} else if targetVal < min-(spanIncrement*rangeMinPaddingPercentMax) {
				break expoLoop // no match possible, use the min padding as a default
			}

			spanIncrementMultiplier = rangeMinPaddingPercentMin
			for ; spanIncrementMultiplier < rangeMinPaddingPercentMax; spanIncrementMultiplier++ {
				adjustedMin := min - (spanIncrement * spanIncrementMultiplier)
				if adjustedMin <= targetVal { // we found our target value between the min and adjustedMin
					// set the min to the target value and the multiplier so the max can be set
					spanIncrementMultiplier = (min - targetVal) / spanIncrement
					minResult = targetVal
					break rootLoop
				}
			}
		}
	}
	if minResult == min {
		minResult, spanIncrementMultiplier =
			friendlyRound(min, spanIncrement, rangeDefaultPaddingPercent,
				rangeMinPaddingPercentMin, rangeMinPaddingPercentMax, false)
	}
	// update max and match based off the ideal padding
	maxResult, _ =
		friendlyRound(max, spanIncrement, spanIncrementMultiplier,
			rangeMaxPaddingPercentMin, rangeMaxPaddingPercentMax, true)
	// adjust max so that the intervals and labels are also round if possible
	interval := (maxResult - minResult) / float64(labelCount-1)
	maxIntervalIncrease := ((rangeMaxPaddingPercentMax * spanIncrement) - (maxResult - max)) / float64(labelCount-1)
	roundedInterval, _ := friendlyRound(interval, 1.0, 0.0, 0.0, maxIntervalIncrease, true)

	if roundedInterval != interval {
		maxResult = minResult + (roundedInterval * float64(labelCount-1))
	}

	return minResult, maxResult
}

func friendlyRound(val, increment, defaultMultiplier, minMultiplier, maxMultiplier float64, add bool) (float64, float64) {
	absVal := math.Abs(val)
	for orderOfMagnitude := math.Floor(math.Log10(absVal)); orderOfMagnitude > 0; orderOfMagnitude-- {
		roundValue := math.Pow(10, orderOfMagnitude)
		var proposedVal float64
		var proposedMultiplier float64
		for roundAdjust := 0.0; roundAdjust < 9.0; roundAdjust++ {
			if add {
				proposedVal = (math.Ceil(absVal/roundValue) * roundValue) + (roundValue * roundAdjust)
			} else {
				proposedVal = (math.Floor(absVal/roundValue) * roundValue) + (roundValue * roundAdjust)
			}
			if val < 0 { // Apply the original sign back to proposedVal
				proposedVal = -proposedVal
			}
			if add {
				proposedMultiplier = (proposedVal - val) / increment
			} else {
				proposedMultiplier = (val - proposedVal) / increment
			}

			if proposedMultiplier > maxMultiplier {
				break // shortcut inner loop as multiplier will only go up
			} else if proposedMultiplier > minMultiplier {
				return proposedVal, proposedMultiplier
			}
		}
		if proposedMultiplier <= minMultiplier {
			break // shortcut outer loop if multiplier is below the min after adjust check, as this will only get smaller
		}
	}
	// No match found, let's see if we can just round to the next whole number
	if (increment*maxMultiplier) >= 1.0 && val != math.Trunc(val) {
		var proposedVal float64
		var proposedMultiplier float64
		if add {
			proposedVal = math.Ceil(val)
			proposedMultiplier = (proposedVal - val) / increment
		} else {
			proposedVal = math.Floor(val)
			proposedMultiplier = (val - proposedVal) / increment
		}
		return proposedVal, proposedMultiplier
	}
	// No rounder alternative found, just adjust based off default multiplier
	if add {
		return val + (increment * defaultMultiplier), defaultMultiplier
	} else {
		return val - (increment * defaultMultiplier), defaultMultiplier
	}
}

// Values returns values of range
func (r axisRange) Values() []string {
	offset := (r.max - r.min) / float64(r.divideCount-1)
	formatter := commafWithDigits
	if r.p != nil && r.p.valueFormatter != nil {
		formatter = r.p.valueFormatter
	}
	values := make([]string, r.divideCount)
	for i := 0; i < r.divideCount; i++ {
		v := r.min + float64(i)*offset
		values[i] = formatter(v)
	}
	return values
}

func (r axisRange) getHeight(value float64) int {
	if r.max <= r.min {
		return 0
	}
	v := (value - r.min) / (r.max - r.min)
	return int(v * float64(r.size))
}

func (r axisRange) getRestHeight(value float64) int {
	return r.size - r.getHeight(value)
}

// GetRange returns a range of index
func (r axisRange) GetRange(index int) (float64, float64) {
	unit := float64(r.size) / float64(r.divideCount)
	return unit * float64(index), unit * float64(index+1)
}

// AutoDivide divides the axis
func (r axisRange) AutoDivide() []int {
	return autoDivide(r.size, r.divideCount)
}

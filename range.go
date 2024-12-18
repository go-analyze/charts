package charts

import (
	"math"
)

const rangeMinPaddingPercentMin = 0.0 // increasing could result in forced negative y-axis minimum
const rangeMinPaddingPercentMax = 20.0
const rangeMaxPaddingPercentMin = 5.0 // set minimum spacing at the top of the graph
const rangeMaxPaddingPercentMax = 20.0
const zeroSpanAdjustment = 1 // Adjustment

type axisRange struct {
	p           *Painter
	divideCount int
	min         float64
	max         float64
	size        int
}

// NewRange returns a range of data for an axis, this range will have padding to better present the data.
func NewRange(painter *Painter, size, divideCount int, min, max, minPaddingScale, maxPaddingScale float64) axisRange {
	min, max = padRange(divideCount, min, max, minPaddingScale, maxPaddingScale)
	return axisRange{
		p:           painter,
		divideCount: divideCount,
		min:         min,
		max:         max,
		size:        size,
	}
}

func padRange(divideCount int, min, max, minPaddingScale, maxPaddingScale float64) (float64, float64) {
	if minPaddingScale <= 0.0 && maxPaddingScale <= 0.0 {
		return min, max
	}
	// scale percents for min value
	scaledMinPadPercentMin := rangeMinPaddingPercentMin * minPaddingScale
	scaledMinPadPercentMax := rangeMinPaddingPercentMax * minPaddingScale
	// scale percents for max value
	scaledMaxPadPercentMin := rangeMaxPaddingPercentMin * maxPaddingScale
	scaledMaxPadPercentMax := rangeMaxPaddingPercentMax * maxPaddingScale
	minResult := min
	spanIncrement := (max - min) * 0.01 // must be 1% of the span
	var spanIncrementMultiplier float64
	// find a min value to start our range from
	// we prefer (in order, negative if necessary), 0, 1, 10, 100, ..., 2, 20, ..., 5, 50, ...
	updatedMin := false
rootLoop:
	for _, multiple := range []float64{1.0, 2.0, 5.0} {
		if min < 0 {
			multiple *= -1 // convert multiple sign to adjust targetVal correctly
		}
	expoLoop:
		for expo := -1.0; expo < 6; expo++ {
			if expo == -1.0 && multiple != 1.0 {
				continue expoLoop // we only want to test targetVal 0 once
			}
			// use 10^expo so that we prefer 0, 10, 100, etc numbers
			targetVal := math.Floor(math.Pow(10, expo)) * multiple // Math.Floor to convert 0.1 from -1 expo into 0
			if targetVal < min-(spanIncrement*scaledMinPadPercentMax) {
				break expoLoop // no match possible, target value will only get further from start
			} else if targetVal <= min-(spanIncrement*scaledMinPadPercentMin) {
				// targetVal can be between our span increment increases, calculate and set result
				updatedMin = true
				spanIncrementMultiplier = (min - targetVal) / spanIncrement
				minResult = targetVal
				break rootLoop
			} // else try again to meet minimum padding requirements
		}
	}
	if !updatedMin {
		minResult, spanIncrementMultiplier =
			friendlyRound(min, spanIncrement, scaledMinPadPercentMin,
				scaledMinPadPercentMin, scaledMinPadPercentMax, false)
	}
	if minTrunk := math.Trunc(minResult); minTrunk <= min-(spanIncrement*scaledMinPadPercentMin) {
		minResult = minTrunk // remove possible float multiplication inaccuracies
	}

	if max == minResult {
		// no adjustment was made and there is no span, because of that the max calculation below can't function
		// for that reason we apply a default constant span, still wanting to prefer a zero start if possible
		if minResult == 0 {
			return minResult, minResult + (2 * zeroSpanAdjustment)
		}
		return minResult - zeroSpanAdjustment, minResult + zeroSpanAdjustment
	} else if maxPaddingScale <= 0.0 {
		return minResult, max
	} else if math.Abs(max) < 10 {
		return minResult, math.Ceil(max) + 1
	}

	// update max to provide ideal padding and human friendly intervals
	interval := (max - minResult) / float64(divideCount-1)
	roundedInterval, _ := friendlyRound(interval, spanIncrement/float64(divideCount-1),
		math.Max(spanIncrementMultiplier, scaledMaxPadPercentMin),
		scaledMaxPadPercentMin, scaledMaxPadPercentMax, true)
	maxResult := minResult + (roundedInterval * float64(divideCount-1))
	if maxTrunk := math.Trunc(maxResult); maxTrunk >= max+(spanIncrement*scaledMaxPadPercentMin) {
		maxResult = maxTrunk // remove possible float multiplication inaccuracies
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
				proposedVal *= -1
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
	formatter := defaultValueFormatter
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

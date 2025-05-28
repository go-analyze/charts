package charts

import (
	"math"
	"strconv"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/matrix"
)

const rangeMinPaddingPercentMin = 0.0 // increasing could result in forced negative y-axis minimum
const rangeMinPaddingPercentMax = 20.0
const rangeMaxPaddingPercentMin = 5.0 // set minimum spacing at the top of the graph
const rangeMaxPaddingPercentMax = 20.0
const zeroSpanAdjustment = 1

// axisRange represents the calculated range for the axis, as well as values for fitting labels on the range.
type axisRange struct {
	isCategory bool
	// labels are the rendered labels, 1:1 for categories or the range value labels to render
	labels []string
	// dataStartIndex specifies what index the label values should start from.
	dataStartIndex int
	tickCount      int
	divideCount    int
	labelCount     int
	min, max       float64 // only valid if !isCategory
	size           int
	textMaxWidth   int
	textMaxHeight  int
	labelRotation  float64
	labelFontStyle FontStyle
}

// calculateValueAxisRange centralizes the logic for numeric axes, picking a scale and label count that is human friendly.
func calculateValueAxisRange(p *Painter, isVertical bool, axisSize int,
	minCfg, maxCfg, rangeValuePaddingScale *float64,
	labelsCfg []string, dataStartIndex int,
	labelCountCfg int, labelUnit float64, labelCountAdjustment int,
	seriesList seriesList, yAxisIndex int, stackSeries bool,
	valueFormatter ValueFormatter,
	labelRotation float64, fontStyle FontStyle) axisRange {
	// calculate the range
	minVal, maxVal, sumMax := getSeriesMinMaxSumMax(seriesList, yAxisIndex, stackSeries)
	if stackSeries { // If stacked, maxVal should be the maxVal data point of all series summed together
		if minVal > 0 {
			minVal-- // subtract to ensure that all series are represented as a small stacked bar (may otherwise have 0 height)
		}
		maxVal = sumMax
	}
	minPadScale, maxPadScale := 1.0, 1.0
	if rangeValuePaddingScale != nil {
		minPadScale = *rangeValuePaddingScale
		maxPadScale = minPadScale
	}
	if minCfg != nil && *minCfg < minVal {
		minVal = *minCfg
		minPadScale = 0.0
	}
	if maxCfg != nil && *maxCfg > maxVal {
		maxVal = *maxCfg
		maxPadScale = 0.0
	}
	decimalData := minVal != math.Floor(minVal) || (maxVal-minVal) != math.Floor(maxVal-minVal)

	// Label counts and range padding are linked together to produce a user-friendly graph.
	// First when considering padding we want to prefer a zero axis start if reasonable, and add a slight
	// padding to the maxVal so there is a little space at the top of the graph. We also want to pick
	// a maxVal value and label count that will result in round intervals on the axis, or match the user
	// supplied unit (if provided).
	//
	// In order to accomplish this, we estimate the label count (if necessary), produce some labels to measure,
	// calculate our label limit, pad the range, then once the label count is finalized produce the final labels.
	initialLabelCount := labelCountCfg
	if initialLabelCount < 1 {
		if labelUnit > 0 {
			initialLabelCount = int((maxVal-minVal)/labelUnit) + 1
		} else {
			initialLabelCount =
				chartdraw.MinInt(chartdraw.MaxInt(int(maxVal-minVal)+1, defaultYAxisLabelCountLow),
					defaultYAxisLabelCountHigh)
			if decimalData { // if there is a decimal, we double our labels to provide more detail
				initialLabelCount = chartdraw.MinInt(initialLabelCount*2, defaultYAxisLabelCountHigh)
			}
		}
	}
	initialLabelCount = chartdraw.MaxInt(initialLabelCount+labelCountAdjustment, minimumAxisLabels)
	labels := valueLabels(labelsCfg, valueFormatter, minVal, maxVal, initialLabelCount)
	labelW, labelH := p.measureTextMaxWidthHeight(labels, labelRotation, fontStyle)

	// If user gave an explicit LabelCount, then we do NOT do a collision check
	// For default logic we want to make sure we choose a label count that is visually appealing
	padLabelCount := initialLabelCount
	maxLabelCount := padLabelCount
	if labelCountCfg == 0 {
		if isVertical {
			if labelH > 0 { // avoid divide by zero
				maxLabelCount = axisSize / labelH
			}
		} else {
			if labelW > 0 {
				// add to the label width to give good spacing
				maxLabelCount = axisSize / (labelW + chartdraw.MinInt(20, labelW))
			}
		}
		if maxLabelCount < padLabelCount {
			padLabelCount = chartdraw.MaxInt(maxLabelCount, minimumAxisLabels)
		}
		if labelUnit > 0 && padLabelCount > minimumAxisLabels {
			// reduce padLabelCount to ensure it remains within the max count if we have to add one to meet unit expectations
			padLabelCount--
		}
	}
	minPadded, maxPadded := padRange(padLabelCount, minVal, maxVal, minPadScale, maxPadScale)
	labelCount := padLabelCount
	// if the user set only a unit, we may need to refine again after padding to meet the unit
	if labelCountCfg == 0 && labelUnit > 0 {
		// Quick exit when the unit itself is larger than the whole span
		if dataSpan := maxPadded - minPadded; labelUnit >= dataSpan {
			labelCount = minimumAxisLabels
		} else {
			// Snap helpers to ensure we maintain a multiple of `labelUnit`
			down := func(v float64) float64 { return math.Floor(v/labelUnit) * labelUnit }
			up := func(v float64) float64 { return math.Ceil(v/labelUnit) * labelUnit }

			var bestCount int
			bestMin, bestMax := minPadded, maxPadded
			bestPad := math.Inf(1)
			bestDeltaAbs := math.MaxInt

			// Helper that records the “best so far”
			accept := func(c int, mn, mx float64) {
				deltaAbs := int(math.Abs(float64(padLabelCount - c)))
				pad := (minPadded - mn) + (mx - maxPadded)
				if pad < bestPad-matrix.DefaultEpsilon ||
					(math.Abs(pad-bestPad) < matrix.DefaultEpsilon && deltaAbs < bestDeltaAbs) {
					if bestCount == 0 || down(mn)-mn < matrix.DefaultEpsilon && mx-up(mx) < matrix.DefaultEpsilon {
						bestPad = pad
						bestMin = mn
						bestMax = mx
						bestCount = c
						bestDeltaAbs = deltaAbs
					}
				}
			}

			// The search expands symmetrically around padLabelCount
			maxDelta := chartdraw.MaxInt(padLabelCount-minimumAxisLabels, maxLabelCount-padLabelCount)
			for delta := 0; delta <= maxDelta; delta++ {
				try := func(c int) {
					if c < minimumAxisLabels || c > maxLabelCount {
						return
					}
					spanCount := float64(c - 1)
					span := spanCount * labelUnit

					// Attempt in the order of preferred match method

					// snapped min and max to the current label count
					snappedMin := down(minPadded)
					snappedMax := up(maxPadded)
					snappedInterval := up((snappedMax - snappedMin) / spanCount)
					flip := true // flip to alternate which boundary we expand to meet the snapped interval
					for snappedMin+(snappedInterval*spanCount)-snappedMax > matrix.DefaultEpsilon {
						if snappedMin-snappedInterval >= 0 && flip {
							flip = false
							snappedMin -= snappedInterval
						} else {
							flip = true
							snappedMax += snappedInterval
						}
					}
					// final max adjustment to ensure we meet the snapped interval
					snappedMax = math.Ceil(snappedMax/snappedInterval) * snappedInterval
					accept(c, snappedMin, snappedMax)

					// shift MIN downward
					if minCfg == nil {
						candMax := up(maxPadded) // snapped top
						candMin := candMax - span
						if (minPadded < 0 || candMin >= 0.0-matrix.DefaultEpsilon) &&
							candMin <= minPadded+matrix.DefaultEpsilon {
							accept(c, candMin, candMax)
						}
					}

					// split padding (both free)
					if minCfg == nil && maxCfg == nil {
						// center the span around the data as much as multiples allow
						candMin := down(minPadded - (span-dataSpan)/2)
						candMax := candMin + span
						if (minPadded < 0 || candMin >= 0.0-matrix.DefaultEpsilon) &&
							candMin <= minPadded+matrix.DefaultEpsilon &&
							candMax >= maxPadded-matrix.DefaultEpsilon {
							accept(c, candMin, candMax)
						}
					}

					// grow the MAX upward
					if maxCfg == nil {
						candMin := down(minPadded) // snapped bottom
						candMax := candMin + span
						if candMax >= maxPadded-matrix.DefaultEpsilon {
							accept(c, candMin, candMax)
						}
					}
				}

				// Try padLabelCount-delta and padLabelCount+delta in that order
				if cand := padLabelCount - delta; cand >= minimumAxisLabels {
					try(cand)
				}
				if delta != 0 {
					if cand := padLabelCount + delta; cand <= maxLabelCount {
						try(cand)
					}
				}

				if bestPad < matrix.DefaultEpsilon {
					break // perfect fit with zero extra padding
				}
			}

			if bestCount > 0 {
				labelCount = bestCount
				minPadded = bestMin
				maxPadded = bestMax
			} // else, could not match the unit inside range, fallback to the original padLabelCount
		}
	}

	if len(labels) != labelCount || minVal-minPadded > matrix.DefaultEpsilon || maxPadded-maxVal > matrix.DefaultEpsilon {
		// regenerate labels to meet new scale
		labels = valueLabels(labelsCfg, valueFormatter, minPadded, maxPadded, labelCount)
		labelW, labelH = p.measureTextMaxWidthHeight(labels, labelRotation, fontStyle)
	}

	return axisRange{
		isCategory:     false,
		labels:         labels,
		dataStartIndex: dataStartIndex,
		divideCount:    len(labels),
		tickCount:      labelCount,
		labelCount:     labelCount,
		min:            minPadded,
		max:            maxPadded,
		size:           axisSize,
		textMaxWidth:   labelW,
		textMaxHeight:  labelH,
		labelRotation:  labelRotation,
		labelFontStyle: fontStyle,
	}
}

// calculateCategoryAxisRange does the same for category axes (common for X-axis in line/bar charts).
func calculateCategoryAxisRange(p *Painter, axisSize int, isVertical bool, extraSpace bool,
	labels []string, dataStartIndex int,
	labelCountCfg int, labelCountAdjustment int, labelUnit float64,
	seriesList seriesList, labelRotation float64, fontStyle FontStyle) axisRange {
	// If user provided no labels, use series names.
	// If provided only partially, fill in the remaining labels.
	for i := len(labels); i < getSeriesMaxDataCount(seriesList); i++ {
		labels = append(labels, strconv.Itoa(i+1))
	}
	dataCount := len(labels)

	textW, textH := p.measureTextMaxWidthHeight(labels, labelRotation, fontStyle)

	labelCount := labelCountCfg
	if labelCount <= 0 {
		labelCount = dataCount
	} else if labelCount > dataCount {
		labelCount = dataCount
	}
	labelCount = chartdraw.MaxInt(labelCount+labelCountAdjustment, minimumAxisLabels)
	// validate the labels fit, otherwise reduce the count
	if labelCountCfg == 0 {
		maxLabelCount := labelCount
		if isVertical {
			if textH > 0 {
				var extra int
				if extraSpace {
					extra = 10
				}
				maxLabelCount = chartdraw.MaxInt(axisSize/(textH+extra), minimumAxisLabels)
			}
		} else {
			if textW > 0 {
				// add a little extra padding for horizontal layouts
				extra := textW
				if !extraSpace {
					extra /= 2
				}
				maxLabelCount = chartdraw.MaxInt(axisSize/(textW+extra), minimumAxisLabels)
			}
		}
		if labelUnit > 0 {
			// If the user gave a 'unit', figure out how many 'units' fit
			multiplier := 1.0
			for {
				count := ceilFloatToInt(float64(dataCount) / (labelUnit * multiplier))
				if count > maxLabelCount {
					multiplier++
				} else {
					labelCount = chartdraw.MaxInt(count, minimumAxisLabels)
					break
				}
			}
		} else if maxLabelCount < labelCount {
			// Instead of a slight reduction, we choose a skip factor (step) so that we skip every other label until
			// we are within our limit.
			step := 1
			candidateCount := 2 + (dataCount-2)/step
			for candidateCount > maxLabelCount {
				step++
				candidateCount = 2 + (dataCount-2)/step
			}
			labelCount = chartdraw.MaxInt(candidateCount, minimumAxisLabels)
		}
	}
	// ensure there are not too many ticks, we want them relative and related to the label positions
	tickCount := dataCount
	if tickCount > labelCount*2 {
		// it's difficult to choose a tick count that allows multiple ticks while staying lined up with the labels
		// TODO - I would like to improve this, but for simplicity we will match the label count if ticks are too dense
		tickCount = labelCount
	}

	return axisRange{
		isCategory:     true,
		labels:         labels,
		dataStartIndex: dataStartIndex,
		divideCount:    dataCount,
		tickCount:      tickCount,
		labelCount:     labelCount,
		size:           axisSize,
		textMaxWidth:   textW,
		textMaxHeight:  textH,
		labelRotation:  labelRotation,
		labelFontStyle: fontStyle,
	}
}

func valueLabels(labelsCfg []string, valueFormatter ValueFormatter, min, max float64, labelCount int) []string {
	labels := make([]string, labelCount)
	offset := (max - min) / float64(labelCount-1)
	for i := range labels {
		if i < len(labelsCfg) {
			labels[i] = labelsCfg[i]
		} else {
			labels[i] = valueFormatter(min + float64(i)*offset)
		}
	}
	return labels
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
	// else, just adjust based off default multiplier
	if add {
		return val + (increment * defaultMultiplier), defaultMultiplier
	} else {
		return val - (increment * defaultMultiplier), defaultMultiplier
	}
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

// getRange returns a range at a given index.
func (r axisRange) getRange(index int) (float64, float64) {
	unit := float64(r.size) / float64(r.divideCount)
	return unit * float64(index), unit * float64(index+1)
}

// autoDivide divides the axis size by the configured count.
func (r axisRange) autoDivide() []int {
	return autoDivide(r.size, r.divideCount)
}

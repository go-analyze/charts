package charts

import (
	"math"
	"strconv"

	"github.com/go-analyze/charts/chartdraw"
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
	labels []string, dataStartIndex int,
	labelCountCfg int, labelUnit float64, labelCountAdjustment int,
	seriesList seriesList, yAxisIndex int, stackSeries bool,
	valueFormatter ValueFormatter,
	labelRotation float64, fontStyle FontStyle) axisRange {
	// calculate the range
	minVal, maxVal, sumMax := getSeriesMinMaxSumMax(seriesList, yAxisIndex, stackSeries)
	if stackSeries { // If stacked, maxVal should be the maxVal data point of all series summed together
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
	decimalData := minVal != math.Floor(minVal) || (maxVal-minVal) != math.Floor(maxVal-minVal) // TODO - does this need to be done after padding

	// Label counts and y-axis padding are linked together to produce a user-friendly graph.
	// First when considering padding we want to prefer a zero axis start if reasonable, and add a slight
	// padding to the maxVal so there is a little space at the top of the graph. In addition, we want to pick
	// a maxVal value that will result in round intervals on the axis. These details are in range.go.
	// But in order to produce round intervals we need to have an idea of how many intervals there are.
	// In addition, if the user specified a `Unit` value we may need to adjust our label count calculation
	// based on the padded range.
	//
	// In order to accomplish this, we estimate the label count (if necessary), pad the range, then precisely
	// calculate the label count.
	padLabelCount := labelCountCfg
	if padLabelCount < 1 {
		if labelUnit > 0 {
			padLabelCount = int((maxVal-minVal)/labelUnit) + 1
		} else {
			padLabelCount = defaultYAxisLabelCount(maxVal-minVal, decimalData)
		}
	}
	padLabelCount = chartdraw.MaxInt(padLabelCount+labelCountAdjustment, minimumAxisLabels)
	// we call padRange directly because we need to do this padding before we can calculate the final labelCount for the axisRange
	minPadded, maxPadded := padRange(padLabelCount, minVal, maxVal, minPadScale, maxPadScale)
	labelCount := labelCountCfg
	if labelCount == 0 { // If user never explicitly set labelCount, refine again after padding
		if labelUnit > 0 {
			// Round up max to a multiple of the unit
			if maxCfg == nil {
				maxPadded = math.Trunc(math.Ceil(maxPadded/labelUnit) * labelUnit)
			}
			// set labelCount so that unit is respected
			labelCount = int((maxPadded-minPadded)/labelUnit) + 1
		} else {
			labelCount = defaultYAxisLabelCount(maxPadded-minPadded, decimalData)
		}
	}
	if labelCountAdjustment != 0 {
		if labelUnit > 0 {
			// in order to maintain unit we need to expand the scale range to add the labels
			for i := 0; i < labelCountAdjustment; i++ {
				if i%2 == 0 || minPadded == 0 {
					maxPadded += labelUnit
				} else {
					minPadded -= labelUnit
				}
			}
		} else {
			labelCount = chartdraw.MaxInt(labelCount+labelCountAdjustment, minimumAxisLabels)
		}
	}

	// ensure labels are set for the full range is met
	if len(labels) < labelCount {
		offset := (maxPadded - minPadded) / float64(labelCount-1)
		for i := len(labels); i < labelCount; i++ {
			labels = append(labels, valueFormatter(minPadded+float64(i)*offset))
		}
	}

	textW, textH := p.measureTextMaxWidthHeight(labels, labelRotation, fontStyle)

	// If user gave an explicit .LabelCount, then we do NOT do a collision check.
	// For default logic we want to make sure we choose a label count that is visually appealing
	if labelCountCfg == 0 {
		if maxLabels := calculateLabelFixMax(isVertical, labelUnit, labelCount, axisSize, textW, textH); maxLabels < labelCount {
			labelCount = maxLabels
		}
	}

	return axisRange{
		isCategory:     false,
		labels:         labels,
		dataStartIndex: dataStartIndex, // TODO - ensure implemented correct
		divideCount:    len(labels),
		labelCount:     labelCount,
		min:            minPadded,
		max:            maxPadded,
		size:           axisSize,
		textMaxWidth:   textW,
		textMaxHeight:  textH,
		labelRotation:  labelRotation,
		labelFontStyle: fontStyle,
	}
}

// calculateCategoryAxisRange does the same for category axes (common for X-axis in line/bar charts).
func calculateCategoryAxisRange(p *Painter, axisSize int, isVertical bool,
	labels []string, dataStartIndex int,
	labelCountCfg int, labelCountAdjustment int, labelUnit float64,
	seriesList seriesList, labelRotation float64, fontStyle FontStyle) axisRange {
	// If user provided no labels, use series names.
	// If provided only partially, fill in the remaining labels.
	if len(labels) == 0 {
		labels = seriesList.names()
	} else {
		for i := len(labels); i < seriesList.len(); i++ {
			seriesName := seriesList.getSeriesName(i)
			if seriesName == "" {
				seriesName = strconv.Itoa(i + 1) // TODO - use value formatter?
			}
			labels = append(labels, seriesName)
		}
	}
	dataCount := len(labels)

	textW, textH := p.measureTextMaxWidthHeight(labels, labelRotation, fontStyle)

	labelCount := labelCountCfg
	if labelCount <= 0 {
		labelCount = dataCount
	} else if labelCount > dataCount {
		labelCount = dataCount
	}
	labelCount += labelCountAdjustment
	if labelCount < minimumAxisLabels {
		labelCount = minimumAxisLabels
	}
	if labelCountCfg == 0 {
		if maxLabels := calculateLabelFixMax(isVertical, labelUnit, dataCount, axisSize, textW, textH); maxLabels < labelCount {
			labelCount = maxLabels
		}
	}

	return axisRange{
		isCategory:     true,
		labels:         labels,
		dataStartIndex: dataStartIndex, // TODO - ensure implemented correct
		divideCount:    dataCount,
		labelCount:     labelCount,
		size:           axisSize,
		textMaxWidth:   textW,
		textMaxHeight:  textH,
		labelRotation:  labelRotation,
		labelFontStyle: fontStyle,
	}
}

func defaultYAxisLabelCount(span float64, decimalData bool) int {
	result := chartdraw.MinInt(chartdraw.MaxInt(int(span)+1, defaultYAxisLabelCountLow), defaultYAxisLabelCountHigh)
	if decimalData {
		// if there is a decimal, we double our labels to provide more detail
		result = chartdraw.MinInt(result*2, defaultYAxisLabelCountHigh)
	}
	return result
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

func calculateLabelFixMax(isVertical bool, unit float64, dataCount, space, textW, textH int) int {
	// set minimum values to avoid divide by zero
	textW = chartdraw.MaxInt(textW, 1)
	textH = chartdraw.MaxInt(textH, 1)

	var maxLabelCount int
	if isVertical {
		maxLabelCount = space / textH
	} else {
		// add a little extra padding for horizontal layouts
		maxLabelCount = space / (textW + 10)
	}
	if maxLabelCount < minimumAxisLabels {
		return minimumAxisLabels // required to prevent infinite loop if less than zero
	}

	if unit > 0 {
		// If the user gave a 'unit', figure out how many 'units' fit
		multiplier := 1.0
		for {
			count := ceilFloatToInt(float64(dataCount) / (unit * multiplier))
			if count > maxLabelCount {
				multiplier++
			} else {
				return count
			}
		}
	}
	// TODO - check if a small adjustment allows for a better unit
	return maxLabelCount
}

// Values returns values of range.
// For a value axis, if no custom labels are set, it uses the valueFormatter to compute friendly numbers.
// For a category axis, it returns the stored labels.
func (r axisRange) Values() []string {
	// TODO - copy necessary?
	result := make([]string, len(r.labels))
	copy(result, r.labels)
	return result
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

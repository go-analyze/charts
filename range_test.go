package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRange(size, divideCount int, min, max, minPaddingScale, maxPaddingScale float64) axisRange {
	min, max = padRange(divideCount, min, max, minPaddingScale, maxPaddingScale)
	return axisRange{
		divideCount: divideCount,
		min:         min,
		max:         max,
		size:        size,
	}
}

func newTestRangeForLabels(labels []string, rotation float64, style FontStyle) axisRange {
	p := NewPainter(PainterOptions{})
	style = fillFontStyleDefaults(style, defaultFontSize, ColorBlack)
	width, height := p.measureTextMaxWidthHeight(labels, rotation, style)
	return axisRange{
		isCategory:     true,
		labels:         labels,
		divideCount:    len(labels),
		tickCount:      len(labels),
		labelCount:     len(labels),
		size:           800,
		textMaxWidth:   width,
		textMaxHeight:  height,
		labelFontStyle: style,
		labelRotation:  rotation,
	}
}

// testSeries implements the series interface.
type testSeries struct {
	yAxisIndex int
	values     []float64
}

func (s testSeries) getType() string {
	return "fake"
}

func (s testSeries) getYAxisIndex() int {
	return s.yAxisIndex
}

func (s testSeries) getValues() []float64 {
	return s.values
}

// testSeriesList implements the seriesList interface.
type testSeriesList []testSeries

func (tsl testSeriesList) len() int {
	return len(tsl)
}

func (tsl testSeriesList) getSeries(index int) series {
	return tsl[index]
}

func (tsl testSeriesList) getSeriesName(index int) string {
	return "series:" + strconv.Itoa(index)
}

func (tsl testSeriesList) getSeriesValues(index int) []float64 {
	return tsl[index].values
}

func (tsl testSeriesList) getSeriesLen(i int) int {
	return len(tsl[i].values)
}

func (tsl testSeriesList) names() []string {
	result := make([]string, tsl.len())
	for i := range result {
		result[i] = tsl.getSeriesName(i)
	}
	return result
}

func (tsl testSeriesList) hasMarkPoint() bool {
	return false
}

func (tsl testSeriesList) setSeriesName(_ int, _ string) {
	panic("not implemented")
}

func (tsl testSeriesList) sortByNameIndex(_ map[string]int) {
	panic("not implemented")
}

func (tsl testSeriesList) getSeriesSymbol(_ int) Symbol {
	panic("not implemented")
}

func TestCalculateValueAxisRange(t *testing.T) {
	fs := FontStyle{FontSize: 16, FontColor: ColorGray}

	t.Run("label_count", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, false, 800, nil, nil, Ptr(0.0),
			nil, 0, 3, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Len(t, ar.labels, 3)
		assert.Equal(t, []string{"10", "20", "30"}, ar.labels)
		assert.Equal(t, 3, ar.divideCount)
		assert.Equal(t, 3, ar.labelCount)
	})

	t.Run("label_unit", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{0, 50}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 5, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, 12, ar.labelCount)
		assert.Equal(t, []string{"0", "5", "10", "15", "20", "25", "30", "35", "40", "45", "50", "55"}, ar.labels)
	})

	t.Run("label_unit_adjusted_positive", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 1200, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 100}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, false, 1200, nil, nil, nil,
			nil, 0, 0, 5, 2,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, 8, ar.labelCount)
		assert.InDelta(t, 0.0, ar.min, 0.0)
		assert.InDelta(t, 105.0, ar.max, 0.0)
	})

	t.Run("label_unit_adjusted_negative", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 2400, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{-10, 100}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, false, 2400, nil, nil, nil,
			nil, 0, 0, 5, 4,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, 25, ar.labelCount)
		assert.InDelta(t, -10.0, ar.min, 0.0)
		assert.InDelta(t, 110.0, ar.max, 0.0)
	})

	t.Run("stacked_series", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{
			{values: []float64{1, 2, 3}},
			{values: []float64{4, 5, 6}},
		}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, true, defaultValueFormatter, 0, fs, nil)

		assert.InDelta(t, 0.0, ar.min, 0.0)
		assert.InDelta(t, 10.0, ar.max, 0.0)
	})

	t.Run("min_max_set", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20}}
		tsl := testSeriesList{series}

		min := Ptr(5.0)
		max := Ptr(25.0)
		ar := calculateValueAxisRange(p, true, 800, min, max,
			nil, []string{}, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.InDelta(t, 5.0, ar.min, 0.0)
		assert.InDelta(t, 25.0, ar.max, 0.0)
	})

	t.Run("decimal_range", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{1.1, 2.2, 3.3}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.InDelta(t, 1.0, ar.min, 0.0)
		assert.InDelta(t, 5.0, ar.max, 0.0)
		assert.Equal(t, 6, ar.labelCount)
	})

	t.Run("long_horizontal_labels", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series}

		fs := FontStyle{FontSize: 28, FontColor: ColorGray}
		inputLabels := []string{"ThisIsAVeryLongLabelThatExceedsNormal", "AnotherVeryLongLabelThatExceedsNormal",
			"WowLookAtTheseLabels!", "AndHereIsAnotherReallyLongLabel"}
		ar := calculateValueAxisRange(p, false, 800, nil, nil, nil,
			inputLabels, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, 810, ar.textMaxWidth)
		assert.Equal(t, 41, ar.textMaxHeight)
		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("zero_span_nonzero", func(t *testing.T) {
		// Series with a single nonzero value should trigger the zero–span adjustment.
		// When the data point is nonzero, we expect: min = value - zeroSpanAdjustment, max = value + zeroSpanAdjustment.
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{50}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.InDelta(t, 49.0, ar.min, 0.0)
		assert.InDelta(t, 51.0, ar.max, 0.0)
	})

	t.Run("vertical_label_rotation", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series}

		rotation := DegreesToRadians(45.0)
		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			[]string{"Label One", "Label Two", "Label Three", "Label Four"}, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, rotation, fs, nil)

		assert.Equal(t, 103, ar.textMaxWidth)
		assert.Equal(t, 103, ar.textMaxHeight)
		assert.InDelta(t, rotation, ar.labelRotation, 0.0)
	})

	t.Run("provided_labels_excess", func(t *testing.T) {
		// If the caller supplies more labels than the explicit labelCount,
		// the provided labels should be distributed across the range
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{5, 15, 25}}
		tsl := testSeriesList{series}

		providedLabels := []string{"Label1", "Label2", "Label3", "Label4", "Label5"}
		explicitLabelCount := 3
		ar := calculateValueAxisRange(p, false, 800, nil, nil, nil,
			providedLabels, 0, explicitLabelCount, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, []string{"Label1", "Label2", "Label3"}, ar.labels)
		assert.Equal(t, 3, ar.divideCount)
		assert.Equal(t, explicitLabelCount, ar.labelCount)
	})

	t.Run("label_unit_search_loop_zero_origin", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{0, 100}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, false, 800,
			nil, nil, Ptr(0.0), // force no padding
			nil, 0, 0, 7, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, 6, ar.labelCount)
		assert.InDelta(t, 0.0, ar.min, 0.0)
		assert.InDelta(t, 105.0, ar.max, 0.0)
		assert.Equal(t, []string{"0", "21", "42", "63", "84", "105"}, ar.labels)
	})

	t.Run("label_unit_search_loop_nonzero_min", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{9, 30}}
		tsl := testSeriesList{series}

		ar := calculateValueAxisRange(p, false, 800,
			nil, nil, Ptr(0.0), // force no padding
			nil, 0, 0, 9, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)

		assert.Equal(t, 4, ar.labelCount)
		assert.InDelta(t, 9.0, ar.min, 0.0)
		assert.InDelta(t, 36.0, ar.max, 0.0)
		assert.Equal(t, []string{"9", "18", "27", "36"}, ar.labels)
	})

	t.Run("prefer_nice_intervals", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{0, 50}}
		tsl := testSeriesList{series}

		// without PreferNiceIntervals the flex logic is not triggered
		arDefault := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, nil)
		// with PreferNiceIntervals the flex logic produces nicer intervals
		arNice := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs, Ptr(true))

		assert.NotEqual(t, arDefault.max, arNice.max)
		// verify the nice interval produces a round interval
		interval := (arNice.max - arNice.min) / float64(arNice.labelCount-1)
		niceInterval := niceNum(interval)
		assert.InDelta(t, niceInterval, interval, 1e-10)
	})

	t.Run("label_unit_infinite_loop", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 462, Height: 400})

		tsl := testSeriesList{
			{values: []float64{20, 46}},
		}

		ar := calculateValueAxisRange(
			p, false, 462, // isVertical, axisSize
			nil, nil, nil, nil, // minCfg, maxCfg, rangeValuePaddingScale, labelsCfg
			0,            // dataStartIndex
			0,            // labelCountCfg
			100000,       // labelUnit (much larger than the data span)
			0,            // labelCountAdjustment
			tsl, 0, true, // seriesList, yAxisIndex, stackSeries
			defaultValueFormatter,
			0, fs, // labelRotation, fontStyle
			nil, // preferNiceIntervals
		)

		assert.Equal(t, 2, ar.labelCount)
		assert.InDelta(t, 19.0, ar.min, 0.0)
		assert.InDelta(t, 49, ar.max, 0.0)
		assert.Equal(t, []string{"19", "49"}, ar.labels)
	})
}

func TestCalculateCategoryAxisRange(t *testing.T) {
	fs := FontStyle{FontSize: 16, FontColor: ColorGray}

	t.Run("no_labels_provided_uses_series_names", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}

		ar := calculateCategoryAxisRange(p, 800, false, false, nil, 0,
			0, 0, 0, tsl, 0, fs)

		expectedLabels := []string{"1"}
		assert.Equal(t, expectedLabels, ar.labels)
		assert.Equal(t, 1, ar.divideCount)
		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("provided_labels_filled_to_series_length", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		// Provide fewer labels than the number of series.
		providedLabels := []string{"CustomLabel"}
		tsl := testSeriesList{
			{values: []float64{1, 1}},
			{values: []float64{2, 1}},
			{values: []float64{3, 1}},
		}

		ar := calculateCategoryAxisRange(p, 800, false, false, providedLabels, 0,
			0, 0, 0, tsl, 0, fs)

		assert.Equal(t, []string{"CustomLabel", "2"}, ar.labels)
		assert.Equal(t, 2, ar.divideCount)
		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("explicit_label_count_cfg", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
			{values: []float64{4}},
		}

		ar := calculateCategoryAxisRange(p, 800, false, false, nil, 0,
			2, 1, 0, tsl, 0, fs)

		assert.Equal(t, 1, ar.divideCount)
		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("label_rotation", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
		}

		rotation := DegreesToRadians(30.0)
		ar := calculateCategoryAxisRange(p, 800, true, false, []string{}, 0,
			0, 0, 0, tsl, rotation, fs)

		assert.Equal(t, 17, ar.textMaxWidth)
		assert.Equal(t, 20, ar.textMaxHeight)
		assert.InDelta(t, rotation, ar.labelRotation, 0.0)
	})

	t.Run("negative_label_count_adjustment", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}

		ar := calculateCategoryAxisRange(p, 800, false, false, []string{}, 0,
			0, -2, 0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("label_count_exceeds_series_count", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
		}

		ar := calculateCategoryAxisRange(p, 800, false, false, []string{}, 0,
			5, 0, 0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("long_horizontal_labels", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}

		inputLabels := []string{"ThisIsAVeryLongLabelThatExceedsNormal", "AnotherVeryLongLabelThatExceedsNormal",
			"WowLookAtTheseLabels!", "AndHereIsAnotherReallyLongLabel"}
		ar := calculateCategoryAxisRange(p, 600, false, false, inputLabels, 0,
			0, 0, 0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("label_unit", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
			{values: []float64{4}},
			{values: []float64{5}},
			{values: []float64{6}},
			{values: []float64{7}},
			{values: []float64{8}},
			{values: []float64{9}},
			{values: []float64{10}},
		}

		ar := calculateCategoryAxisRange(p, 800, false, false, []string{}, 0,
			0, 0, 4.0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("empty_series_list", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{}
		ar := calculateCategoryAxisRange(p, 800, false, false, nil, 0,
			0, 0, 0, tsl, 0, fs)

		assert.Empty(t, ar.labels)
		assert.Equal(t, 0, ar.divideCount)
		assert.Equal(t, 2, ar.labelCount)
	})
}

func TestNiceNum(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"zero", 0, 0},
		{"negative", -5, 0},
		{"exact_one", 1, 1},
		{"exact_two", 2, 2},
		{"exact_two_point_five", 2.5, 2.5},
		{"exact_five", 5, 5},
		{"exact_ten", 10, 10},
		{"exact_twenty_five", 25, 25},
		{"sub_unit_small", 0.07, 0.1},
		{"sub_unit_mid", 0.3, 0.5},
		{"just_above_one", 0.7, 1.0},
		{"three", 3, 5},
		{"seven", 7, 10},
		{"fifteen_point_six", 15.6, 20},
		{"twenty_one", 21, 25},
		{"large_seven_hundred", 700, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, niceNum(tc.input), 1e-10)
		})
	}
}

func TestNiceNumFrom(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"three", 3, 3},
		{"three_point_five", 3.5, 4},
		{"four_point_five", 4.5, 5},
		{"five_point_five", 5.5, 6},
		{"seven", 7, 8},
		{"eight_point_five", 8.5, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, niceNumFrom(tc.input, extendedNiceNums[:]), 1e-10)
		})
	}
}

func TestUnitAlignmentTier(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		interval float64
		unit     float64
		expected int
	}{
		{"interval_multiple_of_unit", 80, 40, 0},
		{"interval_equals_unit", 40, 40, 0},
		{"unit_multiple_of_interval", 20, 60, 1},
		{"no_alignment", 33, 40, 2},
		{"zero_unit", 40, 0, 2},
		{"zero_interval", 0, 40, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, unitAlignmentTier(tc.interval, tc.unit))
		})
	}
}

func TestPadRange(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		expectedMinValue float64
		expectedMaxValue float64
		minValue         float64
		maxValue         float64
		labelCount       int
	}{
		{
			name:             "pad_max_only",
			expectedMinValue: 0.0,
			expectedMaxValue: 10.5,
			minValue:         0.0,
			maxValue:         10.0,
			labelCount:       10,
		},
		{
			name:             "pad_min_to_zero",
			expectedMinValue: 0.0,
			expectedMaxValue: 21.0,
			minValue:         1.0,
			maxValue:         20.0,
			labelCount:       10,
		},
		{
			name:             "pad_negative_min_positive_max",
			expectedMinValue: -5.0,
			expectedMaxValue: 12.0,
			minValue:         -3.0,
			maxValue:         10.0,
			labelCount:       10,
		},
		{
			name:             "pad_negative_min_negative_max",
			expectedMinValue: -20.0,
			expectedMaxValue: -9.0,
			minValue:         -20.0,
			maxValue:         -10.0,
			labelCount:       10,
		},
		{
			name:             "pad_positive_min_positive_max",
			expectedMinValue: 100.0,
			expectedMaxValue: 214.0,
			minValue:         100.0,
			maxValue:         200.0,
			labelCount:       20,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			min, max := padRange(tc.labelCount, tc.minValue, tc.maxValue, 1.0, 1.0)

			assert.InDelta(t, tc.expectedMinValue, min, 0, "Unexpected value rounding %v", tc.minValue)
			assert.InDelta(t, tc.expectedMaxValue, max, 0, "Unexpected value rounding %v", tc.maxValue)
		})
	}

	t.Run("flex_capped_by_baseline", func(t *testing.T) {
		// niceNum jumps from 2.5→5 inflating the range; flex should be capped by the friendlyRound baseline
		prep := valueAxisPrep{
			minVal: 1, maxVal: 26, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 5, maxLabelCount: 5,
		}
		_, flexMax, _ := resolveValueAxisRange(&prep, true, 0)
		_, baselineMax, _ := resolveValueAxisRange(&prep, false, 0)
		// the flex result should not significantly exceed the baseline
		assert.LessOrEqual(t, flexMax, baselineMax+((baselineMax-26)*1.0)+1e-10)
	})

	t.Run("tier2_flex_fallback", func(t *testing.T) {
		// range where tier-1 niceNums (1,2,2.5,5) can't find a match in ±1 delta,
		// but tier-2 extended set can find one at the original divideCount
		prep := valueAxisPrep{
			minVal: 0, maxVal: 340, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 8, maxLabelCount: 8,
		}
		_, flexMax, flexCount := resolveValueAxisRange(&prep, true, 0)
		interval := flexMax / float64(flexCount-1)
		// the flex path should produce a clean interval
		assert.InDelta(t, math.Round(interval), interval, 0.5)
		// verify reasonable range
		assert.GreaterOrEqual(t, flexMax, 340.0)
	})
}

func TestFriendlyRound(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedValue float64
		value         float64
		minMultiplier float64
		maxMultiplier float64
		add           bool
	}{
		{
			name:          "original_zero_sub",
			expectedValue: 0.0,
			value:         0.0,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "original_zero_add",
			expectedValue: 0.0,
			value:         0.0,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "round_fraction_sub",
			expectedValue: -2.0,
			value:         -1.2,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "round_fraction_add",
			expectedValue: 2.0,
			value:         1.2,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "round_very_close_to_zero_sub",
			expectedValue: -1.0,
			value:         -0.01,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "round_very_close_to_zero_add",
			expectedValue: 0.0,
			value:         -0.01,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "original_large_sub",
			expectedValue: 1337,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "original_large_add",
			expectedValue: 1337,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "round_thousand_large_sub",
			expectedValue: 1000,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           false,
		},
		{
			name:          "round_thousand_large_add",
			expectedValue: 2000,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           true,
		},
		{
			name:          "round_hundred_large_sub",
			expectedValue: 1300,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
		{
			name:          "round_hundred_large_add",
			expectedValue: 1400,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           true,
		},
		{
			name:          "round_negative_small_sub",
			expectedValue: -1.0,
			value:         -0.5,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "round_halfway_point_sub",
			expectedValue: 100.0,
			value:         150.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
		{
			name:          "round_halfway_point_add",
			expectedValue: 200.0,
			value:         150.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           true,
		},
		{
			name:          "round_thousands_negative_large_sub",
			expectedValue: -2000.0,
			value:         -1337.0,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           false,
		},
		{
			name:          "round_hundreds_negative_large_sub",
			expectedValue: -1400.0,
			value:         -1337.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, _ := friendlyRound(tc.value, 1.0, 0.0,
				tc.minMultiplier, tc.maxMultiplier, tc.add)

			assert.InDelta(t, tc.expectedValue, val, 0, "Unexpected value rounding %v", tc.value)
		})
	}
}

func TestResolveValueAxisRange(t *testing.T) {
	t.Run("basic_no_flex", func(t *testing.T) {
		prep := valueAxisPrep{
			minVal: 10, maxVal: 90, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 5, maxLabelCount: 10,
		}
		minP, maxP, count := resolveValueAxisRange(&prep, false, 0)
		// should match padRange baseline
		expMin, expMax := padRange(5, 10, 90, 1.0, 1.0)
		assert.InDelta(t, expMin, minP, 1e-10)
		assert.InDelta(t, expMax, maxP, 1e-10)
		assert.Equal(t, 5, count)
	})

	t.Run("target_overrides_count", func(t *testing.T) {
		prep := valueAxisPrep{
			minVal: 0, maxVal: 100, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 5, maxLabelCount: 10,
		}
		_, _, count := resolveValueAxisRange(&prep, true, 8)
		assert.Equal(t, 8, count)
	})

	t.Run("flex_produces_nice_interval", func(t *testing.T) {
		prep := valueAxisPrep{
			minVal: 0, maxVal: 47, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 6, maxLabelCount: 10,
		}
		minP, maxP, count := resolveValueAxisRange(&prep, true, 0)
		require.Greater(t, count, 1)
		interval := (maxP - minP) / float64(count-1)
		assert.InDelta(t, niceNum(interval), interval, 1e-10)
		assert.GreaterOrEqual(t, maxP, 47.0)
		assert.LessOrEqual(t, minP, 0.0)
	})

	t.Run("unit_snaps_range", func(t *testing.T) {
		prep := valueAxisPrep{
			minVal: 0, maxVal: 50, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 6, maxLabelCount: 15,
			labelUnit: 10,
		}
		minP, maxP, count := resolveValueAxisRange(&prep, false, 0)
		require.Greater(t, count, 1)
		interval := (maxP - minP) / float64(count-1)
		// interval should be a multiple of labelUnit
		ratio := interval / 10
		assert.InDelta(t, math.Round(ratio), ratio, 1e-9)
		assert.GreaterOrEqual(t, maxP, 50.0)
	})

	t.Run("unit_with_target_count", func(t *testing.T) {
		prep := valueAxisPrep{
			minVal: 0, maxVal: 100, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 6, maxLabelCount: 15,
			labelUnit: 10,
		}
		_, _, count := resolveValueAxisRange(&prep, false, 6)
		assert.Equal(t, 6, count)
	})

	t.Run("unit_exceeds_span", func(t *testing.T) {
		prep := valueAxisPrep{
			minVal: 10, maxVal: 20, minPadScale: 1.0, maxPadScale: 1.0,
			padLabelCount: 5, maxLabelCount: 10,
			labelUnit: 100,
		}
		_, _, count := resolveValueAxisRange(&prep, false, 0)
		assert.Equal(t, minimumAxisLabels, count)
	})

	t.Run("unit_with_min_cfg", func(t *testing.T) {
		minCfg := 0.0
		prep := valueAxisPrep{
			minVal: 0, maxVal: 100, minPadScale: 0.0, maxPadScale: 1.0,
			padLabelCount: 6, maxLabelCount: 15,
			labelUnit: 20, minCfg: &minCfg,
		}
		minP, maxP, count := resolveValueAxisRange(&prep, false, 0)
		// min should stay at or near 0 since minCfg blocks shifting downward
		assert.InDelta(t, 0.0, minP, 1e-10)
		require.Greater(t, count, 1)
		interval := (maxP - minP) / float64(count-1)
		ratio := interval / 20
		assert.InDelta(t, math.Round(ratio), ratio, 1e-9)
	})

	t.Run("unit_with_max_cfg", func(t *testing.T) {
		maxCfg := 120.0
		prep := valueAxisPrep{
			minVal: 0, maxVal: 100, minPadScale: 1.0, maxPadScale: 0.0,
			padLabelCount: 6, maxLabelCount: 15,
			labelUnit: 20, maxCfg: &maxCfg,
		}
		minP, maxP, count := resolveValueAxisRange(&prep, false, 0)
		require.Greater(t, count, 1)
		interval := (maxP - minP) / float64(count-1)
		ratio := interval / 20
		assert.InDelta(t, math.Round(ratio), ratio, 1e-9)
	})
}

func TestCoordinateValueAxisRanges(t *testing.T) {
	fs := FontStyle{FontSize: 16, FontColor: ColorGray}

	t.Run("empty_preps", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		results := coordinateValueAxisRanges(p, nil)
		assert.Nil(t, results)
	})

	t.Run("single_axis", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeriesList{{values: []float64{0, 100}}}
		prep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			series, 0, false, defaultValueFormatter, 0, fs, Ptr(true))
		results := coordinateValueAxisRanges(p, []*valueAxisPrep{&prep})

		assert.Len(t, results, 1)
		assert.GreaterOrEqual(t, results[0].max, 100.0)
	})

	t.Run("forced_count", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		leftSeries := testSeriesList{{values: []float64{0, 100}}}
		leftPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 5, 0, 0,
			leftSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))
		rightSeries := testSeriesList{{values: []float64{0, 200}}}
		rightPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			rightSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))

		results := coordinateValueAxisRanges(p, []*valueAxisPrep{&leftPrep, &rightPrep})

		assert.Len(t, results, 2)
		assert.Equal(t, 5, results[0].labelCount)
		assert.Equal(t, 5, results[1].labelCount)
	})

	t.Run("conflicting_counts", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		leftSeries := testSeriesList{{values: []float64{0, 100}}}
		leftPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 5, 0, 0,
			leftSeries, 0, false, defaultValueFormatter, 0, fs, nil)
		rightSeries := testSeriesList{{values: []float64{0, 200}}}
		rightPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 8, 0, 0,
			rightSeries, 0, false, defaultValueFormatter, 0, fs, nil)

		results := coordinateValueAxisRanges(p, []*valueAxisPrep{&leftPrep, &rightPrep})

		assert.Len(t, results, 2)
		// resolved independently, each gets its own count
		assert.Equal(t, 5, results[0].labelCount)
		assert.Equal(t, 8, results[1].labelCount)
	})

	t.Run("no_secondary_nice", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		leftSeries := testSeriesList{{values: []float64{0, 100}}}
		leftPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			leftSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))
		rightSeries := testSeriesList{{values: []float64{0, 200}}}
		rightPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			rightSeries, 0, false, defaultValueFormatter, 0, fs, nil)

		results := coordinateValueAxisRanges(p, []*valueAxisPrep{&leftPrep, &rightPrep})

		assert.Len(t, results, 2)
		// secondary adopts primary's count
		assert.Equal(t, results[0].labelCount, results[1].labelCount)
	})

	t.Run("natural_counts_match", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		// use similar data ranges so natural counts are likely the same
		leftSeries := testSeriesList{{values: []float64{0, 100}}}
		leftPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			leftSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))
		rightSeries := testSeriesList{{values: []float64{0, 100}}}
		rightPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			rightSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))

		results := coordinateValueAxisRanges(p, []*valueAxisPrep{&leftPrep, &rightPrep})

		assert.Len(t, results, 2)
		assert.Equal(t, results[0].labelCount, results[1].labelCount)
	})

	t.Run("dual_axis_with_unit", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		// left axis: data 0-250, no unit
		leftSeries := testSeriesList{{values: []float64{0, 250}}}
		leftPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			leftSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))
		// right axis: data 0-680, unit=40
		rightSeries := testSeriesList{{values: []float64{0, 680}}}
		rightPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 40, 0,
			rightSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))

		preps := []*valueAxisPrep{&leftPrep, &rightPrep}
		results := coordinateValueAxisRanges(p, preps)

		assert.Len(t, results, 2)
		// right axis interval should be compatible with unit=40
		rightRange := results[1]
		require.Greater(t, rightRange.labelCount, 1)
		interval := (rightRange.max - rightRange.min) / float64(rightRange.labelCount-1)
		tier := unitAlignmentTier(interval, 40)
		assert.LessOrEqual(t, tier, 1, "interval %.1f should align with unit 40", interval)
		// both axes should have matching label counts
		assert.Equal(t, results[0].labelCount, results[1].labelCount)
	})

	t.Run("dual_axis_without_unit", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		// both axes without units — existing scoring unchanged
		leftSeries := testSeriesList{{values: []float64{0, 100}}}
		leftPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			leftSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))
		rightSeries := testSeriesList{{values: []float64{0, 200}}}
		rightPrep := prepareValueAxisRange(p, true, 500,
			nil, nil, nil, nil, 0, 0, 0, 0,
			rightSeries, 0, false, defaultValueFormatter, 0, fs, Ptr(true))

		preps := []*valueAxisPrep{&leftPrep, &rightPrep}
		results := coordinateValueAxisRanges(p, preps)

		assert.Len(t, results, 2)
		// both should have the same label count
		assert.Equal(t, results[0].labelCount, results[1].labelCount)
		// intervals should be nice numbers
		for _, r := range results {
			require.Greater(t, r.labelCount, 1)
			interval := (r.max - r.min) / float64(r.labelCount-1)
			ni := niceNum(interval)
			assert.InDelta(t, ni, interval, 1e-10, "interval should be a nice number")
		}
	})
}

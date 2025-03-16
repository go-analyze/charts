package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

func newTestRangeForLabels(labels []string, style FontStyle) axisRange {
	p := NewPainter(PainterOptions{})
	style = fillFontStyleDefaults(style, defaultFontSize, ColorBlack)
	width, height := p.measureTextMaxWidthHeight(labels, 0, style)
	return axisRange{
		isCategory:     true,
		labels:         labels,
		divideCount:    len(labels),
		labelCount:     len(labels),
		size:           800,
		textMaxWidth:   width,
		textMaxHeight:  height,
		labelFontStyle: style,
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
type testSeriesList struct {
	series []testSeries
}

func (tsl testSeriesList) len() int {
	return len(tsl.series)
}

func (tsl testSeriesList) getSeries(index int) series {
	return tsl.series[index]
}

func (tsl testSeriesList) getSeriesName(index int) string {
	return "series:" + strconv.Itoa(index)
}

func (tsl testSeriesList) getSeriesValues(index int) []float64 {
	return tsl.series[index].values
}

func (tsl testSeriesList) getSeriesLen(i int) int {
	return len(tsl.series[i].values)
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

func (tsl testSeriesList) setSeriesName(index int, name string) {
	panic("not implemented")
}

func (tsl testSeriesList) sortByNameIndex(dict map[string]int) {
	panic("not implemented")
}

func (tsl testSeriesList) getSeriesSymbol(index int) Symbol {
	panic("not implemented")
}

func TestCalculateValueAxisRange(t *testing.T) {
	fs := FontStyle{FontSize: 16}

	t.Run("lable_rotation_and_font_defaults", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series: []testSeries{series}}
		fs := FontStyle{} // Missing Font and FontSize.

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 5, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.NotNil(t, ar.labelFontStyle.Font)
		assert.InDelta(t, defaultFontSize, ar.labelFontStyle.FontSize, 0.0)
	})

	t.Run("label_count", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series: []testSeries{series}}

		ar := calculateValueAxisRange(p, false, 800, nil, nil, Ptr(0.0),
			nil, 0, 3, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.Len(t, ar.labels, 3)
		assert.Equal(t, []string{"10", "20", "30"}, ar.labels)
		assert.Equal(t, 3, ar.divideCount)
		assert.Equal(t, 3, ar.labelCount)
	})

	t.Run("label_unit", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{0, 50}}
		tsl := testSeriesList{series: []testSeries{series}}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 5, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.Equal(t, 3, ar.labelCount)
		assert.Equal(t, []string{"0", "5", "10", "15", "20", "25", "30", "35", "40", "45", "50", "55"}, ar.labels)
	})

	t.Run("label_unit_adjusted_positive", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 100}}
		tsl := testSeriesList{series: []testSeries{series}}

		ar := calculateValueAxisRange(p, false, 800, nil, nil, nil,
			nil, 0, 0, 5, 2,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.Equal(t, 5, ar.labelCount)
		assert.InDelta(t, 0.0, ar.min, 0.0)
		assert.InDelta(t, 120.0, ar.max, 0.0)
	})

	t.Run("label_unit_adjusted_negative", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{-10, 100}}
		tsl := testSeriesList{series: []testSeries{series}}

		ar := calculateValueAxisRange(p, false, 800, nil, nil, nil,
			nil, 0, 0, 5, 4,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.Equal(t, 5, ar.labelCount)
		assert.InDelta(t, -20.0, ar.min, 0.0)
		assert.InDelta(t, 120.0, ar.max, 0.0)
	})

	t.Run("stacked_series", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1, 2, 3}},
			{values: []float64{4, 5, 6}},
		}}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, true, defaultValueFormatter, 0, fs)

		assert.InDelta(t, 0.0, ar.min, 0.0)
		assert.InDelta(t, 10.0, ar.max, 0.0)
	})

	t.Run("min_max_set", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20}}
		tsl := testSeriesList{series: []testSeries{series}}

		min := Ptr(5.0)
		max := Ptr(25.0)
		ar := calculateValueAxisRange(p, true, 800, min, max,
			nil, []string{}, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.InDelta(t, 5.0, ar.min, 0.0)
		assert.InDelta(t, 25.0, ar.max, 0.0)
	})

	t.Run("decimal_range", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{1.1, 2.2, 3.3}}
		tsl := testSeriesList{series: []testSeries{series}}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.InDelta(t, 1.0, ar.min, 0.0)
		assert.InDelta(t, 5.0, ar.max, 0.0)
		assert.Equal(t, 10, ar.labelCount)
	})

	t.Run("long_horizontal_labels", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series: []testSeries{series}}

		fs := FontStyle{FontSize: 28}
		inputLabels := []string{"ThisIsAVeryLongLabelThatExceedsNormal", "AnotherVeryLongLabelThatExceedsNormal",
			"WowLookAtTheseLabels!", "AndHereIsAnotherReallyLongLabel"}
		ar := calculateValueAxisRange(p, false, 800, nil, nil, nil,
			inputLabels, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.Equal(t, 811, ar.textMaxWidth)
		assert.Equal(t, 41, ar.textMaxHeight)
		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("zero_span_nonzero", func(t *testing.T) {
		// Series with a single nonzero value should trigger the zero–span adjustment.
		// When the data point is nonzero, we expect: min = value - zeroSpanAdjustment, max = value + zeroSpanAdjustment.
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{50}}
		tsl := testSeriesList{series: []testSeries{series}}

		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			nil, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.InDelta(t, 49.0, ar.min, 0.0)
		assert.InDelta(t, 51.0, ar.max, 0.0)
	})

	t.Run("vertical_label_rotation", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{10, 20, 30}}
		tsl := testSeriesList{series: []testSeries{series}}

		rotation := DegreesToRadians(45.0)
		ar := calculateValueAxisRange(p, true, 800, nil, nil, nil,
			[]string{"Label One", "Label Two", "Label Three", "Label Four"}, 0, 0, 0, 0,
			tsl, 0, false, defaultValueFormatter, rotation, fs)

		assert.Equal(t, 103, ar.textMaxWidth)
		assert.Equal(t, 103, ar.textMaxHeight)
		assert.InDelta(t, rotation, ar.labelRotation, 0.0)
	})

	t.Run("provided_labels_excess", func(t *testing.T) {
		// If the caller supplies more labels than the explicit labelCount,
		// the provided labels should be distributed across the range
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		series := testSeries{yAxisIndex: 0, values: []float64{5, 15, 25}}
		tsl := testSeriesList{series: []testSeries{series}}

		providedLabels := []string{"Label1", "Label2", "Label3", "Label4", "Label5"}
		explicitLabelCount := 3
		ar := calculateValueAxisRange(p, false, 800, nil, nil, nil,
			providedLabels, 0, explicitLabelCount, 0, 0,
			tsl, 0, false, defaultValueFormatter, 0, fs)

		assert.Equal(t, providedLabels, ar.labels)
		assert.Equal(t, len(providedLabels), ar.divideCount)
		assert.Equal(t, explicitLabelCount, ar.labelCount)
	})
}

func TestCalculateCategoryAxisRange(t *testing.T) {
	fs := FontStyle{FontSize: 16}

	t.Run("no_labels_provided_uses_series_names", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}}

		ar := calculateCategoryAxisRange(p, 800, false, nil, 0,
			0, 0, 0, tsl, 0, fs)

		expectedLabels := []string{"series:0", "series:1", "series:2"}
		assert.Equal(t, expectedLabels, ar.labels)
		assert.Equal(t, 3, ar.divideCount)
		assert.Equal(t, 3, ar.labelCount)
	})

	t.Run("provided_labels_filled_to_series_length", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		// Provide fewer labels than the number of series.
		providedLabels := []string{"CustomLabel"}
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}}

		ar := calculateCategoryAxisRange(p, 800, false, providedLabels, 0,
			0, 0, 0, tsl, 0, fs)

		assert.Equal(t, []string{"CustomLabel", "series:1", "series:2"}, ar.labels)
		assert.Equal(t, 3, ar.divideCount)
		assert.Equal(t, 3, ar.labelCount)
	})

	t.Run("explicit_label_count_cfg", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
			{values: []float64{4}},
		}}

		ar := calculateCategoryAxisRange(p, 800, false, nil, 0,
			2, 1, 0, tsl, 0, fs)

		assert.Equal(t, []string{"series:0", "series:1", "series:2", "series:3"}, ar.labels)
		assert.Equal(t, 4, ar.divideCount)
		assert.Equal(t, 3, ar.labelCount)
	})

	t.Run("font_defaults", func(t *testing.T) {
		fsEmpty := FontStyle{}
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
		}}

		ar := calculateCategoryAxisRange(p, 800, true, []string{}, 0,
			0, 0, 0, tsl, 0, fsEmpty)

		assert.NotNil(t, ar.labelFontStyle.Font)
		assert.InDelta(t, defaultFontSize, ar.labelFontStyle.FontSize, 0.0)
	})

	t.Run("label_rotation", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
		}}

		rotation := DegreesToRadians(30.0)
		ar := calculateCategoryAxisRange(p, 800, true, []string{}, 0,
			0, 0, 0, tsl, rotation, fs)

		assert.Equal(t, 81, ar.textMaxWidth)
		assert.Equal(t, 57, ar.textMaxHeight)
		assert.InDelta(t, rotation, ar.labelRotation, 0.0)
	})

	t.Run("negative_label_count_adjustment", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}}

		ar := calculateCategoryAxisRange(p, 800, false, []string{}, 0,
			0, -2, 0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("label_count_exceeds_series_count", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
		}}

		ar := calculateCategoryAxisRange(p, 800, false, []string{}, 0,
			5, 0, 0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("long_horizontal_labels", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		tsl := testSeriesList{series: []testSeries{
			{values: []float64{1}},
			{values: []float64{2}},
			{values: []float64{3}},
		}}

		inputLabels := []string{"ThisIsAVeryLongLabelThatExceedsNormal", "AnotherVeryLongLabelThatExceedsNormal",
			"WowLookAtTheseLabels!", "AndHereIsAnotherReallyLongLabel"}
		ar := calculateCategoryAxisRange(p, 600, false, inputLabels, 0,
			0, 0, 0, tsl, 0, fs)

		assert.Equal(t, 2, ar.labelCount)
	})

	t.Run("label_unit", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{
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
		}}

		ar := calculateCategoryAxisRange(p, 800, false, []string{}, 0,
			0, 0, 4.0, tsl, 0, fs)

		assert.Equal(t, 3, ar.labelCount)
	})

	t.Run("empty_series_list", func(t *testing.T) {
		p := NewPainter(PainterOptions{Width: 800, Height: 600})
		tsl := testSeriesList{series: []testSeries{}}
		ar := calculateCategoryAxisRange(p, 800, false, nil, 0,
			0, 0, 0, tsl, 0, fs)

		assert.Equal(t, []string{}, ar.labels)
		assert.Equal(t, 0, ar.divideCount)
		assert.Equal(t, 2, ar.labelCount)
	})
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

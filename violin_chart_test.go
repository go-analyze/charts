package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicViolinChartOption() ViolinChartOption {
	seriesList := NewSeriesListViolin([][][2]float64{
		{{0.22, 0.18}, {0.40, 0.34}, {0.54, 0.62}, {0.20, 0.16}},
	}, ViolinSeriesOption{
		Names: []string{"A"},
	})
	return ViolinChartOption{
		Theme:      GetDefaultTheme(),
		Padding:    NewBoxEqual(10),
		SeriesList: seriesList,
		ValueAxis:  ViolinAxis{},
		Legend: LegendOption{
			Show:        Ptr(false),
			SeriesNames: []string{"A"},
		},
	}
}

func makeDualViolinChartOption() ViolinChartOption {
	seriesList := NewSeriesListViolin([][][2]float64{
		{{0.15, 0.20}, {0.32, 0.40}, {0.48, 0.50}, {0.28, 0.22}},
		{{0.22, 0.18}, {0.40, 0.34}, {0.54, 0.62}, {0.20, 0.16}},
	}, ViolinSeriesOption{
		Names: []string{"A", "B"},
	})
	return ViolinChartOption{
		Theme:      GetDefaultTheme(),
		Padding:    NewBoxEqual(10),
		SeriesList: seriesList,
		ValueAxis:  ViolinAxis{},
		Legend: LegendOption{
			Show:        Ptr(false),
			SeriesNames: []string{"A", "B"},
		},
	}
}

func makeAllConfigOption() ViolinChartOption {
	opt := makeDualViolinChartOption()
	opt.Horizontal = true
	opt.ViolinWidth = "64"
	opt.ShowSpine = Ptr(true)
	opt.Theme = opt.Theme.WithAxisSplitLineColor(ColorBlack)
	opt.SpineWidth = 2.0
	opt.Title = TitleOption{Text: "all config"}
	opt.Legend = LegendOption{
		Show:        Ptr(true),
		SeriesNames: []string{"left", "right"},
	}
	return opt
}

func makeViolinAxisVerticalOption() ViolinChartOption {
	opt := makeBasicViolinChartOption()
	opt.ValueAxis.Show = Ptr(true)
	opt.ValueAxis.Title = "value axis"
	opt.ValueAxis.TitleFontStyle = FontStyle{
		FontColor: ColorBlue,
		FontSize:  14,
	}
	opt.ValueAxis.LabelFontStyle = FontStyle{
		FontColor: ColorRed,
		FontSize:  10,
	}
	opt.ValueAxis.LabelRotation = math.Pi / 9
	opt.ValueAxis.Limit = Ptr(1.4)
	opt.ValueAxis.Unit = 0.2
	opt.ValueAxis.LabelCount = 6
	opt.ValueAxis.PreferNiceIntervals = Ptr(true)
	opt.ValueFormatter = func(f float64) string {
		return strconv.FormatFloat(f, 'f', 3, 64)
	}
	return opt
}

func TestNewViolinChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewViolinChartOptionWithData([][][2]float64{{{0.2, 0.2}, {0.5, 0.5}}})

	require.Len(t, opt.SeriesList, 1)
	assert.Equal(t, ChartTypeViolin, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.ViolinChart(opt))
}

func TestViolinChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		makeOptions func() ViolinChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_vertical",
			makeOptions: makeBasicViolinChartOption,
			pngCRC:      0x1dfaf01b,
		},
		{
			name:        "basic_vertical",
			makeOptions: makeDualViolinChartOption,
			pngCRC:      0x3fa74c84,
		},
		{
			name: "horizontal_split",
			makeOptions: func() ViolinChartOption {
				opt := NewViolinChartOptionWithData([][][2]float64{
					{{0.2, 0.4}, {0.3, 0.8}, {0.2, 0.5}, {0.1, 0.2}},
					{{0.4, 0.2}, {0.8, 0.3}, {0.5, 0.2}, {0.2, 0.1}},
				})
				opt.Padding = NewBoxEqual(10)
				opt.Horizontal = true
				opt.ShowSpine = Ptr(true)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x3262eb4f,
		},
		{
			name: "width",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ViolinWidth = "22"
				return opt
			},
			pngCRC: 0xa94e1d4,
		},
		{
			name: "hide_spine",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ShowSpine = Ptr(false)
				return opt
			},
			pngCRC: 0xe24ea399,
		},
		{
			name: "value_formatter_axis_limit",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ValueFormatter = func(f float64) string {
					return strconv.FormatFloat(f, 'f', 2, 64)
				}
				opt.ValueAxis.Limit = Ptr(10.0)
				return opt
			},
			pngCRC: 0x48d36ad,
		},
		{
			name: "negative_width",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ViolinWidth = "-8"
				return opt
			},
			pngCRC: 0x1dfaf01b,
		},
		{
			name: "nan_inf_null_extents",
			makeOptions: func() ViolinChartOption {
				opt := makeDualViolinChartOption()
				opt.SeriesList[0].Data = [][2]float64{
					{math.NaN(), math.Inf(1)},
					{GetNullValue(), 0.2},
				}
				return opt
			},
			pngCRC: 0x1545de49,
		},
		{
			name:        "all_config",
			makeOptions: makeAllConfigOption,
			pngCRC:      0x38f15b81,
		},
		{
			name:        "vertical_axis_all_fields",
			makeOptions: makeViolinAxisVerticalOption,
			pngCRC:      0xb8b76901,
		},
		{
			name: "horizontal_axis_all_fields",
			makeOptions: func() ViolinChartOption {
				opt := makeViolinAxisVerticalOption()
				opt.Horizontal = true
				opt.ValueAxis.LabelRotation = -math.Pi / 12
				opt.ValueAxis.LabelCountAdjustment = 1
				opt.ValueAxis.Unit = 0.25
				return opt
			},
			pngCRC: 0x46687e85,
		},
		{
			name: "axis_hidden",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ValueAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xb7444e6c,
		},
		{
			name: "label_adjustment_1",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ValueAxis.LabelCountAdjustment = 1
				return opt
			},
			pngCRC: 0xee42246e,
		},
		{
			name: "label_adjustment_-1",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.ValueAxis.LabelCountAdjustment = -1
				return opt
			},
			pngCRC: 0x1dfaf01b,
		},
		{
			name: "mark_line_dual",
			makeOptions: func() ViolinChartOption {
				samples := [][]float64{
					{2, 4, 6, 8, 10, 12},
					{1, 3, 5, 7, 9, 11, 13},
				}
				opt, err := NewViolinChartOptionWithSamples(samples, 20)
				if err != nil {
					panic(err)
				}
				opt.SeriesList[0].MarkLine = NewMarkLine(
					SeriesMarkTypeAverage,
				)
				opt.SeriesList[1].MarkLine = NewMarkLine(
					SeriesMarkTypeAverage,
					SeriesMarkTypeMedian,
				)
				return opt
			},
			pngCRC: 0x5fb61b23,
		},
		{
			name: "mark_line_horizontal",
			makeOptions: func() ViolinChartOption {
				samples := [][]float64{
					{2, 4, 6, 8, 10, 12},
					{1, 3, 5, 7, 9, 11, 13},
				}
				opt, err := NewViolinChartOptionWithSamples(samples, 20)
				if err != nil {
					panic(err)
				}
				opt.Horizontal = true
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeAverage, SeriesMarkTypeMedian)
				return opt
			},
			pngCRC: 0xb4e26ec0,
		},
		{
			name: "nil_data",
			makeOptions: func() ViolinChartOption {
				return ViolinChartOption{
					SeriesList: ViolinSeriesList{{Name: "A"}},
					Legend:     LegendOption{Show: Ptr(false)},
				}
			},
			pngCRC: 0x8ff80d58,
		},
		{
			name: "single_pair",
			makeOptions: func() ViolinChartOption {
				return ViolinChartOption{
					SeriesList: ViolinSeriesList{{Data: [][2]float64{{0.3, 0.5}}, Name: "A"}},
					Legend:     LegendOption{Show: Ptr(false)},
				}
			},
			pngCRC: 0xa74aff81,
		},
		{
			name: "mixed_empty_and_populated",
			makeOptions: func() ViolinChartOption {
				return ViolinChartOption{
					SeriesList: ViolinSeriesList{
						{Name: "empty"},
						{Data: [][2]float64{{0.2, 0.4}, {0.3, 0.5}}, Name: "populated"},
					},
					Legend: LegendOption{Show: Ptr(false)},
				}
			},
			pngCRC: 0x568dfb38,
		},
		{
			name: "horizontal_nil_data",
			makeOptions: func() ViolinChartOption {
				return ViolinChartOption{
					Horizontal: true,
					SeriesList: ViolinSeriesList{{Name: "A"}, {Name: "B"}},
					Legend:     LegendOption{Show: Ptr(false)},
				}
			},
			pngCRC: 0xe18bf168,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			svgPainter := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})
			pngPainter := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        600,
				Height:       400,
			})
			validateViolinChartRender(t, svgPainter, pngPainter, tt.makeOptions(), tt.pngCRC)
		})
	}
}

func TestViolinResolveAxisBounds(t *testing.T) {
	t.Parallel()

	one := 1.0
	five := 5.0
	pointTwo := 0.2

	tests := []struct {
		name       string
		absMax     float64
		limitOpt   *float64
		labelCount int
		unit       float64
		wantMin    float64
		wantMax    float64
	}{
		{
			name:       "empty_defaults",
			absMax:     0,
			labelCount: 0,
			wantMin:    -1,
			wantMax:    1,
		},
		{
			name:       "data_defaults",
			absMax:     4,
			labelCount: 0,
			wantMin:    -4,
			wantMax:    4,
		},
		{
			name:       "decimal_data_rounds_to_nice_extent",
			absMax:     0.62,
			labelCount: 5,
			wantMin:    -0.8,
			wantMax:    0.8,
		},
		{
			name:       "unit_guides_extent",
			absMax:     0.62,
			labelCount: 5,
			unit:       pointTwo,
			wantMin:    -0.8,
			wantMax:    0.8,
		},
		{
			name:       "limit_applied",
			absMax:     0,
			limitOpt:   &five,
			labelCount: 0,
			wantMin:    -5,
			wantMax:    5,
		},
		{
			name:       "limit_does_not_clip_data",
			absMax:     9,
			limitOpt:   &one,
			labelCount: 0,
			wantMin:    -9,
			wantMax:    9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := violinResolveAxisBounds(tt.absMax, tt.limitOpt, tt.labelCount, tt.unit)
			assert.InDelta(t, tt.wantMin, gotMin, 1e-12)
			assert.InDelta(t, tt.wantMax, gotMax, 1e-12)
		})
	}
}

func TestViolinEnsureZeroLabelCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     int
		wantCount int
	}{
		{name: "default_when_zero", input: 0, wantCount: 5},
		{name: "default_when_negative", input: -2, wantCount: 5},
		{name: "minimum_when_one", input: 1, wantCount: 3},
		{name: "minimum_when_two", input: 2, wantCount: 3},
		{name: "odd_kept", input: 7, wantCount: 7},
		{name: "even_incremented", input: 8, wantCount: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantCount, violinEnsureZeroLabelCount(tt.input))
		})
	}
}

func TestViolinSingleSeriesUsesFullCategoryArea(t *testing.T) {
	t.Parallel()

	seriesList := NewSeriesListViolin([][][2]float64{
		{{0.2, 0.4}, {0.3, 0.5}},
	})
	axisSize := 320
	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})

	r := calculateCategoryAxisRange(
		p,
		axisSize,
		true,
		false,
		[]string{"solo"},
		0,
		0,
		0,
		0,
		seriesList,
		0,
		FontStyle{},
	)

	assert.Equal(t, 1, r.divideCount)
	start, end := r.getRange(0)
	assert.InDelta(t, float64(axisSize), end-start, 1e-12)
}

func TestViolinChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() ViolinChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() ViolinChartOption {
				return NewViolinChartOptionWithSeries(ViolinSeriesList{})
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "invalid_yaxis_index",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.SeriesList[0].YAxisIndex = 1
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
		{
			name: "mark_line_without_stats",
			makeOptions: func() ViolinChartOption {
				opt := makeBasicViolinChartOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeAverage)
				opt.SeriesList[0].Stats = nil
				return opt
			},
			errorMsgContains: "mark line requires series stats",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})
			err := p.ViolinChart(tt.makeOptions())
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsgContains)
		})
	}
}

func TestNewViolinChartOptionWithSamples(t *testing.T) {
	t.Parallel()

	samples := [][]float64{
		{1, 2, 3, 4, 5, 6, 7, 8},
		{2, 3, 4, 5, 6, 7, 8, 9},
	}
	opt, err := NewViolinChartOptionWithSamples(samples, 24)
	require.NoError(t, err)
	require.Len(t, opt.SeriesList, 2)
	require.NotEmpty(t, opt.SeriesList[0].Data)

	for _, pair := range opt.SeriesList[0].Data {
		assert.InDelta(t, pair[0], pair[1], 1e-12)
	}

	expectedA := summarizePopulationData(samples[0])
	expectedB := summarizePopulationData(samples[1])
	require.NotNil(t, opt.SeriesList[0].Stats)
	require.NotNil(t, opt.SeriesList[1].Stats)
	assert.InDelta(t, expectedA.Average, opt.SeriesList[0].Stats.Average, 1e-12)
	assert.InDelta(t, expectedA.Median, opt.SeriesList[0].Stats.Median, 1e-12)
	assert.InDelta(t, expectedA.Min, opt.SeriesList[0].Stats.Min, 1e-12)
	assert.InDelta(t, expectedA.Max, opt.SeriesList[0].Stats.Max, 1e-12)
	assert.InDelta(t, expectedB.Average, opt.SeriesList[1].Stats.Average, 1e-12)
	assert.InDelta(t, expectedB.Median, opt.SeriesList[1].Stats.Median, 1e-12)
	assert.InDelta(t, expectedB.Min, opt.SeriesList[1].Stats.Min, 1e-12)
	assert.InDelta(t, expectedB.Max, opt.SeriesList[1].Stats.Max, 1e-12)
}

func TestViolinSamplesZeroVariance(t *testing.T) {
	t.Parallel()

	samples := [][]float64{
		{5, 5, 5, 5},
		{1, 2, 3, 4, 5},
	}
	opt, err := NewViolinChartOptionWithSamples(samples, 20)
	require.NoError(t, err)
	require.Len(t, opt.SeriesList, 2)
	assert.Empty(t, opt.SeriesList[0].Data)
	assert.NotEmpty(t, opt.SeriesList[1].Data)
}

func TestViolinSamplesAllGapsAllowed(t *testing.T) {
	t.Parallel()

	samples := [][]float64{
		{5, 5, 5, 5},
		{9, 9, 9, 9},
	}
	opt, err := NewViolinChartOptionWithSamples(samples, 20)
	require.NoError(t, err)
	require.Len(t, opt.SeriesList, 2)
	assert.Empty(t, opt.SeriesList[0].Data)
	assert.Empty(t, opt.SeriesList[1].Data)

	// all-empty-data option renders without error
	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})
	require.NoError(t, p.ViolinChart(opt))
	data, err := p.Bytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestViolinSamplesNormalization(t *testing.T) {
	t.Parallel()

	samples := [][]float64{
		{1, 2, 3, 4, 5, 6, 7, 8},
		{1, 1.5, 2, 6, 6.5, 7, 10, 12},
	}
	perSeries, err := NewViolinChartOptionWithSamples(samples, 24, ViolinNormalizationPerSeries)
	require.NoError(t, err)
	global, err := NewViolinChartOptionWithSamples(samples, 24, ViolinNormalizationGlobal)
	require.NoError(t, err)

	perMaxA := maxViolinExtent(perSeries.SeriesList[0].Data)
	perMaxB := maxViolinExtent(perSeries.SeriesList[1].Data)
	assert.InDelta(t, 1.0, perMaxA, 1e-9)
	assert.InDelta(t, 1.0, perMaxB, 1e-9)

	globalMaxA := maxViolinExtent(global.SeriesList[0].Data)
	globalMaxB := maxViolinExtent(global.SeriesList[1].Data)
	assert.InDelta(t, 1.0, math.Max(globalMaxA, globalMaxB), 1e-9)
	assert.True(t, globalMaxA < 1.0 || globalMaxB < 1.0)
}

func TestViolinSamplesBandwidthOverride(t *testing.T) {
	t.Parallel()

	samples := [][]float64{
		{1, 1.5, 2, 3, 4, 5, 6, 8, 9},
	}

	defaultOpt, err := NewViolinChartOptionWithSamples(samples, 25)
	require.NoError(t, err)
	smallBWOpt, err := NewViolinChartOptionWithSamples(samples, 25, "bandwidth=0.2")
	require.NoError(t, err)
	largeBWOpt, err := NewViolinChartOptionWithSamples(samples, 25, "bw=3.0")
	require.NoError(t, err)

	require.NotEmpty(t, defaultOpt.SeriesList[0].Data)
	require.NotEmpty(t, smallBWOpt.SeriesList[0].Data)
	require.NotEmpty(t, largeBWOpt.SeriesList[0].Data)

	center := len(defaultOpt.SeriesList[0].Data) / 2
	defaultCenter := defaultOpt.SeriesList[0].Data[center][0]
	smallBWCenter := smallBWOpt.SeriesList[0].Data[center][0]
	largeBWCenter := largeBWOpt.SeriesList[0].Data[center][0]

	assert.NotEqual(t, defaultCenter, smallBWCenter)
	assert.NotEqual(t, defaultCenter, largeBWCenter)
	assert.NotEqual(t, smallBWCenter, largeBWCenter)
}

func TestViolinSamplesInvalidInput(t *testing.T) {
	t.Parallel()

	_, err := NewViolinChartOptionWithSamples([][]float64{{1, 2}}, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "point count")

	_, err = NewViolinChartOptionWithSamples([][]float64{{1, 2}}, 12, "invalid")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "normalization")

	_, err = NewViolinChartOptionWithSamples([][]float64{{1, 2}}, 12, "bandwidth=abc")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bandwidth")

	_, err = NewViolinChartOptionWithSamples([][]float64{{1, 2}}, 12, "bw=0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bandwidth")

	_, err = NewViolinChartOptionWithSamples([][]float64{{GetNullValue(), GetNullValue()}}, 12)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no valid data")
}

func TestRenderViolinChart(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: NewSeriesListViolin([][][2]float64{
			{{0.1, 0.2}, {0.4, 0.5}, {0.2, 0.3}},
			{{0.2, 0.1}, {0.5, 0.4}, {0.3, 0.2}},
		}).ToGenericSeriesList(),
		XAxis: XAxisOption{
			Labels: []string{"A", "B"},
		},
		Legend: LegendOption{Show: Ptr(false)},
	}

	svgPainter, err := Render(opt, SVGOutputOptionFunc())
	require.NoError(t, err)
	svgData, err := svgPainter.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, svgData)

	pngPainter, err := Render(opt, PNGOutputOptionFunc())
	require.NoError(t, err)
	pngData, err := pngPainter.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, 0x7c42daa2, pngData)

	// empty data via ChartOption API renders without error
	emptyOpt := ChartOption{
		SeriesList: GenericSeriesList{
			{Type: ChartTypeViolin, Name: "A"},
			{Type: ChartTypeViolin, Values: []float64{0.1, 0.2}, Name: "B"},
		},
		Legend: LegendOption{Show: Ptr(false)},
	}
	p, err := Render(emptyOpt, SVGOutputOptionFunc())
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestViolinRender(t *testing.T) {
	t.Parallel()

	svgPainter, err := ViolinRender([][][2]float64{
		{{0.1, 0.2}, {0.4, 0.5}, {0.2, 0.3}},
		{{0.2, 0.1}, {0.5, 0.4}, {0.3, 0.2}},
	},
		SVGOutputOptionFunc(),
		XAxisLabelsOptionFunc([]string{"A", "B"}),
	)
	require.NoError(t, err)
	svgData, err := svgPainter.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, svgData)

	pngPainter, err := ViolinRender([][][2]float64{
		{{0.1, 0.2}, {0.4, 0.5}, {0.2, 0.3}},
		{{0.2, 0.1}, {0.5, 0.4}, {0.3, 0.2}},
	},
		PNGOutputOptionFunc(),
		XAxisLabelsOptionFunc([]string{"A", "B"}),
	)
	require.NoError(t, err)
	pngData, err := pngPainter.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, 0x7c42daa2, pngData)
}

func TestRenderHorizontalViolinChart(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: GenericSeriesList{
			{
				Type:   ChartTypeHorizontalViolin,
				Values: []float64{0.1, 0.2, 0.4, 0.5, 0.2, 0.3},
				Name:   "A",
			},
			{
				Type:   ChartTypeHorizontalViolin,
				Values: []float64{0.2, 0.1, 0.5, 0.4, 0.3, 0.2},
				Name:   "B",
			},
		},
		YAxis: []YAxisOption{
			{
				Labels: []string{"A", "B"},
			},
		},
		Legend: LegendOption{Show: Ptr(false)},
	}

	p, err := Render(opt, SVGOutputOptionFunc())
	require.NoError(t, err)
	svgData, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, svgData)
}

func TestRenderViolinCanNotMix(t *testing.T) {
	t.Parallel()

	violinSeries := NewSeriesListViolin([][][2]float64{{{0.1, 0.2}, {0.3, 0.4}}}).ToGenericSeriesList()
	lineSeries := NewSeriesListLine([][]float64{{1, 2, 3}}).ToGenericSeriesList()
	mixed := append(violinSeries, lineSeries...)

	_, err := Render(ChartOption{SeriesList: mixed})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "violin can not mix other charts")
}

func validateViolinChartRender(t *testing.T, svgP, pngP *Painter, opt ViolinChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.ViolinChart(opt)
	require.NoError(t, err)
	svgData, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, svgData)

	err = pngP.ViolinChart(opt)
	require.NoError(t, err)
	pngData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, pngData)
}

func maxViolinExtent(data [][2]float64) float64 {
	maxVal := 0.0
	for _, pair := range data {
		if pair[0] > maxVal {
			maxVal = pair[0]
		}
		if pair[1] > maxVal {
			maxVal = pair[1]
		}
	}
	return maxVal
}

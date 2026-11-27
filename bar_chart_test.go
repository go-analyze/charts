package charts

import (
	"strconv"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicBarChartOption() BarChartOption {
	seriesList := NewSeriesListBar([][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3},
	})
	for index := range seriesList {
		seriesList[index].Label.Show = Ptr(true)
	}
	return BarChartOption{
		Padding:    NewBoxEqual(10),
		SeriesList: seriesList,
		CategoryAxis: CategoryAxisOption{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		},
		ValueAxis: []ValueAxisOption{
			{
				Show:                Ptr(true),
				PreferNiceIntervals: Ptr(true),
			},
		},
	}
}

func makeMinimalBarChartOption() BarChartOption {
	opt := NewBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})
	opt.CategoryAxis = CategoryAxisOption{
		Show:   Ptr(false),
		Labels: []string{"A", "B"},
	}
	opt.ValueAxis[0].Show = Ptr(false)
	return opt
}

func makeFullBarChartStackedOption() BarChartOption {
	seriesList := NewSeriesListBar([][]float64{
		{4.9, 23.2, 25.6, 102.6, 142.2, 32.6, 20.0, 3.3},
		{9.0, 26.4, 28.7, 144.6, 122.2, 48.7, 18.8, 2.3},
		{80.0, 40.4, 28.4, 28.8, 24.4, 24.2, 40.8, 80.8},
	}, BarSeriesOption{
		Label: SeriesLabel{
			Show: Ptr(false),
		},
		MarkPoint: NewMarkPoint("max"),
	})
	return BarChartOption{
		Padding:     NewBoxEqual(20),
		SeriesList:  seriesList,
		StackSeries: Ptr(true),
		CategoryAxis: CategoryAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		},
		Legend: LegendOption{
			SeriesNames: []string{"A", "B", "C"},
			Symbol:      SymbolDot,
		},
		ValueAxis: []ValueAxisOption{
			{
				RangeValuePaddingScale: Ptr(1.0),
			},
		},
	}
}

func TestNewBarChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeBar, opt.SeriesList[0].getType())
	assert.Len(t, opt.ValueAxis, 1)
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.BarChart(opt))
}

func TestBarChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() BarChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicBarChartOption,
			pngCRC:      0x84501d88,
		},
		{
			name: "rounded_caps",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.RoundedBarCaps = Ptr(true)
				return opt
			},
			pngCRC: 0xb8509c94,
		},
		{
			name: "custom_font",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.CategoryAxis.LabelFontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0x7a38623d,
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.CategoryAxis.BoundaryGap = Ptr(true)
				opt.ValueAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xcf5dd6d,
		},
		{
			name: "boundary_gap_disable",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.CategoryAxis.BoundaryGap = Ptr(false)
				return opt
			},
			pngCRC: 0x98f53836,
		},
		{
			name: "value_formatter",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.ValueFormatter = func(f float64) string {
					return "f"
				}
				return opt
			},
			pngCRC: 0x39525f8f,
		},
		{
			name: "bar_width_truncate",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.CategoryAxis.Show = Ptr(false)
				opt.ValueAxis[0].Show = Ptr(false)
				opt.BarSize = 2.0 // exceeds the slot, limited to fit
				return opt
			},
			pngCRC: 0xd248993a,
		},
		{
			name: "bar_width_thin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarSize = 0.018 // ~2px wide
				return opt
			},
			pngCRC: 0x83ffbda8,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0xcf9984cd,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarMargin = Ptr(5.0) // will be limited to fit graph
				return opt
			},
			pngCRC: 0x3ba0d40d,
		},
		{
			name: "bar_width_and_narrow_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarSize = 0.075 // ~10px wide
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x86798cdc,
		},
		{
			name: "bar_width_and_wide_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarSize = 0.075      // ~10px wide
				opt.BarMargin = Ptr(5.0) // will be limited for readability
				return opt
			},
			pngCRC: 0xfec5eff9,
		},
		{
			name: "dual_yaxis",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.Title.Text = "T"
				opt.SeriesList[1].YAxisIndex = 1
				opt.ValueAxis = append(opt.ValueAxis, opt.ValueAxis[0])
				opt.ValueAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.ValueAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				return opt
			},
			pngCRC: 0xc833fe47,
		},
		{
			name: "mark_line",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].MarkLine = NewMarkLine("min", "max", "average")
				}
				return opt
			},
			pngCRC: 0x8db6641c,
		},
		{
			name: "mark_point",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].MarkPoint = NewMarkPoint("min", "max")
				}
				return opt
			},
			pngCRC: 0x87ca50d8,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullBarChartStackedOption,
			pngCRC:      0x8f3f3642,
		},
		{
			name: "stack_series_capped_bar",
			makeOptions: func() BarChartOption {
				opt := makeFullBarChartStackedOption()
				opt.RoundedBarCaps = Ptr(true)
				// disable some extra visuals
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
					opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
						return strconv.Itoa(int(f))
					}
					opt.SeriesList[i].MarkLine.Lines = nil
					opt.SeriesList[i].MarkPoint.Points = nil
				}
				opt.ValueAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xa1b5f660,
		},
		{
			name: "stack_series_global_mark_point",
			makeOptions: func() BarChartOption {
				opt := makeFullBarChartStackedOption()
				// add global point configurations
				opt.SeriesList[len(opt.SeriesList)-1].MarkPoint.AddGlobalPoints(SeriesMarkTypeMax)
				// disable some extra visuals
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
					opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
						return strconv.Itoa(int(f))
					}
					opt.SeriesList[i].MarkLine.Lines = nil
					if i != len(opt.SeriesList)-1 {
						opt.SeriesList[i].MarkPoint.Points = nil
					}
				}
				opt.ValueAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xcd381eeb,
		},
		{
			name: "stack_series_simple",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{{4.0}, {1.0}})
				opt.StackSeries = Ptr(true)
				// disable extra
				opt.CategoryAxis.Show = Ptr(false)
				opt.ValueAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x5f9d8d95,
		},
		{
			name: "stack_series_global_mark_line",
			makeOptions: func() BarChartOption {
				opt := makeFullBarChartStackedOption()
				// add global line configurations
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine.AddGlobalLines(SeriesMarkTypeAverage,
					SeriesMarkTypeMin, SeriesMarkTypeMax)
				// disable some extra visuals
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(false)
					opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
						return strconv.Itoa(int(f))
					}
					opt.SeriesList[i].MarkPoint.Points = nil
				}
				opt.ValueAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xd2bbd8a5,
		},
		{
			name: "empty_series",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{})
				opt.Padding = NewBoxEqual(10)
				opt.Legend = LegendOption{
					Show:        Ptr(true),
					SeriesNames: []string{"Series A", "Series B"},
				}
				opt.CategoryAxis = CategoryAxisOption{Labels: []string{"Jan", "Feb", "Mar"}}
				opt.ValueAxis = []ValueAxisOption{{Show: Ptr(true)}, {Show: Ptr(true)}}
				return opt
			},
			pngCRC: 0x374571ca,
		},
		{
			name: "null_values",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.SeriesList[0].Values[1] = GetNullValue()
				opt.SeriesList[1].Values[7] = GetNullValue()
				return opt
			},
			pngCRC: 0xf422fba1,
		},
		{
			name: "stack_series_null",
			makeOptions: func() BarChartOption {
				opt := makeFullBarChartStackedOption()
				opt.SeriesList[0].Values[3] = GetNullValue()
				return opt
			},
			pngCRC: 0x8bd0d4fd,
		},
		{
			name: "stack_series_global_mark_point_null_last",
			makeOptions: func() BarChartOption {
				// combined max (index 1) and min (index 3) fall where the last series is null
				opt := NewBarChartOptionWithData([][]float64{
					{10, 40, 10, 5, 10},
					{10, 40, 10, 5, 10},
					{5, GetNullValue(), 5, GetNullValue(), 5},
				})
				opt.StackSeries = Ptr(true)
				opt.SeriesList[len(opt.SeriesList)-1].MarkPoint.AddGlobalPoints(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x3dff262e,
		},
		{
			name: "stack_series_global_mark_point_null_series",
			makeOptions: func() BarChartOption {
				// last stacked series is entirely null, marks anchor to the combined stack
				opt := NewBarChartOptionWithData([][]float64{
					{10, 40, 10, 5, 10},
					{10, 40, 10, 5, 10},
					{GetNullValue(), GetNullValue(), GetNullValue(), GetNullValue(), GetNullValue()},
				})
				opt.StackSeries = Ptr(true)
				opt.SeriesList[len(opt.SeriesList)-1].MarkPoint.AddGlobalPoints(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x75bb6a49,
		},
		{
			name: "stack_series_all_secondary_axis",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.StackSeries = Ptr(true)
				for i := range opt.SeriesList {
					opt.SeriesList[i].YAxisIndex = 1
				}
				opt.ValueAxis = append(opt.ValueAxis, opt.ValueAxis[0])
				return opt
			},
			pngCRC: 0x9810c3b3,
		},
		{
			name: "stack_series_dual_yaxis",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{4.9, 23.2, 25.6, 102.6},
					{9.0, 26.4, 28.7, 144.6},
					{80.0, 40.4, 28.4, 28.8},
				})
				opt.StackSeries = Ptr(true)
				opt.SeriesList[2].YAxisIndex = 1
				opt.ValueAxis = append(opt.ValueAxis, opt.ValueAxis[0])
				return opt
			},
			pngCRC: 0xe3ef4961,
		},
		{
			name: "stack_series_single_secondary",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{{4.0, 2.0, 3.0}})
				opt.StackSeries = Ptr(true)
				opt.SeriesList[0].YAxisIndex = 1
				opt.ValueAxis = append(opt.ValueAxis, opt.ValueAxis[0])
				return opt
			},
			pngCRC: 0x1d666b75,
		},
		{
			name: "negative_values",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{-12, -5, -2, -5, -12, -19, -22, -19},
					{-16, -9, -6, -9, -16, -23, -26, -23},
				})
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0xfa8d8497,
		},
		{
			name: "mixed_sign",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{0, 7, 10, 7, 0, -7, -10, -7},
					{10, 7, 0, -7, -10, -7, 0, 7},
				})
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0x66a23407,
		},
		{
			name: "stack_series_negative",
			makeOptions: func() BarChartOption {
				// mixed signs so stacks diverge above and below zero
				opt := NewBarChartOptionWithData([][]float64{
					{0, 7, 10, 7, 0, -7, -10, -7},
					{-7, -10, -7, 0, 7, 10, 7, 0},
					{3, 5, 3, 0, -3, -5, -3, 0},
				})
				opt.StackSeries = Ptr(true)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0xf4284353,
		},
		{
			name: "stack_series_all_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{-12, -5, -2, -5, -12, -19, -22, -19},
					{-6, -2, -1, -2, -6, -9, -11, -9},
				})
				opt.StackSeries = Ptr(true)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0x56b30499,
		},
		{
			name: "rounded_caps_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{5, 10, 5, -5, -10, -5},
					{-10, -5, 5, 10, 5, -5},
				})
				opt.RoundedBarCaps = Ptr(true)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F"}
				return opt
			},
			pngCRC: 0x8e5bd82,
		},
		{
			name: "label_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{5, 10, 5, -5, -10, -5},
					{-10, -5, 5, 10, 5, -5},
				})
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F"}
				return opt
			},
			pngCRC: 0x715cc9d8,
		},
		{
			name: "mark_point_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{5, 10, 5, -5, -10, -5},
				})
				opt.SeriesList[0].MarkPoint = NewMarkPoint(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F"}
				return opt
			},
			pngCRC: 0x8b32f238,
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		rasterOptions := PainterOptions{
			OutputFormat: ChartOutputPNG,
			Width:        600,
			Height:       400,
		}
		if tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				r := NewPainter(rasterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateBarChartRender(t, p, r, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				r := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateBarChartRender(t, p, r, opt, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				r := NewPainter(rasterOptions)

				validateBarChartRender(t, p, r, tt.makeOptions(), tt.pngCRC)
			})
		}
	}
}

func validateBarChartRender(t *testing.T, svgP, pngP *Painter, opt BarChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.BarChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.BarChart(opt)
	require.NoError(t, err)
	rdata, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rdata)
}

func TestEnsureBarSpan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		lo, hi         int
		baseLow        bool
		gap, limit     int
		wantLo, wantHi int
	}{
		{name: "non_empty_high_base", lo: 10, hi: 20, limit: 100, wantLo: 10, wantHi: 20},
		{name: "non_empty_low_base", lo: 10, hi: 20, baseLow: true, limit: 100, wantLo: 10, wantHi: 20},
		{name: "empty_high_base", lo: 50, hi: 50, limit: 100, wantLo: 49, wantHi: 50},
		{name: "empty_low_base", lo: 50, hi: 50, baseLow: true, limit: 100, wantLo: 50, wantHi: 51},
		{name: "empty_high_base_gap", lo: 49, hi: 50, gap: 1, limit: 100, wantLo: 49, wantHi: 50},
		{name: "empty_low_base_gap", lo: 50, hi: 51, baseLow: true, gap: 1, limit: 100, wantLo: 50, wantHi: 51},
		{name: "inverted_high_base", lo: 51, hi: 49, limit: 100, wantLo: 48, wantHi: 49},
		{name: "inverted_low_base", lo: 50, hi: 49, baseLow: true, limit: 100, wantLo: 50, wantHi: 51},
		// no pixel fits on the baseline side, so the span stays empty
		{name: "high_base_at_origin", lo: 1, hi: 0, limit: 100, wantLo: 0, wantHi: 0},
		{name: "low_base_at_limit", lo: 100, hi: 100, baseLow: true, limit: 100, wantLo: 100, wantHi: 100},
		{name: "low_base_beyond_limit", lo: 101, hi: 100, baseLow: true, limit: 100, wantLo: 101, wantHi: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo, hi := ensureBarSpan(tt.lo, tt.hi, tt.baseLow, tt.gap, tt.limit)
			assert.Equal(t, tt.wantLo, lo)
			assert.Equal(t, tt.wantHi, hi)
		})
	}
}

func TestStackedBarLanes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		axes          []int
		expectedLanes []int
		expectedCount int
	}{
		{
			name:          "mixed_axes",
			axes:          []int{0, 1},
			expectedLanes: []int{0, 1},
			expectedCount: 2,
		},
		{
			name:          "all_secondary",
			axes:          []int{1, 1},
			expectedLanes: []int{0, 1},
			expectedCount: 2,
		},
		{
			name:          "single_secondary",
			axes:          []int{1},
			expectedLanes: []int{0},
			expectedCount: 1,
		},
		{
			name:          "all_primary",
			axes:          []int{0, 0, 0},
			expectedLanes: []int{0, 0, 0},
			expectedCount: 1,
		},
		{
			name:          "interleaved",
			axes:          []int{1, 0, 1},
			expectedLanes: []int{1, 0, 2},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := make(BarSeriesList, len(tt.axes))
			for i, axis := range tt.axes {
				sl[i] = BarSeries{Values: []float64{1}, YAxisIndex: axis}
			}
			lanes, count := stackedBarLanes(sl)
			assert.Equal(t, tt.expectedLanes, lanes)
			assert.Equal(t, tt.expectedCount, count)
		})
	}
}

func TestCalculateGroupMarginsAndSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		seriesCount     int
		space           int
		size            int
		margin          *float64
		expectedMargin  int
		expectedSpacing int
		expectedSize    int
	}{
		{
			name:            "large_space_default",
			seriesCount:     3,
			space:           100,
			expectedMargin:  10,
			expectedSpacing: 5,
			expectedSize:    23,
		},
		{
			name:            "small_space_multiple",
			seriesCount:     2,
			space:           18,
			expectedMargin:  2,
			expectedSpacing: 2,
			expectedSize:    6,
		},
		{
			name:            "custom_width_and_margin",
			seriesCount:     2,
			space:           90,
			size:            20,
			margin:          Ptr(4.0),
			expectedMargin:  23,
			expectedSpacing: 4,
			expectedSize:    20,
		},
		{
			name:            "width_too_wide",
			seriesCount:     2,
			space:           40,
			size:            20,
			margin:          Ptr(1.0),
			expectedMargin:  5,
			expectedSpacing: 3,
			expectedSize:    13,
		},
		{
			name:            "single_series_default",
			seriesCount:     1,
			space:           50,
			expectedMargin:  10,
			expectedSpacing: 5,
			expectedSize:    30,
		},
		{
			name:            "single_series_custom",
			seriesCount:     1,
			space:           50,
			size:            10,
			margin:          Ptr(0.0),
			expectedMargin:  20,
			expectedSpacing: 0,
			expectedSize:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			margin, spacing, size := calculateGroupMarginsAndSize(tt.seriesCount, tt.space, tt.size, tt.margin)
			assert.Equal(t, tt.expectedMargin, margin)
			assert.Equal(t, tt.expectedSpacing, spacing)
			assert.Equal(t, tt.expectedSize, size)
		})
	}
}

func TestBarChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() BarChartOption
		errorMsgContains string
	}{
		{
			name: "invalid_yaxis_index",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{{1, 2, 3}})
				opt.SeriesList[0].YAxisIndex = 2
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
		{
			name: "negative_yaxis_index",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{{1, 2, 3}})
				opt.SeriesList[0].YAxisIndex = -1
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				Width:  600,
				Height: 400,
			})

			err := p.BarChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

func makeBasicHorizontalBarOption() BarChartOption {
	return BarChartOption{
		Horizontal: true,
		Padding:    NewBoxEqual(10),
		SeriesList: NewSeriesListBar([][]float64{
			{18203, 23489, 29034, 104970, 131744, 630230},
			{19325, 23438, 31000, 121594, 134141, 681807},
		}),
		Title: TitleOption{
			Text: "World Population",
		},
		Legend: LegendOption{
			SeriesNames: []string{"2011", "2012"},
			Symbol:      SymbolDot,
		},
		CategoryAxis: CategoryAxisOption{
			Labels: []string{"Brazil", "Indonesia", "USA", "India", "China", "World"},
		},
	}
}

func makeMinimalHorizontalBarOption() BarChartOption {
	opt := NewBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})
	opt.Horizontal = true
	opt.CategoryAxis = CategoryAxisOption{
		Show:   Ptr(false),
		Labels: []string{"A", "B"},
	}
	opt.ValueAxis[0].Show = Ptr(false)
	return opt
}

func makeFullHorizontalBarStackedOption() BarChartOption {
	seriesList := NewSeriesListBar([][]float64{
		{4.9, 23.2, 25.6, 102.6, 142.2, 32.6, 20.0, 3.3},
		{19.0, 26.4, 28.7, 144.6, 122.2, 48.7, 28.8, 22.3},
		{80.0, 40.4, 28.4, 28.8, 24.4, 24.2, 40.8, 80.8},
	}, BarSeriesOption{
		Label: SeriesLabel{
			Show: Ptr(true),
			ValueFormatter: func(f float64) string {
				return strconv.Itoa(int(f))
			},
		},
	})
	return BarChartOption{
		Horizontal:  true,
		Padding:     NewBoxEqual(20),
		SeriesList:  seriesList,
		StackSeries: Ptr(true),
		Legend: LegendOption{
			Symbol: SymbolDot,
		},
		CategoryAxis: CategoryAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		},
	}
}

func TestBarChartHorizontal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() BarChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicHorizontalBarOption,
			pngCRC:      0x7a0bd5c0,
		},
		{
			name: "custom_fonts",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.ValueAxis = []ValueAxisOption{{LabelFontStyle: customFont}}
				opt.CategoryAxis.LabelFontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0x30aadf4d,
		},
		{
			name: "value_labels",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
					series[i].Label.ValueFormatter = func(f float64) string {
						return humanize.FtoaWithDigits(f, 2)
					}
				}
				return opt
			},
			pngCRC: 0xf8d3d830,
		},
		{
			name: "value_formatter",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.ValueFormatter = func(f float64) string {
					return "f"
				}
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
					series[i].Label.ValueFormatter = opt.ValueFormatter
				}
				return opt
			},
			pngCRC: 0x9049e1b2,
		},
		{
			name: "bar_size_truncate",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.Title.Show = Ptr(false)
				opt.ValueAxis = []ValueAxisOption{{Show: Ptr(false)}}
				opt.CategoryAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				opt.BarSize = 2.0 // exceeds the slot, limited to fit
				return opt
			},
			pngCRC: 0x60fc4ac0,
		},
		{
			name: "mark_line",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				opt.CategoryAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x77e4f95,
		},
		{
			name: "bar_size_thin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarSize = 0.028 // ~2px tall
				return opt
			},
			pngCRC: 0x74310edf,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x8f92fa71,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarMargin = Ptr(5.0) // will be limited to fit graph
				return opt
			},
			pngCRC: 0x44f506e1,
		},
		{
			name: "bar_size_and_narrow_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarSize = 0.117 // ~10px tall
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0xd5c8fd54,
		},
		{
			name: "bar_size_and_wide_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarSize = 0.117      // ~10px tall
				opt.BarMargin = Ptr(5.0) // will be limited for readability
				return opt
			},
			pngCRC: 0xe0e52271,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullHorizontalBarStackedOption,
			pngCRC:      0xf3227c9a,
		},
		{
			name: "stack_series_simple",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{{4.0}, {1.0}})
				opt.Horizontal = true
				opt.StackSeries = Ptr(true)
				opt.ValueAxis[0].Unit = 1
				opt.CategoryAxis = CategoryAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0xc9877e11,
		},
		{
			name: "stack_series_with_mark",
			makeOptions: func() BarChartOption {
				opt := makeFullHorizontalBarStackedOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine = NewMarkLine(SeriesMarkTypeMax)
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine.Lines[0].Global = true
				opt.CategoryAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x5c06be7,
		},
		{
			name: "empty_series",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{})
				opt.Horizontal = true
				opt.Padding = NewBoxEqual(10)
				opt.Legend = LegendOption{
					Show:        Ptr(true),
					SeriesNames: []string{"Series A", "Series B"},
				}
				opt.CategoryAxis = CategoryAxisOption{
					Show:   Ptr(true),
					Labels: []string{"A", "B", "C"},
				}
				opt.ValueAxis = []ValueAxisOption{{Show: Ptr(true)}}
				return opt
			},
			pngCRC: 0xb4c84169,
		},
		{
			name: "rounded_caps",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.RoundedBarCaps = Ptr(true)
				return opt
			},
			pngCRC: 0x3bf203d6,
		},
		{
			name: "mark_point",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.SeriesList[0].MarkPoint = NewMarkPoint(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x8d96333a,
		},
		{
			name: "mark_point_axis_right",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.CategoryAxis.Position = PositionRight
				opt.SeriesList[0].MarkPoint = NewMarkPoint(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xf8c5e60,
		},
		{
			name: "category_axis_right",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.CategoryAxis.Position = PositionRight
				return opt
			},
			pngCRC: 0x199a0d02,
		},
		{
			name: "category_axis_right_with_mark_line",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.CategoryAxis.Position = PositionRight
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				return opt
			},
			pngCRC: 0x6d5fe69a,
		},
		{
			name: "null_values",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				opt.SeriesList[0].Values[2] = GetNullValue()
				opt.SeriesList[1].Values[4] = GetNullValue()
				return opt
			},
			pngCRC: 0x13855691,
		},
		{
			name: "stack_series_null",
			makeOptions: func() BarChartOption {
				opt := makeFullHorizontalBarStackedOption()
				opt.SeriesList[1].Values[3] = GetNullValue()
				return opt
			},
			pngCRC: 0x73ac3e1,
		},
		{
			name: "stack_series_global_mark_point_null_last",
			makeOptions: func() BarChartOption {
				// combined max (index 1) and min (index 3) fall where the last series is null
				opt := NewBarChartOptionWithData([][]float64{
					{10, 40, 10, 5, 10},
					{10, 40, 10, 5, 10},
					{5, GetNullValue(), 5, GetNullValue(), 5},
				})
				opt.Horizontal = true
				opt.StackSeries = Ptr(true)
				opt.SeriesList[len(opt.SeriesList)-1].MarkPoint.AddGlobalPoints(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x51906f8e,
		},
		{
			name: "stack_series_global_mark_point_null_series",
			makeOptions: func() BarChartOption {
				// last stacked series is entirely null, marks anchor to the combined stack
				opt := NewBarChartOptionWithData([][]float64{
					{10, 40, 10, 5, 10},
					{10, 40, 10, 5, 10},
					{GetNullValue(), GetNullValue(), GetNullValue(), GetNullValue(), GetNullValue()},
				})
				opt.Horizontal = true
				opt.StackSeries = Ptr(true)
				opt.SeriesList[len(opt.SeriesList)-1].MarkPoint.AddGlobalPoints(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xf3199ab4,
		},
		{
			name: "negative_values",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{-12, -5, -2, -5, -12, -19, -22, -19},
					{-16, -9, -6, -9, -16, -23, -26, -23},
				})
				opt.Horizontal = true
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0x5cbfb55b,
		},
		{
			name: "mixed_sign",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{0, 7, 10, 7, 0, -7, -10, -7},
					{10, 7, 0, -7, -10, -7, 0, 7},
				})
				opt.Horizontal = true
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0x5c9cece1,
		},
		{
			name: "mixed_sign_reversed",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{0, 7, 10, 7, 0, -7, -10, -7},
					{10, 7, 0, -7, -10, -7, 0, 7},
				})
				opt.Horizontal = true
				opt.CategoryAxis.Position = PositionRight
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0x4e0755a0,
		},
		{
			name: "stack_series_negative",
			makeOptions: func() BarChartOption {
				// mixed signs so stacks diverge above and below zero
				opt := NewBarChartOptionWithData([][]float64{
					{0, 7, 10, 7, 0, -7, -10, -7},
					{-7, -10, -7, 0, 7, 10, 7, 0},
					{3, 5, 3, 0, -3, -5, -3, 0},
				})
				opt.Horizontal = true
				opt.StackSeries = Ptr(true)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H"}
				return opt
			},
			pngCRC: 0xfecbb170,
		},
		{
			name: "rounded_caps_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{5, 10, 5, -5, -10, -5},
					{-10, -5, 5, 10, 5, -5},
				})
				opt.Horizontal = true
				opt.RoundedBarCaps = Ptr(true)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F"}
				return opt
			},
			pngCRC: 0xfafd3a0d,
		},
		{
			name: "label_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{5, 10, 5, -5, -10, -5},
					{-10, -5, 5, 10, 5, -5},
				})
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				opt.Horizontal = true
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F"}
				return opt
			},
			pngCRC: 0x5c2d58f3,
		},
		{
			name: "mark_point_negative",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{
					{5, 10, 5, -5, -10, -5},
				})
				opt.Horizontal = true
				opt.SeriesList[0].MarkPoint = NewMarkPoint(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.CategoryAxis.Labels = []string{"A", "B", "C", "D", "E", "F"}
				return opt
			},
			pngCRC: 0xf4c06245,
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		rasterOptions := PainterOptions{
			OutputFormat: ChartOutputPNG,
			Width:        600,
			Height:       400,
		}
		if tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				r := NewPainter(rasterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateBarChartRender(t, p, r, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				r := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateBarChartRender(t, p, r, opt, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				r := NewPainter(rasterOptions)

				validateBarChartRender(t, p, r, tt.makeOptions(), tt.pngCRC)
			})
		}
	}
}

func TestBarChartOptionNotMutated(t *testing.T) {
	t.Parallel()

	t.Run("vertical_value_axis", func(t *testing.T) {
		opt := NewBarChartOptionWithData([][]float64{{1, 2, 3}})
		opt.ValueAxis = []ValueAxisOption{{Position: PositionTop}}

		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		require.NoError(t, p.BarChart(opt))

		assert.Equal(t, PositionTop, opt.ValueAxis[0].Position)
		assert.Nil(t, opt.ValueAxis[0].Theme)
	})
	t.Run("horizontal_value_axis", func(t *testing.T) {
		opt := NewBarChartOptionWithData([][]float64{{1, 2, 3}})
		opt.Horizontal = true
		opt.ValueAxis = []ValueAxisOption{{Position: PositionLeft}}

		p := NewPainter(PainterOptions{Width: 600, Height: 400})
		require.NoError(t, p.BarChart(opt))

		assert.Equal(t, PositionLeft, opt.ValueAxis[0].Position)
		assert.Nil(t, opt.ValueAxis[0].Theme)
	})
}

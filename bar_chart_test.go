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
			pngCRC:      0x4fb909a1,
		},
		{
			name: "rounded_caps",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.RoundedBarCaps = Ptr(true)
				return opt
			},
			pngCRC: 0x8523cb0f,
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
			pngCRC: 0x99797c4f,
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.CategoryAxis.BoundaryGap = Ptr(true)
				opt.ValueAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xad9a8d00,
		},
		{
			name: "boundary_gap_disable",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.CategoryAxis.BoundaryGap = Ptr(false)
				return opt
			},
			pngCRC: 0x348f762e,
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
			pngCRC: 0x844c7378,
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
			pngCRC: 0x78126afc,
		},
		{
			name: "bar_width_thin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarSize = 0.018 // ~2px wide
				return opt
			},
			pngCRC: 0x5561e961,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x625488de,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarMargin = Ptr(5.0) // will be limited to fit graph
				return opt
			},
			pngCRC: 0x7ea8b5bb,
		},
		{
			name: "bar_width_and_narrow_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarSize = 0.075 // ~10px wide
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x2cd0c5e6,
		},
		{
			name: "bar_width_and_wide_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarSize = 0.075      // ~10px wide
				opt.BarMargin = Ptr(5.0) // will be limited for readability
				return opt
			},
			pngCRC: 0x11f1be5c,
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
			pngCRC: 0x25fa0e30,
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
			pngCRC: 0x825725a1,
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
			pngCRC: 0x46d0902,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullBarChartStackedOption,
			pngCRC:      0x855ffac5,
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
			pngCRC: 0xe32e452a,
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
			pngCRC: 0x14521e2c,
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
			pngCRC: 0xa03afde3,
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
			pngCRC: 0xa7c8dbb8,
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
			pngCRC: 0x135c35e9,
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
			pngCRC:      0xf27d8d42,
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
			pngCRC: 0x216be35d,
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
			pngCRC: 0xaaa09e64,
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
			pngCRC: 0xd2977ce7,
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
			pngCRC: 0x2aa51c50,
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
			pngCRC: 0xb04ecd33,
		},
		{
			name: "bar_size_thin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarSize = 0.028 // ~2px tall
				return opt
			},
			pngCRC: 0x1c0932ef,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0xc0b82a32,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarMargin = Ptr(5.0) // will be limited to fit graph
				return opt
			},
			pngCRC: 0x76c61fe0,
		},
		{
			name: "bar_size_and_narrow_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarSize = 0.117 // ~10px tall
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x134a9d3c,
		},
		{
			name: "bar_size_and_wide_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.BarSize = 0.117      // ~10px tall
				opt.BarMargin = Ptr(5.0) // will be limited for readability
				return opt
			},
			pngCRC: 0xd92bcab2,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullHorizontalBarStackedOption,
			pngCRC:      0xec9a6543,
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
			pngCRC: 0x2f4f3f65,
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
			pngCRC: 0x9425c0a7,
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
			pngCRC: 0x82eca088,
		},
		{
			name: "rounded_caps",
			makeOptions: func() BarChartOption {
				opt := makeMinimalHorizontalBarOption()
				opt.RoundedBarCaps = Ptr(true)
				return opt
			},
			pngCRC: 0xc758a0a0,
		},
		{
			name: "mark_point",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.SeriesList[0].MarkPoint = NewMarkPoint(SeriesMarkTypeMin, SeriesMarkTypeMax)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xde15ffec,
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
			pngCRC: 0xa4c035f6,
		},
		{
			name: "category_axis_right",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.CategoryAxis.Position = PositionRight
				return opt
			},
			pngCRC: 0xc3a78173,
		},
		{
			name: "category_axis_right_with_mark_line",
			makeOptions: func() BarChartOption {
				opt := makeBasicHorizontalBarOption()
				opt.CategoryAxis.Position = PositionRight
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				return opt
			},
			pngCRC: 0x6f4575e6,
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

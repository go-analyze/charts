package charts

import (
	"strconv"
	"testing"

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
		XAxis: XAxisOption{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		},
		YAxis: []YAxisOption{
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
	opt.XAxis = XAxisOption{
		Show:   Ptr(false),
		Labels: []string{"A", "B"},
	}
	opt.YAxis[0].Show = Ptr(false)
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
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		},
		Legend: LegendOption{
			SeriesNames: []string{"A", "B", "C"},
			Symbol:      SymbolDot,
		},
		YAxis: []YAxisOption{
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
	assert.Len(t, opt.YAxis, 1)
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
			pngCRC:      0xf99b3b59,
		},
		{
			name: "rounded_caps",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.RoundedBarCaps = Ptr(true)
				return opt
			},
			pngCRC: 0xb177bec7,
		},
		{
			name: "custom_font",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.XAxis.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0xdc9be388,
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.XAxis.BoundaryGap = Ptr(true)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xab7b0063,
		},
		{
			name: "boundary_gap_disable",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.XAxis.BoundaryGap = Ptr(false)
				return opt
			},
			pngCRC: 0xd56b1106,
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
			pngCRC: 0xae4b0308,
		},
		{
			name: "bar_width_truncate",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				opt.BarWidth = 1000
				return opt
			},
			pngCRC: 0x78126afc,
		},
		{
			name: "bar_width_thin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarWidth = 2
				return opt
			},
			pngCRC: 0xa9ef8762,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x83daa059,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarMargin = Ptr(1000.0) // will be limited to fit graph
				return opt
			},
			pngCRC: 0x2f5d0556,
		},
		{
			name: "bar_width_and_narrow_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarWidth = 10
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x882b6b42,
		},
		{
			name: "bar_width_and_wide_margin",
			makeOptions: func() BarChartOption {
				opt := makeMinimalBarChartOption()
				opt.BarWidth = 10
				opt.BarMargin = Ptr(1000.0) // will be limited for readability
				return opt
			},
			pngCRC: 0xea22413e,
		},
		{
			name: "dual_yaxis",
			makeOptions: func() BarChartOption {
				opt := makeBasicBarChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.Title.Text = "T"
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				return opt
			},
			pngCRC: 0xa34f80b9,
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
			pngCRC: 0x88db720a,
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
			pngCRC: 0x71230483,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullBarChartStackedOption,
			pngCRC:      0x86614a2,
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
				opt.YAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x58590bd4,
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
				opt.YAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x6ed1ee54,
		},
		{
			name: "stack_series_simple",
			makeOptions: func() BarChartOption {
				opt := NewBarChartOptionWithData([][]float64{{4.0}, {1.0}})
				opt.StackSeries = Ptr(true)
				// disable extra
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
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
				opt.YAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x75d281c8,
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

func TestCalculateBarMarginsAndSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		seriesCount     int
		space           int
		barWidth        int
		barMargin       *float64
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
			barWidth:        20,
			barMargin:       Ptr(4.0),
			expectedMargin:  23,
			expectedSpacing: 4,
			expectedSize:    20,
		},
		{
			name:            "width_too_wide",
			seriesCount:     2,
			space:           40,
			barWidth:        20,
			barMargin:       Ptr(1.0),
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
			barWidth:        10,
			barMargin:       Ptr(0.0),
			expectedMargin:  20,
			expectedSpacing: 0,
			expectedSize:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			margin, spacing, size := calculateBarMarginsAndSize(tt.seriesCount, tt.space, tt.barWidth, tt.barMargin)
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
			name: "empty_series",
			makeOptions: func() BarChartOption {
				return NewBarChartOptionWithData([][]float64{})
			},
			errorMsgContains: "empty series list",
		},
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

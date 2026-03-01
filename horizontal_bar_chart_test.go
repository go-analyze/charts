package charts

import (
	"strconv"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicHorizontalBarChartOption() HorizontalBarChartOption {
	return HorizontalBarChartOption{
		Padding: NewBoxEqual(10),
		SeriesList: NewSeriesListHorizontalBar([][]float64{
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
		YAxis: YAxisOption{
			Labels: []string{"Brazil", "Indonesia", "USA", "India", "China", "World"},
		},
	}
}

func makeMinimalHorizontalBarChartOption() HorizontalBarChartOption {
	opt := NewHorizontalBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})
	opt.YAxis = YAxisOption{
		Show:   Ptr(false),
		Labels: []string{"A", "B"},
	}
	opt.XAxis.Show = Ptr(false)
	return opt
}

func makeFullHorizontalBarChartStackedOption() HorizontalBarChartOption {
	seriesList := NewSeriesListHorizontalBar([][]float64{
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
	return HorizontalBarChartOption{
		Padding:     NewBoxEqual(20),
		SeriesList:  seriesList,
		StackSeries: Ptr(true),
		Legend: LegendOption{
			Symbol: SymbolDot,
		},
		YAxis: YAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		},
	}
}

func TestNewHorizontalBarChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewHorizontalBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeHorizontalBar, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.HorizontalBarChart(opt))
}

func TestHorizontalBarChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() HorizontalBarChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicHorizontalBarChartOption,
			pngCRC:      0xfb168012,
		},
		{
			name: "custom_fonts",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.XAxis.FontStyle = customFont
				opt.YAxis.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0x7d773c06,
		},
		{
			name: "value_labels",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
					series[i].Label.ValueFormatter = func(f float64) string {
						return humanize.FtoaWithDigits(f, 2)
					}
				}
				return opt
			},
			pngCRC: 0xdda92f8,
		},
		{
			name: "value_formatter",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
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
			pngCRC: 0xd72c917f,
		},
		{
			name: "bar_height_truncate",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				opt.BarHeight = 1000
				return opt
			},
			pngCRC: 0xb35224f4,
		},
		{
			name: "mark_line",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				opt.YAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x56386ea2,
		},
		{
			name: "bar_height_thin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 2
				return opt
			},
			pngCRC: 0xedf8c602,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x62e361ae,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = Ptr(1000.0) // will be limited to fit graph
				return opt
			},
			pngCRC: 0xe22b0ccd,
		},
		{
			name: "bar_height_and_narrow_margin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			pngCRC: 0x8a6043ab,
		},
		{
			name: "bar_height_and_wide_margin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = Ptr(1000.0) // will be limited for readability
				return opt
			},
			pngCRC: 0x56436af8,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullHorizontalBarChartStackedOption,
			pngCRC:      0xf75e263c,
		},
		{
			name: "stack_series_simple",
			makeOptions: func() HorizontalBarChartOption {
				opt := NewHorizontalBarChartOptionWithData([][]float64{{4.0}, {1.0}})
				opt.StackSeries = Ptr(true)
				opt.XAxis.Unit = 1
				// disable extra
				opt.YAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x2f4f3f65,
		},
		{
			name: "stack_series_with_mark",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeFullHorizontalBarChartStackedOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine = NewMarkLine(SeriesMarkTypeMax)
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine.Lines[0].Global = true
				opt.YAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x2e00befd,
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
				rp := NewPainter(rasterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateHorizontalBarChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateHorizontalBarChartRender(t, p, rp, opt, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateHorizontalBarChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
		}
	}
}

func validateHorizontalBarChartRender(t *testing.T, svgP, pngP *Painter, opt HorizontalBarChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.HorizontalBarChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.HorizontalBarChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestHorizontalBarChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() HorizontalBarChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() HorizontalBarChartOption {
				return NewHorizontalBarChartOptionWithData([][]float64{})
			},
			errorMsgContains: "empty series list",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.HorizontalBarChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

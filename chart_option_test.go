package charts

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChartOption(t *testing.T) {
	t.Parallel()

	fns := []OptionFunc{
		SVGOutputOptionFunc(),
		FontOptionFunc(GetDefaultFont()),
		ThemeNameOptionFunc(ThemeVividDark),
		TitleTextOptionFunc("title"),
		LegendLabelsOptionFunc([]string{"label"}),
		XAxisLabelsOptionFunc([]string{"xaxis"}),
		YAxisLabelsOptionFunc([]string{"yaxis"}),
		DimensionsOptionFunc(800, 600),
		PaddingOptionFunc(NewBoxEqual(10)),
	}
	opt := ChartOption{}
	for _, fn := range fns {
		fn(&opt)
	}
	require.Equal(t, ChartOption{
		OutputFormat: ChartOutputSVG,
		Font:         GetDefaultFont(),
		Theme:        GetTheme(ThemeVividDark),
		Title: TitleOption{
			Text: "title",
		},
		Legend: LegendOption{
			SeriesNames: []string{"label"},
		},
		XAxis: XAxisOption{
			Labels: []string{"xaxis"},
		},
		YAxis: []YAxisOption{
			{
				Labels: []string{"yaxis"},
			},
		},
		Width:   800,
		Height:  600,
		Padding: NewBoxEqual(10),
	}, opt)

	makeInvalidYAxisOption := func(axis int) ChartOption {
		invalidOpt := ChartOption{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
			SeriesList: NewSeriesListLine([][]float64{{1, 2, 3}}, LineSeriesOption{
				Names: []string{"Series A"},
			}).ToGenericSeriesList(),
		}
		invalidOpt.SeriesList[0].YAxisIndex = axis
		return invalidOpt
	}

	t.Run("invalid_yaxis_index_returns_error", func(t *testing.T) {
		_, err := Render(makeInvalidYAxisOption(2))
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid y-axis index")
	})

	t.Run("negative_yaxis_index_returns_error", func(t *testing.T) {
		_, err := Render(makeInvalidYAxisOption(-1))
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid y-axis index")
	})
}

func TestChartOptionSeriesShowLabel(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: NewSeriesListPie([]float64{1, 2}).ToGenericSeriesList(),
	}
	SeriesShowLabel(true)(&opt)
	assert.True(t, flagIs(true, opt.SeriesList[0].Label.Show))

	SeriesShowLabel(false)(&opt)
	assert.True(t, flagIs(false, opt.SeriesList[0].Label.Show))
}

func newNoTypeSeriesListFromValues(values [][]float64) GenericSeriesList {
	return NewSeriesListGeneric(values, "")
}

func TestChartOptionMarkLine(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: newNoTypeSeriesListFromValues([][]float64{{1, 2}}),
	}
	MarkLineOptionFunc(0, "min", "max")(&opt)
	assert.Equal(t, NewMarkLine("min", "max"), opt.SeriesList[0].MarkLine)
}

func TestChartOptionMarkPoint(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: newNoTypeSeriesListFromValues([][]float64{{1, 2}}),
	}
	MarkPointOptionFunc(0, "min", "max")(&opt)
	assert.Equal(t, NewMarkPoint("min", "max"), opt.SeriesList[0].MarkPoint)
}

func TestLineRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	p, err := LineRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Line"),
		XAxisLabelsOptionFunc([]string{
			"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
		}),
		LegendLabelsOptionFunc([]string{
			"Email", "Union Ads", "Video Ads", "Direct", "Search Engine",
		}),
		func(opt *ChartOption) {
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestScatterRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	p, err := ScatterRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Scatter"),
		XAxisLabelsOptionFunc([]string{
			"1", "2", "3", "4", "5", "6", "7",
		}),
		LegendLabelsOptionFunc([]string{
			"A", "B", "C", "D", "E",
		}),
		func(opt *ChartOption) {
			opt.Symbol = SymbolSquare
			opt.Legend.Symbol = SymbolSquare
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestBarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3},
	}
	p, err := BarRender(
		values,
		SVGOutputOptionFunc(),
		XAxisLabelsOptionFunc([]string{
			"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		}),
		LegendLabelsOptionFunc([]string{
			"Rainfall", "Evaporation",
		}),
		MarkLineOptionFunc(0, SeriesMarkTypeAverage),
		MarkPointOptionFunc(0, SeriesMarkTypeMax, SeriesMarkTypeMin),
		// custom option func
		func(opt *ChartOption) {
			opt.Legend.Offset = OffsetRight
			opt.Legend.OverlayChart = Ptr(true)
			opt.SeriesList[1].MarkPoint = NewMarkPoint(SeriesMarkTypeMax, SeriesMarkTypeMin)
			opt.SeriesList[1].MarkLine = NewMarkLine(SeriesMarkTypeAverage)
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestHorizontalBarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{18203, 23489, 29034, 104970, 131744, 630230},
		{19325, 23438, 31000, 121594, 134141, 681807},
	}
	p, err := HorizontalBarRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("World Population"),
		PaddingOptionFunc(Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		LegendLabelsOptionFunc([]string{"2011", "2012"}),
		YAxisLabelsOptionFunc([]string{"Brazil", "Indonesia", "USA", "India", "China", "World"}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestPieRender(t *testing.T) {
	t.Parallel()

	values := []float64{1048, 735, 580, 484, 300}
	p, err := PieRender(
		values,
		SVGOutputOptionFunc(),
		TitleOptionFunc(TitleOption{
			Text:    "Rainfall vs Evaporation",
			Subtext: "Fake Data",
			Offset:  OffsetCenter,
		}),
		PaddingOptionFunc(NewBoxEqual(20)),
		LegendOptionFunc(LegendOption{
			Vertical:    Ptr(true),
			SeriesNames: []string{"Search Engine", "Direct", "Email", "Union Ads", "Video Ads"},
			Offset:      OffsetLeft,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestDoughnutRender(t *testing.T) {
	t.Parallel()

	values := []float64{1048, 735, 580, 484, 300}
	p, err := DoughnutRender(
		values,
		SVGOutputOptionFunc(),
		TitleOptionFunc(TitleOption{
			Text:    "Title",
			Subtext: "Fake Data",
			Offset:  OffsetCenter,
		}),
		PaddingOptionFunc(NewBoxEqual(20)),
		LegendOptionFunc(LegendOption{
			Vertical:    Ptr(true),
			SeriesNames: []string{"A", "B", "C", "D", "E"},
			Offset:      OffsetLeft,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestRadarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{4200, 3000, 20000, 35000, 50000, 18000},
		{5000, 14000, 28000, 26000, 42000, 21000},
	}
	p, err := RadarRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Basic Radar Chart"),
		LegendLabelsOptionFunc([]string{
			"Allocated Budget", "Actual Spending",
		}),
		RadarIndicatorOptionFunc([]string{
			"Sales",
			"Administration",
			"Information Technology",
			"Customer Support",
			"Development",
			"Marketing",
		}, []float64{
			6500, 16000, 30000, 38000, 52000, 25000,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestFunnelRender(t *testing.T) {
	t.Parallel()

	values := []float64{
		100, 80, 60, 40, 20,
	}
	p, err := FunnelRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Funnel"),
		LegendLabelsOptionFunc([]string{
			"Show", "Click", "Visit", "Inquiry", "Order",
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestChildRender(t *testing.T) {
	p, err := LineRender(
		[][]float64{
			{120, 132, 101, 134, 90, 230, 210},
			{150, 232, 201, 154, 190, 330, 410},
			{320, 332, 301, 334, 390, 330, 320},
		},
		SVGOutputOptionFunc(),
		XAxisLabelsOptionFunc([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}),
		ChildOptionFunc(ChartOption{
			Box: NewBox(200, 10, 500, 200),
			SeriesList: NewSeriesListHorizontalBar([][]float64{
				{70, 90, 110, 130},
				{80, 100, 120, 140},
			}).ToGenericSeriesList(),
			Legend: LegendOption{
				SeriesNames: []string{"2011", "2012"},
			},
			YAxis: []YAxisOption{
				{
					Labels: []string{"USA", "India", "China", "World"},
				},
			},
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestChartCombos(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		optFunc func() ChartOption
	}{
		{
			name: "line+scatter",
			optFunc: func() ChartOption {
				lineValues := [][]float64{
					{120, 132, 101, 134, 90, 230, 210},
				}
				scatterValues := [][]float64{
					{180, 200, 150, 180, 160, 280, 250},
				}

				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				scatterSeries := NewSeriesListScatter(scatterValues, ScatterSeriesOption{
					Names: []string{"Scatter Data"},
				}).ToGenericSeriesList()

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(500.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Line Data", "Scatter Data"},
					},
					SeriesList: append(lineSeries, scatterSeries...),
				}
			},
		},
		{
			name: "line+bar",
			optFunc: func() ChartOption {
				lineValues := [][]float64{
					{120, 132, 101, 134, 90, 230, 210},
				}
				barValues := [][]float64{
					{70, 90, 110, 130, 80, 100, 120},
				}

				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				barSeries := NewSeriesListBar(barValues, BarSeriesOption{
					Names: []string{"Bar Data"},
				}).ToGenericSeriesList()

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(400.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Line Data", "Bar Data"},
					},
					SeriesList: append(lineSeries, barSeries...),
				}
			},
		},
		{
			name: "line+candlestick",
			optFunc: func() ChartOption {
				lineValues := [][]float64{
					{25, 28, 23, 30, 27, 32, 29},
				}
				candlestickData := [][]OHLCData{
					{
						{Open: 24, High: 30, Low: 20, Close: 25},
						{Open: 25, High: 32, Low: 22, Close: 28},
						{Open: 28, High: 30, Low: 18, Close: 23},
						{Open: 23, High: 35, Low: 25, Close: 30},
						{Open: 30, High: 33, Low: 24, Close: 27},
						{Open: 27, High: 36, Low: 28, Close: 32},
						{Open: 32, High: 34, Low: 26, Close: 29},
					},
				}

				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				candlestickSeries := NewSeriesListCandlestick(candlestickData, CandlestickSeriesOption{
					Names: []string{"OHLC"},
				}).ToGenericSeriesList()

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(40.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Line Data", "OHLC"},
					},
					SeriesList: append(lineSeries, candlestickSeries...),
				}
			},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			opt := tc.optFunc()
			opt.OutputFormat = ChartOutputSVG

			p, err := Render(opt)
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

func TestDualAxisChartCombos(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		optFunc func() ChartOption
	}{
		{
			name: "scatter+line",
			optFunc: func() ChartOption {
				scatterValues := [][]float64{
					{180, 200, 150, 180, 160, 280, 250},
				}
				lineValues := [][]float64{
					{1.2, 1.5, 1.1, 1.8, 1.4, 2.1, 1.7},
				}

				scatterSeries := NewSeriesListScatter(scatterValues, ScatterSeriesOption{
					Names: []string{"Scatter Data"},
				}).ToGenericSeriesList()
				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				lineSeries[0].YAxisIndex = 1

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(300.0),
						},
						{
							Max: Ptr(3.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Scatter Data", "Line Data"},
					},
					SeriesList: append(scatterSeries, lineSeries...),
				}
			},
		},
		{
			name: "bar+line",
			optFunc: func() ChartOption {
				barValues := [][]float64{
					{70, 90, 110, 130, 80, 100, 120},
				}
				lineValues := [][]float64{
					{1.2, 1.8, 2.3, 1.9, 2.5, 1.6, 2.1},
				}

				barSeries := NewSeriesListBar(barValues, BarSeriesOption{
					Names: []string{"Bar Data"},
				}).ToGenericSeriesList()
				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				lineSeries[0].YAxisIndex = 1

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(150.0),
						},
						{
							Max: Ptr(3.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Bar Data", "Line Data"},
					},
					SeriesList: append(barSeries, lineSeries...),
				}
			},
		},
		{
			name: "candlestick+line",
			optFunc: func() ChartOption {
				candlestickData := [][]OHLCData{
					{
						{Open: 24, High: 30, Low: 20, Close: 25},
						{Open: 25, High: 32, Low: 22, Close: 28},
						{Open: 28, High: 30, Low: 18, Close: 23},
						{Open: 23, High: 35, Low: 25, Close: 30},
						{Open: 30, High: 33, Low: 24, Close: 27},
						{Open: 27, High: 36, Low: 28, Close: 32},
						{Open: 32, High: 34, Low: 26, Close: 29},
					},
				}
				lineValues := [][]float64{
					{1200, 1400, 1100, 1600, 1300, 1800, 1500},
				}

				candlestickSeries := NewSeriesListCandlestick(candlestickData, CandlestickSeriesOption{
					Names: []string{"OHLC"},
				}).ToGenericSeriesList()
				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Volume"},
				}).ToGenericSeriesList()
				lineSeries[0].YAxisIndex = 1

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(40.0),
						},
						{
							Max: Ptr(2000.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"OHLC", "Volume"},
					},
					SeriesList: append(candlestickSeries, lineSeries...),
				}
			},
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			opt := tc.optFunc()
			opt.OutputFormat = ChartOutputSVG

			p, err := Render(opt)
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

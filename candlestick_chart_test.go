package charts

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicCandlestickData() []OHLCData {
	return []OHLCData{
		{Open: 100, High: 110, Low: 95, Close: 105},
		{Open: 105, High: 115, Low: 100, Close: 112},
		{Open: 112, High: 118, Low: 108, Close: 115},
		{Open: 115, High: 120, Low: 105, Close: 108}, // bearish
		{Open: 108, High: 113, Low: 105, Close: 109},
	}
}

func makeBasicCandlestickChartOption() CandlestickChartOption {
	return CandlestickChartOption{
		Title: TitleOption{
			Text: "Candlestick Chart",
		},
		Theme:   GetTheme(ThemeVividLight),
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May"},
		},
		YAxis: []YAxisOption{
			{
				PreferNiceIntervals: Ptr(true),
			},
		},
		Legend: LegendOption{
			SeriesNames: []string{"Price"},
		},
		SeriesList: CandlestickSeriesList{{Data: makeBasicCandlestickData()}},
	}
}

func makeMinimalCandlestickChartOption() CandlestickChartOption {
	return CandlestickChartOption{
		Theme:   GetTheme(ThemeVividLight),
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5"},
			Show:   Ptr(false),
		},
		YAxis: []YAxisOption{
			{
				PreferNiceIntervals: Ptr(true),
			},
		},
		SeriesList: CandlestickSeriesList{{Data: makeBasicCandlestickData()}},
	}
}

func TestCandlestickChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		makeOptions func() CandlestickChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic",
			makeOptions: makeBasicCandlestickChartOption,
			pngCRC:      0x85c13d24,
		},
		{
			name:        "minimal",
			makeOptions: makeMinimalCandlestickChartOption,
			pngCRC:      0x6f23a53d,
		},
		{
			name: "custom_style",
			makeOptions: func() CandlestickChartOption {
				opt := makeBasicCandlestickChartOption()
				opt.CandleWidth = 0.5
				opt.WickWidth = 2.0
				opt.SeriesList[0].ShowWicks = Ptr(false)
				return opt
			},
			pngCRC: 0xbe8ea1f8,
		},
		{
			name: "doji",
			makeOptions: func() CandlestickChartOption {
				opt := makeBasicCandlestickChartOption()
				data := makeBasicCandlestickData()
				if len(data) > 0 {
					data[0] = OHLCData{Open: 100, High: 110, Low: 95, Close: 100}
				}
				opt.SeriesList[0] = CandlestickSeries{Data: data}
				return opt
			},
			pngCRC: 0x930d420a,
		},
		{
			name: "dual_axis",
			makeOptions: func() CandlestickChartOption {
				opt := makeBasicCandlestickChartOption()
				secondSeriesData := []OHLCData{
					{Open: 200, High: 220, Low: 190, Close: 210},
					{Open: 210, High: 230, Low: 200, Close: 225},
					{Open: 225, High: 240, Low: 215, Close: 230},
					{Open: 230, High: 245, Low: 210, Close: 215},
					{Open: 215, High: 225, Low: 205, Close: 220},
				}
				// Add second Y axis and second series
				opt.YAxis = append(opt.YAxis, YAxisOption{
					PreferNiceIntervals: Ptr(true),
				})
				opt.SeriesList = append(opt.SeriesList, CandlestickSeries{
					Data:       secondSeriesData,
					YAxisIndex: 1,
					Name:       "Volume",
				})
				// Update legend to show both series
				opt.Legend.SeriesNames = []string{"Price", "Volume"}
				return opt
			},
			pngCRC: 0x8df4ee78,
		},
		{
			name: "trend_lines",
			makeOptions: func() CandlestickChartOption {
				opt := makeBasicCandlestickChartOption()
				opt.SeriesList[0].OpenTrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeSMA, Period: 2, LineColor: ColorRed},
				}
				opt.SeriesList[0].HighTrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeEMA, Period: 4, LineColor: ColorRedAlt1},
				}
				opt.SeriesList[0].LowTrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeEMA, Period: 4, LineColor: ColorBlue},
				}
				opt.SeriesList[0].CloseTrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeSMA, Period: 2, LineColor: ColorGreen},
				}
				return opt
			},
			pngCRC: 0xd2f00aca,
		},
		{
			name: "mark_lines",
			makeOptions: func() CandlestickChartOption {
				series := CandlestickSeries{
					Data: makeBasicCandlestickData(),
					OpenMarkLine: SeriesMarkLine{
						Lines: []SeriesMark{
							{Type: SeriesMarkTypeMax}, // Open resistance
						},
					},
					HighMarkLine: SeriesMarkLine{
						Lines: []SeriesMark{
							{Type: SeriesMarkTypeMax}, // Absolute high
						},
					},
					LowMarkLine: SeriesMarkLine{
						Lines: []SeriesMark{
							{Type: SeriesMarkTypeMin}, // Absolute low
						},
					},
					CloseMarkLine: SeriesMarkLine{
						Lines: []SeriesMark{
							{Type: SeriesMarkTypeAverage}, // Close average
						},
					},
				}
				return CandlestickChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"Jan", "Feb", "Mar", "Apr", "May"},
					},
					YAxis:      make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{series},
				}
			},
			pngCRC: 0x577c3666,
		},
		{
			name: "mark_point",
			makeOptions: func() CandlestickChartOption {
				series := CandlestickSeries{
					Data: makeBasicCandlestickData(),
					OpenMarkPoint: SeriesMarkPoint{
						Points: []SeriesMark{
							{Type: SeriesMarkTypeMax}, // Highest open
						},
					},
					HighMarkPoint: SeriesMarkPoint{
						Points: []SeriesMark{
							{Type: SeriesMarkTypeMax}, // Absolute maximum
						},
					},
					LowMarkPoint: SeriesMarkPoint{
						Points: []SeriesMark{
							{Type: SeriesMarkTypeMin}, // Absolute minimum
						},
					},
					CloseMarkPoint: SeriesMarkPoint{
						Points: []SeriesMark{
							{Type: SeriesMarkTypeMin}, // Lowest close
						},
					},
				}
				return CandlestickChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"Jan", "Feb", "Mar", "Apr", "May"},
					},
					YAxis:      make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{series},
				}
			},
			pngCRC: 0x9aa8a7fd,
		},
		{
			name: "patterns_replace_mode",
			makeOptions: func() CandlestickChartOption {
				series := CandlestickSeries{
					Data: []OHLCData{
						{Open: 100, High: 110, Low: 95, Close: 100.01}, // Doji: 0.01/15 = 0.0007 < 0.001
						{Open: 105, High: 115, Low: 100, Close: 112},   // Normal
					},
					Label: SeriesLabel{
						Show: Ptr(true),
						LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
							return "User Label", nil
						},
					},
					PatternConfig: (&CandlestickPatternConfig{}).
						WithPreferPatternLabels(true).
						WithDoji().
						WithDojiThreshold(0.001).
						WithShadowTolerance(0.01).
						WithShadowRatio(2.0).
						WithEngulfingMinSize(0.8),
				}
				return CandlestickChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"A", "B"},
					},
					YAxis:      make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{series},
				}
			},
			pngCRC: 0x161b609d,
		},
		{
			name: "large_dataset",
			makeOptions: func() CandlestickChartOption {
				var data []OHLCData
				for i := 0; i < 50; i++ {
					basePrice := 100.0 + float64(i)*0.5
					data = append(data, OHLCData{
						Open:  basePrice,
						High:  basePrice + 5,
						Low:   basePrice - 3,
						Close: basePrice + 2,
					})
				}
				return CandlestickChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Show: Ptr(false), // Hide labels for large dataset
					},
					YAxis:      make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{{Data: data}},
				}
			},
			pngCRC: 0x1d24c9ce,
		},
		{
			name: "multiple_series",
			makeOptions: func() CandlestickChartOption {
				series1Data := []OHLCData{
					{Open: 100.0, High: 110.0, Low: 95.0, Close: 105.0},
					{Open: 105.0, High: 115.0, Low: 100.0, Close: 112.0},
					{Open: 112.0, High: 118.0, Low: 108.0, Close: 115.0},
					{Open: 115.0, High: 120.0, Low: 110.0, Close: 108.0},
					{Open: 108.0, High: 113.0, Low: 105.0, Close: 109.0},
				}
				series2Data := []OHLCData{
					{Open: 120.0, High: 130.0, Low: 115.0, Close: 125.0},
					{Open: 125.0, High: 135.0, Low: 120.0, Close: 132.0},
					{Open: 132.0, High: 138.0, Low: 128.0, Close: 135.0},
					{Open: 135.0, High: 140.0, Low: 130.0, Close: 128.0},
					{Open: 128.0, High: 133.0, Low: 125.0, Close: 129.0},
				}
				series3Data := []OHLCData{
					{Open: 80.0, High: 110.0, Low: 45.0, Close: 85.0},
					{Open: 85.0, High: 115.0, Low: 40.0, Close: 82.0},
					{Open: 82.0, High: 118.0, Low: 48.0, Close: 85.0},
					{Open: 85.0, High: 120.0, Low: 40.0, Close: 88.0},
					{Open: 88.0, High: 113.0, Low: 45.0, Close: 89.0},
				}
				return CandlestickChartOption{
					XAxis: XAxisOption{
						Labels: []string{"Day 1", "Day 2", "Day 3", "Day 4", "Day 5"},
					},
					YAxis: make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{
						{Data: series1Data, Name: "Stock A"},
						{Data: series2Data, Name: "Stock B"},
						{Data: series3Data, Name: "Stock C"},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Stock A", "Stock B", "Stock C"},
						Show:        Ptr(true),
					},
					Padding: NewBoxEqual(10),
				}
			},
			pngCRC: 0x7ed89e2b,
		},
		{
			name: "bollinger_bands",
			makeOptions: func() CandlestickChartOption {
				// Create longer dataset for meaningful Bollinger Bands
				ohlcData := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},
					{Open: 105, High: 115, Low: 100, Close: 112},
					{Open: 112, High: 118, Low: 108, Close: 115},
					{Open: 115, High: 120, Low: 110, Close: 118},
					{Open: 118, High: 125, Low: 115, Close: 122},
					{Open: 122, High: 128, Low: 119, Close: 125},
					{Open: 125, High: 130, Low: 122, Close: 127},
					{Open: 127, High: 132, Low: 124, Close: 129},
					{Open: 129, High: 135, Low: 126, Close: 131},
					{Open: 131, High: 138, Low: 128, Close: 135},
				}
				return CandlestickChartOption{
					XAxis: XAxisOption{
						Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
					},
					YAxis: []YAxisOption{
						{
							PreferNiceIntervals: Ptr(true),
						},
					},
					SeriesList: CandlestickSeriesList{{
						Data: ohlcData,
						CloseTrendLine: []SeriesTrendLine{
							{Type: SeriesTrendTypeBollingerUpper, Period: 5},
							{Type: SeriesTrendTypeSMA, Period: 5},
							{Type: SeriesTrendTypeBollingerLower, Period: 5},
						},
					}},
					Legend:  LegendOption{Show: Ptr(false)},
					Padding: NewBoxEqual(10),
				}
			},
			pngCRC: 0x36e96980,
		},
		{
			name: "aggregation",
			makeOptions: func() CandlestickChartOption {
				// Create longer dataset to test aggregation
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Period 1
					{Open: 105, High: 115, Low: 100, Close: 112}, // Period 1
					{Open: 112, High: 118, Low: 108, Close: 115}, // Period 2
					{Open: 115, High: 120, Low: 110, Close: 118}, // Period 2
					{Open: 118, High: 125, Low: 115, Close: 122}, // Period 3
					{Open: 122, High: 128, Low: 119, Close: 125}, // Period 3
				}
				series := CandlestickSeries{Data: data, Name: "1-Period"}
				aggregated := AggregateCandlestick(series, 2) // Aggregate into 2-period candles

				return CandlestickChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"Period 1", "Period 2", "Period 3"},
					},
					YAxis: []YAxisOption{
						{
							PreferNiceIntervals: Ptr(true),
						},
					},
					SeriesList: CandlestickSeriesList{aggregated},
				}
			},
			pngCRC: 0x7090b691,
		},
		{
			name: "large_series_count",
			makeOptions: func() CandlestickChartOption {
				const seriesCount = 10
				const dataPointsPerSeries = 5

				var seriesList CandlestickSeriesList
				var seriesNames []string
				for i := 0; i < seriesCount; i++ {
					basePrice := 100.0 + float64(i*20) // Different price ranges
					data := make([]OHLCData, dataPointsPerSeries)

					for j := 0; j < dataPointsPerSeries; j++ {
						open := basePrice + float64(j*5)
						high := open + 10.0
						low := open - 5.0
						close := open + float64((j%2)*10-5) // Alternating up/down

						data[j] = OHLCData{
							Open:  open,
							High:  high,
							Low:   low,
							Close: close,
						}
					}

					series := CandlestickSeries{
						Data: data,
						Name: fmt.Sprintf("Series %d", i+1),
					}
					seriesList = append(seriesList, series)
					seriesNames = append(seriesNames, series.Name)
				}

				return CandlestickChartOption{
					XAxis: XAxisOption{
						Labels: []string{"T1", "T2", "T3", "T4", "T5"},
					},
					YAxis:      make([]YAxisOption, 1),
					SeriesList: seriesList,
					Legend: LegendOption{
						SeriesNames: seriesNames,
						Show:        Ptr(true),
					},
					Padding: NewBoxEqual(10),
				}
			},
			pngCRC: 0x31d4c785,
		},
		{
			name: "series_styles",
			makeOptions: func() CandlestickChartOption {
				// Create different datasets for multiple series with different styles
				series1Data := []OHLCData{
					{Open: 100.0, High: 110.0, Low: 95.0, Close: 105.0},
					{Open: 105.0, High: 115.0, Low: 100.0, Close: 112.0},
					{Open: 112.0, High: 118.0, Low: 108.0, Close: 115.0},
				}
				series2Data := []OHLCData{
					{Open: 150.0, High: 160.0, Low: 145.0, Close: 155.0},
					{Open: 155.0, High: 165.0, Low: 150.0, Close: 162.0},
					{Open: 162.0, High: 168.0, Low: 158.0, Close: 165.0},
				}
				return CandlestickChartOption{
					XAxis: XAxisOption{
						Labels: []string{"Day 1", "Day 2", "Day 3"},
					},
					YAxis: make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{
						{
							Data:        series1Data,
							Name:        "Stock A (Filled)",
							CandleStyle: CandleStyleFilled,
						},
						{
							Data:        series2Data,
							Name:        "Stock B (Traditional)",
							CandleStyle: CandleStyleTraditional,
						},
					},
					Legend:  LegendOption{Show: Ptr(false)},
					Padding: NewBoxEqual(10),
				}
			},
			pngCRC: 0x60a5d5cd,
		},
		{
			name: "candle_margin_zero",
			makeOptions: func() CandlestickChartOption {
				// Create multiple series with similar price ranges to show candlesticks side-by-side
				series1Data := []OHLCData{
					{Open: 100.0, High: 110.0, Low: 95.0, Close: 105.0},
					{Open: 105.0, High: 115.0, Low: 100.0, Close: 112.0},
					{Open: 112.0, High: 118.0, Low: 108.0, Close: 115.0},
				}
				series2Data := []OHLCData{
					{Open: 102.0, High: 112.0, Low: 97.0, Close: 107.0},
					{Open: 107.0, High: 117.0, Low: 102.0, Close: 114.0},
					{Open: 114.0, High: 120.0, Low: 110.0, Close: 117.0},
				}
				return CandlestickChartOption{
					XAxis: XAxisOption{
						Labels: []string{"Day 1", "Day 2", "Day 3"},
					},
					YAxis: make([]YAxisOption, 1),
					SeriesList: CandlestickSeriesList{
						{Data: series1Data, Name: "Stock A"},
						{Data: series2Data, Name: "Stock B"},
					},
					CandleMargin: Ptr(0.0),
					Legend:       LegendOption{Show: Ptr(false)},
					Padding:      NewBoxEqual(10),
				}
			},
			pngCRC: 0x9bb012da,
		},
		{
			name: "null_values",
			makeOptions: func() CandlestickChartOption {
				opt := makeBasicCandlestickChartOption()
				opt.SeriesList[0] = CandlestickSeries{Data: []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105}, // Valid
					{Open: GetNullValue(), High: GetNullValue(), Low: GetNullValue(), Close: GetNullValue()},
					{Open: 112, High: 118, Low: 108, Close: 115},            // Valid
					{Open: GetNullValue(), High: 108, Low: 105, Close: 108}, // Partial null (invalid)
					{Open: 115, High: GetNullValue(), Low: 105, Close: 108}, // Partial null (invalid)
					{Open: 108, High: 113, Low: 105, Close: 109},            // Valid
				}}
				return opt
			},
			pngCRC: 0xe7a6d8e7,
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        800,
				Height:       600,
			})
			r := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        800,
				Height:       600,
			})

			opt := tc.makeOptions()

			validateCandlestickChartRender(t, p, r, opt, tc.pngCRC)
		})
	}
}

func TestCandlestickChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() CandlestickChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series_list",
			makeOptions: func() CandlestickChartOption {
				return NewCandlestickOptionWithSeries()
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "empty_series",
			makeOptions: func() CandlestickChartOption {
				return NewCandlestickOptionWithSeries(CandlestickSeries{})
			},
			errorMsgContains: "no data in any series",
		},
		{
			name: "invalid_yaxis_index",
			makeOptions: func() CandlestickChartOption {
				opt := NewCandlestickOptionWithSeries(CandlestickSeries{
					Data: []OHLCData{{Open: 10, High: 12, Low: 9, Close: 11}},
				})
				opt.SeriesList[0].YAxisIndex = 2
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
		{
			name: "negative_yaxis_index",
			makeOptions: func() CandlestickChartOption {
				opt := NewCandlestickOptionWithSeries(CandlestickSeries{
					Data: []OHLCData{{Open: 10, High: 12, Low: 9, Close: 11}},
				})
				opt.SeriesList[0].YAxisIndex = -1
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			painterOptions := PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}
			p := NewPainter(painterOptions)
			opt := tt.makeOptions()

			err := p.CandlestickChart(opt)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsgContains)
		})
	}
}

func TestOHLCDataValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ohlc     OHLCData
		expected bool
	}{
		{
			name:     "valid_bullish",
			ohlc:     OHLCData{Open: 100, High: 110, Low: 95, Close: 105},
			expected: true,
		},
		{
			name:     "valid_bearish",
			ohlc:     OHLCData{Open: 110, High: 115, Low: 100, Close: 105},
			expected: true,
		},
		{
			name:     "valid_doji",
			ohlc:     OHLCData{Open: 100, High: 105, Low: 95, Close: 100},
			expected: true,
		},
		{
			name:     "invalid_high_too_low",
			ohlc:     OHLCData{Open: 100, High: 98, Low: 95, Close: 105},
			expected: false,
		},
		{
			name:     "invalid_low_too_high",
			ohlc:     OHLCData{Open: 100, High: 110, Low: 102, Close: 105},
			expected: false,
		},
		{
			name:     "invalid_null_values",
			ohlc:     OHLCData{Open: GetNullValue(), High: 110, Low: 95, Close: 105},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, validateOHLCData(tt.ohlc))
		})
	}
}

func TestCandlestickSeriesListMethods(t *testing.T) {
	t.Parallel()

	data := makeBasicCandlestickData()
	series := CandlestickSeries{Data: data, Name: "Test Series"}
	seriesList := CandlestickSeriesList{series}

	assert.Equal(t, 1, seriesList.len())
	assert.Equal(t, "Test Series", seriesList.names()[0])
	assert.Len(t, series.Data, len(data))
	assert.Empty(t, string(seriesList.getSeriesSymbol(0)))

	// Test ToGenericSeriesList
	genericList := seriesList.ToGenericSeriesList()
	assert.Len(t, genericList, 1)
	assert.Equal(t, ChartTypeCandlestick, genericList[0].Type)
	// OHLC data is encoded as 4 values per candlestick
	assert.Len(t, genericList[0].Values, len(data)*4)
}

func TestRenderCandlestickChart(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: CandlestickSeriesList{{Data: makeBasicCandlestickData()}}.ToGenericSeriesList(),
		Title:      TitleOption{Text: "Price"},
		XAxis: XAxisOption{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May"},
		},
		YAxis: []YAxisOption{
			{
				PreferNiceIntervals: Ptr(true),
			},
		},
		Legend: LegendOption{SeriesNames: []string{"Price"}},
	}

	painter, err := Render(opt, SVGOutputOptionFunc())
	require.NoError(t, err)
	data, err := painter.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestCandlestickCandleWidthClamped(t *testing.T) {
	t.Parallel()

	opt := makeMinimalCandlestickChartOption()
	opt.CandleWidth = 1
	base := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 800, Height: 600})
	require.NoError(t, base.CandlestickChart(opt))
	expected, err := base.Bytes()
	require.NoError(t, err)

	opt.CandleWidth = 2
	over := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 800, Height: 600})
	require.NoError(t, over.CandlestickChart(opt))
	actual, err := over.Bytes()
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual))
}

func validateCandlestickChartRender(t *testing.T, svgP, pngP *Painter, opt CandlestickChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.CandlestickChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.CandlestickChart(opt)
	require.NoError(t, err)
	rdata, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rdata)
}

func TestCandlestickYAxisScalingExtremes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		data []OHLCData
		yMin float64
		yMax float64
	}{
		{
			name: "very_large_millions",
			data: []OHLCData{
				{Open: 1_100_000, High: 1_250_000, Low: 1_000_000, Close: 1_200_000},
				{Open: 1_200_000, High: 1_300_000, Low: 1_150_000, Close: 1_280_000},
			},
			yMin: 1_000_000,
			yMax: 1_500_000,
		},
		{
			name: "small_decimals",
			data: []OHLCData{
				{Open: 0.015, High: 0.02, Low: 0.01, Close: 0.018},
				{Open: 0.017, High: 0.025, Low: 0.012, Close: 0.02},
			},
			yMin: 0.0,
			yMax: 0.05,
		},
		{
			name: "negative_values",
			data: []OHLCData{
				{Open: -1_000, High: -900, Low: -1_100, Close: -950},
				{Open: -950, High: -800, Low: -1_050, Close: -820},
			},
			yMin: -2_000,
			yMax: -500,
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i)+"-"+c.name, func(t *testing.T) {
			opt := CandlestickChartOption{
				XAxis: XAxisOption{Labels: []string{"A", "B"}},
				YAxis: []YAxisOption{{Min: Ptr(c.yMin), Max: Ptr(c.yMax)}},
				SeriesList: CandlestickSeriesList{{
					Data: c.data,
				}},
				Padding: NewBoxEqual(10),
			}

			p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 800, Height: 600})
			require.NoError(t, p.CandlestickChart(opt))
			svg, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, svg)
		})
	}
}

func TestCandlestickTrendLineAlignmentSingleSeries(t *testing.T) {
	p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 600, Height: 400})
	data := []OHLCData{{Open: 100, High: 110, Low: 90, Close: 105}, {Open: 105, High: 115, Low: 95, Close: 108}, {Open: 108, High: 118, Low: 100, Close: 112}}
	opt := CandlestickChartOption{
		Theme:   GetDefaultTheme(),
		Padding: NewBoxEqual(0),
		XAxis:   XAxisOption{Labels: []string{"A", "B", "C"}, Show: Ptr(false)},
		YAxis:   make([]YAxisOption, 1),
		SeriesList: CandlestickSeriesList{{
			Data:           data,
			CloseTrendLine: []SeriesTrendLine{{Type: SeriesTrendTypeLinear, DashedLine: Ptr(true)}},
		}},
		ShowWicks: Ptr(false),
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     &opt.SeriesList,
		xAxis:          &opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	require.NoError(t, err)

	_, err = newCandlestickChart(p, opt).renderChart(renderResult)
	require.NoError(t, err)
	svgBytes, err := p.Bytes()
	require.NoError(t, err)

	paths := extractDashedPathXCoords(string(svgBytes))
	require.Len(t, paths, 1)
	expected := computeCenters(renderResult, opt, 0)
	assert.Equal(t, expected, paths[0])
}

func TestCandlestickTrendLineAlignmentMultiSeries(t *testing.T) {
	p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 600, Height: 400})
	data := []OHLCData{{Open: 100, High: 110, Low: 90, Close: 105}, {Open: 105, High: 115, Low: 95, Close: 108}, {Open: 108, High: 118, Low: 100, Close: 112}}
	opt := CandlestickChartOption{
		Theme:   GetDefaultTheme(),
		Padding: NewBoxEqual(0),
		XAxis:   XAxisOption{Labels: []string{"A", "B", "C"}, Show: Ptr(false)},
		YAxis:   make([]YAxisOption, 1),
		SeriesList: CandlestickSeriesList{
			{Data: data, CloseTrendLine: []SeriesTrendLine{{Type: SeriesTrendTypeLinear, DashedLine: Ptr(true)}}},
			{Data: data, CloseTrendLine: []SeriesTrendLine{{Type: SeriesTrendTypeLinear, DashedLine: Ptr(true)}}},
		},
		ShowWicks: Ptr(false),
	}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     &opt.SeriesList,
		xAxis:          &opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	require.NoError(t, err)

	_, err = newCandlestickChart(p, opt).renderChart(renderResult)
	require.NoError(t, err)
	svgBytes, err := p.Bytes()
	require.NoError(t, err)

	paths := extractDashedPathXCoords(string(svgBytes))
	require.Len(t, paths, 2)
	expected0 := computeCenters(renderResult, opt, 0)
	expected1 := computeCenters(renderResult, opt, 1)
	assert.Equal(t, expected0, paths[0])
	assert.Equal(t, expected1, paths[1])
}

// extractDashedPathXCoords provides dashed trend line x coordinates from svg.
func extractDashedPathXCoords(svg string) [][]int {
	rePath := regexp.MustCompile(`<path[^>]*stroke-dasharray[^>]*d="([^"]+)"`)
	paths := rePath.FindAllStringSubmatch(svg, -1)
	coordRe := regexp.MustCompile(`[ML] ([0-9]+) [0-9]+`)
	result := make([][]int, len(paths))
	for i, p := range paths {
		coords := coordRe.FindAllStringSubmatch(p[1], -1)
		xs := make([]int, len(coords))
		for j, c := range coords {
			x, _ := strconv.Atoi(c[1])
			xs[j] = x
		}
		result[i] = xs
	}
	return result
}

// computeCenters computes expected center positions (absolute) for series.
func computeCenters(r *defaultRenderResult, opt CandlestickChartOption, seriesIndex int) []int {
	width := r.seriesPainter.Width()
	seriesCount := opt.SeriesList.len()
	maxDataCount := getSeriesMaxDataCount(opt.SeriesList)
	candleWidthRatio := opt.CandleWidth
	if candleWidthRatio <= 0 {
		candleWidthRatio = 0.8
	}
	candleWidth := int(float64(width) * candleWidthRatio / float64(maxDataCount))
	if candleWidth < 1 {
		candleWidth = 1
	}
	candleWidthPerSeries := candleWidth / seriesCount
	if candleWidthPerSeries < 1 {
		candleWidthPerSeries = 1
	}
	divideValues := r.xaxisRange.autoDivide()
	centers := make([]int, len(opt.SeriesList.getSeries(seriesIndex).(*CandlestickSeries).Data))
	for j := range centers {
		if j >= len(divideValues) {
			continue
		}
		var sectionWidth int
		if j < len(divideValues)-1 {
			sectionWidth = divideValues[j+1] - divideValues[j]
		} else if j > 0 {
			sectionWidth = divideValues[j] - divideValues[j-1]
		} else {
			sectionWidth = width / maxDataCount
		}
		var groupMargin, candleMargin, cWidth int
		if seriesCount == 1 {
			cWidth = candleWidthPerSeries
		} else {
			var candleMarginFloat *float64
			if opt.CandleMargin != nil {
				marginPixels := float64(sectionWidth) * (*opt.CandleMargin)
				candleMarginFloat = &marginPixels
			}
			groupMargin, candleMargin, cWidth = calculateCandleMarginsAndSize(seriesCount, sectionWidth, candleWidthPerSeries, candleMarginFloat)
		}
		var center int
		if seriesCount == 1 {
			center = divideValues[j] + sectionWidth/2
		} else {
			x := divideValues[j] + groupMargin + seriesIndex*(cWidth+candleMargin)
			center = x + cWidth/2
		}
		centers[j] = center + r.seriesPainter.box.Left
	}
	return centers
}

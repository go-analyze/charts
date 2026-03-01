package charts

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeFullScatterChartOption() ScatterChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return ScatterChartOption{
		Title: TitleOption{
			Text: "Scatter",
		},
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		},
		YAxis: make([]YAxisOption, 1),
		Legend: LegendOption{
			SeriesNames: []string{"Email", "Union Ads", "Video Ads", "Direct", "Search Engine"},
		},
		SeriesList: NewSeriesListScatter(values),
	}
}

func makeBasicScatterChartOption() ScatterChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return ScatterChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"A", "B", "C", "D", "E", "F", "G"},
		},
		YAxis: make([]YAxisOption, 1),
		Legend: LegendOption{
			SeriesNames: []string{"1", "2"},
		},
		SeriesList: NewSeriesListScatter(values),
	}
}

func makeMinimalScatterChartOption() ScatterChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return ScatterChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7"},
			Show:   Ptr(false),
		},
		YAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		SeriesList: NewSeriesListScatter(values),
	}
}

func makeMinimalMultiValueScatterChartOption() ScatterChartOption {
	values := [][][]float64{
		{{120, GetNullValue()}, {132}, {101, 20}, {134}, {90, 28}, {230}, {210}},
		{{820, GetNullValue()}, {932}, {901, 600}, {934}, {1290}, {1330}, {1320}},
	}
	return ScatterChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7"},
			Show:   Ptr(false),
		},
		YAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		SeriesList: NewSeriesListScatterMultiValue(values),
	}
}

func generateRandomScatterData(seriesCount int, dataPointCount int, maxVariationPercentage float64) [][][]float64 {
	data := make([][][]float64, seriesCount)
	for i := 0; i < seriesCount; i++ {
		data[i] = make([][]float64, dataPointCount)
	}
	r := rand.New(rand.NewSource(1))

	for i := 0; i < seriesCount; i++ {
		for j := 0; j < dataPointCount; j++ {
			if j == 0 {
				// Set the initial value for the line
				data[i][j] = []float64{r.Float64() * 100}
			} else {
				// Calculate the allowed variation range
				variationRange := data[i][j-1][0] * maxVariationPercentage / 100
				min := data[i][j-1][0] - variationRange
				max := data[i][j-1][0] + variationRange

				// Generate a random value within the allowed range
				values := []float64{min + r.Float64()*(max-min)}
				if j%2 == 0 {
					values = append(values, min+r.Float64()*(max-min))
				}
				if j%10 == 0 {
					values = append(values, min+r.Float64()*(max-min))
				}
				data[i][j] = values
			}
		}
	}

	return data
}

func makeDenseScatterChartOption() ScatterChartOption {
	const dataPointCount = 100
	values := generateRandomScatterData(3, dataPointCount, 10)

	xAxisLabels := make([]string, dataPointCount)
	for i := 0; i < dataPointCount; i++ {
		xAxisLabels[i] = strconv.Itoa(i)
	}

	return ScatterChartOption{
		SeriesList: NewSeriesListScatterMultiValue(values, ScatterSeriesOption{
			TrendLine: NewTrendLine(SeriesMarkTypeAverage),
			Label: SeriesLabel{
				ValueFormatter: func(f float64) string {
					return FormatValueHumanizeShort(f, 0, false)
				},
			},
		}),
		Padding: NewBoxEqual(20),
		Theme:   GetTheme(ThemeLight),
		YAxis: []YAxisOption{
			{
				Min:            Ptr(0.0), // force min to be zero
				Max:            Ptr(200.0),
				Unit:           10,
				LabelSkipCount: 1,
			},
		},
		XAxis: XAxisOption{
			Labels:        xAxisLabels,
			BoundaryGap:   Ptr(false),
			LabelCount:    10,
			LabelRotation: DegreesToRadians(45),
		},
	}
}

func TestNewScatterChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewScatterChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeScatter, opt.SeriesList[0].getType())
	assert.Len(t, opt.YAxis, 1)
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.ScatterChart(opt))
}

func TestScatterChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ignore      string // specified if the test is ignored
		themed      bool
		makeOptions func() ScatterChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeFullScatterChartOption,
			pngCRC:      0xe4ea791c,
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalScatterChartOption()
				opt.XAxis.Show = Ptr(true)
				opt.XAxis.BoundaryGap = Ptr(true)
				return opt
			},
			pngCRC: 0xd9149cb1,
		},
		{
			name: "dual_yaxis",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x295ebe5b,
		},
		{
			name: "no_nice_interval",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(false)
				opt.YAxis[1].PreferNiceIntervals = Ptr(false)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x295ebe5b,
		},
		{
			name: "left_nice_interval",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(true)
				opt.YAxis[1].PreferNiceIntervals = Ptr(false)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x2f46b989,
		},
		{
			name: "right_nice_interval",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(false)
				opt.YAxis[1].PreferNiceIntervals = Ptr(true)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x5d2cb52f,
		},
		{
			name: "data_gap",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalScatterChartOption()
				opt.SeriesList[0].Values[4] = []float64{GetNullValue()}
				opt.SeriesList[1].Values[2] = []float64{GetNullValue()}
				return opt
			},
			pngCRC: 0x2ec0cf1c,
		},
		{
			name: "mark_line",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalMultiValueScatterChartOption()
				opt.Padding = NewBoxEqual(40)
				opt.SymbolSize = 4.5
				for i := range opt.SeriesList {
					markLine := NewMarkLine("min", "max", "average")
					markLine.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 0, false)
					}
					opt.SeriesList[i].MarkLine = markLine
				}
				return opt
			},
			pngCRC: 0x77aa06d5,
		},
		{
			name: "series_label",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalMultiValueScatterChartOption()
				opt.YAxis[0].Show = Ptr(false)
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
					opt.SeriesList[i].Label.FontStyle = FontStyle{
						FontSize:  12.0,
						Font:      GetDefaultFont(),
						FontColor: ColorBlue,
					}
					opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 2, false)
					}
				}
				return opt
			},
			pngCRC: 0x2dd0ba8d,
		},
		{
			name: "symbol_dot",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolDot
				opt.Legend.Symbol = SymbolDot
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x65aa700f,
		},
		{
			name: "symbol_circle",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolCircle
				opt.Legend.Symbol = SymbolCircle
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xdd29c17b,
		},
		{
			name: "symbol_square",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolSquare
				opt.Legend.Symbol = SymbolSquare
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x11891d83,
		},
		{
			name: "symbol_diamond",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolDiamond
				opt.Legend.Symbol = SymbolDiamond
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x5caf6eb3,
		},
		{
			name:   "symbol_mixed",
			ignore: "size", // svg is too big to commit
			makeOptions: func() ScatterChartOption {
				opt := makeFullScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.0
				opt.SeriesList[0].Symbol = SymbolCircle
				opt.SeriesList[1].Symbol = SymbolSquare
				opt.SeriesList[2].Symbol = SymbolDiamond
				opt.SeriesList[3].Symbol = SymbolDot
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
		},
		{
			name:   "dense_trends",
			ignore: "size", // svg is too big to commit
			makeOptions: func() ScatterChartOption {
				opt := makeDenseScatterChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].TrendLine[0].StrokeSmoothingTension = 0.9 // smooth average line
					opt.SeriesList[i].TrendLine[0].Period = 5
					c1 := Color{
						R: uint8(80 + (20 * i)),
						G: uint8(80 + (20 * i)),
						B: uint8(80 + (20 * i)),
						A: 255,
					}
					c2 := c1
					if i%2 == 0 {
						c2.R = 200
					} else {
						c2.B = 200
					}
					trendLine1 := SeriesTrendLine{
						Type:      SeriesTrendTypeCubic,
						LineColor: c1,
					}
					trendLine2 := SeriesTrendLine{
						Type:      SeriesTrendTypeLinear,
						LineColor: c2,
					}
					opt.SeriesList[i].TrendLine = append(opt.SeriesList[i].TrendLine, trendLine1, trendLine2)
				}
				// disable extras
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
		},
		{
			name: "trend_line_dashed",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetDefaultTheme()
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				// Add dashed trend lines to each series
				for i := range opt.SeriesList {
					opt.SeriesList[i].TrendLine = []SeriesTrendLine{
						{
							StrokeSmoothingTension: 0.8,
							Type:                   SeriesTrendTypeSMA,
							DashedLine:             Ptr(true), // Explicitly set to dashed
							LineColor:              opt.Theme.GetSeriesTrendColor(i).WithAdjustHSL(0, .2, -.2),
						},
						{
							Type:       SeriesTrendTypeCubic,
							DashedLine: Ptr(true), // Explicitly set to dashed
							LineColor:  opt.Theme.GetSeriesTrendColor(i).WithAdjustHSL(0, .4, -.4),
						},
					}
				}
				return opt
			},
			pngCRC: 0x8faf9659,
		},
		{
			name: "with_conditional_labels",
			makeOptions: func() ScatterChartOption {
				return ScatterChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"A", "B", "C", "D", "E"},
					},
					YAxis: []YAxisOption{{}},
					SeriesList: NewSeriesListScatter([][]float64{
						{50, 150, 100, 200, 175},
						{75, 125, 90, 160, 140},
					}, ScatterSeriesOption{
						Names: []string{"Dataset1", "Dataset2"},
						Label: SeriesLabel{
							Show: Ptr(true),
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								// Show labels only for values above 120
								if val > 120 {
									switch {
									case val >= 180: // High values - gold styling
										return "⭐ " + strconv.FormatFloat(val, 'f', 0, 64), &LabelStyle{
											FontStyle:       FontStyle{FontColor: ColorBlack, FontSize: 14},
											BackgroundColor: ColorFromHex("#FFD700"), // Gold
											CornerRadius:    6,
										}
									case val >= 150: // Medium-high values - silver styling
										return "📊 " + strconv.FormatFloat(val, 'f', 0, 64), &LabelStyle{
											FontStyle:       FontStyle{FontColor: ColorBlack, FontSize: 12},
											BackgroundColor: ColorFromHex("#C0C0C0"), // Silver
											CornerRadius:    4,
										}
									default: // Values above 120 but below 150 - simple styling
										return strconv.FormatFloat(val, 'f', 0, 64), &LabelStyle{
											FontStyle: FontStyle{FontColor: ColorBlue, FontSize: 10},
										}
									}
								}
								// Hide labels for values <= 120
								return "", nil
							},
						},
					}),
					Title: TitleOption{
						Show: Ptr(false),
					},
					Legend: LegendOption{
						Show: Ptr(false),
					},
				}
			},
			pngCRC: 0xce19f6cc,
		},
		{
			name: "bollinger",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.SeriesList[0].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeBollingerLower, Period: 3},
				}
				opt.SeriesList[1].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeBollingerUpper, Period: 3},
				}
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x2731f499,
		},
		{
			name: "rsi",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.SeriesList[0].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeRSI, Period: 3},
				}
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x6a3d4fd4,
		},
	}

	for i, tt := range tests {
		if tt.ignore != "" {
			continue
		}
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
		if !tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateScatterChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
		} else {
			theme := GetTheme(ThemeVividDark)
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(theme))
				rp := NewPainter(rasterOptions, PainterThemeOption(theme))

				validateScatterChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_opt", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = theme

				validateScatterChartRender(t, p, rp, opt, tt.pngCRC)
			})
		}
	}
}

func validateScatterChartRender(t *testing.T, svgP, pngP *Painter, opt ScatterChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.ScatterChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.ScatterChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestScatterChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() ScatterChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() ScatterChartOption {
				return NewScatterChartOptionWithData([][]float64{})
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "invalid_yaxis_index",
			makeOptions: func() ScatterChartOption {
				opt := NewScatterChartOptionWithData([][]float64{{1, 2, 3}})
				opt.SeriesList[0].YAxisIndex = 2
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
		{
			name: "negative_yaxis_index",
			makeOptions: func() ScatterChartOption {
				opt := NewScatterChartOptionWithData([][]float64{{1, 2, 3}})
				opt.SeriesList[0].YAxisIndex = -1
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.ScatterChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

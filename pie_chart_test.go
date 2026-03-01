package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicPieChartOption() PieChartOption {
	values := []float64{
		1048, 735, 580, 484, 300,
	}
	return PieChartOption{
		SeriesList: NewSeriesListPie(values),
		Title: TitleOption{
			Text:    "Title",
			Subtext: "Sub",
			Offset:  OffsetCenter,
		},
		Padding: NewBoxEqual(20),
		Legend: LegendOption{
			Vertical:    Ptr(true),
			SeriesNames: []string{"Series-A", "Series-B", "Series-C", "Series-D", "Series-E"},
			Offset:      OffsetLeft,
			Symbol:      SymbolDot,
		},
	}
}

func TestNewPieChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewPieChartOptionWithData([]float64{12, 24, 48})

	assert.Len(t, opt.SeriesList, 3)
	assert.Equal(t, ChartTypePie, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.PieChart(opt))
}

func TestPieChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		width       int
		height      int
		themed      bool
		makeOptions func() PieChartOption
		pngCRC      uint32
	}{
		{
			name: "defaults",
			makeOptions: func() PieChartOption {
				opt := makeBasicPieChartOption()
				opt.Title.Show = Ptr(false)
				opt.Legend.Offset = OffsetStr{}
				opt.Legend.Symbol = ""
				opt.Legend.Vertical = nil
				return opt
			},
			pngCRC: 0x52188ed5,
		},
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicPieChartOption,
			pngCRC:      0x2380709a,
		},
		{
			name:   "lots_labels-sortedDescending",
			width:  1000,
			height: 800,
			makeOptions: func() PieChartOption {
				values := []float64{
					84358845, 68070697, 58850717, 48059777, 36753736, 19051562, 17947406, 11754004,
					10827529, 10521556, 10467366, 10394055, 9597085, 9104772, 6447710, 5932654,
					5563970, 5428792, 5194336, 3850894, 2857279, 2116792, 1883008, 1373101,
					920701, 660809, 542051,
				}

				return PieChartOption{
					SeriesList: NewSeriesListPie(values, PieSeriesOption{
						Label: SeriesLabel{
							Show:           Ptr(true),
							FormatTemplate: "{b} ({c} ≅ {d})",
							ValueFormatter: func(f float64) string {
								return humanize.FtoaWithDigits(f, 2)
							},
						},
					}),
					Radius:  "200",
					Padding: NewBoxEqual(20),
					Legend: LegendOption{
						SeriesNames: []string{
							"Germany",
							"France",
							"Italy",
							"Spain",
							"Poland",
							"Romania",
							"Netherlands",
							"Belgium",
							"Czech Republic",
							"Sweden",
							"Portugal",
							"Greece",
							"Hungary",
							"Austria",
							"Bulgaria",
							"Denmark",
							"Finland",
							"Slovakia",
							"Ireland",
							"Croatia",
							"Lithuania",
							"Slovenia",
							"Latvia",
							"Estonia",
							"Cyprus",
							"Luxembourg",
							"Malta",
						},
						Show: Ptr(false),
					},
				}
			},
			pngCRC: 0x6c21894d,
		},
		{
			name:   "lots_labels-unsorted",
			width:  1000,
			height: 800,
			makeOptions: func() PieChartOption {
				values := []float64{
					9104772, 11754004, 6447710, 3850894, 920701, 10827529, 5932654, 1373101,
					5563970, 68070697, 84358845, 10394055, 9597085, 5194336, 58850717, 1883008,
					2857279, 660809, 542051, 17947406, 36753736, 10467366, 19051562, 5428792,
					2116792, 48059777, 10521556,
				}

				return PieChartOption{
					SeriesList: NewSeriesListPie(values, PieSeriesOption{
						Label: SeriesLabel{
							Show:           Ptr(true),
							FormatTemplate: "{b} ({c} ≅ {d})",
							ValueFormatter: func(f float64) string {
								return humanize.FtoaWithDigits(f, 2)
							},
						},
					}),
					Radius:  "200",
					Padding: NewBoxEqual(20),
					Legend: LegendOption{
						SeriesNames: []string{
							"Austria",
							"Belgium",
							"Bulgaria",
							"Croatia",
							"Cyprus",
							"Czech Republic",
							"Denmark",
							"Estonia",
							"Finland",
							"France",
							"Germany",
							"Greece",
							"Hungary",
							"Ireland",
							"Italy",
							"Latvia",
							"Lithuania",
							"Luxembourg",
							"Malta",
							"Netherlands",
							"Poland",
							"Portugal",
							"Romania",
							"Slovakia",
							"Slovenia",
							"Spain",
							"Sweden",
						},
						Show: Ptr(false),
					},
				}
			},
			pngCRC: 0x28ad9867,
		},
		{
			name:   "100labels-sorted",
			width:  1000,
			height: 900,
			makeOptions: func() PieChartOption {
				var values []float64
				var labels []string
				for i := 1; i <= 100; i++ {
					values = append(values, float64(1))
					labels = append(labels, "Label "+strconv.Itoa(i))
				}

				return PieChartOption{
					SeriesList: NewSeriesListPie(values),
					Radius:     "200",
					Padding:    NewBoxEqual(20),
					Legend: LegendOption{
						SeriesNames: labels,
						Show:        Ptr(false),
					},
				}
			},
			pngCRC: 0x65b9fb22,
		},
		{
			name:   "fix_label_pos",
			width:  1150,
			height: 550,
			makeOptions: func() PieChartOption {
				values := []float64{
					397594, 185596, 149086, 144258, 120194, 117514, 99412, 91135,
					87282, 76790, 72586, 58818, 58270, 56306, 55486, 54792,
					53746, 51460, 41242, 39476, 37414, 36644, 33784, 32788,
					32566, 29608, 29558, 29384, 28166, 26998, 26948, 26054,
					25804, 25730, 24438, 23782, 22896, 21404, 428978,
				}
				return PieChartOption{
					SeriesList: NewSeriesListPie(values, PieSeriesOption{
						Label: SeriesLabel{
							Show:           Ptr(true),
							FormatTemplate: "{b} ({c} ≅ {d})",
							ValueFormatter: func(f float64) string {
								return humanize.FtoaWithDigits(f, 2)
							},
						},
					}),
					Radius: "150",
					Title: TitleOption{
						Text:   "Fix label K (72586)",
						Offset: OffsetRight,
					},
					Padding: NewBoxEqual(20),
					Legend: LegendOption{
						SeriesNames: []string{
							"A", "B", "C", "D", "E", "F", "G", "H",
							"I", "J", "K", "L", "M", "N", "O", "P",
							"Q", "R", "S", "T", "U", "V", "W", "X",
							"Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF",
							"AG", "AH", "AI", "AJ", "AK", "AL", "AM",
						},
						Show: Ptr(false),
					},
				}
			},
			pngCRC: 0xcf524a67,
		},
		{
			name: "pie_chart_with_value_formatter",
			makeOptions: func() PieChartOption {
				return PieChartOption{
					SeriesList: NewSeriesListPie([]float64{
						1048, 735, 580,
					}, PieSeriesOption{
						Names: []string{"A", "B", "C"},
						Label: SeriesLabel{
							Show: Ptr(true),
							ValueFormatter: func(f float64) string {
								return "ValueFormatter: " + humanize.FtoaWithDigits(f, 0)
							},
						},
					}),
				}
			},
			pngCRC: 0x3030453f,
		},
		{
			name: "pie_chart_with_label_formatter_precedence",
			makeOptions: func() PieChartOption {
				return PieChartOption{
					SeriesList: NewSeriesListPie([]float64{
						1048, 735, 580,
					}, PieSeriesOption{
						Names: []string{"A", "B", "C"},
						Label: SeriesLabel{
							Show: Ptr(true),
							ValueFormatter: func(f float64) string {
								return "ValueFormatter: " + humanize.FtoaWithDigits(f, 0)
							},
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								return "LabelFormatter: " + name + "=" + humanize.FtoaWithDigits(val, 0), nil
							},
						},
					}),
				}
			},
			pngCRC: 0x3603a38e,
		},
		{
			name: "custom_fonts",
			makeOptions: func() PieChartOption {
				opt := makeBasicPieChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.SeriesList[0].Label.FontStyle = customFont
				opt.Legend.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0xc0ee21e3,
		},
		{
			name: "legend_bottom_right",
			makeOptions: func() PieChartOption {
				opt := makeBasicPieChartOption()
				opt.Legend.Offset = OffsetStr{
					Top:  PositionBottom,
					Left: PositionRight,
				}
				return opt
			},
			pngCRC: 0x33f71ce5,
		},
		{
			name: "variable_series_radius",
			makeOptions: func() PieChartOption {
				opt := makeBasicPieChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].Radius = strconv.Itoa((i+1)*10) + "%"
				}
				// disable extras
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xaf1b9e0b,
		},
		{
			name: "segment_gap",
			makeOptions: func() PieChartOption {
				opt := makeBasicPieChartOption()
				opt.SegmentGap = 4.0
				// disable extras
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(false)
				}
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xde32ca61,
		},
		{
			name: "mixed_label_style",
			makeOptions: func() PieChartOption {
				return PieChartOption{
					SeriesList: NewSeriesListPie([]float64{
						1048, 735, 580, 484, 300,
					}, PieSeriesOption{
						Names: []string{"Visible", "Hidden", "Styled", "Hidden", "Custom"},
						Label: SeriesLabel{
							Show: Ptr(true),
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								switch index {
								case 0: // first - simple visible
									return name + ": " + strconv.FormatFloat(val, 'f', 0, 64), nil
								case 1, 3: // second and fourth - hidden
									return "", nil
								case 2: // third - styled with background
									return "📊 " + name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorWhite, FontSize: 13},
										BackgroundColor: ColorPurple,
										CornerRadius:    4,
									}
								default: // last - custom color only
									return name + " ★", &LabelStyle{
										FontStyle: FontStyle{FontColor: ColorGreen, FontSize: 14},
									}
								}
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
			pngCRC: 0x95d1ed60,
		},
		{
			name: "border_styling",
			makeOptions: func() PieChartOption {
				values := []float64{100, 200, 150}
				opt := NewPieChartOptionWithData(values)
				opt.SeriesList[0].Label = SeriesLabel{
					Show: Ptr(true),
					LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
						return "background", &LabelStyle{
							BackgroundColor: ColorBlue,
							BorderColor:     ColorRed,
							BorderWidth:     2,
							CornerRadius:    5,
						}
					},
				}
				opt.SeriesList[1].Label = SeriesLabel{
					Show: Ptr(true),
					LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
						return "transparent", &LabelStyle{
							BorderColor: ColorGreen,
							BorderWidth: 1.5,
						}
					},
				}
				return opt
			},
			pngCRC: 0x4c6d0d53,
		},
		{
			name: "label_title_collision",
			makeOptions: func() PieChartOption {
				values := []float64{640, 360}
				labels := []string{"Used", "Free"}
				return PieChartOption{
					SeriesList: NewSeriesListPie(values),
					Legend: LegendOption{
						SeriesNames: labels,
						Offset:      OffsetStr{Top: "20"},
					},
					Title: TitleOption{
						Text: "Test title very long text - my pie chart",
					},
				}
			},
			pngCRC: 0xbafdc837,
		},
		{
			name: "title_legend_overlap",
			makeOptions: func() PieChartOption {
				opt := NewPieChartOptionWithData([]float64{1048, 735, 580, 484, 300})
				opt.Title = TitleOption{
					Text:    "Pie Chart",
					Subtext: "Title Subtext",
					Offset:  OffsetCenter,
				}
				opt.Legend.SeriesNames = []string{
					"Search Engine", "Direct", "Email", "Union Ads", "Video Ads",
				}
				return opt
			},
			pngCRC: 0xd2c9c46b,
		},
		{
			name: "vertical_legend_title_collision",
			makeOptions: func() PieChartOption {
				opt := NewPieChartOptionWithData([]float64{1048, 735, 580})
				opt.Title = TitleOption{
					Text:   "Chart Title",
					Offset: OffsetLeft,
				}
				opt.Legend.SeriesNames = []string{"Search", "Direct", "Email"}
				opt.Legend.Vertical = Ptr(true)
				return opt
			},
			pngCRC: 0xe775743f,
		},
		{
			name: "explicit_offset_no_reposition",
			makeOptions: func() PieChartOption {
				opt := NewPieChartOptionWithData([]float64{1048, 735, 580, 484, 300})
				opt.Title = TitleOption{
					Text:   "Pie Chart",
					Offset: OffsetCenter,
				}
				opt.Legend.SeriesNames = []string{"A", "B", "C", "D", "E"}
				opt.Legend.Offset = OffsetStr{Top: "10"}
				return opt
			},
			pngCRC: 0x2d186954,
		},
	}

	for i, tt := range tests {
		if tt.width == 0 {
			tt.width = 600
		}
		if tt.height == 0 {
			tt.height = 400
		}
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        tt.width,
			Height:       tt.height,
		}
		rasterOptions := PainterOptions{
			OutputFormat: ChartOutputPNG,
			Width:        tt.width,
			Height:       tt.height,
		}
		if tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				rp := NewPainter(rasterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validatePieChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validatePieChartRender(t, p, rp, opt, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validatePieChartRender(t, p.Child(PainterPaddingOption(NewBoxEqual(20))),
					rp.Child(PainterPaddingOption(NewBoxEqual(20))), tt.makeOptions(), tt.pngCRC)
			})
		}
	}
}

func validatePieChartRender(t *testing.T, svgP, pngP *Painter, opt PieChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.PieChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.PieChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestPieChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() PieChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() PieChartOption {
				return NewPieChartOptionWithData([]float64{})
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "zero_sum",
			makeOptions: func() PieChartOption {
				return NewPieChartOptionWithData([]float64{0.0, 0.0})
			},
			errorMsgContains: "greater than 0",
		},
		{
			name: "negative_values",
			makeOptions: func() PieChartOption {
				return NewPieChartOptionWithData([]float64{10.0, -1.0})
			},
			errorMsgContains: "unsupported negative value",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.PieChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

func TestCircleChartPosition(t *testing.T) {
	t.Parallel()

	tests := []struct {
		width, height int
		cx, cy        int
		diameter      float64
	}{
		{200, 100, 100, 50, 100},
		{120, 180, 60, 90, 120},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := NewPainter(PainterOptions{Width: tt.width, Height: tt.height})
			cx, cy, d := circleChartPosition(p)
			assert.Equal(t, tt.cx, cx)
			assert.Equal(t, tt.cy, cy)
			assert.InDelta(t, tt.diameter, d, 0.0)
		})
	}
}

func TestSectorOuterLabelLines(t *testing.T) {
	t.Parallel()

	sec := sector{midAngle: 0, quadrant: 1}
	lsX, lsY, lbX, lbY, leX, leY := sec.calculateOuterLabelLines(50, 50, 20, 30, 15)
	assert.Equal(t, 70, lsX)
	assert.Equal(t, 50, lsY)
	assert.Equal(t, 80, lbX)
	assert.Equal(t, 50, lbY)
	assert.Equal(t, 95, leX)
	assert.Equal(t, 50, leY)
}

func TestSectorAdjustedOuterLabelPosition(t *testing.T) {
	t.Parallel()

	secTop := sector{midAngle: 0, quadrant: 1}
	textBox := NewBox(0, 0, 20, 10)
	lsX, lsY, lbX, adjY, leX, leY, textX, textY := secTop.calculateAdjustedOuterLabelPosition(50, 50, 20, 30, 15, 60, 10, textBox)
	assert.Equal(t, 70, lsX)
	assert.Equal(t, 50, lsY)
	assert.Equal(t, 80, lbX)
	assert.Equal(t, 44, adjY)
	assert.Equal(t, 95, leX)
	assert.Equal(t, 44, leY)
	assert.Equal(t, 98, textX)
	assert.Equal(t, 48, textY)

	secBottom := sector{midAngle: math.Pi, quadrant: 3}
	lsX, lsY, lbX, adjY, leX, leY, textX, textY = secBottom.calculateAdjustedOuterLabelPosition(50, 50, 20, 30, 15, 40, 10, textBox)
	assert.Equal(t, 30, lsX)
	assert.Equal(t, 50, lsY)
	assert.Equal(t, 20, lbX)
	assert.Equal(t, 56, adjY)
	assert.Equal(t, 5, leX)
	assert.Equal(t, 56, leY)
	assert.Equal(t, -18, textX)
	assert.Equal(t, 60, textY)
}

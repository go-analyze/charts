package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicFunnelChartOption() FunnelChartOption {
	return FunnelChartOption{
		SeriesList: NewSeriesListFunnel([]float64{
			100, 80, 60, 40, 20,
		}),
		Legend: LegendOption{
			SeriesNames: []string{"Show", "Click", "Visit", "Inquiry", "Order"},
		},
		Title: TitleOption{
			Text: "Funnel",
		},
	}
}

func TestNewFunnelChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewFunnelChartOptionWithData([]float64{12, 24, 48})

	assert.Len(t, opt.SeriesList, 3)
	assert.Equal(t, ChartTypeFunnel, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.FunnelChart(opt))
}

func TestFunnelChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() FunnelChartOption
		svg         string
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicFunnelChartOption,
			svg:         "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><text x=\"0\" y=\"16\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path d=\"M 86 3\nL 116 3\nL 116 16\nL 86 16\nL 86 3\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"118\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path d=\"M 176 3\nL 206 3\nL 206 16\nL 176 16\nL 176 3\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"208\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path d=\"M 262 3\nL 292 3\nL 292 16\nL 262 16\nL 262 3\" style=\"stroke:none;fill:rgb(100,180,210)\"/><text x=\"294\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path d=\"M 345 3\nL 375 3\nL 375 16\nL 345 16\nL 345 3\" style=\"stroke:none;fill:rgb(64,160,110)\"/><text x=\"377\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path d=\"M 444 3\nL 474 3\nL 474 16\nL 444 16\nL 444 3\" style=\"stroke:none;fill:rgb(154,96,180)\"/><text x=\"476\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><path d=\"M 0 36\nL 600 36\nL 540 107\nL 60 107\nL 0 36\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"264\" y=\"71\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Show(100%)</text><path d=\"M 60 109\nL 540 109\nL 480 180\nL 120 180\nL 60 109\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"269\" y=\"144\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Click(80%)</text><path d=\"M 120 182\nL 480 182\nL 420 253\nL 180 253\nL 120 182\" style=\"stroke:none;fill:rgb(100,180,210)\"/><text x=\"271\" y=\"217\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Visit(60%)</text><path d=\"M 180 255\nL 420 255\nL 360 326\nL 240 326\nL 180 255\" style=\"stroke:none;fill:rgb(64,160,110)\"/><text x=\"264\" y=\"290\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Inquiry(40%)</text><path d=\"M 240 328\nL 360 328\nL 300 399\nL 300 399\nL 240 328\" style=\"stroke:none;fill:rgb(154,96,180)\"/><text x=\"268\" y=\"363\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Order(20%)</text></svg>",
			pngCRC:      0x3c56896,
		},
		{
			name: "custom_legend",
			makeOptions: func() FunnelChartOption {
				opt := makeBasicFunnelChartOption()
				opt.Legend.Symbol = SymbolDot
				opt.Legend.FontStyle = NewFontStyleWithSize(4.0)
				opt.Legend.Vertical = Ptr(true)
				opt.Legend.Offset = OffsetStr{
					Left: PositionRight,
					Top:  PositionBottom,
				}
				opt.Title.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 552 304\nL 582 304\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"567\" cy=\"304\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"584\" y=\"310\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Show</text><path d=\"M 552 324\nL 582 324\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"567\" cy=\"324\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"584\" y=\"330\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Click</text><path d=\"M 552 344\nL 582 344\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"567\" cy=\"344\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"584\" y=\"350\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Visit</text><path d=\"M 552 364\nL 582 364\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"567\" cy=\"364\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"584\" y=\"370\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path d=\"M 552 384\nL 582 384\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"567\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"584\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Order</text><path d=\"M 0 0\nL 600 0\nL 540 78\nL 60 78\nL 0 0\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"264\" y=\"39\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Show(100%)</text><path d=\"M 60 80\nL 540 80\nL 480 158\nL 120 158\nL 60 80\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"269\" y=\"119\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Click(80%)</text><path d=\"M 120 160\nL 480 160\nL 420 238\nL 180 238\nL 120 160\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"271\" y=\"199\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Visit(60%)</text><path d=\"M 180 240\nL 420 240\nL 360 318\nL 240 318\nL 180 240\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"264\" y=\"279\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Inquiry(40%)</text><path d=\"M 240 320\nL 360 320\nL 300 398\nL 300 398\nL 240 320\" style=\"stroke:none;fill:rgb(115,192,222)\"/><text x=\"268\" y=\"359\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Order(20%)</text></svg>",
			pngCRC: 0x8f29abb4,
		},
		{
			name: "with_label_formatter",
			makeOptions: func() FunnelChartOption {
				return FunnelChartOption{
					SeriesList: NewSeriesListFunnel([]float64{
						100, 80, 60, 40, 20,
					}, FunnelSeriesOption{
						Names: []string{"Show", "Click", "Visit", "Inquiry", "Order"},
						Label: SeriesLabel{
							Show: Ptr(true),
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								if index == 1 || index == 3 { // highlight 2nd and 4th items
									return "⭐ " + name + ": " + strconv.FormatFloat(val, 'f', 0, 64), nil
								}
								return "", nil // hide other labels
							},
						},
					}),
					Legend: LegendOption{
						Show: Ptr(false),
					},
				}
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 0 0\nL 600 0\nL 540 78\nL 60 78\nL 0 0\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 60 80\nL 540 80\nL 480 158\nL 120 158\nL 60 80\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"269\" y=\"119\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">⭐ Click: 80</text><path d=\"M 120 160\nL 480 160\nL 420 238\nL 180 238\nL 120 160\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 180 240\nL 420 240\nL 360 318\nL 240 318\nL 180 240\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"263\" y=\"279\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">⭐ Inquiry: 40</text><path d=\"M 240 320\nL 360 320\nL 300 398\nL 300 398\nL 240 320\" style=\"stroke:none;fill:rgb(115,192,222)\"/></svg>",
			pngCRC: 0x6d3a879,
		},
		{
			name: "with_styled_labels",
			makeOptions: func() FunnelChartOption {
				return FunnelChartOption{
					SeriesList: NewSeriesListFunnel([]float64{
						100, 80, 60, 40, 20,
					}, FunnelSeriesOption{
						Names: []string{"Show", "Click", "Visit", "Inquiry", "Order"},
						Label: SeriesLabel{
							Show: Ptr(true),
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								switch index {
								case 0: // first item - red background with rounded corners
									return name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorWhite, FontSize: 14},
										BackgroundColor: ColorRed,
										CornerRadius:    5,
									}
								case 1: // second item - blue background, larger font
									return name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorWhite, FontSize: 16},
										BackgroundColor: ColorBlue,
										CornerRadius:    3,
										BorderColor:     ColorPurple,
										BorderWidth:     2,
									}
								case 2: // third item - green background, square corners
									return "🟢 " + name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorBlack, FontSize: 12},
										BackgroundColor: ColorLime,
										BorderColor:     ColorRed,
										BorderWidth:     2,
									}
								case 3: // fourth item - no background, custom color
									return "⭐ " + name, &LabelStyle{
										FontStyle: FontStyle{FontColor: ColorOrange, FontSize: 15},
									}
								default: // last item - no label
									return "", nil
								}
							},
						},
					}),
					Legend: LegendOption{
						Show: Ptr(false),
					},
				}
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 0 0\nL 600 0\nL 540 78\nL 60 78\nL 0 0\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 279 17\nL 322 17\nL 322 17\nA 5 5 90.00 0 1 327 22\nL 327 38\nL 327 38\nA 5 5 90.00 0 1 322 43\nL 279 43\nL 279 43\nA 5 5 90.00 0 1 274 38\nL 274 22\nL 274 22\nA 5 5 90.00 0 1 279 17\nZ\" style=\"stroke:none;fill:red\"/><text x=\"278\" y=\"39\" style=\"stroke:none;fill:white;font-size:17.9px;font-family:'Roboto Medium',sans-serif\">Show</text><path d=\"M 60 80\nL 540 80\nL 480 158\nL 120 158\nL 60 80\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 276 94\nL 324 94\nL 324 94\nA 3 3 90.00 0 1 327 97\nL 327 120\nL 327 120\nA 3 3 90.00 0 1 324 123\nL 276 123\nL 276 123\nA 3 3 90.00 0 1 273 120\nL 273 97\nL 273 97\nA 3 3 90.00 0 1 276 94\nZ\" style=\"stroke-width:2;stroke:purple;fill:blue\"/><text x=\"277\" y=\"119\" style=\"stroke:none;fill:white;font-size:20.4px;font-family:'Roboto Medium',sans-serif\">Click</text><path d=\"M 120 160\nL 480 160\nL 420 238\nL 180 238\nL 120 160\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 273 179\nL 328 179\nL 328 203\nL 273 203\nL 273 179\" style=\"stroke-width:2;stroke:red;fill:lime\"/><text x=\"277\" y=\"199\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">🟢 Visit</text><path d=\"M 180 240\nL 420 240\nL 360 318\nL 240 318\nL 180 240\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"261\" y=\"279\" style=\"stroke:none;fill:rgb(255,165,0);font-size:19.2px;font-family:'Roboto Medium',sans-serif\">⭐ Inquiry</text><path d=\"M 240 320\nL 360 320\nL 300 398\nL 300 398\nL 240 320\" style=\"stroke:none;fill:rgb(115,192,222)\"/></svg>",
			pngCRC: 0x24a84c4e,
		},
		{
			name: "border_without_background",
			makeOptions: func() FunnelChartOption {
				values := []float64{100, 50}
				opt := NewFunnelChartOptionWithData(values)
				opt.SeriesList[0].Label = SeriesLabel{
					Show: Ptr(true),
					LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
						return "test label", &LabelStyle{
							BorderColor: ColorRed,
							BorderWidth: 2.5,
						}
					},
				}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 20 20\nL 580 20\nL 440 199\nL 160 199\nL 20 20\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 270 92\nL 331 92\nL 331 113\nL 270 113\nL 270 92\" style=\"stroke-width:2.5;stroke:red;fill:none\"/><text x=\"274\" y=\"109\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">test label</text><path d=\"M 160 201\nL 440 201\nL 300 380\nL 300 380\nL 160 201\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"284\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(50%)</text></svg>",
			pngCRC: 0x94752d3e,
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

				validateFunnelChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateFunnelChartRender(t, p, rp, opt, tt.svg, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateFunnelChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
		}
	}
}

func validateFunnelChartRender(t *testing.T, svgP, pngP *Painter, opt FunnelChartOption, expectedSVG string, expectedCRC uint32) {
	t.Helper()

	err := svgP.FunnelChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedSVG, data)

	err = pngP.FunnelChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestFunnelChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() FunnelChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() FunnelChartOption {
				return NewFunnelChartOptionWithData([]float64{})
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

			err := p.FunnelChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

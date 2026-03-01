package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicDoughnutChartOption() DoughnutChartOption {
	values := []float64{
		1048, 735, 580, 484, 300,
	}
	return DoughnutChartOption{
		SeriesList: NewSeriesListDoughnut(values),
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

func makeMinimalDoughnutChartOption() DoughnutChartOption {
	opt := makeBasicDoughnutChartOption()
	// disable extras
	for i := range opt.SeriesList {
		opt.SeriesList[i].Label.Show = Ptr(false)
	}
	opt.Title.Show = Ptr(false)
	opt.Legend.Show = Ptr(false)
	return opt
}

func TestNewDoughnutChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewDoughnutChartOptionWithData([]float64{12, 24, 48})

	assert.Len(t, opt.SeriesList, 3)
	assert.Equal(t, ChartTypeDoughnut, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.DoughnutChart(opt))
}

func TestDoughnutChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() DoughnutChartOption
		pngCRC      uint32
	}{
		{
			name: "defaults",
			makeOptions: func() DoughnutChartOption {
				opt := makeBasicDoughnutChartOption()
				opt.Title.Show = Ptr(false)
				opt.Legend.Offset = OffsetStr{}
				opt.Legend.Symbol = ""
				opt.Legend.Vertical = nil
				return opt
			},
			pngCRC: 0xc2978e0c,
		},
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicDoughnutChartOption,
			pngCRC:      0xeed84398,
		},
		{
			name: "custom_fonts",
			makeOptions: func() DoughnutChartOption {
				opt := makeBasicDoughnutChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.SeriesList[0].Label.FontStyle = customFont
				opt.Legend.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0x3baaad19,
		},
		{
			name: "variable_series_radius",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].Radius = strconv.Itoa((i+1)*10) + "%"
				}
				return opt
			},
			pngCRC: 0x35cb31f8,
		},
		{
			name: "center_radius_small",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.RadiusCenter = "20"
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				return opt
			},
			pngCRC: 0x7c33df8b,
		},
		{
			name: "center_radius_large",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.RadiusRing = "42%"
				opt.RadiusCenter = "40%"
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				return opt
			},
			pngCRC: 0x88346256,
		},
		{
			name: "segment_gap",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.SegmentGap = 20.0
				return opt
			},
			pngCRC: 0xc5e85f35,
		},
		{
			name: "center_sum",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.CenterValues = "sum"
				opt.CenterValuesFontStyle.FontSize = 24.0
				opt.CenterValuesFontStyle.FontColor = ColorNavy
				return opt
			},
			pngCRC: 0xff9189f4,
		},
		{
			name: "center_labels",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.CenterValues = "labels"
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				return opt
			},
			pngCRC: 0x5c692a2,
		},
		{
			name: "center_lots_labels",
			makeOptions: func() DoughnutChartOption {
				values := []float64{
					9104772, 11754004, 10827529, 10394055, 9597085,
					17947406, 36753736, 10467366, 19051562, 10521556,
				}

				return DoughnutChartOption{
					SeriesList:   NewSeriesListDoughnut(values),
					CenterValues: "labels",
					Legend: LegendOption{
						SeriesNames: []string{
							"Cyprus", "Denmark", "Estonia", "Finland", "France",
							"Germany", "Greece", "Hungary", "Ireland", "Italy",
						},
						Show: Ptr(false),
					},
				}
			},
			pngCRC: 0xc9f5c160,
		},
		{
			name: "styled_custom_labels",
			makeOptions: func() DoughnutChartOption {
				return DoughnutChartOption{
					SeriesList: NewSeriesListDoughnut([]float64{
						1048, 735, 580, 484, 300,
					}, DoughnutSeriesOption{
						Names: []string{"Analytics", "Marketing", "Sales", "Support", "Development"},
						Label: SeriesLabel{
							Show: Ptr(true),
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								switch index {
								case 0: // analytics - data icon with blue background
									return "📊 " + name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorWhite, FontSize: 12},
										BackgroundColor: ColorBlue,
										CornerRadius:    4,
									}
								case 1: // marketing - no icon with red background
									return name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorWhite, FontSize: 12},
										BackgroundColor: ColorRed,
										CornerRadius:    6,
									}
								case 2: // sales - money icon with green background
									return "💰 " + name, &LabelStyle{
										FontStyle:       FontStyle{FontColor: ColorBlack, FontSize: 11},
										BackgroundColor: ColorLime,
										CornerRadius:    3,
										BorderColor:     ColorRed,
										BorderWidth:     2,
									}
								case 3: // support - help icon with orange color
									return "❓ " + name, &LabelStyle{
										FontStyle: FontStyle{FontColor: ColorOrange, FontSize: 13},
									}
								default: // development, no label
									return "", nil
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
			pngCRC: 0x5d5f7cfb,
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

				validateDoughnutChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateDoughnutChartRender(t, p, rp, opt, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateDoughnutChartRender(t, p.Child(PainterPaddingOption(NewBoxEqual(20))),
					rp.Child(PainterPaddingOption(NewBoxEqual(20))), tt.makeOptions(), tt.pngCRC)
			})
		}
	}
}

func validateDoughnutChartRender(t *testing.T, svgP, pngP *Painter, opt DoughnutChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.DoughnutChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.DoughnutChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestDoughnutChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() DoughnutChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() DoughnutChartOption {
				return NewDoughnutChartOptionWithData([]float64{})
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "zero_sum",
			makeOptions: func() DoughnutChartOption {
				return NewDoughnutChartOptionWithData([]float64{0.0, 0.0})
			},
			errorMsgContains: "greater than 0",
		},
		{
			name: "negative_values",
			makeOptions: func() DoughnutChartOption {
				return NewDoughnutChartOptionWithData([]float64{10.0, -1.0})
			},
			errorMsgContains: "unsupported negative value",
		},
		{
			name: "invalid_radius_center",
			makeOptions: func() DoughnutChartOption {
				opt := NewDoughnutChartOptionWithData([]float64{10.0, 20.0})
				opt.RadiusCenter = "abc"
				return opt
			},
			errorMsgContains: "invalid RadiusCenter",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.DoughnutChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

func TestClampAngleToSector(t *testing.T) {
	t.Parallel()

	cases := []struct {
		angle, start, end float64
		want              float64
	}{
		{0.5, 0, math.Pi / 2, 0.5},
		{-0.1, 0, math.Pi / 2, math.Pi / 2},
		{math.Pi, 0, math.Pi / 2, math.Pi / 2},
		{math.Pi, 3 * math.Pi / 2, math.Pi / 2, 3 * math.Pi / 2},
		{0, 3 * math.Pi / 2, math.Pi / 2, 0},
	}
	for _, tt := range cases {
		got := clampAngleToSector(tt.angle, tt.start, tt.end)
		assert.InDelta(t, tt.want, got, 1e-6)
	}
}

func TestConnectionPoint(t *testing.T) {
	t.Parallel()

	s := sector{startAngle: 0, delta: math.Pi / 2}
	x, y := s.connectionPoint(0, 0, 5, 0, 10)
	assert.Equal(t, 0, x)
	assert.Equal(t, 5, y)

	x, y = s.connectionPoint(0, 0, 5, 5, 5)
	assert.Equal(t, 3, x)
	assert.Equal(t, 3, y)

	x, y = s.connectionPoint(0, 0, 5, -10, 0)
	assert.Equal(t, 0, x)
	assert.Equal(t, 5, y)
}

func TestIsInsideCircle(t *testing.T) {
	t.Parallel()

	c := isInsideCircle(NewBox(-1, -1, 1, 1), 0, 0, 5)
	assert.True(t, c)
	c = isInsideCircle(NewBox(4, 4, 6, 6), 0, 0, 5)
	assert.False(t, c)
}

func TestClampInsideCircle(t *testing.T) {
	t.Parallel()

	lp := &labelPlacement{box: NewBox(6, -1, 8, 1)}
	clampInsideCircle(lp, 0, 0, 5)
	assert.True(t, isInsideCircle(lp.box, 0, 0, 5))
}

func TestMinimalRadialPush(t *testing.T) {
	t.Parallel()

	p := &labelPlacement{box: NewBox(1, 0, 3, 2)}
	q := &labelPlacement{box: NewBox(0, 0, 2, 2)}
	dx, dy := minimalRadialPush(p, q, 0, 0)
	assert.Equal(t, 3, dx)
	assert.Equal(t, 2, dy)
}

func TestProjectBoxRadially(t *testing.T) {
	t.Parallel()

	b := NewBox(0, 0, 2, 2)
	min, max := projectBoxRadially(b, 1, 0)
	assert.InDelta(t, 0.0, min, 1e-6)
	assert.InDelta(t, 2.0, max, 1e-6)

	min, max = projectBoxRadially(b, 0, 1)
	assert.InDelta(t, 0.0, min, 1e-6)
	assert.InDelta(t, 2.0, max, 1e-6)
}

func TestShiftLabelHorizontallyTowardSector(t *testing.T) {
	t.Parallel()

	lp := &labelPlacement{box: NewBox(-2, -1, 0, 1)}
	shiftLabelHorizontallyTowardSector(lp, 0, 0)
	assert.Equal(t, NewBox(0, -1, 2, 1), lp.box)

	lp = &labelPlacement{box: NewBox(2, -1, 4, 1)}
	shiftLabelHorizontallyTowardSector(lp, 0, math.Pi)
	assert.Equal(t, NewBox(-2, -1, 0, 1), lp.box)
}

func TestAnyLabelCollision(t *testing.T) {
	t.Parallel()

	placed := []*labelPlacement{{box: NewBox(1, 1, 3, 3)}}
	c := anyLabelCollision(NewBox(0, 0, 2, 2), placed)
	assert.True(t, c)
	c = anyLabelCollision(NewBox(3, 3, 5, 5), placed)
	assert.False(t, c)
}

func TestComputeLabelBox(t *testing.T) {
	t.Parallel()

	b := computeLabelBox(0, 0, 5, 0, NewBox(0, 0, 2, 2))
	assert.Equal(t, NewBox(3, -2, 5, 0), b)

	b = computeLabelBox(0, 0, 5, math.Pi, NewBox(0, 0, 2, 2))
	assert.Equal(t, NewBox(-5, -2, -3, 0), b)
}

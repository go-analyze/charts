package charts

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allThemes = []string{ThemeLight, ThemeDark, ThemeVividLight, ThemeVividDark, ThemeAnt, ThemeGrafana,
	ThemeNatureLight, ThemeNatureDark, ThemeRetro, ThemeOcean, ThemeSlate, ThemeGray,
	ThemeWinter, ThemeSpring, ThemeSummer, ThemeFall}

func TestInstallGetTheme(t *testing.T) {
	t.Parallel()

	name := "TestInstallGetTheme"
	InstallTheme(name, ThemeOption{IsDarkMode: true})
	getThemeResult := GetTheme(name)
	assert.NotNil(t, getThemeResult)
	assert.NotEqual(t, GetDefaultTheme(), getThemeResult)
	assert.Equal(t, name, fmt.Sprintf("%s", getThemeResult))
}

func TestGetPreferredTheme(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		assert.Equal(t, getPreferredTheme(nil), GetDefaultTheme())
	})
	t.Run("provided", func(t *testing.T) {
		theme := GetTheme(ThemeGrafana)
		assert.Equal(t, getPreferredTheme(theme), theme)
	})
}

func TestDefaultTheme(t *testing.T) {
	t.Parallel()

	assert.Equal(t, GetTheme(ThemeLight), GetDefaultTheme())
	assert.Equal(t, GetTheme("Unknown Theme"), GetDefaultTheme())
}

func TestSetDefaultThemeError(t *testing.T) {
	t.Parallel()

	assert.Error(t, SetDefaultTheme("not a theme"))
}

func renderTestLineChartWithThemeName(t *testing.T, fullChart bool, themeName string) []byte {
	t.Helper()

	return renderTestLineChartWithTheme(t, fullChart, GetTheme(themeName))
}

func renderTestLineChartWithTheme(t *testing.T, fullChart bool, theme ColorPalette) []byte {
	t.Helper()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})
	opt := makeFullLineChartOption()
	if len(opt.YAxis) == 0 {
		opt.YAxis = []YAxisOption{{}}
	}
	opt.YAxis[0].LabelCount = 5
	opt.YAxis[0].LabelSkipCount = 1
	if !fullChart {
		opt.Title.Show = Ptr(false)
		opt.Symbol = SymbolNone
	}
	opt.Theme = theme

	err := p.LineChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	return data
}

func TestThemeLight(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeLight)
	assertTestdataSVG(t, svg)
}

func TestThemeDark(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeDark)
	assertTestdataSVG(t, svg)
}

func TestThemeVividLight(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeVividLight)
	assertTestdataSVG(t, svg)
}

func TestThemeVividDark(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeVividDark)
	assertTestdataSVG(t, svg)
}

func TestThemeAnt(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeAnt)
	assertTestdataSVG(t, svg)
}

func TestThemeGrafana(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeGrafana)
	assertTestdataSVG(t, svg)
}

func TestThemeNatureLight(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeNatureLight)
	assertTestdataSVG(t, svg)
}

func TestThemeNatureDark(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeNatureDark)
	assertTestdataSVG(t, svg)
}

func TestThemeRetro(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeRetro)
	assertTestdataSVG(t, svg)
}

func TestThemeOcean(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeOcean)
	assertTestdataSVG(t, svg)
}

func TestThemeSlate(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeSlate)
	assertTestdataSVG(t, svg)
}

func TestThemeGray(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeGray)
	assertTestdataSVG(t, svg)
}

func TestThemeWinter(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeWinter)
	assertTestdataSVG(t, svg)
}

func TestThemeSpring(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeSpring)
	assertTestdataSVG(t, svg)
}

func TestThemeSummer(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeSummer)
	assertTestdataSVG(t, svg)
}

func TestThemeFall(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeFall)
	assertTestdataSVG(t, svg)
}

func TestLightThemeSeriesRepeat(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithTheme(t, false,
		&themeColorPalette{
			name:               t.Name(),
			isDarkMode:         false,
			xaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			yaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			axisSplitLineColor: Color{R: 224, G: 230, B: 242, A: 255},
			backgroundColor:    ColorWhite,
			legendTextColor:    Color{R: 70, G: 70, B: 70, A: 255},
			seriesColors: []Color{
				{R: 50, G: 50, B: 50, A: 255},
				{R: 200, G: 50, B: 50, A: 255},
			},
		})
	assertTestdataSVG(t, svg)
}

func TestDarkThemeSeriesRepeat(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithTheme(t, false,
		&themeColorPalette{
			name:               t.Name(),
			isDarkMode:         true,
			xaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			yaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			axisSplitLineColor: Color{R: 40, G: 40, B: 60, A: 255},
			backgroundColor:    Color{R: 40, G: 40, B: 60, A: 255},
			legendTextColor:    Color{R: 70, G: 70, B: 70, A: 255},
			seriesColors: []Color{
				{R: 250, G: 250, B: 250, A: 255},
				{R: 200, G: 50, B: 50, A: 255},
			},
		})
	assertTestdataSVG(t, svg)
}

func TestWithAxisColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithXAxisColor(ColorBlack).WithYAxisColor(ColorBlack)

	assert.Equal(t, ColorWhite, whiteCP.xaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.yaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.backgroundColor)
	assert.Equal(t, ColorWhite, whiteCP.seriesColors[0])

	assert.NotEqual(t, ColorBlack, blackCP.GetAxisSplitLineColor())
	assert.Equal(t, ColorBlack, blackCP.GetXAxisStrokeColor())
	assert.Equal(t, ColorBlack, blackCP.GetYAxisStrokeColor())
	assert.Equal(t, ColorWhite, blackCP.GetBackgroundColor())
	assert.Equal(t, ColorWhite, blackCP.GetSeriesColor(0))
}

func TestWithTextColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		titleTextColor:     ColorWhite,
		markTextColor:      ColorWhite,
		labelTextColor:     ColorWhite,
		legendTextColor:    ColorWhite,
		xaxisTextColor:     ColorWhite,
		yaxisTextColor:     ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithTextColor(ColorBlack)

	assert.Equal(t, ColorWhite, whiteCP.axisSplitLineColor)
	assert.Equal(t, ColorWhite, whiteCP.xaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.yaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.backgroundColor)
	assert.Equal(t, ColorWhite, whiteCP.GetTitleTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetLabelTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetLegendTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetMarkTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetXAxisTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetYAxisTextColor())
	assert.Equal(t, ColorWhite, whiteCP.seriesColors[0])

	assert.Equal(t, ColorBlack, blackCP.GetTitleTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetLabelTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetLegendTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetMarkTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetXAxisTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetYAxisTextColor())
	assert.Equal(t, ColorWhite, blackCP.GetAxisSplitLineColor())
	assert.Equal(t, ColorWhite, blackCP.GetXAxisStrokeColor())
	assert.Equal(t, ColorWhite, blackCP.GetYAxisStrokeColor())
	assert.Equal(t, ColorWhite, blackCP.GetBackgroundColor())
	assert.Equal(t, ColorWhite, blackCP.GetSeriesColor(0))

	redCP := blackCP.WithTitleTextColor(ColorRed).WithLabelTextColor(ColorRed).
		WithMarkTextColor(ColorRed).WithLegendTextColor(ColorRed).
		WithXAxisTextColor(ColorRed).WithYAxisTextColor(ColorRed)

	assert.Equal(t, ColorRed, redCP.GetTitleTextColor())
	assert.Equal(t, ColorRed, redCP.GetLabelTextColor())
	assert.Equal(t, ColorRed, redCP.GetLegendTextColor())
	assert.Equal(t, ColorRed, redCP.GetMarkTextColor())
	assert.Equal(t, ColorRed, redCP.GetXAxisTextColor())
	assert.Equal(t, ColorRed, redCP.GetYAxisTextColor())
}

func TestWithSeriesColors(t *testing.T) {
	t.Parallel()

	yellowCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorYellow,
		xaxisStrokeColor:   ColorYellow,
		yaxisStrokeColor:   ColorYellow,
		backgroundColor:    ColorYellow,
		seriesColors:       []Color{ColorYellow},
	}

	t.Run("ignored", func(t *testing.T) {
		invalidCP := yellowCP.WithSeriesColors([]Color{})

		assert.Equal(t, yellowCP.GetSeriesColor(0), invalidCP.GetSeriesColor(0))
	})
	t.Run("updated", func(t *testing.T) {
		blackCP := yellowCP.WithSeriesColors([]Color{ColorBlack})

		assert.Equal(t, ColorYellow, blackCP.GetAxisSplitLineColor())
		assert.Equal(t, ColorYellow, blackCP.GetXAxisStrokeColor())
		assert.Equal(t, ColorYellow, blackCP.GetYAxisStrokeColor())
		assert.Equal(t, ColorYellow, blackCP.GetBackgroundColor())
		assert.Equal(t, ColorBlack, blackCP.GetSeriesColor(0))
	})
}

func TestWithBackgroundColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithBackgroundColor(ColorBlack)

	require.False(t, whiteCP.IsDark())
	assert.True(t, blackCP.IsDark())

	yellowCP := blackCP.WithBackgroundColor(ColorYellow)

	require.False(t, yellowCP.IsDark())
}

func TestWithLegendBorderColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		legendBorderColor:  ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithLegendBorderColor(ColorBlack)
	assert.Equal(t, ColorBlack, blackCP.GetLegendBorderColor())
	assert.Equal(t, ColorWhite, whiteCP.GetLegendBorderColor())
}

func getThemeSeriesColors(name string) []Color {
	if value, ok := palettes.Load(name); ok {
		if cp, ok := value.(*themeColorPalette); ok {
			return cp.seriesColors
		}
	}
	return nil
}

const makeSeriesLoopColorSamples = false

func TestThemeLightSeriesLoopColors(t *testing.T) {
	t.Parallel()
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	for _, theme := range allThemes {
		t.Run(theme+"-loop", func(t *testing.T) {
			testColorShadeVariation(t, func(c Color) Color {
				return adjustSeriesColor(c, 1, false)
			}, getThemeSeriesColors(theme)...)
		})
		t.Run(theme+"-pie", func(t *testing.T) {
			svg := renderTestPieChartWithThemeName(t, theme)
			assertEqualSVG(t, nil, svg)
		})
	}
}

func renderTestPieChartWithThemeName(t *testing.T, themeName string) []byte {
	t.Helper()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        800,
		Height:       600,
	})
	values := []float64{
		64358845, 48070697, 40850717, 40059777, 36753736, 19051562, 17947406, 11754004,
		10827529, 10521556, 10467366, 10394055, 9597085, 9104772, 6447710, 5932654,
		5563970, 5428792, 5194336, 3850894, 2857279, 2116792, 1883008, 1373101,
		920701, 660809, 542051,
	}
	opt := PieChartOption{
		Theme: GetTheme(themeName),
		Title: TitleOption{
			Text: "Example " + themeName,
		},
		SeriesList: NewSeriesListPie(values, PieSeriesOption{
			Label: SeriesLabel{
				Show: Ptr(false),
			},
		}),
		Radius:  "200",
		Padding: NewBoxEqual(20),
		Legend: LegendOption{
			Show: Ptr(false),
		},
	}

	err := p.PieChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	return data
}

const makeTrendColorSamples = false

func TestThemeTrendColors(t *testing.T) {
	t.Parallel()
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	for _, theme := range allThemes {
		t.Run(theme, func(t *testing.T) {
			testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(theme)...)
		})
	}
}

func testColorShadeVariation(t *testing.T, mutateFunc func(Color) Color, colors ...Color) {
	t.Helper()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})

	sampleWidth := p.Width() / len(colors)
	for i, c := range colors {
		startX := i * sampleWidth
		endX := (i + 1) * sampleWidth
		if i == len(colors)-1 {
			endX = p.Width() // ensure edge is painted
		}
		p.FilledRect(startX, 0, endX, p.Height(),
			c, ColorTransparent, 0.0)
		margin := (endX - startX) / 4
		c2 := mutateFunc(c)
		p.FilledRect(startX+margin, 0, endX-margin, p.Height(),
			c2, ColorTransparent, 0.0)
	}

	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, nil, data)
}

func TestCandlestickThemes(t *testing.T) {
	t.Parallel()

	// Create large candlestick series to test color cycling
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

	for i, theme := range allThemes {
		t.Run(strconv.Itoa(i)+"-"+theme, func(t *testing.T) {
			opt := CandlestickChartOption{
				Theme: GetTheme(theme),
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

			painterOptions := PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        1200,
				Height:       800,
			}
			p := NewPainter(painterOptions)

			err := p.CandlestickChart(opt)
			require.NoError(t, err)

			buf, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, buf)
		})
	}
}

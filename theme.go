package charts

import (
	"fmt"
	"sync"

	"github.com/wcharczuk/go-chart/v2/drawing"
)

// ThemeLight is the default theme used, with series colors from echarts.
const ThemeLight = "light"

// ThemeDark is a dark alternative to the default theme 'light, with series colors from echarts'.
const ThemeDark = "dark"

// ThemeVividLight is an alternative light theme that has red, yellow, and other bright colors initially in the series.
// It can be a good option when you want the first few series items to grab the most attention.
const ThemeVividLight = "vivid-light"

// ThemeVividDark is a dark alternative to 'ThemeVividLight', with the same bright initial series colors.
const ThemeVividDark = "vivid-dark"

// ThemeGrafana is a grafana styled theme.
const ThemeGrafana = "grafana"

// ThemeAnt is an ant styled theme.
const ThemeAnt = "ant"

type ColorPalette interface {
	IsDark() bool
	GetAxisStrokeColor() Color
	GetAxisSplitLineColor() Color
	GetSeriesColor(int) Color
	GetBackgroundColor() Color
	GetTextColor() Color
	// WithAxisColor will provide a new ColorPalette that uses the specified color for axis values.
	// This includes the Axis Stroke, Split Line, and Text Color.
	WithAxisColor(Color) ColorPalette
}

type themeColorPalette struct {
	name               string
	isDarkMode         bool
	axisStrokeColor    Color
	axisSplitLineColor Color
	backgroundColor    Color
	textColor          Color
	seriesColors       []Color
}

type ThemeOption struct {
	IsDarkMode         bool
	AxisStrokeColor    Color
	AxisSplitLineColor Color
	BackgroundColor    Color
	TextColor          Color
	SeriesColors       []Color
}

var palettes = sync.Map{}

const defaultTheme = "default"

var defaultLightFontColor = drawing.Color{R: 70, G: 70, B: 70, A: 255}
var defaultDarkFontColor = drawing.Color{R: 238, G: 238, B: 238, A: 255}

func init() {
	darkGray := Color{R: 40, G: 40, B: 40, A: 255}
	echartSeriesColors := []Color{
		parseColor("#5470c6"),
		parseColor("#91cc75"),
		parseColor("#fac858"),
		parseColor("#ee6666"),
		parseColor("#73c0de"),
		parseColor("#3ba272"),
		parseColor("#fc8452"),
		parseColor("#9a60b4"),
		parseColor("#ea7ccc"),
	}
	InstallTheme(
		ThemeLight,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 110, G: 112, B: 121, A: 255},
			AxisSplitLineColor: Color{R: 224, G: 230, B: 242, A: 255},
			BackgroundColor:    drawing.ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			SeriesColors:       echartSeriesColors,
		},
	)
	InstallTheme(
		ThemeDark,
		ThemeOption{
			IsDarkMode:         true,
			AxisStrokeColor:    Color{R: 185, G: 184, B: 206, A: 255},
			AxisSplitLineColor: Color{R: 72, G: 71, B: 83, A: 255},
			BackgroundColor:    darkGray,
			TextColor:          Color{R: 238, G: 238, B: 238, A: 255},
			SeriesColors:       echartSeriesColors,
		},
	)
	vividSeriesColors := []Color{
		{ // red
			R: 255, G: 100, B: 100, A: 255,
		},
		{ // yellow
			R: 255, G: 210, B: 100, A: 255,
		},
		{ // blue
			R: 100, G: 180, B: 210, A: 255,
		},
		{ // green
			R: 64, G: 160, B: 110, A: 255,
		},
		{ // purple
			R: 154, G: 100, B: 180, A: 255,
		},
		{ // light red
			R: 250, G: 128, B: 80, A: 255,
		},
		{ // light green
			R: 90, G: 210, B: 110, A: 255,
		},
		{ // light purple
			R: 220, G: 150, B: 210, A: 255,
		},
		{ // dark blue
			R: 90, G: 118, B: 140, A: 255,
		},
	}
	InstallTheme(
		ThemeVividLight,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 110, G: 112, B: 121, A: 255},
			AxisSplitLineColor: Color{R: 224, G: 230, B: 242, A: 255},
			BackgroundColor:    drawing.ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			SeriesColors:       vividSeriesColors,
		},
	)
	InstallTheme(
		ThemeVividDark,
		ThemeOption{
			IsDarkMode:         true,
			AxisStrokeColor:    Color{R: 185, G: 184, B: 206, A: 255},
			AxisSplitLineColor: Color{R: 72, G: 71, B: 83, A: 255},
			BackgroundColor:    darkGray,
			TextColor:          Color{R: 238, G: 238, B: 238, A: 255},
			SeriesColors:       vividSeriesColors,
		},
	)
	InstallTheme(
		ThemeAnt,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 110, G: 112, B: 121, A: 255},
			AxisSplitLineColor: Color{R: 224, G: 230, B: 242, A: 255},
			BackgroundColor:    drawing.ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			SeriesColors: []Color{
				parseColor("#5b8ff9"),
				parseColor("#5ad8a6"),
				parseColor("#5d7092"),
				parseColor("#f6bd16"),
				parseColor("#6f5ef9"),
				parseColor("#6dc8ec"),
				parseColor("#945fb9"),
				parseColor("#ff9845"),
			},
		},
	)
	InstallTheme(
		ThemeGrafana,
		ThemeOption{
			IsDarkMode:         true,
			AxisStrokeColor:    Color{R: 185, G: 184, B: 206, A: 255},
			AxisSplitLineColor: Color{R: 68, G: 67, B: 67, A: 255},
			BackgroundColor:    Color{R: 31, G: 29, B: 29, A: 255},
			TextColor:          Color{R: 216, G: 217, B: 218, A: 255},
			SeriesColors: []Color{
				parseColor("#7EB26D"),
				parseColor("#EAB839"),
				parseColor("#6ED0E0"),
				parseColor("#EF843C"),
				parseColor("#E24D42"),
				parseColor("#1F78C1"),
				parseColor("#705DA0"),
				parseColor("#508642"),
			},
		},
	)
	if err := SetDefaultTheme(ThemeLight); err != nil {
		panic(fmt.Errorf("could not setup default theme %s", ThemeLight))
	}
}

// SetDefaultTheme sets default theme by name.
func SetDefaultTheme(name string) error {
	if value, ok := palettes.Load(name); ok {
		palettes.Store(defaultTheme, value)
		return nil
	}
	return fmt.Errorf("theme not found: %s", name)
}

func getPreferredTheme(t ...ColorPalette) ColorPalette {
	for _, theme := range t {
		if theme != nil {
			return theme
		}
	}
	return GetDefaultTheme()
}

// GetDefaultTheme returns the default theme.
func GetDefaultTheme() ColorPalette {
	return GetTheme(defaultTheme)
}

// InstallTheme adds a theme to the catalog.
func InstallTheme(name string, opt ThemeOption) {
	palettes.Store(name, &themeColorPalette{
		name:               name,
		isDarkMode:         opt.IsDarkMode,
		axisStrokeColor:    opt.AxisStrokeColor,
		axisSplitLineColor: opt.AxisSplitLineColor,
		backgroundColor:    opt.BackgroundColor,
		textColor:          opt.TextColor,
		seriesColors:       opt.SeriesColors,
	})
}

// GetTheme returns an installed theme by name, or the default if the theme is not installed.
func GetTheme(name string) ColorPalette {
	if value, ok := palettes.Load(name); ok {
		if cp, ok := value.(ColorPalette); ok {
			return cp
		}
	}
	return GetDefaultTheme()
}

func (t *themeColorPalette) String() string {
	return t.name
}

func (t *themeColorPalette) IsDark() bool {
	return t.isDarkMode
}

func (t *themeColorPalette) GetAxisStrokeColor() Color {
	return t.axisStrokeColor
}

func (t *themeColorPalette) GetAxisSplitLineColor() Color {
	return t.axisSplitLineColor
}

func (t *themeColorPalette) GetSeriesColor(index int) Color {
	colors := t.seriesColors
	return colors[index%len(colors)]
}

func (t *themeColorPalette) GetBackgroundColor() Color {
	return t.backgroundColor
}

func (t *themeColorPalette) GetTextColor() Color {
	return t.textColor
}

func (t *themeColorPalette) WithAxisColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-modified"
	copy.axisStrokeColor = c
	copy.axisSplitLineColor = c
	copy.textColor = c
	return &copy
}

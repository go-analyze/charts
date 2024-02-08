package charts

import (
	"github.com/golang/freetype/truetype"
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
	GetFontSize() float64
	GetFont() *truetype.Font
}

type themeColorPalette struct {
	isDarkMode         bool
	axisStrokeColor    Color
	axisSplitLineColor Color
	backgroundColor    Color
	textColor          Color
	seriesColors       []Color
	fontSize           float64
	font               *truetype.Font
}

type ThemeOption struct {
	IsDarkMode         bool
	AxisStrokeColor    Color
	AxisSplitLineColor Color
	BackgroundColor    Color
	TextColor          Color
	SeriesColors       []Color
}

var palettes = map[string]*themeColorPalette{}

const defaultFontSize = 12.0

var defaultTheme ColorPalette

var defaultLightFontColor = drawing.Color{
	R: 70,
	G: 70,
	B: 70,
	A: 255,
}
var defaultDarkFontColor = drawing.Color{
	R: 238,
	G: 238,
	B: 238,
	A: 255,
}

func init() {
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
	AddTheme(
		ThemeLight,
		ThemeOption{
			IsDarkMode: false,
			AxisStrokeColor: Color{
				R: 110,
				G: 112,
				B: 121,
				A: 255,
			},
			AxisSplitLineColor: Color{
				R: 224,
				G: 230,
				B: 242,
				A: 255,
			},
			BackgroundColor: drawing.ColorWhite,
			TextColor: Color{
				R: 70,
				G: 70,
				B: 70,
				A: 255,
			},
			SeriesColors: echartSeriesColors,
		},
	)
	AddTheme(
		ThemeDark,
		ThemeOption{
			IsDarkMode: true,
			AxisStrokeColor: Color{
				R: 185,
				G: 184,
				B: 206,
				A: 255,
			},
			AxisSplitLineColor: Color{
				R: 72,
				G: 71,
				B: 83,
				A: 255,
			},
			BackgroundColor: Color{
				R: 16,
				G: 12,
				B: 42,
				A: 255,
			},
			TextColor: Color{
				R: 238,
				G: 238,
				B: 238,
				A: 255,
			},
			SeriesColors: echartSeriesColors,
		},
	)
	vividSeriesColors := []Color{
		drawing.ColorFromAlphaMixedRGBA(255, 100, 100, 255), // red
		drawing.ColorFromAlphaMixedRGBA(255, 210, 100, 255), // yellow
		drawing.ColorFromAlphaMixedRGBA(100, 180, 210, 255), // blue
		drawing.ColorFromAlphaMixedRGBA(64, 160, 110, 255),  // green
		drawing.ColorFromAlphaMixedRGBA(154, 100, 180, 255), // purple
		drawing.ColorFromAlphaMixedRGBA(220, 150, 210, 255), // light purple
		drawing.ColorFromAlphaMixedRGBA(250, 128, 80, 255),  // light red
		drawing.ColorFromAlphaMixedRGBA(90, 210, 110, 255),  // light green
		drawing.ColorFromAlphaMixedRGBA(90, 110, 140, 255),  // dark blue
	}
	AddTheme(
		ThemeVividLight,
		ThemeOption{
			IsDarkMode: false,
			AxisStrokeColor: Color{
				R: 110,
				G: 112,
				B: 121,
				A: 255,
			},
			AxisSplitLineColor: Color{
				R: 224,
				G: 230,
				B: 242,
				A: 255,
			},
			BackgroundColor: drawing.ColorWhite,
			TextColor: Color{
				R: 70,
				G: 70,
				B: 70,
				A: 255,
			},
			SeriesColors: vividSeriesColors,
		},
	)
	AddTheme(
		ThemeVividDark,
		ThemeOption{
			IsDarkMode: true,
			AxisStrokeColor: Color{
				R: 185,
				G: 184,
				B: 206,
				A: 255,
			},
			AxisSplitLineColor: Color{
				R: 72,
				G: 71,
				B: 83,
				A: 255,
			},
			BackgroundColor: Color{
				R: 40,
				G: 40,
				B: 40,
				A: 255,
			},
			TextColor: Color{
				R: 238,
				G: 238,
				B: 238,
				A: 255,
			},
			SeriesColors: vividSeriesColors,
		},
	)
	AddTheme(
		ThemeAnt,
		ThemeOption{
			IsDarkMode: false,
			AxisStrokeColor: Color{
				R: 110,
				G: 112,
				B: 121,
				A: 255,
			},
			AxisSplitLineColor: Color{
				R: 224,
				G: 230,
				B: 242,
				A: 255,
			},
			BackgroundColor: drawing.ColorWhite,
			TextColor: drawing.Color{
				R: 70,
				G: 70,
				B: 70,
				A: 255,
			},
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
	AddTheme(
		ThemeGrafana,
		ThemeOption{
			IsDarkMode: true,
			AxisStrokeColor: Color{
				R: 185,
				G: 184,
				B: 206,
				A: 255,
			},
			AxisSplitLineColor: Color{
				R: 68,
				G: 67,
				B: 67,
				A: 255,
			},
			BackgroundColor: drawing.Color{
				R: 31,
				G: 29,
				B: 29,
				A: 255,
			},
			TextColor: Color{
				R: 216,
				G: 217,
				B: 218,
				A: 255,
			},
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
	SetDefaultTheme(ThemeLight)
}

// SetDefaultTheme sets default theme
func SetDefaultTheme(name string) {
	defaultTheme = NewTheme(name)
}

func AddTheme(name string, opt ThemeOption) {
	palettes[name] = &themeColorPalette{
		isDarkMode:         opt.IsDarkMode,
		axisStrokeColor:    opt.AxisStrokeColor,
		axisSplitLineColor: opt.AxisSplitLineColor,
		backgroundColor:    opt.BackgroundColor,
		textColor:          opt.TextColor,
		seriesColors:       opt.SeriesColors,
	}
}

func NewTheme(name string) ColorPalette {
	p, ok := palettes[name]
	if !ok {
		p = palettes[ThemeLight]
	}
	clone := *p
	return &clone
}

func NewThemeWithFont(name string, font *truetype.Font) ColorPalette {
	p, ok := palettes[name]
	if !ok {
		p = palettes[ThemeLight]
	}
	clone := *p
	clone.font = font
	return &clone
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

func (t *themeColorPalette) GetFontSize() float64 {
	if t.fontSize != 0 {
		return t.fontSize
	}
	return defaultFontSize
}

func (t *themeColorPalette) GetFont() *truetype.Font {
	if t.font != nil {
		return t.font
	}
	f, _ := GetDefaultFont()
	return f
}

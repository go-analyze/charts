package charts

import (
	"fmt"
	"hash/crc32"
	"sync"

	"github.com/go-analyze/charts/chartdraw"
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
	// WithTextColor will provide a new ColorPalette that uses the specified color for text.
	// This is generally recommended over using the FontColor config values.
	WithTextColor(Color) ColorPalette
	// WithSeriesColors will provide a new ColorPalette that uses the specified series colors.
	WithSeriesColors([]Color) ColorPalette
	// WithBackgroundColor will provide a new ColorPalette that uses the specified color for the background.
	WithBackgroundColor(Color) ColorPalette
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

var defaultLightFontColor = Color{R: 70, G: 70, B: 70, A: 255}
var defaultDarkFontColor = Color{R: 238, G: 238, B: 238, A: 255}
var defaultGlobalMarkFillColor = ColorLightGray

func init() {
	echartSeriesColors := []Color{
		{ // blue
			R: 84, G: 112, B: 198, A: 255,
		},
		{ // green
			R: 145, G: 204, B: 117, A: 255,
		},
		ColorOrangeAlt2,
		{ // red
			R: 238, G: 102, B: 102, A: 255,
		},
		{ // aqua
			R: 115, G: 192, B: 222, A: 255,
		},
		ColorGreenAlt3,
		{ // dark orange
			R: 252, G: 132, B: 82, A: 255,
		},
		{ // dark purple
			R: 154, G: 96, B: 180, A: 255,
		},
		{ // light purple
			R: 234, G: 124, B: 204, A: 255,
		},
	}
	InstallTheme(
		ThemeLight,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 110, G: 112, B: 121, A: 255},
			AxisSplitLineColor: Color{R: 224, G: 230, B: 242, A: 255},
			BackgroundColor:    ColorWhite,
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
			BackgroundColor:    ColorDarkGray,
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
			BackgroundColor:    ColorWhite,
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
			BackgroundColor:    ColorDarkGray,
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
			BackgroundColor:    ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			SeriesColors: []Color{
				{ // light blue
					R: 91, G: 143, B: 249, A: 255,
				},
				{ // light green
					R: 90, G: 216, B: 166, A: 255,
				},
				{ // dark blue
					R: 93, G: 112, B: 146, A: 255,
				},
				{ // dark yellow
					R: 246, G: 189, B: 22, A: 255,
				},
				{ // blue
					R: 111, G: 94, B: 249, A: 255,
				},
				{ // aqua
					R: 109, G: 200, B: 236, A: 255,
				},
				{ // purple
					R: 148, G: 95, B: 185, A: 255,
				},
				ColorOrangeAlt3,
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
				{ // dark green
					R: 126, G: 178, B: 109, A: 255,
				},
				{ // orange
					R: 234, G: 184, B: 57, A: 255,
				},
				{ // aqua
					R: 110, G: 208, B: 224, A: 255,
				},
				{ // orange
					R: 239, G: 132, B: 60, A: 255,
				},
				ColorRedAlt2,
				{ // dark blue
					R: 31, G: 120, B: 193, A: 255,
				},
				{ // dark purple
					R: 112, G: 93, B: 160, A: 255,
				},
				ColorGreenAlt4,
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

// MakeTheme constructs a one-off theme without installing it into the catalog.
func MakeTheme(opt ThemeOption) ColorPalette {
	optStr := fmt.Sprintf("%v", opt)
	optId := crc32.ChecksumIEEE([]byte(optStr))
	return &themeColorPalette{
		name:               fmt.Sprintf("custom-%x", optId),
		isDarkMode:         opt.IsDarkMode,
		axisStrokeColor:    opt.AxisStrokeColor,
		axisSplitLineColor: opt.AxisSplitLineColor,
		backgroundColor:    opt.BackgroundColor,
		textColor:          opt.TextColor,
		seriesColors:       opt.SeriesColors,
	}
}

// InstallTheme adds a theme to the catalog which can later be retrieved using GetTheme.
func InstallTheme(name string, opt ThemeOption) {
	cp := &themeColorPalette{
		name:               name,
		isDarkMode:         opt.IsDarkMode,
		axisStrokeColor:    opt.AxisStrokeColor,
		axisSplitLineColor: opt.AxisSplitLineColor,
		backgroundColor:    opt.BackgroundColor,
		textColor:          opt.TextColor,
		seriesColors:       opt.SeriesColors,
	}
	palettes.Store(name, cp)
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
	colorCount := len(colors)
	if index < colorCount {
		return colors[index]
	} else {
		result := colors[index%colorCount]
		// adjust the color shade automatically
		rMax, gMax, bMax := 200, 200, 200
		var rMin, gMin, bMin int
		// the adjustment amount and mod count must be balanced to ensure colors don't hit their limits quickly
		adjustment := 40 * ((index / colorCount) % 3)
		if t.IsDark() { // adjust the shade darker for dark themes
			adjustment *= -1
			rMax, gMax, bMax = 255, 255, 255
			rMin, gMin, bMin = 40, 40, 40
		}
		if result.R != result.G || result.R != result.B {
			// try to ensure the brightest channel maintains emphasis
			if result.R >= result.G && result.R >= result.B {
				rMin += 80
				gMax -= 20
				bMax -= 20
			} else if result.G >= result.R && result.G >= result.B {
				gMin += 80
				rMax -= 20
				bMax -= 20
			} else {
				bMin += 80
				rMax -= 20
				gMax -= 20
			}
		}

		result.R = uint8(chartdraw.MaxInt(chartdraw.MinInt(int(result.R)+adjustment, rMax), rMin))
		result.G = uint8(chartdraw.MaxInt(chartdraw.MinInt(int(result.G)+adjustment, gMax), gMin))
		result.B = uint8(chartdraw.MaxInt(chartdraw.MinInt(int(result.B)+adjustment, bMax), bMin))

		return result
	}
}

func (t *themeColorPalette) GetBackgroundColor() Color {
	return t.backgroundColor
}

func (t *themeColorPalette) GetTextColor() Color {
	return t.textColor
}

func (t *themeColorPalette) WithAxisColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-axis_mod"
	copy.axisStrokeColor = c
	copy.textColor = c
	return &copy
}

func (t *themeColorPalette) WithTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-text_mod"
	copy.textColor = c
	return &copy
}

func (t *themeColorPalette) WithSeriesColors(colors []Color) ColorPalette {
	copy := *t
	if len(colors) == 0 { // ignore invalid input rather than panic later
		copy.name += "-ignored_invalid_series_mod"
		return &copy
	}
	copy.name += "-series_mod"
	copy.seriesColors = colors
	return &copy
}

func (t *themeColorPalette) WithBackgroundColor(color Color) ColorPalette {
	copy := *t
	copy.name += "-background_mod"
	copy.backgroundColor = color
	updatedDark := !isLightColor(color)
	if copy.isDarkMode != updatedDark {
		copy.isDarkMode = updatedDark
		if copy.isDarkMode {
			copy.name += "_dark"
		} else {
			copy.name += "_light"
		}
	}
	return &copy
}

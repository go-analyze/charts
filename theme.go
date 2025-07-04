package charts

import (
	"fmt"
	"hash/crc32"
	"sync"

	"github.com/go-analyze/charts/chartdraw"
)

const (
	// ThemeLight is the default theme used, with series colors from echarts.
	ThemeLight = "light"
	// ThemeDark is a dark alternative to the default theme 'light, with series colors from echarts'.
	ThemeDark = "dark"
	// ThemeVividLight is an alternative light theme that has red, yellow, and other bright colors initially in the series.
	// It can be a good option when you want the first few series items to grab the most attention.
	ThemeVividLight = "vivid-light"
	// ThemeVividDark is a dark alternative to 'ThemeVividLight', with the same bright initial series colors.
	ThemeVividDark = "vivid-dark"
	// ThemeGrafana is a grafana styled theme.
	ThemeGrafana = "grafana"
	// ThemeAnt is an ant styled theme.
	ThemeAnt = "ant"
	// ThemeNatureLight provides earthy color tones.
	ThemeNatureLight = "nature-light"
	// ThemeNatureDark provides earthy color tones with a dark background.
	ThemeNatureDark = "nature-dark"
	// ThemeRetro provides colors from the 50's and 60's, silver, maroon, tan, and other vintage colors.
	ThemeRetro = "retro"
	// ThemeOcean is a light colored theme that focuses on shades of green, blue, and other ocean colors.
	ThemeOcean = "ocean"
	// ThemeSlate is a dark theme with a slate background, and light pastel series colors.
	ThemeSlate = "slate"
	// ThemeGray is a light theme that only contains shades of gray.
	ThemeGray = "gray"
	// ThemeWinter is a light theme with shades of white and blue, some light purple.
	ThemeWinter = "winter"
	// ThemeSpring is a light theme with bright greens, yellows, and blues.
	ThemeSpring = "spring"
	// ThemeSummer is a light theme with red, orange, and yellow shades.
	ThemeSummer = "summer"
	// ThemeFall is a dark theme with shades of yellow, orange and brown.
	ThemeFall = "fall"
)

// ColorPalette provides the theming for the chart.
type ColorPalette interface {
	IsDark() bool
	GetXAxisStrokeColor() Color
	GetYAxisStrokeColor() Color
	GetAxisSplitLineColor() Color
	GetSeriesColor(int) Color
	GetSeriesTrendColor(int) Color
	GetBackgroundColor() Color
	GetTitleTextColor() Color
	GetMarkTextColor() Color
	GetLabelTextColor() Color
	GetLegendTextColor() Color
	GetXAxisTextColor() Color
	GetYAxisTextColor() Color
	GetTitleBorderColor() Color
	GetLegendBorderColor() Color
	// WithXAxisColor will provide a new ColorPalette that uses the specified color for X axis. To adjust the text
	// color invoke WithXAxisTextColor following this.
	WithXAxisColor(Color) ColorPalette
	// WithYAxisColor will provide a new ColorPalette that uses the specified color for Y axis. To adjust the text
	// color invoke WithYAxisTextColor following this.
	WithYAxisColor(Color) ColorPalette
	// WithYAxisSeriesColor will provide a new ColorPalette that uses the specified series index color for Y axis and values.
	WithYAxisSeriesColor(int) ColorPalette
	// WithTitleTextColor will provide a new ColorPalette that uses the specified color for the title text.
	WithTitleTextColor(Color) ColorPalette
	// WithMarkTextColor will provide a new ColorPalette that uses the specified color for mark point and mark line labels.
	WithMarkTextColor(Color) ColorPalette
	// WithLabelTextColor will provide a new ColorPalette that uses the specified color for value labels.
	WithLabelTextColor(Color) ColorPalette
	// WithLegendTextColor will provide a new ColorPalette that uses the specified color for the legend labels
	WithLegendTextColor(Color) ColorPalette
	// WithXAxisTextColor will provide a new ColorPalette that uses the specified color for the x-axis labels.
	WithXAxisTextColor(Color) ColorPalette
	// WithYAxisTextColor will provide a new ColorPalette that uses the specified color for the y-axis labels.
	WithYAxisTextColor(Color) ColorPalette
	// WithSeriesColors will provide a new ColorPalette that uses the specified series colors. This will default the
	// trend line colors to be related to the series colors provided. If you want to customize them further use
	// WithSeriesTrendColors.
	WithSeriesColors([]Color) ColorPalette
	// WithSeriesTrendColors will provide a new ColorPalette that uses the specified series trend line colors.
	WithSeriesTrendColors([]Color) ColorPalette
	// WithBackgroundColor will provide a new ColorPalette that uses the specified color for the background.
	WithBackgroundColor(Color) ColorPalette
	// WithTitleBorderColor will provide a new ColorPalette that uses the specified color for the title border.
	WithTitleBorderColor(Color) ColorPalette
	// WithLegendBorderColor will provide a new ColorPalette that uses the specified color for the legend border.
	WithLegendBorderColor(Color) ColorPalette
}

type themeColorPalette struct {
	name               string
	isDarkMode         bool
	xaxisStrokeColor   Color
	yaxisStrokeColor   Color
	axisSplitLineColor Color
	backgroundColor    Color
	titleTextColor     Color
	markTextColor      Color
	labelTextColor     Color
	legendTextColor    Color
	xaxisTextColor     Color
	yaxisTextColor     Color
	titleBorderColor   Color
	legendBorderColor  Color
	seriesColors       []Color
	seriesTrendColors  []Color
}

func (t *themeColorPalette) GetTitleTextColor() Color {
	return t.titleTextColor
}

func (t *themeColorPalette) GetMarkTextColor() Color {
	return t.markTextColor
}

func (t *themeColorPalette) GetLabelTextColor() Color {
	return t.labelTextColor
}

func (t *themeColorPalette) GetLegendTextColor() Color {
	return t.legendTextColor
}

func (t *themeColorPalette) GetXAxisTextColor() Color {
	return t.xaxisTextColor
}

func (t *themeColorPalette) GetYAxisTextColor() Color {
	return t.yaxisTextColor
}

// ThemeOption defines color options for a theme.
type ThemeOption struct {
	// IsDarkMode indicates whether the theme is designed for dark backgrounds, affecting color adjustments and text visibility.
	IsDarkMode bool
	// AxisStrokeColor is the default stroke color for both axes.
	AxisStrokeColor Color
	// XAxisStrokeColor overrides AxisStrokeColor for the x-axis.
	XAxisStrokeColor Color
	// YAxisStrokeColor overrides AxisStrokeColor for the y-axis.
	YAxisStrokeColor Color
	// AxisSplitLineColor sets the color of grid lines drawn between ticks.
	AxisSplitLineColor Color
	// BackgroundColor sets the chart background.
	BackgroundColor Color
	// TextColor is the default font color applied if specific text colors are unset.
	TextColor Color
	// TextColorTitle sets the title text color.
	TextColorTitle Color
	// TextColorMark sets the color of mark line and point labels.
	TextColorMark Color
	// TextColorLabel defines the color of series labels.
	TextColorLabel Color
	// TextColorLegend defines the legend text color.
	TextColorLegend Color
	// TextColorXAxis defines the x-axis label text color.
	TextColorXAxis Color
	// TextColorYAxis defines the y-axis label text color.
	TextColorYAxis Color
	// TitleBorderColor draws an optional border around the title.
	TitleBorderColor Color
	// LegendBorderColor draws an optional border around the legend.
	LegendBorderColor Color
	// SeriesColors provides the color palette used for series data.
	SeriesColors []Color
	// SeriesTrendColors provides the palette for rendered trend lines.
	SeriesTrendColors []Color
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
	natureSeriesColors := []Color{
		ColorSageGreen,
		{ // Terracotta
			R: 242, G: 153, B: 119, A: 255,
		},
		{ // Sky blue alt
			R: 130, G: 175, B: 222, A: 255,
		},
		{ // Moss green
			R: 121, G: 191, B: 127, A: 255,
		},
		{ // Light Forest brown
			R: 171, G: 136, B: 100, A: 255,
		},
		ColorDesertSand,
		{ // Ocean blue
			R: 100, G: 150, B: 180, A: 255,
		},
		{ // Clay red
			R: 203, G: 134, B: 115, A: 255,
		},
		{ // Earthy olive
			R: 135, G: 164, B: 112, A: 255,
		},
		{ // Driftwood gray-brown
			R: 145, G: 133, B: 116, A: 255,
		},
		{ // River stone gray
			R: 128, G: 146, B: 140, A: 255,
		},
	}
	greenHeaderText := ColorGreenAlt3.WithAdjustHSL(0, 0, -0.2)
	InstallTheme(
		ThemeNatureLight,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 138, G: 142, B: 146, A: 255},
			AxisSplitLineColor: Color{R: 200, G: 203, B: 208, A: 255},
			BackgroundColor:    ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			TextColorTitle:     greenHeaderText,
			TextColorLegend:    greenHeaderText,
			SeriesColors:       natureSeriesColors,
		},
	)
	InstallTheme(
		ThemeNatureDark,
		ThemeOption{
			IsDarkMode:         true,
			AxisStrokeColor:    Color{R: 138, G: 142, B: 146, A: 255},
			AxisSplitLineColor: Color{R: 72, G: 71, B: 83, A: 255},
			BackgroundColor:    ColorDarkGray,
			TextColor:          Color{R: 238, G: 238, B: 238, A: 255},
			TextColorTitle:     ColorWhite,
			SeriesColors:       natureSeriesColors,
		},
	)
	InstallTheme(
		ThemeRetro,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 140, G: 135, B: 130, A: 255},
			AxisSplitLineColor: Color{R: 190, G: 190, B: 190, A: 255},
			BackgroundColor:    Color{R: 255, G: 250, B: 240, A: 255},
			TextColor:          Color{R: 51, G: 36, B: 16, A: 255},
			SeriesColors: []Color{
				ColorMaroon,
				{ // sage olive green
					R: 145, G: 150, B: 99, A: 255,
				},
				ColorTan,
				{ // dark orange
					R: 184, G: 90, B: 0, A: 255,
				},
				{ // brown
					R: 101, G: 67, B: 33, A: 255,
				},
				ColorMustardYellow,
				ColorTeal,
				{ // light red
					R: 200, G: 80, B: 50, A: 255,
				},
				{ // navy blue
					R: 25, G: 42, B: 64, A: 255,
				},
			},
		},
	)
	blueHeaderText := ColorBlue.WithAdjustHSL(0, 0, -0.2)
	InstallTheme(
		ThemeOcean,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 120, G: 130, B: 140, A: 255},
			AxisSplitLineColor: Color{R: 200, G: 210, B: 220, A: 255},
			BackgroundColor:    ColorWhite,
			TextColor:          Color{R: 0, G: 45, B: 72, A: 255},
			TextColorTitle:     blueHeaderText,
			TextColorLegend:    blueHeaderText,
			SeriesColors: []Color{
				ColorSkyBlue,
				{ // seafoam green
					R: 90, G: 210, B: 160, A: 255,
				},
				{ // coral pink
					R: 252, G: 163, B: 148, A: 255,
				},
				{ // light purple
					R: 180, G: 140, B: 210, A: 255,
				},
				ColorBlueAlt1,
				{ // light teal
					R: 110, G: 220, B: 210, A: 255,
				},
				ColorPink,
				ColorPlum,
			},
		},
	)
	InstallTheme(
		ThemeSlate,
		ThemeOption{
			IsDarkMode:         true,
			AxisStrokeColor:    Color{R: 128, G: 129, B: 132, A: 255},
			AxisSplitLineColor: Color{R: 72, G: 71, B: 73, A: 255},
			BackgroundColor:    Color{R: 51, G: 53, B: 60, A: 255}, // Slate with slight blue tint
			TextColor:          ColorWhite,
			TextColorLegend:    ColorSlateGray.WithAdjustHSL(0, 0, 0.2),
			SeriesColors: []Color{
				ColorLightCoral,
				{ // pale aqua
					R: 125, G: 210, B: 196, A: 255,
				},
				{ // pale yellow
					R: 250, G: 240, B: 120, A: 255,
				},
				{ // pale blue
					R: 170, G: 190, B: 255, A: 255,
				},
				{ // pale purple
					R: 190, G: 164, B: 240, A: 255,
				},
				{ // pale salmon
					R: 250, G: 160, B: 140, A: 255,
				},
				{ // pale sage green
					R: 180, G: 200, B: 165, A: 255,
				},
				{ // dusty yellow
					R: 238, G: 214, B: 63, A: 255,
				},
				{ // slate blue
					R: 100, G: 130, B: 170, A: 255,
				},
				{ // dusty purple
					R: 178, G: 132, B: 173, A: 255,
				},
				{ // muted green
					R: 120, G: 160, B: 140, A: 255,
				},
			},
		},
	)
	InstallTheme(
		ThemeGray,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 142, G: 142, B: 142, A: 255},
			AxisSplitLineColor: Color{R: 204, G: 204, B: 204, A: 255},
			BackgroundColor:    ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			TextColorTitle:     ColorBlack,
			SeriesColors: []Color{
				ColorDarkGray,
				{R: 88, G: 88, B: 88, A: 255},
				ColorGray,
				{R: 160, G: 160, B: 160, A: 255},
				{R: 190, G: 190, B: 190, A: 255},
				ColorLightGray,
				{R: 228, G: 228, B: 228, A: 255},
				{R: 248, G: 248, B: 248, A: 255},
			},
		},
	)
	InstallTheme(
		ThemeWinter,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 140, G: 145, B: 150, A: 255},
			AxisSplitLineColor: Color{R: 200, G: 210, B: 220, A: 255},
			BackgroundColor:    ColorAzure,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			TextColorTitle:     blueHeaderText,
			TextColorLegend:    blueHeaderText,
			SeriesColors: []Color{
				{ // Light blue
					R: 150, G: 190, B: 255, A: 255,
				},
				ColorPlum,
				{ // Frosty blue
					R: 110, G: 150, B: 240, A: 255,
				},
				{ // Pale lavender
					R: 210, G: 180, B: 230, A: 255,
				},
				{ // Ice blue
					R: 90, G: 130, B: 210, A: 255,
				},
				{ // Soft purple
					R: 190, G: 160, B: 220, A: 255,
				},
				{ // Glacier blue
					R: 80, G: 110, B: 190, A: 255,
				},
			},
		},
	)
	InstallTheme(
		ThemeSpring,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 120, G: 130, B: 140, A: 255},
			AxisSplitLineColor: Color{R: 200, G: 210, B: 220, A: 255},
			BackgroundColor:    ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			TextColorTitle:     greenHeaderText,
			TextColorLegend:    greenHeaderText,
			SeriesColors: []Color{
				{ // Lime green
					R: 120, G: 200, B: 130, A: 255,
				},
				{ // Golden yellow
					R: 220, G: 210, B: 100, A: 255,
				},
				ColorSkyBlue,
				{ // Mint green
					R: 150, G: 230, B: 180, A: 255,
				},
				{ // Sun yellow
					R: 240, G: 220, B: 120, A: 255,
				},
				{ // Light teal
					R: 110, G: 210, B: 200, A: 255,
				},
			},
		},
	)
	InstallTheme(
		ThemeSummer,
		ThemeOption{
			IsDarkMode:         false,
			AxisStrokeColor:    Color{R: 150, G: 140, B: 130, A: 255},
			AxisSplitLineColor: Color{R: 220, G: 210, B: 200, A: 255},
			BackgroundColor:    ColorWhite,
			TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
			TextColorTitle:     ColorBlack,
			SeriesColors: []Color{
				ColorSalmon,
				{ // Bright yellow
					R: 230, G: 220, B: 110, A: 255,
				},
				ColorOrange,
				ColorOrangeAlt1,
				{ // Coral pink
					R: 250, G: 160, B: 140, A: 255,
				},
				{ // pale yellow
					R: 250, G: 240, B: 120, A: 255,
				},
				ColorOrangeAlt2,
			},
		},
	)
	InstallTheme(
		ThemeFall,
		ThemeOption{
			IsDarkMode:         true,
			AxisStrokeColor:    Color{R: 130, G: 90, B: 60, A: 255},
			AxisSplitLineColor: Color{R: 160, G: 120, B: 80, A: 255},
			BackgroundColor:    ColorDarkGray,
			TextColor:          ColorWhite,
			TextColorLegend:    ColorLightGray,
			SeriesColors: []Color{
				{ // Golden yellow
					R: 220, G: 190, B: 110, A: 255,
				},
				{ // Copper orange
					R: 200, G: 130, B: 70, A: 255,
				},
				{ // Burnt brown
					R: 180, G: 120, B: 60, A: 255,
				},
				{ // Amber yellow
					R: 240, G: 190, B: 90, A: 255,
				},
				{ // Rust orange
					R: 230, G: 140, B: 50, A: 255,
				},
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
	cp := makeColorPalette(opt)
	cp.name = fmt.Sprintf("custom-%x", crc32.ChecksumIEEE([]byte(fmt.Sprintf("%v", opt))))
	return cp
}

// InstallTheme adds a theme to the catalog which can later be retrieved using GetTheme.
func InstallTheme(name string, opt ThemeOption) {
	cp := makeColorPalette(opt)
	cp.name = name
	palettes.Store(name, cp)
}

func makeColorPalette(o ThemeOption) *themeColorPalette {
	if o.XAxisStrokeColor.IsZero() {
		o.XAxisStrokeColor = o.AxisStrokeColor
	}
	if o.YAxisStrokeColor.IsZero() {
		o.YAxisStrokeColor = o.AxisStrokeColor
	}
	if o.TextColor.IsZero() {
		o.TextColor = ColorBlack
	}
	if o.TextColorLabel.IsZero() {
		o.TextColorLabel = o.TextColor
	}
	if o.TextColorTitle.IsZero() {
		o.TextColorTitle = o.TextColor
	}
	if o.TextColorMark.IsZero() {
		o.TextColorMark = o.TextColor
	}
	if o.TextColorLegend.IsZero() {
		o.TextColorLegend = o.TextColor
	}
	if o.TextColorXAxis.IsZero() {
		o.TextColorXAxis = o.TextColor
	}
	if o.TextColorYAxis.IsZero() {
		o.TextColorYAxis = o.TextColor
	}
	if o.LegendBorderColor.IsZero() {
		o.LegendBorderColor = ColorBlack
	}
	if o.TitleBorderColor.IsZero() {
		o.TitleBorderColor = ColorBlack
	}
	for i := len(o.SeriesTrendColors); i < len(o.SeriesColors); i++ {
		o.SeriesTrendColors = append(o.SeriesTrendColors, autoSeriesTrendColor(o.SeriesColors[i]))
	}
	return &themeColorPalette{
		isDarkMode:         o.IsDarkMode,
		xaxisStrokeColor:   o.XAxisStrokeColor,
		yaxisStrokeColor:   o.YAxisStrokeColor,
		axisSplitLineColor: o.AxisSplitLineColor,
		backgroundColor:    o.BackgroundColor,
		titleBorderColor:   o.TitleBorderColor,
		legendBorderColor:  o.LegendBorderColor,
		titleTextColor:     o.TextColorTitle,
		markTextColor:      o.TextColorMark,
		labelTextColor:     o.TextColorLabel,
		legendTextColor:    o.TextColorLegend,
		xaxisTextColor:     o.TextColorXAxis,
		yaxisTextColor:     o.TextColorYAxis,
		seriesColors:       o.SeriesColors,
		seriesTrendColors:  o.SeriesTrendColors,
	}
}

func autoSeriesTrendColor(color Color) Color {
	if color.IsTransparent() {
		return color
	}
	c := color.WithAdjustHSL(0.0, 0.1, -0.1)
	if c.A < 255 {
		c.A += (255 - c.A) / 2
	}
	return c
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

func (t *themeColorPalette) GetXAxisStrokeColor() Color {
	return t.xaxisStrokeColor
}

func (t *themeColorPalette) GetYAxisStrokeColor() Color {
	return t.yaxisStrokeColor
}

func (t *themeColorPalette) GetAxisSplitLineColor() Color {
	return t.axisSplitLineColor
}

func (t *themeColorPalette) GetSeriesColor(index int) Color {
	return getSeriesColor(t.seriesColors, t.isDarkMode, index)
}

func (t *themeColorPalette) GetSeriesTrendColor(index int) Color {
	return getSeriesColor(t.seriesTrendColors, t.isDarkMode, index)
}

func getSeriesColor(colors []Color, darkTheme bool, index int) Color {
	colorCount := len(colors)
	if index < colorCount {
		return colors[index]
	} else {
		return adjustSeriesColor(colors[index%colorCount], index/colorCount, darkTheme)
	}
}

func adjustSeriesColor(c Color, loopCount int, darkTheme bool) Color {
	impact := ((loopCount - 1) % 3) + 1
	satAdj := float64(impact) * -0.1
	ltAdj := 0.08
	if chartdraw.AbsInt(int(c.R)-int(c.G)) < 20 && chartdraw.AbsInt(int(c.R)-int(c.B)) < 20 {
		ltAdj += 0.1 // more adjustment if close to gray
	}
	ltAdj *= float64(impact)
	if darkTheme {
		ltAdj *= -1
	}
	return c.WithAdjustHSL(0.0, satAdj, ltAdj)
}

func (t *themeColorPalette) GetBackgroundColor() Color {
	return t.backgroundColor
}

func (t *themeColorPalette) GetTitleBorderColor() Color {
	return t.titleBorderColor
}

func (t *themeColorPalette) GetLegendBorderColor() Color {
	return t.legendBorderColor
}

func (t *themeColorPalette) WithXAxisColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-xaxis_stroke_mod"
	copy.xaxisStrokeColor = c
	return &copy
}

func (t *themeColorPalette) WithYAxisColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-yaxis_stroke_mod"
	copy.yaxisStrokeColor = c
	return &copy
}

func (t *themeColorPalette) WithYAxisSeriesColor(series int) ColorPalette {
	copy := *t
	copy.name += "-yaxis_mod"
	seriesColor := t.GetSeriesColor(series)
	copy.yaxisStrokeColor = seriesColor
	copy.yaxisTextColor = seriesColor
	return &copy
}

func (t *themeColorPalette) WithTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-text_mod"
	copy.titleTextColor = c
	copy.markTextColor = c
	copy.labelTextColor = c
	copy.legendTextColor = c
	copy.xaxisTextColor = c
	copy.yaxisTextColor = c
	return &copy
}

func (t *themeColorPalette) WithTitleTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-title_mod"
	copy.titleTextColor = c
	return &copy
}

func (t *themeColorPalette) WithMarkTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-mark_mod"
	copy.markTextColor = c
	return &copy
}

func (t *themeColorPalette) WithLabelTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-label_mod"
	copy.labelTextColor = c
	return &copy
}

func (t *themeColorPalette) WithLegendTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-legend_text_mod"
	copy.legendTextColor = c
	return &copy
}

func (t *themeColorPalette) WithXAxisTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-xaxis_text_mod"
	copy.xaxisTextColor = c
	return &copy
}

func (t *themeColorPalette) WithYAxisTextColor(c Color) ColorPalette {
	copy := *t
	copy.name += "-yaxis_text_mod"
	copy.yaxisTextColor = c
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
	for i, c := range colors {
		trendColor := autoSeriesTrendColor(c)
		if i < len(copy.seriesTrendColors) {
			copy.seriesTrendColors[i] = trendColor
		} else {
			copy.seriesTrendColors = append(copy.seriesTrendColors, trendColor)
		}
	}
	return &copy
}

func (t *themeColorPalette) WithSeriesTrendColors(colors []Color) ColorPalette {
	copy := *t
	if len(colors) == 0 { // ignore invalid input rather than panic later
		copy.name += "-ignored_invalid_series_mod"
		return &copy
	}
	copy.name += "-trend_mod"
	copy.seriesTrendColors = colors
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

func (t *themeColorPalette) WithTitleBorderColor(color Color) ColorPalette {
	copy := *t
	copy.name += "-title_border_mod"
	copy.titleBorderColor = color
	return &copy
}

func (t *themeColorPalette) WithLegendBorderColor(color Color) ColorPalette {
	copy := *t
	copy.name += "-legend_border_mod"
	copy.legendBorderColor = color
	return &copy
}

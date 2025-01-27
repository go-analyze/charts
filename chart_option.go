package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"
)

type ChartOption struct {
	// OutputFormat specifies the output type of chart, "svg" or "png", default value is "png".
	OutputFormat string
	// Width is the width of chart, default width is 600.
	Width int
	// Height is the height of chart, default height is 400.
	Height int
	// Theme specifies the colors used for the chart. Built in themes can be loaded using GetTheme with
	// "light", "dark", "vivid-light", "vivid-dark", "ant" or "grafana".
	Theme ColorPalette
	// Padding specifies the padding for chart, default padding is [20, 10, 10, 10].
	Padding Box
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// Font is the font to use for rendering the chart.
	Font *truetype.Font
	// Box specifies the canvas box for the chart.
	Box Box
	// SeriesList provides the data series.
	SeriesList SeriesList
	// StackSeries if set to *true the lines will be layered or stacked. This option significantly changes the chart
	// visualization, please see the specific chart docs for full details.
	StackSeries *bool
	// RadarIndicators are radar indicator list for radar charts.
	RadarIndicators []RadarIndicator
	// Symbol specifies the symbols to draw at the data points. Empty (default) will vary based on the dataset.
	// Specify 'none' to enforce no symbol, or specify a desired symbol: 'circle', 'dot', 'square', 'diamond'.
	Symbol Symbol
	// LineStrokeWidth is the stroke width for line charts.
	LineStrokeWidth float64
	// FillArea set to *true to fill the area under the line in line charts
	FillArea *bool
	// FillOpacity is the opacity (alpha) of the area fill in line charts.
	FillOpacity uint8
	// BarWidth is the width of the bars for bar charts.
	BarWidth int
	// BarHeight is the height of the bars for horizontal bar charts.
	BarHeight int
	// BarMargin specifies the margin between bars grouped together. BarWidth or BarHeight takes priority over the margin.
	BarMargin *float64
	// Radius default radius for pie and radar charts e.g.: 40%, default is "40%"
	Radius string
	// Children are Child charts to render together.
	Children []ChartOption
	parent   *Painter
	// ValueFormatter to format numeric values into labels.
	ValueFormatter ValueFormatter
}

// OptionFunc option function.
type OptionFunc func(opt *ChartOption)

// SVGOutputOptionFunc set svg type of chart's output.
func SVGOutputOptionFunc() OptionFunc {
	return outputFormatOptionFunc(ChartOutputSVG)
}

// PNGOutputOptionFunc set png type of chart's output.
func PNGOutputOptionFunc() OptionFunc {
	return outputFormatOptionFunc(ChartOutputPNG)
}

// outputFormatOptionFunc set type of chart's output.
func outputFormatOptionFunc(t string) OptionFunc {
	return func(opt *ChartOption) {
		opt.OutputFormat = t
	}
}

// FontOptionFunc set the default font of the chart.
func FontOptionFunc(font *truetype.Font) OptionFunc {
	return func(opt *ChartOption) {
		opt.Font = font
	}
}

// ThemeNameOptionFunc set them of chart by name.
func ThemeNameOptionFunc(theme string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Theme = GetTheme(theme)
	}
}

// ThemeOptionFunc sets the theme of the chart.
func ThemeOptionFunc(theme ColorPalette) OptionFunc {
	return func(opt *ChartOption) {
		opt.Theme = theme
	}
}

// TitleOptionFunc sets the title of the chart.
func TitleOptionFunc(title TitleOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Title = title
	}
}

// TitleTextOptionFunc sets the title text of chart.
func TitleTextOptionFunc(text string, subtext ...string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Title.Text = text
		if len(subtext) != 0 {
			opt.Title.Subtext = subtext[0]
		}
	}
}

// LegendOptionFunc sets the legend of the chart.
func LegendOptionFunc(legend LegendOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = legend
	}
}

// LegendLabelsOptionFunc sets the legend series name labels of the chart.
func LegendLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = LegendOption{
			SeriesNames: labels,
		}
	}
}

// XAxisOptionFunc sets the x-axis of the chart.
func XAxisOptionFunc(xAxisOption XAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = xAxisOption
	}
}

// XAxisLabelsOptionFunc sets the x-axis labels of the chart.
func XAxisLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = XAxisOption{
			Labels: labels,
		}
	}
}

// YAxisOptionFunc sets the y-axis of chart, supports up to two y-axis.
func YAxisOptionFunc(yAxisOption ...YAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxis = yAxisOption
	}
}

// YAxisLabelsOptionFunc sets the y-axis labels of the chart.
func YAxisLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxis = []YAxisOption{
			{
				Labels: labels,
			},
		}
	}
}

// DimensionsOptionFunc sets the width and height dimensions of the chart.
func DimensionsOptionFunc(width, height int) OptionFunc {
	return func(opt *ChartOption) {
		opt.Width = width
		opt.Height = height
	}
}

// PaddingOptionFunc sets the padding of the chart.
func PaddingOptionFunc(padding Box) OptionFunc {
	return func(opt *ChartOption) {
		opt.Padding = padding
	}
}

// SeriesShowLabel sets the series show label state for all series.
func SeriesShowLabel(show bool) OptionFunc {
	return func(opt *ChartOption) {
		for index := range opt.SeriesList {
			opt.SeriesList[index].Label.Show = BoolPointer(show)
		}
	}
}

// ChildOptionFunc adds a Child chart on top of the current one. Use Padding and Box for positioning.
func ChildOptionFunc(child ...ChartOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Children = append(opt.Children, child...)
	}
}

// RadarIndicatorOptionFunc sets the radar indicator of chart
func RadarIndicatorOptionFunc(names []string, values []float64) OptionFunc {
	return func(opt *ChartOption) {
		opt.RadarIndicators = NewRadarIndicators(names, values)
	}
}

// MarkLineOptionFunc sets the mark line for series of the chart.
func MarkLineOptionFunc(seriesIndex int, markLineTypes ...string) OptionFunc {
	return func(opt *ChartOption) {
		if len(opt.SeriesList) <= seriesIndex {
			return
		}
		opt.SeriesList[seriesIndex].MarkLine = NewMarkLine(markLineTypes...)
	}
}

// MarkPointOptionFunc sets the mark point for series of the chart.
func MarkPointOptionFunc(seriesIndex int, markPointTypes ...string) OptionFunc {
	return func(opt *ChartOption) {
		if len(opt.SeriesList) <= seriesIndex {
			return
		}
		opt.SeriesList[seriesIndex].MarkPoint = NewMarkPoint(markPointTypes...)
	}
}

func (o *ChartOption) fillDefault() error {
	o.Width = getDefaultInt(o.Width, defaultChartWidth)
	o.Height = getDefaultInt(o.Height, defaultChartHeight)

	yaxisCount := o.SeriesList.getYAxisCount()
	if yaxisCount < 0 {
		return errors.New("series specified invalid y-axis index")
	}
	if len(o.YAxis) < yaxisCount {
		yAxisOptions := make([]YAxisOption, yaxisCount)
		copy(yAxisOptions, o.YAxis)
		o.YAxis = yAxisOptions
	}

	if o.Font == nil {
		o.Font = GetDefaultFont()
	}
	if o.Theme == nil {
		o.Theme = GetDefaultTheme()
	}
	fillThemeDefaults(o.Theme, &o.Title, &o.Legend, &o.XAxis)

	if o.Padding.IsZero() {
		o.Padding = defaultPadding
	}
	return nil
}

func fillThemeDefaults(defaultTheme ColorPalette, title *TitleOption, legend *LegendOption, xaxis *XAxisOption) {
	if title.Theme == nil {
		title.Theme = defaultTheme
	}
	if legend.Theme == nil {
		legend.Theme = defaultTheme
	}
	if xaxis.Theme == nil {
		xaxis.Theme = defaultTheme
	}
}

// LineRender line chart render.
func LineRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListLine(values),
	}, opts...)
}

// BarRender bar chart render.
func BarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListBar(values),
	}, opts...)
}

// HorizontalBarRender horizontal bar chart render.
func HorizontalBarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListHorizontalBar(values),
	}, opts...)
}

// PieRender pie chart render.
func PieRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListPie(values),
	}, opts...)
}

// RadarRender radar chart render.
func RadarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListRadar(values),
	}, opts...)
}

// FunnelRender funnel chart render.
func FunnelRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListFunnel(values),
	}, opts...)
}

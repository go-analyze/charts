package charts

import (
	"fmt"
	"sort"

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
	// RadarIndicators are radar indicator list for radar charts.
	RadarIndicators []RadarIndicator
	// SymbolShow set this to *false or *true (using False() or True()) to force if the symbols should be shown or hidden.
	SymbolShow *bool
	// LineStrokeWidth is the stroke width for line charts.
	LineStrokeWidth float64
	// FillArea set to true to fill the area under the line in line charts
	FillArea bool
	// FillOpacity is the opacity (alpha) of the area fill in line charts.
	FillOpacity uint8
	// BarWidth is the width of the bars for bar charts.
	BarWidth int
	// BarHeight is the height of the bars for horizontal bar charts.
	BarHeight int
	// Children are Child charts to render together.
	Children []ChartOption
	parent   *Painter
	// ValueFormatter to format numeric values into labels.
	ValueFormatter ValueFormatter
}

// OptionFunc option function
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

// FontOptionFunc set font of chart.
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

// ThemeOptionFunc set them of chart
func ThemeOptionFunc(theme ColorPalette) OptionFunc {
	return func(opt *ChartOption) {
		opt.Theme = theme
	}
}

// TitleOptionFunc set title of chart
func TitleOptionFunc(title TitleOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Title = title
	}
}

// TitleTextOptionFunc set title text of chart
func TitleTextOptionFunc(text string, subtext ...string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Title.Text = text
		if len(subtext) != 0 {
			opt.Title.Subtext = subtext[0]
		}
	}
}

// LegendOptionFunc set legend of chart
func LegendOptionFunc(legend LegendOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = legend
	}
}

// LegendLabelsOptionFunc set legend labels of chart
func LegendLabelsOptionFunc(labels []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = LegendOption{
			Data: labels,
		}
	}
}

// XAxisOptionFunc set x-axis of chart
func XAxisOptionFunc(xAxisOption XAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = xAxisOption
	}
}

// XAxisDataOptionFunc set x-axis data of chart
func XAxisDataOptionFunc(data []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = XAxisOption{
			Data: data,
		}
	}
}

// YAxisOptionFunc set y-axis of chart, supports up to two y-axis.
func YAxisOptionFunc(yAxisOption ...YAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxis = yAxisOption
	}
}

// YAxisDataOptionFunc set y-axis data of chart
func YAxisDataOptionFunc(data []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxis = []YAxisOption{
			{
				Data: data,
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

// PaddingOptionFunc set padding of chart
func PaddingOptionFunc(padding Box) OptionFunc {
	return func(opt *ChartOption) {
		opt.Padding = padding
	}
}

// PieSeriesShowLabel set pie series show label.
func PieSeriesShowLabel() OptionFunc {
	return func(opt *ChartOption) {
		for index := range opt.SeriesList {
			opt.SeriesList[index].Label.Show = true
		}
	}
}

// ChildOptionFunc adds a Child chart on top of the current one. Use Padding and Box for positioning.
func ChildOptionFunc(child ...ChartOption) OptionFunc {
	return func(opt *ChartOption) {
		if opt.Children == nil {
			opt.Children = make([]ChartOption, 0)
		}
		opt.Children = append(opt.Children, child...)
	}
}

// RadarIndicatorOptionFunc set radar indicator of chart
func RadarIndicatorOptionFunc(names []string, values []float64) OptionFunc {
	return func(opt *ChartOption) {
		opt.RadarIndicators = NewRadarIndicators(names, values)
	}
}

// MarkLineOptionFunc set mark line for series of chart
func MarkLineOptionFunc(seriesIndex int, markLineTypes ...string) OptionFunc {
	return func(opt *ChartOption) {
		if len(opt.SeriesList) <= seriesIndex {
			return
		}
		opt.SeriesList[seriesIndex].MarkLine = NewMarkLine(markLineTypes...)
	}
}

// MarkPointOptionFunc set mark point for series of chart
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

	yaxisCount := 1
	for _, series := range o.SeriesList {
		if series.YAxisIndex == 1 {
			yaxisCount++
			break
		} else if series.YAxisIndex > 1 {
			return fmt.Errorf("series '%s' specified invalid y-axis index: %v", series.Name, series.YAxisIndex)
		}
	}
	if len(o.YAxis) < yaxisCount {
		yAxisOptions := make([]YAxisOption, yaxisCount)
		copy(yAxisOptions, o.YAxis)
		o.YAxis = yAxisOptions
	}
	// TODO - this is a hack, we need to update the yaxis based on the markpoint state
	// TODO - but can't do this earlier due to needing the axis initialized
	// TODO - we should reconsider the API for configuration
	for _, sl := range o.SeriesList {
		if len(sl.MarkPoint.Data) > 0 { // if graph has markpoint
			// adjust padding scale to give space for mark point (if not specified by user)
			for i := range o.YAxis {
				if o.YAxis[i].RangeValuePaddingScale == nil {
					o.YAxis[i].RangeValuePaddingScale = FloatPointer(2.5)
				}
			}
			break
		}
	}

	if o.Font == nil {
		o.Font = GetDefaultFont()
	}
	if o.Theme == nil {
		o.Theme = GetDefaultTheme()
	}
	fillThemeDefaults(o.Theme, &o.Title, &o.Legend, &o.XAxis)

	if o.Padding.IsZero() {
		o.Padding = Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}
	}
	// association between legend and series name
	if len(o.Legend.Data) == 0 {
		o.Legend.Data = o.SeriesList.Names()
	} else {
		seriesCount := len(o.SeriesList)
		for index, name := range o.Legend.Data {
			if index < seriesCount && len(o.SeriesList[index].Name) == 0 {
				o.SeriesList[index].Name = name
			}
		}
		nameIndexDict := map[string]int{}
		for index, name := range o.Legend.Data {
			nameIndexDict[name] = index
		}
		// ensure order of series is consistent with legend
		sort.Slice(o.SeriesList, func(i, j int) bool {
			return nameIndexDict[o.SeriesList[i].Name] < nameIndexDict[o.SeriesList[j].Name]
		})
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

// LineRender line chart render
func LineRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListDataFromValues(values, ChartTypeLine),
	}, opts...)
}

// BarRender bar chart render
func BarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListDataFromValues(values, ChartTypeBar),
	}, opts...)
}

// HorizontalBarRender horizontal bar chart render
func HorizontalBarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListDataFromValues(values, ChartTypeHorizontalBar),
	}, opts...)
}

// PieRender pie chart render
func PieRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewPieSeriesList(values),
	}, opts...)
}

// RadarRender radar chart render
func RadarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewSeriesListDataFromValues(values, ChartTypeRadar),
	}, opts...)
}

// FunnelRender funnel chart render
func FunnelRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	return Render(ChartOption{
		SeriesList: NewFunnelSeriesList(values),
	}, opts...)
}

// TableRender table chart render
func TableRender(header []string, data [][]string, spanMaps ...map[int]int) (*Painter, error) {
	opt := TableChartOption{
		Header: header,
		Data:   data,
	}
	if len(spanMaps) != 0 {
		spanMap := spanMaps[0]
		spans := make([]int, len(opt.Header))
		for index := range spans {
			v, ok := spanMap[index]
			if !ok {
				v = 1
			}
			spans[index] = v
		}
		opt.Spans = spans
	}
	return TableOptionRender(opt)
}

// TableOptionRender table render with option
func TableOptionRender(opt TableChartOption) (*Painter, error) {
	if opt.OutputFormat == "" {
		opt.OutputFormat = chartDefaultOutputFormat
	}
	if opt.Width <= 0 {
		opt.Width = defaultChartWidth
	}

	p, err := NewPainter(PainterOptions{
		OutputFormat: opt.OutputFormat,
		Width:        opt.Width,
		Height:       100, // is only used to calculate the height of the table
		Font:         opt.FontStyle.Font,
	})
	if err != nil {
		return nil, err
	}
	info, err := NewTableChart(p, opt).render()
	if err != nil {
		return nil, err
	}

	p, err = NewPainter(PainterOptions{
		OutputFormat: opt.OutputFormat,
		Width:        info.width,
		Height:       info.height,
		Font:         opt.FontStyle.Font,
	})
	if err != nil {
		return nil, err
	}
	if _, err = NewTableChart(p, opt).renderWithInfo(info); err != nil {
		return nil, err
	}
	return p, nil
}

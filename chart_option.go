package charts

import (
	"sort"

	"github.com/golang/freetype/truetype"
)

type ChartOption struct {
	// The font to use for rendering the chart
	Font *truetype.Font
	// The output type of chart, "svg" or "png", default value is "svg"
	Type string
	// The theme of chart.  Built in themes can be loaded using `GetTheme` with
	//"light", "dark", "vivid-light", "vivid-dark", "ant" or "grafana".
	Theme ColorPalette
	// The title option
	Title TitleOption
	// The legend option
	Legend LegendOption
	// The x-axis options
	XAxis XAxisOption
	// The y-axis option list
	YAxisOptions []YAxisOption
	// The width of chart, default width is 600
	Width int
	// The height of chart, default height is 400
	Height int
	Parent *Painter
	// The padding for chart, default padding is [20, 10, 10, 10]
	Padding Box
	// The canvas box for chart
	Box Box
	// The series list
	SeriesList SeriesList
	// The radar indicator list
	RadarIndicators []RadarIndicator
	// The background color of chart
	BackgroundColor Color
	// The flag for show symbol of line, set this to *false will hide symbol
	SymbolShow *bool
	// The stroke width of line chart
	LineStrokeWidth float64
	// The bar with of bar chart
	BarWidth int
	// The bar height of horizontal bar chart
	BarHeight int
	// Fill in the area under the line
	FillArea bool
	// background fill (alpha) opacity
	Opacity uint8
	// The child charts
	Children []ChartOption
	// The value formatter
	ValueFormatter ValueFormatter
}

// OptionFunc option function
type OptionFunc func(opt *ChartOption)

// SVGTypeOption set svg type of chart's output
func SVGTypeOption() OptionFunc {
	return TypeOptionFunc(ChartOutputSVG)
}

// PNGTypeOption set png type of chart's output
func PNGTypeOption() OptionFunc {
	return TypeOptionFunc(ChartOutputPNG)
}

// TypeOptionFunc set type of chart's output
func TypeOptionFunc(t string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Type = t
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
func LegendLabelsOptionFunc(labels []string, left ...string) OptionFunc {
	return func(opt *ChartOption) {
		opt.Legend = NewLegendOption(labels, left...)
	}
}

// XAxisOptionFunc set x-axis of chart
func XAxisOptionFunc(xAxisOption XAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = xAxisOption
	}
}

// XAxisDataOptionFunc set x-axis data of chart
func XAxisDataOptionFunc(data []string, boundaryGap ...*bool) OptionFunc {
	return func(opt *ChartOption) {
		opt.XAxis = NewXAxisOption(data, boundaryGap...)
	}
}

// YAxisOptionFunc set y-axis of chart, supports up to two y-axis.
func YAxisOptionFunc(yAxisOption ...YAxisOption) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxisOptions = yAxisOption
	}
}

// YAxisDataOptionFunc set y-axis data of chart
func YAxisDataOptionFunc(data []string) OptionFunc {
	return func(opt *ChartOption) {
		opt.YAxisOptions = NewYAxisOptions(data)
	}
}

// WidthOptionFunc set width of chart
func WidthOptionFunc(width int) OptionFunc {
	return func(opt *ChartOption) {
		opt.Width = width
	}
}

// HeightOptionFunc set height of chart
func HeightOptionFunc(height int) OptionFunc {
	return func(opt *ChartOption) {
		opt.Height = height
	}
}

// PaddingOptionFunc set padding of chart
func PaddingOptionFunc(padding Box) OptionFunc {
	return func(opt *ChartOption) {
		opt.Padding = padding
	}
}

// BoxOptionFunc set box of chart
func BoxOptionFunc(box Box) OptionFunc {
	return func(opt *ChartOption) {
		opt.Box = box
	}
}

// PieSeriesShowLabel set pie series show label
func PieSeriesShowLabel() OptionFunc {
	return func(opt *ChartOption) {
		for index := range opt.SeriesList {
			opt.SeriesList[index].Label.Show = true
		}
	}
}

// ChildOptionFunc add child chart
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

// BackgroundColorOptionFunc set background color of chart
func BackgroundColorOptionFunc(color Color) OptionFunc {
	return func(opt *ChartOption) {
		opt.BackgroundColor = color
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

func (o *ChartOption) fillDefault() {
	axisCount := 1
	for _, series := range o.SeriesList {
		if series.AxisIndex >= axisCount {
			axisCount++
		}
	}
	o.Width = getDefaultInt(o.Width, defaultChartWidth)
	o.Height = getDefaultInt(o.Height, defaultChartHeight)

	yAxisOptions := make([]YAxisOption, axisCount)
	copy(yAxisOptions, o.YAxisOptions)
	o.YAxisOptions = yAxisOptions
	// TODO - this is a hack, we need to update the yaxis based on the markpoint state
	// TODO - but can't do this earlier due to needing the axis initialized
	// TODO - we should reconsider the API for configuration
	hasMarkpoint := false
	for _, sl := range o.SeriesList {
		if len(sl.MarkPoint.Data) > 0 {
			hasMarkpoint = true
			break
		}
	}
	if hasMarkpoint {
		for i := range o.YAxisOptions {
			if o.YAxisOptions[i].RangeValuePaddingScale == nil {
				defaultPadding := 2.5 // default a larger padding to give space for the mark point
				o.YAxisOptions[i].RangeValuePaddingScale = &defaultPadding
			}
		}
	}

	if o.Font == nil {
		o.Font = GetDefaultFont()
	}
	if o.Theme == nil {
		o.Theme = GetDefaultTheme()
	}

	if o.BackgroundColor.IsZero() {
		o.BackgroundColor = o.Theme.GetBackgroundColor()
	}
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
			if index < seriesCount &&
				len(o.SeriesList[index].Name) == 0 {
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
}

// LineRender line chart render
func LineRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	seriesList := NewSeriesListDataFromValues(values, ChartTypeLine)
	return Render(ChartOption{
		SeriesList: seriesList,
	}, opts...)
}

// BarRender bar chart render
func BarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	seriesList := NewSeriesListDataFromValues(values, ChartTypeBar)
	return Render(ChartOption{
		SeriesList: seriesList,
	}, opts...)
}

// HorizontalBarRender horizontal bar chart render
func HorizontalBarRender(values [][]float64, opts ...OptionFunc) (*Painter, error) {
	seriesList := NewSeriesListDataFromValues(values, ChartTypeHorizontalBar)
	return Render(ChartOption{
		SeriesList: seriesList,
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
	seriesList := NewSeriesListDataFromValues(values, ChartTypeRadar)
	return Render(ChartOption{
		SeriesList: seriesList,
	}, opts...)
}

// FunnelRender funnel chart render
func FunnelRender(values []float64, opts ...OptionFunc) (*Painter, error) {
	seriesList := NewFunnelSeriesList(values)
	return Render(ChartOption{
		SeriesList: seriesList,
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
	if opt.Type == "" {
		opt.Type = ChartOutputPNG
	}
	if opt.Width <= 0 {
		opt.Width = defaultChartWidth
	}

	p, err := NewPainter(PainterOptions{
		Type:  opt.Type,
		Width: opt.Width,
		// is only used to calculate the hight of the table
		Height: 100,
		Font:   opt.Font,
	})
	if err != nil {
		return nil, err
	}
	info, err := NewTableChart(p, opt).render()
	if err != nil {
		return nil, err
	}

	p, err = NewPainter(PainterOptions{
		Type:   opt.Type,
		Width:  info.Width,
		Height: info.Height,
		Font:   opt.Font,
	})
	if err != nil {
		return nil, err
	}
	if _, err = NewTableChart(p, opt).renderWithInfo(info); err != nil {
		return nil, err
	}
	return p, nil
}

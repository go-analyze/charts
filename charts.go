package charts

import (
	"errors"
	"math"
	"sort"

	"github.com/wcharczuk/go-chart/v2"
)

const labelFontSize = 10
const smallLabelFontSize = 8
const defaultDotWidth = 2.0
const defaultStrokeWidth = 2.0

var defaultChartWidth = 600
var defaultChartHeight = 400

// SetDefaultWidth sets default width of chart
func SetDefaultWidth(width int) {
	if width > 0 {
		defaultChartWidth = width
	}
}

// SetDefaultHeight sets default height of chart
func SetDefaultHeight(height int) {
	if height > 0 {
		defaultChartHeight = height
	}
}

var nullValue = math.MaxFloat64

// SetNullValue sets the null value, default is MaxFloat64
func SetNullValue(v float64) {
	nullValue = v
}

// GetNullValue gets the null value
func GetNullValue() float64 {
	return nullValue
}

type Renderer interface {
	Render() (Box, error)
}

type renderHandler struct {
	list []func() error
}

func (rh *renderHandler) Add(fn func() error) {
	list := rh.list
	if len(list) == 0 {
		list = make([]func() error, 0)
	}
	rh.list = append(list, fn)
}

func (rh *renderHandler) Do() error {
	for _, fn := range rh.list {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

type defaultRenderOption struct {
	Theme      ColorPalette
	Padding    Box
	SeriesList SeriesList
	// The y axis option
	YAxisOptions []YAxisOption
	// The x axis option
	XAxis XAxisOption
	// The title option
	TitleOption TitleOption
	// The legend option
	LegendOption LegendOption
	// background is filled
	backgroundIsFilled bool
	// x y axis is reversed
	axisReversed bool
}

type defaultRenderResult struct {
	axisRanges map[int]axisRange
	// legend area
	seriesPainter *Painter
}

func defaultRender(p *Painter, opt defaultRenderOption) (*defaultRenderResult, error) {
	seriesList := opt.SeriesList
	seriesList.init()
	if !opt.backgroundIsFilled {
		p.SetBackground(p.Width(), p.Height(), opt.Theme.GetBackgroundColor())
	}

	if !opt.Padding.IsZero() {
		p = p.Child(PainterPaddingOption(opt.Padding))
	}

	legendHeight := 0
	if len(opt.LegendOption.Data) != 0 {
		if opt.LegendOption.Theme == nil {
			opt.LegendOption.Theme = opt.Theme
		}
		legendResult, err := NewLegendPainter(p, opt.LegendOption).Render()
		if err != nil {
			return nil, err
		}
		legendHeight = legendResult.Height()
	}

	if opt.TitleOption.Text != "" {
		if opt.TitleOption.Theme == nil {
			opt.TitleOption.Theme = opt.Theme
		}
		titlePainter := NewTitlePainter(p, opt.TitleOption)

		titleBox, err := titlePainter.Render()
		if err != nil {
			return nil, err
		}

		top := chart.MaxInt(legendHeight, titleBox.Height())
		// if in verticle mode, the legend height is not calculated
		if opt.LegendOption.Orient == OrientVertical {
			top = titleBox.Height()
		}
		p = p.Child(PainterPaddingOption(Box{
			// leave space under the title
			Top: top + 20,
		}))
	}

	result := defaultRenderResult{
		axisRanges: make(map[int]axisRange),
	}

	axisIndexList := make([]int, 0)
	for _, series := range opt.SeriesList {
		if containsInt(axisIndexList, series.AxisIndex) {
			continue
		}
		axisIndexList = append(axisIndexList, series.AxisIndex)
	}
	// the height needs to be subtracted from the height of the x-axis
	rangeHeight := p.Height() - defaultXAxisHeight
	rangeWidthLeft := 0
	rangeWidthRight := 0

	sort.Sort(sort.Reverse(sort.IntSlice(axisIndexList)))

	// calculate the axis range
	for _, index := range axisIndexList {
		yAxisOption := YAxisOption{}
		if len(opt.YAxisOptions) > index {
			yAxisOption = opt.YAxisOptions[index]
		}
		divideCount := yAxisOption.DivideCount
		if divideCount <= 0 {
			divideCount = defaultAxisDivideCount
		}
		max, min := opt.SeriesList.GetMaxMin(index)
		r := NewRange(AxisRangeOption{
			Painter: p,
			Min:     min,
			Max:     max,
			// the height needs to be subtracted from the height of the x-axis
			Size: rangeHeight,
			// serparate quantity
			DivideCount: divideCount,
		})
		if yAxisOption.Min != nil && *yAxisOption.Min <= min {
			r.min = *yAxisOption.Min
		}
		if yAxisOption.Max != nil && *yAxisOption.Max >= max {
			r.max = *yAxisOption.Max
		}
		result.axisRanges[index] = r

		if yAxisOption.Theme == nil {
			yAxisOption.Theme = opt.Theme
		}
		if !opt.axisReversed {
			yAxisOption.Data = r.Values()
		} else {
			yAxisOption.isCategoryAxis = true
			// since the x-axis is the value part, it's label is calculated and processed seperately
			opt.XAxis.Data = NewRange(AxisRangeOption{
				Painter: p,
				Min:     min,
				Max:     max,
				// the height needs to be subtracted from the height of the x-axis
				Size: rangeHeight,
				// seperate quantities
				DivideCount: defaultAxisDivideCount,
			}).Values()
			opt.XAxis.isValueAxis = true
		}
		reverseStringSlice(yAxisOption.Data)
		// TODO - generate other positions and y-axis
		var yAxis *axisPainter
		child := p.Child(PainterPaddingOption(Box{
			Left:  rangeWidthLeft,
			Right: rangeWidthRight,
		}))
		if index == 0 {
			yAxis = NewLeftYAxis(child, yAxisOption)
		} else {
			yAxis = NewRightYAxis(child, yAxisOption)
		}
		yAxisBox, err := yAxis.Render()
		if err != nil {
			return nil, err
		}
		if index == 0 {
			rangeWidthLeft += yAxisBox.Width()
		} else {
			rangeWidthRight += yAxisBox.Width()
		}
	}

	if opt.XAxis.Theme == nil {
		opt.XAxis.Theme = opt.Theme
	}
	xAxis := NewBottomXAxis(p.Child(PainterPaddingOption(Box{
		Left:  rangeWidthLeft,
		Right: rangeWidthRight,
	})), opt.XAxis)
	if _, err := xAxis.Render(); err != nil {
		return nil, err
	}

	result.seriesPainter = p.Child(PainterPaddingOption(Box{
		Bottom: defaultXAxisHeight,
		Left:   rangeWidthLeft,
		Right:  rangeWidthRight,
	}))
	return &result, nil
}

func doRender(renderers ...Renderer) error {
	for _, r := range renderers {
		if _, err := r.Render(); err != nil {
			return err
		}
	}
	return nil
}

func Render(opt ChartOption, opts ...OptionFunc) (*Painter, error) {
	for _, fn := range opts {
		fn(&opt)
	}
	opt.fillDefault()

	isChild := true
	if opt.Parent == nil {
		isChild = false
		p, err := NewPainter(PainterOptions{
			Type:   opt.Type,
			Width:  opt.Width,
			Height: opt.Height,
			Font:   opt.font,
		})
		if err != nil {
			return nil, err
		}
		opt.Parent = p
	}
	p := opt.Parent
	if opt.ValueFormatter != nil {
		p.valueFormatter = opt.ValueFormatter
	}
	if !opt.Box.IsZero() {
		p = p.Child(PainterBoxOption(opt.Box))
	}
	if !isChild {
		p.SetBackground(p.Width(), p.Height(), opt.BackgroundColor)
	}
	seriesList := opt.SeriesList
	seriesList.init()

	seriesCount := len(seriesList)

	// line chart
	lineSeriesList := seriesList.Filter(ChartTypeLine)
	barSeriesList := seriesList.Filter(ChartTypeBar)
	horizontalBarSeriesList := seriesList.Filter(ChartTypeHorizontalBar)
	pieSeriesList := seriesList.Filter(ChartTypePie)
	radarSeriesList := seriesList.Filter(ChartTypeRadar)
	funnelSeriesList := seriesList.Filter(ChartTypeFunnel)

	if len(horizontalBarSeriesList) != 0 && len(horizontalBarSeriesList) != seriesCount {
		return nil, errors.New("Horizontal bar can not mix other charts")
	}
	if len(pieSeriesList) != 0 && len(pieSeriesList) != seriesCount {
		return nil, errors.New("Pie can not mix other charts")
	}
	if len(radarSeriesList) != 0 && len(radarSeriesList) != seriesCount {
		return nil, errors.New("Radar can not mix other charts")
	}
	if len(funnelSeriesList) != 0 && len(funnelSeriesList) != seriesCount {
		return nil, errors.New("Funnel can not mix other charts")
	}

	axisReversed := len(horizontalBarSeriesList) != 0
	renderOpt := defaultRenderOption{
		Theme:        opt.theme,
		Padding:      opt.Padding,
		SeriesList:   opt.SeriesList,
		XAxis:        opt.XAxis,
		YAxisOptions: opt.YAxisOptions,
		TitleOption:  opt.Title,
		LegendOption: opt.Legend,
		axisReversed: axisReversed,
		// the background color has been set
		backgroundIsFilled: true,
	}
	if len(pieSeriesList) != 0 ||
		len(radarSeriesList) != 0 ||
		len(funnelSeriesList) != 0 {
		renderOpt.XAxis.Show = FalseFlag()
		renderOpt.YAxisOptions = []YAxisOption{
			{
				Show: FalseFlag(),
			},
		}
	}
	if len(horizontalBarSeriesList) != 0 {
		renderOpt.YAxisOptions[0].DivideCount = len(renderOpt.YAxisOptions[0].Data)
		renderOpt.YAxisOptions[0].Unit = 1
	}

	renderResult, err := defaultRender(p, renderOpt)
	if err != nil {
		return nil, err
	}

	handler := renderHandler{}

	// bar chart
	if len(barSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewBarChart(p, BarChartOption{
				Theme:    opt.theme,
				Font:     opt.font,
				XAxis:    opt.XAxis,
				BarWidth: opt.BarWidth,
			}).render(renderResult, barSeriesList)
			return err
		})
	}

	// horizontal bar chart
	if len(horizontalBarSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewHorizontalBarChart(p, HorizontalBarChartOption{
				Theme:        opt.theme,
				Font:         opt.font,
				BarHeight:    opt.BarHeight,
				YAxisOptions: opt.YAxisOptions,
			}).render(renderResult, horizontalBarSeriesList)
			return err
		})
	}

	// pie chart
	if len(pieSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewPieChart(p, PieChartOption{
				Theme: opt.theme,
				Font:  opt.font,
			}).render(renderResult, pieSeriesList)
			return err
		})
	}

	// line chart
	if len(lineSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewLineChart(p, LineChartOption{
				Theme:       opt.theme,
				Font:        opt.font,
				XAxis:       opt.XAxis,
				SymbolShow:  opt.SymbolShow,
				StrokeWidth: opt.LineStrokeWidth,
				FillArea:    opt.FillArea,
				Opacity:     opt.Opacity,
			}).render(renderResult, lineSeriesList)
			return err
		})
	}

	// radar chart
	if len(radarSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewRadarChart(p, RadarChartOption{
				Theme: opt.theme,
				Font:  opt.font,
				// 相应值
				RadarIndicators: opt.RadarIndicators,
			}).render(renderResult, radarSeriesList)
			return err
		})
	}

	// funnel chart
	if len(funnelSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewFunnelChart(p, FunnelChartOption{
				Theme: opt.theme,
				Font:  opt.font,
			}).render(renderResult, funnelSeriesList)
			return err
		})
	}

	if err = handler.Do(); err != nil {
		return nil, err
	}

	for _, item := range opt.Children {
		item.Parent = p
		if item.Theme == "" {
			item.Theme = opt.Theme
		}
		if item.FontFamily == "" {
			item.FontFamily = opt.FontFamily
		}
		if _, err = Render(item); err != nil {
			return nil, err
		}
	}

	return p, nil
}

package charts

import (
	"errors"
	"math"
	"sort"

	"github.com/go-analyze/charts/chartdraw"
)

const labelFontSize = 10.0
const smallLabelFontSize = 8
const defaultDotWidth = 2.0
const defaultStrokeWidth = 2.0
const defaultYAxisLabelCountHigh = 10
const defaultYAxisLabelCountLow = 3

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

func defaultYAxisLabelCount(span float64, decimalData bool) int {
	result := math.Min(math.Max(span+1, defaultYAxisLabelCountLow), defaultYAxisLabelCountHigh)
	if decimalData {
		// if there is a decimal, we double our labels to provide more detail
		result = math.Min(result*2, defaultYAxisLabelCountHigh)
	}
	return int(result)
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
	// Theme specifies the colors used for the chart.
	Theme ColorPalette
	// Padding specifies the padding of chart.
	Padding Box
	// SeriesList provides the data series.
	SeriesList SeriesList
	// XAxis are options for the x-axis.
	XAxis XAxisOption
	// YAxis are options for the y-axis (at most two).
	YAxis []YAxisOption
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// backgroundIsFilled is true if the background is filled.
	backgroundIsFilled bool
	// axisReversed is true if the x y axis is reversed.
	axisReversed bool
}

type defaultRenderResult struct {
	axisRanges map[int]axisRange
	// legend area
	seriesPainter *Painter
}

func defaultRender(p *Painter, opt defaultRenderOption) (*defaultRenderResult, error) {
	fillThemeDefaults(getPreferredTheme(opt.Theme), &opt.Title, &opt.Legend, &opt.XAxis)

	seriesList := opt.SeriesList
	seriesList.init()
	if !opt.backgroundIsFilled {
		p.SetBackground(p.Width(), p.Height(), opt.Theme.GetBackgroundColor())
	}

	if !opt.Padding.IsZero() {
		p = p.Child(PainterPaddingOption(opt.Padding))
	}

	legendHeight := 0
	if len(opt.Legend.Data) != 0 {
		legendResult, err := NewLegendPainter(p, opt.Legend).Render()
		if err != nil {
			return nil, err
		}
		legendHeight = legendResult.Height()
	}

	if opt.Title.Text != "" {
		titleBox, err := NewTitlePainter(p, opt.Title).Render()
		if err != nil {
			return nil, err
		}

		top := chartdraw.MaxInt(legendHeight, titleBox.Height())
		// if in vertical mode, the legend height is not calculated
		if opt.Legend.Orient == OrientVertical {
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
		if containsInt(axisIndexList, series.YAxisIndex) {
			continue
		}
		axisIndexList = append(axisIndexList, series.YAxisIndex)
	}
	// the height needs to be subtracted from the height of the x-axis
	rangeHeight := p.Height() - defaultXAxisHeight
	rangeWidthLeft := 0
	rangeWidthRight := 0

	sort.Sort(sort.Reverse(sort.IntSlice(axisIndexList)))

	// calculate the axis range
	for _, index := range axisIndexList {
		yAxisOption := YAxisOption{}
		if len(opt.YAxis) > index {
			yAxisOption = opt.YAxis[index]
		}
		minPadRange, maxPadRange := 1.0, 1.0
		if yAxisOption.RangeValuePaddingScale != nil {
			minPadRange = *yAxisOption.RangeValuePaddingScale
			maxPadRange = *yAxisOption.RangeValuePaddingScale
		}
		min, max := opt.SeriesList.GetMinMax(index)
		decimalData := min != math.Floor(min) || (max-min) != math.Floor(max-min)
		if yAxisOption.Min != nil && *yAxisOption.Min < min {
			min = *yAxisOption.Min
			minPadRange = 0.0
		}
		if yAxisOption.Max != nil && *yAxisOption.Max > max {
			max = *yAxisOption.Max
			maxPadRange = 0.0
		}

		// Label counts and y-axis padding are linked together to produce a user-friendly graph.
		// First when considering padding we want to prefer a zero axis start if reasonable, and add a slight
		// padding to the max so there is a little space at the top of the graph.  In addition, we want to pick
		// a max value that will result in round intervals on the axis.  These details are in range.go.
		// But in order to produce round intervals we need to have an idea of how many intervals there are.
		// In addition, if the user specified a `Unit` value we may need to adjust our label count calculation
		// based on the padded range.
		//
		// In order to accomplish this, we estimate the label count (if necessary), pad the range, then precisely
		// calculate the label count.
		// TODO - label counts are also calculated in axis.go, for the X axis, ideally we unify these implementations
		labelCount := yAxisOption.LabelCount
		padLabelCount := labelCount
		if padLabelCount < 1 {
			if yAxisOption.Unit > 0 {
				padLabelCount = int((max-min)/yAxisOption.Unit) + 1
			} else {
				padLabelCount = defaultYAxisLabelCount(max-min, decimalData)
			}
		}
		padLabelCount = chartdraw.MaxInt(padLabelCount+yAxisOption.LabelCountAdjustment, 2)
		// we call padRange directly because we need to do this padding before we can calculate the final labelCount for the axisRange
		min, max = padRange(padLabelCount, min, max, minPadRange, maxPadRange)
		if labelCount <= 0 {
			if yAxisOption.Unit > 0 {
				if yAxisOption.Max == nil {
					max = math.Trunc(math.Ceil(max/yAxisOption.Unit) * yAxisOption.Unit)
				}
				labelCount = int((max-min)/yAxisOption.Unit) + 1
			} else {
				labelCount = defaultYAxisLabelCount(max-min, decimalData)
			}
			yAxisOption.LabelCount = labelCount
		}
		labelCount = chartdraw.MaxInt(labelCount+yAxisOption.LabelCountAdjustment, 2)
		r := axisRange{
			p:           p,
			divideCount: labelCount,
			min:         min,
			max:         max,
			size:        rangeHeight,
		}
		result.axisRanges[index] = r

		if yAxisOption.Theme == nil {
			yAxisOption.Theme = opt.Theme
		}
		if !opt.axisReversed {
			yAxisOption.Data = r.Values()
		} else {
			yAxisOption.isCategoryAxis = true
			// we need to update the range labels or the bars wont be aligned to the Y axis
			r.divideCount = len(seriesList[0].Data)
			result.axisRanges[index] = r
			// since the x-axis is the value part, it's label is calculated and processed separately
			opt.XAxis.Data = r.Values()
			opt.XAxis.isValueAxis = true
		}
		reverseStringSlice(yAxisOption.Data)
		child := p.Child(PainterPaddingOption(Box{
			Left:  rangeWidthLeft,
			Right: rangeWidthRight,
		}))
		var yAxis *axisPainter
		if index == 0 {
			yAxis = NewLeftYAxis(child, yAxisOption)
		} else {
			yAxis = NewRightYAxis(child, yAxisOption)
		}
		if yAxisBox, err := yAxis.Render(); err != nil {
			return nil, err
		} else if index == 0 {
			rangeWidthLeft += yAxisBox.Width()
		} else {
			rangeWidthRight += yAxisBox.Width()
		}
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
	if err := opt.fillDefault(); err != nil {
		return nil, err
	}

	isChild := opt.parent != nil
	if !isChild {
		p, err := NewPainter(PainterOptions{
			OutputFormat: opt.OutputFormat,
			Width:        opt.Width,
			Height:       opt.Height,
			Font:         opt.Font,
		})
		if err != nil {
			return nil, err
		}
		opt.parent = p
	}
	p := opt.parent
	if opt.ValueFormatter != nil {
		p.valueFormatter = opt.ValueFormatter
	}
	if !opt.Box.IsZero() {
		p = p.Child(PainterBoxOption(opt.Box))
	}
	if !isChild {
		p.SetBackground(p.Width(), p.Height(), opt.Theme.GetBackgroundColor())
	}
	seriesList := opt.SeriesList
	seriesList.init()

	lineSeriesList := seriesList.Filter(ChartTypeLine)
	barSeriesList := seriesList.Filter(ChartTypeBar)
	horizontalBarSeriesList := seriesList.Filter(ChartTypeHorizontalBar)
	pieSeriesList := seriesList.Filter(ChartTypePie)
	radarSeriesList := seriesList.Filter(ChartTypeRadar)
	funnelSeriesList := seriesList.Filter(ChartTypeFunnel)

	seriesCount := len(seriesList)
	if len(horizontalBarSeriesList) != 0 && len(horizontalBarSeriesList) != seriesCount {
		return nil, errors.New("horizontal bar can not mix other charts")
	} else if len(pieSeriesList) != 0 && len(pieSeriesList) != seriesCount {
		return nil, errors.New("pie can not mix other charts")
	} else if len(radarSeriesList) != 0 && len(radarSeriesList) != seriesCount {
		return nil, errors.New("radar can not mix other charts")
	} else if len(funnelSeriesList) != 0 && len(funnelSeriesList) != seriesCount {
		return nil, errors.New("funnel can not mix other charts")
	}

	axisReversed := len(horizontalBarSeriesList) != 0
	renderOpt := defaultRenderOption{
		Theme:        opt.Theme,
		Padding:      opt.Padding,
		SeriesList:   opt.SeriesList,
		XAxis:        opt.XAxis,
		YAxis:        opt.YAxis,
		Title:        opt.Title,
		Legend:       opt.Legend,
		axisReversed: axisReversed,
		// the background color has been set
		backgroundIsFilled: true,
	}
	if len(pieSeriesList) != 0 ||
		len(radarSeriesList) != 0 ||
		len(funnelSeriesList) != 0 {
		renderOpt.XAxis.Show = False()
		renderOpt.YAxis = []YAxisOption{
			{
				Show: False(),
			},
		}
	}
	if len(horizontalBarSeriesList) != 0 {
		renderOpt.YAxis[0].Unit = 1
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
				Theme:    opt.Theme,
				Font:     opt.Font,
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
				Theme:     opt.Theme,
				Font:      opt.Font,
				BarHeight: opt.BarHeight,
				YAxis:     opt.YAxis,
			}).render(renderResult, horizontalBarSeriesList)
			return err
		})
	}

	// pie chart
	if len(pieSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewPieChart(p, PieChartOption{
				Theme: opt.Theme,
				Font:  opt.Font,
			}).render(renderResult, pieSeriesList)
			return err
		})
	}

	// line chart
	if len(lineSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewLineChart(p, LineChartOption{
				Theme:       opt.Theme,
				Font:        opt.Font,
				XAxis:       opt.XAxis,
				SymbolShow:  opt.SymbolShow,
				StrokeWidth: opt.LineStrokeWidth,
				FillArea:    opt.FillArea,
				FillOpacity: opt.FillOpacity,
			}).render(renderResult, lineSeriesList)
			return err
		})
	}

	// radar chart
	if len(radarSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewRadarChart(p, RadarChartOption{
				Theme:           opt.Theme,
				Font:            opt.Font,
				RadarIndicators: opt.RadarIndicators,
			}).render(renderResult, radarSeriesList)
			return err
		})
	}

	// funnel chart
	if len(funnelSeriesList) != 0 {
		handler.Add(func() error {
			_, err := NewFunnelChart(p, FunnelChartOption{
				Theme: opt.Theme,
				Font:  opt.Font,
			}).render(renderResult, funnelSeriesList)
			return err
		})
	}

	if err = handler.Do(); err != nil {
		return nil, err
	}

	for _, item := range opt.Children {
		item.parent = p
		if item.Theme == nil {
			item.Theme = opt.Theme
		}
		if item.Font == nil {
			item.Font = opt.Font
		}
		if _, err = Render(item); err != nil {
			return nil, err
		}
	}

	return p, nil
}

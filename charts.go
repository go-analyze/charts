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
var defaultPadding = NewBoxEqual(20)

// SetDefaultChartDimensions sets default width and height of charts if not otherwise specified in their configuration.
func SetDefaultChartDimensions(width, height int) {
	if width > 0 {
		defaultChartWidth = width
	}
	if height > 0 {
		defaultChartHeight = height
	}
}

// GetNullValue gets the null value, allowing you to set a series point with "no" value.
func GetNullValue() float64 {
	return math.MaxFloat64
}

func defaultYAxisLabelCount(span float64, decimalData bool) int {
	result := math.Min(math.Max(span+1, defaultYAxisLabelCountLow), defaultYAxisLabelCountHigh)
	if decimalData {
		// if there is a decimal, we double our labels to provide more detail
		result = math.Min(result*2, defaultYAxisLabelCountHigh)
	}
	return int(result)
}

type renderer interface {
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
	// theme specifies the colors used for the chart.
	theme ColorPalette
	// padding specifies the padding of chart.
	padding Box
	// seriesList provides the data series.
	seriesList SeriesList
	// stackSeries can be set to true if the series data will be stacked (summed).
	stackSeries bool
	// xAxis are options for the x-axis.
	xAxis *XAxisOption
	// yAxis are options for the y-axis (at most two).
	yAxis []YAxisOption
	// title are options for rendering the title.
	title TitleOption
	// legend are options for the data legend.
	legend *LegendOption
	// backgroundIsFilled is true if the background is filled.
	backgroundIsFilled bool
	// axisReversed is true if the x y-axis is reversed.
	axisReversed bool
	// valueFormatter to format numeric values into labels.
	valueFormatter ValueFormatter
}

type defaultRenderResult struct {
	axisRanges map[int]axisRange
	// legend area
	seriesPainter *Painter
}

func defaultRender(p *Painter, opt defaultRenderOption) (*defaultRenderResult, error) {
	fillThemeDefaults(getPreferredTheme(opt.theme, p.theme), &opt.title, opt.legend, opt.xAxis)

	// TODO - this is a hack, we need to update the yaxis based on the markpoint state
	for _, sl := range opt.seriesList {
		if len(sl.MarkPoint.Data) > 0 { // if graph has markpoint
			// adjust padding scale to give space for mark point (if not specified by user)
			for i := range opt.yAxis {
				if opt.yAxis[i].RangeValuePaddingScale == nil {
					opt.yAxis[i].RangeValuePaddingScale = FloatPointer(2.5)
				}
			}
			break
		}
	}

	if !opt.backgroundIsFilled {
		p.drawBackground(opt.theme.GetBackgroundColor())
	}
	if !opt.padding.IsZero() {
		p = p.Child(PainterPaddingOption(opt.padding))
	}

	// association between legend and series name
	if len(opt.legend.Data) == 0 {
		opt.legend.Data = opt.seriesList.Names()
	} else {
		seriesCount := len(opt.seriesList)
		for index, name := range opt.legend.Data {
			if index < seriesCount && len(opt.seriesList[index].Name) == 0 {
				opt.seriesList[index].Name = name
			}
		}
		nameIndexDict := map[string]int{}
		for index, name := range opt.legend.Data {
			nameIndexDict[name] = index
		}
		// ensure order of series is consistent with legend
		sort.Slice(opt.seriesList, func(i, j int) bool {
			return nameIndexDict[opt.seriesList[i].Name] < nameIndexDict[opt.seriesList[j].Name]
		})
	}

	const legendTitlePadding = 15
	legendTopSpacing := 0
	legendResult, err := newLegendPainter(p, *opt.legend).Render()
	if err != nil {
		return nil, err
	}
	if !legendResult.IsZero() && !flagIs(true, opt.legend.Vertical) && !flagIs(true, opt.legend.OverlayChart) {
		legendHeight := legendResult.Height()
		if legendResult.Bottom < p.Height()/2 {
			// horizontal legend at the top, set the spacing based on the height
			legendTopSpacing = legendHeight + legendTitlePadding
		} else {
			// horizontal legend at the bottom, raise the chart above it
			p = p.Child(PainterPaddingOption(Box{
				Bottom: legendHeight + legendTitlePadding,
				IsSet:  true,
			}))
		}
	}

	titleBox, err := newTitlePainter(p, opt.title).Render()
	if err != nil {
		return nil, err
	}
	if !titleBox.IsZero() {
		var top, bottom int
		if titleBox.Bottom < p.Height()/2 {
			top = chartdraw.MaxInt(legendTopSpacing, titleBox.Bottom+legendTitlePadding)
		} else {
			// title is at the bottom, raise the chart to be above the title
			bottom = titleBox.Height()
			top = legendTopSpacing // the legend may still need space on the top, set to whatever the legend requested
		}

		p = p.Child(PainterPaddingOption(Box{
			Top:    top,
			Bottom: bottom,
			IsSet:  true,
		}))
	} else if legendTopSpacing > 0 { // apply chart spacing below legend
		p = p.Child(PainterPaddingOption(Box{
			Top:   legendTopSpacing,
			IsSet: true,
		}))
	}

	result := defaultRenderResult{
		axisRanges: make(map[int]axisRange),
	}

	axisIndexList := make([]int, opt.seriesList.getYAxisCount())
	for i := range axisIndexList {
		axisIndexList[i] = i
	}
	reverseIntSlice(axisIndexList)
	// the height needs to be subtracted from the height of the x-axis
	rangeHeight := p.Height() - defaultXAxisHeight
	rangeWidthLeft := 0
	rangeWidthRight := 0

	// calculate the axis range
	for _, index := range axisIndexList {
		yAxisOption := YAxisOption{}
		if len(opt.yAxis) > index {
			yAxisOption = opt.yAxis[index]
		}
		minPadRange, maxPadRange := 1.0, 1.0
		if yAxisOption.RangeValuePaddingScale != nil {
			minPadRange = *yAxisOption.RangeValuePaddingScale
			maxPadRange = *yAxisOption.RangeValuePaddingScale
		}
		min, max, sumMax := opt.seriesList.getMinMaxSumMax(index, opt.stackSeries)
		decimalData := min != math.Floor(min) || (max-min) != math.Floor(max-min)
		if yAxisOption.Min != nil && *yAxisOption.Min < min {
			min = *yAxisOption.Min
			minPadRange = 0.0
		}
		if opt.stackSeries {
			// If stacking, max should be the highest sum
			max = sumMax
		}
		if yAxisOption.Max != nil && *yAxisOption.Max > max {
			max = *yAxisOption.Max
			maxPadRange = 0.0
		}

		// Label counts and y-axis padding are linked together to produce a user-friendly graph.
		// First when considering padding we want to prefer a zero axis start if reasonable, and add a slight
		// padding to the max so there is a little space at the top of the graph. In addition, we want to pick
		// a max value that will result in round intervals on the axis. These details are in range.go.
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
		r := newRange(p, getPreferredValueFormatter(yAxisOption.ValueFormatter, opt.valueFormatter),
			rangeHeight, labelCount, min, max, 0, 0)
		result.axisRanges[index] = r

		if yAxisOption.Theme == nil {
			yAxisOption.Theme = opt.theme
		}
		if !opt.axisReversed {
			yAxisOption.Data = r.Values()
		} else {
			yAxisOption.isCategoryAxis = true
			// we need to update the range labels or the bars won't be aligned to the Y axis
			r.divideCount = len(opt.seriesList[0].Data)
			result.axisRanges[index] = r
			// since the x-axis is the value part, it's label is calculated and processed separately
			opt.xAxis.Data = r.Values()
			opt.xAxis.isValueAxis = true
		}
		reverseStringSlice(yAxisOption.Data)
		child := p.Child(PainterPaddingOption(Box{
			Left:  rangeWidthLeft,
			Right: rangeWidthRight,
			IsSet: true,
		}))
		var yAxis *axisPainter
		if index == 0 {
			yAxis = newLeftYAxis(child, yAxisOption)
		} else {
			yAxis = newRightYAxis(child, yAxisOption)
		}
		if yAxisBox, err := yAxis.Render(); err != nil {
			return nil, err
		} else if index == 0 {
			rangeWidthLeft += yAxisBox.Width()
		} else {
			rangeWidthRight += yAxisBox.Width()
		}
	}

	xAxis := newBottomXAxis(p.Child(PainterPaddingOption(Box{
		Left:  rangeWidthLeft,
		Right: rangeWidthRight,
		IsSet: true,
	})), *opt.xAxis)
	if _, err := xAxis.Render(); err != nil {
		return nil, err
	}

	result.seriesPainter = p.Child(PainterPaddingOption(Box{
		Left:   rangeWidthLeft,
		Right:  rangeWidthRight,
		Bottom: defaultXAxisHeight,
		IsSet:  true,
	}))
	return &result, nil
}

func doRender(renderers ...renderer) error {
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
		opt.parent = NewPainter(PainterOptions{
			OutputFormat: opt.OutputFormat,
			Width:        opt.Width,
			Height:       opt.Height,
			Font:         opt.Font,
		})
	}
	p := opt.parent
	if !opt.Box.IsZero() {
		p = p.Child(PainterBoxOption(opt.Box))
	}
	if !isChild {
		p.drawBackground(opt.Theme.GetBackgroundColor())
	}

	seriesList := opt.SeriesList
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
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		xAxis:          &opt.XAxis,
		yAxis:          opt.YAxis,
		stackSeries:    flagIs(true, opt.StackSeries),
		title:          opt.Title,
		legend:         &opt.Legend,
		axisReversed:   axisReversed,
		valueFormatter: opt.ValueFormatter,
		// the background color has been set
		backgroundIsFilled: true,
	}
	if len(pieSeriesList) != 0 ||
		len(radarSeriesList) != 0 ||
		len(funnelSeriesList) != 0 {
		renderOpt.xAxis.Show = False()
		renderOpt.yAxis = []YAxisOption{
			{
				Show: False(),
			},
		}
	}
	if len(horizontalBarSeriesList) != 0 {
		renderOpt.yAxis[0].Unit = 1
	}

	renderResult, err := defaultRender(p, renderOpt)
	if err != nil {
		return nil, err
	}

	handler := renderHandler{}

	// bar chart
	if len(barSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newBarChart(p, BarChartOption{
				Theme:       opt.Theme,
				Font:        opt.Font,
				XAxis:       opt.XAxis,
				StackSeries: opt.StackSeries,
				BarWidth:    opt.BarWidth,
				BarMargin:   opt.BarMargin,
			}).render(renderResult, barSeriesList)
			return err
		})
	}

	// horizontal bar chart
	if len(horizontalBarSeriesList) != 0 {
		if len(opt.YAxis) > 0 {
			if len(opt.YAxis) > 1 {
				return nil, errors.New("horizontal bar chart only accepts a single Y-Axis")
			}
			// TODO - handle v0.5 API change for selecting the yAxis here
		}
		handler.Add(func() error {
			_, err := newHorizontalBarChart(p, HorizontalBarChartOption{
				Theme:       opt.Theme,
				Font:        opt.Font,
				BarHeight:   opt.BarHeight,
				BarMargin:   opt.BarMargin,
				YAxis:       opt.YAxis,
				StackSeries: opt.StackSeries,
			}).render(renderResult, horizontalBarSeriesList)
			return err
		})
	}

	// pie chart
	if len(pieSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newPieChart(p, PieChartOption{
				Theme: opt.Theme,
				Font:  opt.Font,
			}).render(renderResult, pieSeriesList)
			return err
		})
	}

	// line chart
	if len(lineSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newLineChart(p, LineChartOption{
				Theme:           opt.Theme,
				Font:            opt.Font,
				XAxis:           opt.XAxis,
				StackSeries:     opt.StackSeries,
				SymbolShow:      opt.SymbolShow,
				LineStrokeWidth: opt.LineStrokeWidth,
				FillArea:        opt.FillArea,
				FillOpacity:     opt.FillOpacity,
			}).render(renderResult, lineSeriesList)
			return err
		})
	}

	// radar chart
	if len(radarSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newRadarChart(p, RadarChartOption{
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
			_, err := newFunnelChart(p, FunnelChartOption{
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

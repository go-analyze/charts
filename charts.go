package charts

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/go-analyze/charts/chartdraw"
)

const defaultLabelFontSize = 10.0
const smallLabelFontSize = 8
const defaultDotWidth = 2.0
const defaultStrokeWidth = 2.0
const defaultYAxisLabelCountHigh = 10
const defaultYAxisLabelCountLow = 3

var defaultChartWidth = 600
var defaultChartHeight = 400
var defaultPadding = NewBoxEqual(20)

// SetDefaultChartDimensions sets the default chart width and height when not otherwise specified in their configuration.
func SetDefaultChartDimensions(width, height int) {
	if width > 0 {
		defaultChartWidth = width
	}
	if height > 0 {
		defaultChartHeight = height
	}
}

// GetNullValue returns the null value for setting series points with "no" or "unknown" value.
func GetNullValue() float64 {
	return math.MaxFloat64
}

type renderer interface {
	Render() (Box, error)
}

type renderHandler struct {
	list []func() error
}

func (rh *renderHandler) Add(fn func() error) {
	rh.list = append(rh.list, fn)
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
	// padding specifies the chart padding.
	padding Box
	// seriesList provides the data series.
	seriesList seriesList
	// stackSeries when true causes series data to be stacked (summed).
	stackSeries bool
	// xAxis contains options for the x-axis.
	xAxis *XAxisOption
	// yAxis contains options for the y-axis. At most two y-axes are supported.
	yAxis []YAxisOption
	// title contains options for rendering the chart title.
	title TitleOption
	// legend contains options for the data legend.
	legend *LegendOption
	// backgroundIsFilled is true if the background is filled.
	backgroundIsFilled bool
	// axisReversed is true if the x-axis and y-axis are reversed.
	axisReversed bool
	// valueFormatter formats numeric values into labels.
	valueFormatter ValueFormatter
}

type defaultRenderResult struct {
	yaxisRanges   map[int]axisRange
	xaxisRange    axisRange
	seriesPainter *Painter
}

func defaultRender(p *Painter, opt defaultRenderOption) (*defaultRenderResult, error) {
	theme := getPreferredTheme(opt.theme, p.theme)
	fillThemeDefaults(theme, &opt.title, opt.legend, opt.xAxis, opt.yAxis)
	opt.xAxis = opt.xAxis.prep(theme)
	// TODO - this is a hack, we need to update the yaxis based on the markpoint state
	if opt.seriesList.hasMarkPoint() {
		// adjust padding scale to give space for mark point (if not specified by user)
		for i := range opt.yAxis {
			if opt.yAxis[i].RangeValuePaddingScale == nil {
				opt.yAxis[i].RangeValuePaddingScale = Ptr(2.5)
			}
		}
	}
	top := p

	if !opt.backgroundIsFilled {
		p.drawBackground(opt.theme.GetBackgroundColor())
	}
	if !opt.padding.IsZero() {
		p = p.Child(PainterPaddingOption(opt.padding))
	}

	// association between legend and series name
	if len(opt.legend.SeriesNames) == 0 {
		opt.legend.SeriesNames = opt.seriesList.names()
	} else {
		seriesCount := opt.seriesList.len()
		setAllSeriesNames := true
		for index, name := range opt.legend.SeriesNames {
			if index >= seriesCount {
				setAllSeriesNames = false
				break
			} else if opt.seriesList.getSeriesName(index) == "" {
				opt.seriesList.setSeriesName(index, name)
			} else {
				setAllSeriesNames = false
			}
		}
		if !setAllSeriesNames {
			// if the series had names already set, make sure they are consistent ordered with the legend
			nameIndexDict := map[string]int{}
			for index, name := range opt.legend.SeriesNames {
				nameIndexDict[name] = index
			}
			opt.seriesList.sortByNameIndex(nameIndexDict)
		}
	}
	opt.legend.seriesSymbols = make([]Symbol, opt.seriesList.len())
	for index := range opt.legend.seriesSymbols {
		opt.legend.seriesSymbols[index] = opt.seriesList.getSeriesSymbol(index)
	}

	const legendTitlePadding = 15
	var legendTopSpacing int
	var titleBox, legendResult Box
	var err error

	// make a local copy of legend options before we modify position during collision handling
	legendOpt := *opt.legend

	// helper to check if legend can be repositioned to avoid title collision
	// repositioning is allowed when no explicit numeric offset is set
	// OverlayChart only controls legend/chart overlap, not legend/title overlap
	legendCanReposition := func() bool {
		return (legendOpt.Offset.Top == "" || legendOpt.Offset.Top == PositionTop ||
			legendOpt.Offset.Top == PositionBottom) &&
			(legendOpt.Offset.Left == "" || legendOpt.Offset.Left == PositionLeft ||
				legendOpt.Offset.Left == PositionCenter || legendOpt.Offset.Left == PositionRight)
	}

	// determine positioning
	titleAtBottom := opt.title.Offset.Top == PositionBottom
	legendAtBottom := legendOpt.Offset.Top == PositionBottom

	// calculate legend box for both-at-bottom check and collision detection
	legendPainter := newLegendPainter(p, legendOpt)
	legendResult, err = legendPainter.calculateBox()
	if err != nil {
		return nil, err
	}

	// when both title and legend are at bottom (without overlay), render title on a reduced canvas
	// so title appears above where legend will render (title above legend in reading order)
	titleCanvas := p
	adjustedForBottom := false
	if titleAtBottom && legendAtBottom && !legendResult.IsZero() &&
		!flagIs(true, legendOpt.OverlayChart) {
		titleCanvas = p.Child(PainterPaddingOption(Box{
			Bottom: legendResult.Height() + legendTitlePadding,
			IsSet:  true,
		}))
		adjustedForBottom = true
	}

	titleBox, err = newTitlePainter(titleCanvas, opt.title).Render()
	if err != nil {
		return nil, err
	}

	// check for collision and reposition if needed
	// skip if both-at-bottom was handled via canvas adjustment
	if !adjustedForBottom &&
		!titleBox.IsZero() && !legendResult.IsZero() &&
		legendCanReposition() && titleBox.Overlaps(legendResult) {
		legendOpt.Offset.Top = strconv.Itoa(titleBox.Bottom)
		legendPainter = newLegendPainter(p, legendOpt)
	}

	// render legend with final position
	legendResult, err = legendPainter.Render()
	if err != nil {
		return nil, err
	}

	// reserve space for the legend if not in overlay mode
	// - horizontal legends reserve space (top or bottom)
	// - vertical legends always overlay from the side
	// - skip for adjustedForBottom since title reservation handles combined space
	if !legendResult.IsZero() && !flagIs(true, legendOpt.OverlayChart) && !adjustedForBottom &&
		!flagIs(true, legendOpt.Vertical) {
		if legendResult.Bottom < p.Height()/2 {
			// horizontal legend at top - reserve top space
			legendTopSpacing = chartdraw.MaxInt(legendResult.Height(), legendResult.Bottom) + legendTitlePadding
		} else {
			// horizontal legend at the bottom, raise the chart above it
			p = p.Child(PainterPaddingOption(Box{
				Bottom: legendResult.Height() + legendTitlePadding,
				IsSet:  true,
			}))
		}
	}

	// apply title and legend spacing to chart area
	if !titleBox.IsZero() {
		titlePadBox := Box{IsSet: true}
		if titleBox.Bottom < p.Height()/2 {
			titlePadBox.Top = chartdraw.MaxInt(legendTopSpacing, titleBox.Bottom+legendTitlePadding)
		} else { // title is at the bottom, raise the chart to be above the title
			titlePadBox.Top = legendTopSpacing
			if adjustedForBottom {
				titlePadBox.Bottom = p.Height() - titleBox.Top
			} else {
				titlePadBox.Bottom = titleBox.Height()
			}
		}
		p = p.Child(PainterPaddingOption(titlePadBox))
	} else if legendTopSpacing > 0 { // apply chart spacing below legend
		p = p.Child(PainterPaddingOption(Box{
			Top:   legendTopSpacing,
			IsSet: true,
		}))
	}

	result := defaultRenderResult{
		yaxisRanges: make(map[int]axisRange),
	}

	// calculate x-axis range and do a dry-render to find height
	// we will render on the actual painter once we know the space the y-axis will occupy
	var xAxisOpts axisOption
	if opt.axisReversed { // X is value axis
		xAxisRange := calculateValueAxisRange(p, false, p.Width(),
			nil, nil, nil,
			opt.xAxis.Labels, opt.xAxis.DataStartIndex,
			opt.xAxis.LabelCount, opt.xAxis.Unit, opt.xAxis.LabelCountAdjustment,
			opt.seriesList, 0, opt.stackSeries,
			getPreferredValueFormatter(opt.xAxis.ValueFormatter, opt.valueFormatter),
			opt.xAxis.LabelRotation, opt.xAxis.LabelFontStyle,
			nil)
		xAxisOpts = opt.xAxis.toAxisOption(xAxisRange)
	} else { //  X is category axis
		xAxisRange := calculateCategoryAxisRange(p, p.Width(), false, flagIs(false, opt.xAxis.BoundaryGap),
			opt.xAxis.Labels, opt.xAxis.DataStartIndex,
			opt.xAxis.LabelCount, opt.xAxis.LabelCountAdjustment, opt.xAxis.Unit,
			opt.seriesList,
			opt.xAxis.LabelRotation, opt.xAxis.LabelFontStyle)
		xAxisOpts = opt.xAxis.toAxisOption(xAxisRange)
	}
	if top.Height() < 100 {
		xAxisOpts.minimumAxisHeight = 0 // don't reserve if chart is too small
	}
	xAxisBox, err := newAxisPainter(
		NewPainter(PainterOptions{
			OutputFormat: p.outputFormat,
			Width:        p.Width(),
			Height:       p.Height(),
			Theme:        p.theme,
			Font:         p.font,
		}),
		xAxisOpts,
	).Render()
	if err != nil {
		return nil, err
	}
	xAxisHeight := xAxisBox.Height()

	rangeHeight := p.Height() - xAxisHeight
	var rangeWidthLeft, rangeWidthRight int

	// prepare all y-axis options and range data (allowing multiple axes to be considered)
	yAxisCount := getSeriesYAxisCount(opt.seriesList)
	type yAxisEntry struct {
		option YAxisOption
		prep   *valueAxisPrep
		r      axisRange
	}
	var entries []yAxisEntry
	var valuePreps []*valueAxisPrep
	var valuePrepIndices []int
	if yAxisCount > 0 {
		entries = make([]yAxisEntry, yAxisCount)
		for yIndex := 0; yIndex < yAxisCount; yIndex++ {
			var yAxisOption YAxisOption
			if len(opt.yAxis) > yIndex {
				yAxisOption = opt.yAxis[yIndex]
			}
			yAxisOption = *yAxisOption.prep(getPreferredTheme(yAxisOption.Theme, theme))
			entries[yIndex].option = yAxisOption
			if opt.axisReversed { // Y is category axis
				entries[yIndex].r = calculateCategoryAxisRange(p, rangeHeight, true, false,
					yAxisOption.Labels, 0,
					yAxisOption.LabelCount, yAxisOption.LabelCountAdjustment, yAxisOption.Unit,
					opt.seriesList,
					yAxisOption.LabelRotation, yAxisOption.LabelFontStyle)
			} else { // Standard Y value axis
				floatFormatter := getPreferredValueFormatter(yAxisOption.ValueFormatter, opt.valueFormatter)
				valueFormatter := floatFormatter
				if yAxisOption.Formatter != "" {
					fmtStr := yAxisOption.Formatter
					valueFormatter = func(f float64) string {
						return strings.ReplaceAll(fmtStr, "{value}", floatFormatter(f))
					}
				}
				prep := prepareValueAxisRange(p, true, rangeHeight,
					yAxisOption.Min, yAxisOption.Max, yAxisOption.RangeValuePaddingScale,
					yAxisOption.Labels, 0,
					yAxisOption.LabelCount, yAxisOption.Unit, yAxisOption.LabelCountAdjustment,
					opt.seriesList, yIndex, opt.stackSeries,
					valueFormatter, yAxisOption.LabelRotation, yAxisOption.LabelFontStyle,
					yAxisOption.PreferNiceIntervals)
				entries[yIndex].prep = &prep
				valuePreps = append(valuePreps, entries[yIndex].prep)
				valuePrepIndices = append(valuePrepIndices, yIndex)
			}
		}
	}

	// coordinate and resolve value axis ranges
	if len(valuePreps) > 0 {
		ranges := coordinateValueAxisRanges(p, valuePreps)
		for i, yIndex := range valuePrepIndices {
			entries[yIndex].r = ranges[i]
		}
	}

	// render y-axes (reverse order so mark lines from left axis don't extend into right axis)
	for yIndex := yAxisCount - 1; yIndex >= 0; yIndex-- {
		entry := entries[yIndex]
		result.yaxisRanges[yIndex] = entry.r

		axisOpt := entry.option.toAxisOption(entry.r)
		if yIndex != 0 {
			axisOpt.splitLineShow = false // only show split lines on primary index axis
		}
		if axisOpt.position == "" {
			if yIndex == 0 {
				axisOpt.position = PositionLeft
			} else {
				axisOpt.position = PositionRight
			}
		}
		yAxisBox, err := newAxisPainter(p.Child(PainterPaddingOption(Box{
			Left:   rangeWidthLeft,
			Right:  rangeWidthRight,
			Bottom: xAxisHeight,
			IsSet:  true,
		})), axisOpt).Render()
		if err != nil {
			return nil, err
		} else if axisOpt.position == PositionRight {
			rangeWidthRight += yAxisBox.Width()
		} else {
			rangeWidthLeft += yAxisBox.Width()
		}
	}

	xAxisPadding := Box{
		Left:  rangeWidthLeft,
		Right: rangeWidthRight,
		IsSet: true,
	}
	if opt.axisReversed {
		xAxisOpts.aRange.size = p.Width() - rangeWidthLeft   // adjust size to match new painter dimensions
		xAxisOpts = opt.xAxis.toAxisOption(xAxisOpts.aRange) // regenerate axis options after value changes above
	} else {
		xAxisOpts.aRange.size -= rangeWidthLeft + rangeWidthRight // adjust size to match new painter dimensions
		xAxisPadding.Top = p.Height() - xAxisHeight
		xAxisOpts.painterPrePositioned = true // we must provide the exact painter position which will meet with the y-axis exactly
	}

	_, err = newAxisPainter(p.Child(PainterPaddingOption(xAxisPadding)), xAxisOpts).Render()
	if err != nil {
		return nil, err
	}

	result.xaxisRange = xAxisOpts.aRange
	result.seriesPainter = p.Child(PainterPaddingOption(Box{
		Left:   rangeWidthLeft,
		Right:  rangeWidthRight,
		Bottom: xAxisHeight,
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

// Render creates and renders a chart based on the provided options.
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
	lineSeriesList := filterSeriesList[LineSeriesList](opt.SeriesList, ChartTypeLine)
	scatterSeriesList := filterSeriesList[ScatterSeriesList](opt.SeriesList, ChartTypeScatter)
	barSeriesList := filterSeriesList[BarSeriesList](opt.SeriesList, ChartTypeBar)
	candlestickSeries := filterSeriesList[CandlestickSeriesList](opt.SeriesList, ChartTypeCandlestick)
	horizontalBarSeriesList := filterSeriesList[HorizontalBarSeriesList](opt.SeriesList, ChartTypeHorizontalBar)
	pieSeriesList := filterSeriesList[PieSeriesList](opt.SeriesList, ChartTypePie)
	doughnutSeriesList := filterSeriesList[DoughnutSeriesList](opt.SeriesList, ChartTypeDoughnut)
	radarSeriesList := filterSeriesList[RadarSeriesList](opt.SeriesList, ChartTypeRadar)
	funnelSeriesList := filterSeriesList[FunnelSeriesList](opt.SeriesList, ChartTypeFunnel)

	// Check if any incompatible chart types are being mixed
	// All compatible chart types need the absIndex field in the series
	seriesCount := len(seriesList)
	if len(horizontalBarSeriesList) != 0 && len(horizontalBarSeriesList) != seriesCount {
		return nil, errors.New("horizontal bar can not mix other charts")
	} else if len(pieSeriesList) != 0 && len(pieSeriesList) != seriesCount {
		return nil, errors.New("pie can not mix other charts")
	} else if len(doughnutSeriesList) != 0 && len(doughnutSeriesList) != seriesCount {
		return nil, errors.New("doughnut can not mix other charts")
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
		len(doughnutSeriesList) != 0 ||
		len(radarSeriesList) != 0 ||
		len(funnelSeriesList) != 0 {
		renderOpt.xAxis.Show = Ptr(false)
		renderOpt.yAxis = []YAxisOption{
			{
				Show: Ptr(false),
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
			width := opt.BarSize
			if width == 0 {
				width = opt.BarWidth
			}
			_, err := newBarChart(p, BarChartOption{
				Theme:       opt.Theme,
				Font:        opt.Font,
				XAxis:       opt.XAxis,
				SeriesList:  barSeriesList,
				StackSeries: opt.StackSeries,
				BarWidth:    width,
				BarMargin:   opt.BarMargin,
			}).renderChart(renderResult)
			return err
		})
	}

	// horizontal bar chart
	if len(horizontalBarSeriesList) != 0 {
		var yAxis YAxisOption
		if len(opt.YAxis) > 0 {
			if len(opt.YAxis) > 1 {
				return nil, errors.New("horizontal bar chart only accepts a single Y-Axis")
			}
			yAxis = opt.YAxis[0]
		}

		handler.Add(func() error {
			height := opt.BarSize
			if height == 0 {
				height = opt.BarHeight
			}
			_, err := newHorizontalBarChart(p, HorizontalBarChartOption{
				Theme:       opt.Theme,
				Font:        opt.Font,
				BarHeight:   height,
				BarMargin:   opt.BarMargin,
				YAxis:       yAxis,
				SeriesList:  horizontalBarSeriesList,
				StackSeries: opt.StackSeries,
			}).renderChart(renderResult)
			return err
		})
	}

	// candlestick chart
	if len(candlestickSeries) != 0 {
		handler.Add(func() error {
			_, err := newCandlestickChart(p, CandlestickChartOption{
				Theme:          opt.Theme,
				XAxis:          opt.XAxis,
				YAxis:          opt.YAxis,
				SeriesList:     candlestickSeries,
				CandleWidth:    0.8, // TODO - v0.6 - Use BarSize when it represents a percentage
				WickWidth:      1.0,
				ValueFormatter: opt.ValueFormatter,
			}).renderChart(renderResult)
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
				SeriesList:      lineSeriesList,
				StackSeries:     opt.StackSeries,
				Symbol:          opt.Symbol,
				LineStrokeWidth: opt.LineStrokeWidth,
				FillArea:        opt.FillArea,
				FillOpacity:     opt.FillOpacity,
			}).renderChart(renderResult)
			return err
		})
	}

	// scatter chart
	if len(scatterSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newScatterChart(p, ScatterChartOption{
				Theme:      opt.Theme,
				Font:       opt.Font,
				XAxis:      opt.XAxis,
				Symbol:     opt.Symbol,
				SeriesList: scatterSeriesList,
			}).renderChart(renderResult)
			return err
		})
	}

	// pie chart
	if len(pieSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newPieChart(p, PieChartOption{
				Theme:      opt.Theme,
				Font:       opt.Font,
				Radius:     opt.Radius,
				SeriesList: pieSeriesList,
			}).renderChart(renderResult)
			return err
		})
	}

	// doughnut chart
	if len(doughnutSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newDoughnutChart(p, DoughnutChartOption{
				Theme:      opt.Theme,
				RadiusRing: opt.Radius,
				SeriesList: doughnutSeriesList,
			}).renderChart(renderResult)
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
				Radius:          opt.Radius,
				SeriesList:      radarSeriesList,
			}).renderChart(renderResult)
			return err
		})
	}

	// funnel chart
	if len(funnelSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newFunnelChart(p, FunnelChartOption{
				Theme:      opt.Theme,
				Font:       opt.Font,
				SeriesList: funnelSeriesList,
			}).renderChart(renderResult)
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

package charts

import (
	"errors"
	"math"
	"strconv"
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
	// categoryAxis configures the chart's single category axis.
	// TODO - this contract supports at most one category axis and up-to-two value axes.
	// Extend when defaultRender supports dual category axes.
	categoryAxis *CategoryAxisOption
	// valueAxis configures one or two value axes. Length 0, 1, or 2.
	// Dual value axes are only supported when categoryY is false.
	valueAxis []ValueAxisOption
	// categoryY selects which physical axis holds the category axis.
	// false (typical): category on X, value on Y. true: category on Y, value on X.
	categoryY bool
	// title contains options for rendering the chart title.
	title TitleOption
	// legend contains options for the data legend.
	legend *LegendOption
	// backgroundIsFilled is true if the background is filled.
	backgroundIsFilled bool
	// valueFormatter formats numeric values into labels.
	valueFormatter ValueFormatter
}

type defaultRenderResult struct {
	valueAxisRanges   map[int]axisRange
	categoryAxisRange axisRange
	seriesPainter     *Painter
}

func (r *defaultRenderResult) renderNoData(theme ColorPalette) {
	p := r.seriesPainter
	w := p.Width()
	h := p.Height()
	if w <= 0 || h <= 0 {
		return
	}
	// draw empty set symbol (∅): circle with a diagonal line
	dim := float64(min(w, h))
	radius := dim / 5.0
	if radius < 8 {
		radius = 8
	} else if radius > 60 {
		radius = 60
	}
	strokeWidth := radius / 5.0
	if strokeWidth < 2 {
		strokeWidth = 2
	}
	cx := w / 2
	cy := h / 2
	strokeColor := theme.GetLabelTextColor()
	p.Circle(radius, cx, cy, ColorTransparent, strokeColor, strokeWidth)
	// diagonal line from bottom-left to top-right, extending beyond the circle
	extend := int(radius * 1.1)
	p.LineStroke([]Point{
		{X: cx - extend, Y: cy + extend},
		{X: cx + extend, Y: cy - extend},
	}, strokeColor, strokeWidth)
}

// defaultRender handles axis layout, title, and legend rendering for all chart types.
func defaultRender(p *Painter, opt defaultRenderOption) (*defaultRenderResult, error) {
	if len(opt.valueAxis) > 1 && opt.categoryY {
		return nil, errors.New("multiple value axes with categoryY is not supported") // TODO - future support for two continuous value axes
	}

	theme := getPreferredTheme(opt.theme, p.theme)
	fillThemeDefaults(theme, &opt.title, opt.legend, opt.categoryAxis, opt.valueAxis)
	opt.categoryAxis = opt.categoryAxis.prep(theme, opt.categoryY)
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
	opt.legend.seriesSymbols = make([]SymbolShape, opt.seriesList.len())
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

	// snapshot the pre-title frame the legend was measured against; the collision test below
	// translates the legend box from this frame into the later, title-childed plot frame.
	// flag a vertical overlay legend in the top area covering <=1/4 of the chart as a collision
	// candidate; the position-aware test runs after ranging.
	legendFrameTop, legendFrameLeft := p.box.Top, p.box.Left
	legendFrameW, legendFrameH := p.Width(), p.Height()
	legendTopOverlay := flagIs(true, legendOpt.Vertical) && !flagIs(false, legendOpt.OverlayChart) &&
		!opt.categoryY && !legendResult.IsZero() && !legendAtBottom &&
		legendResult.Bottom < legendFrameH/2 &&
		legendResult.Width()*legendResult.Height() <= legendFrameW*legendFrameH/4

	// reserve space for the legend if not in overlay mode
	// - horizontal legends reserve space (top or bottom)
	// - vertical legends always overlay from the side
	// - skip for adjustedForBottom since title reservation handles combined space
	if !legendResult.IsZero() && !flagIs(true, legendOpt.OverlayChart) && !adjustedForBottom &&
		!flagIs(true, legendOpt.Vertical) {
		if legendResult.Bottom < p.Height()/2 {
			// horizontal legend at top - reserve top space
			legendTopSpacing = max(legendResult.Height(), legendResult.Bottom) + legendTitlePadding
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
			titlePadBox.Top = max(legendTopSpacing, titleBox.Bottom+legendTitlePadding)
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
		valueAxisRanges: make(map[int]axisRange),
	}

	// calculate x-axis range and do a dry-render to find height
	// we will render on the actual painter once we know the space the y-axis will occupy
	var xAxisOpts axisOption
	var xValueAxis ValueAxisOption // prepped X-slot value axis; only populated when categoryY
	if opt.categoryY {             // X is value axis
		xValueAxis = opt.valueAxis[0] // value copy avoids mutating caller's slice
		xValueAxis.prep(getPreferredTheme(xValueAxis.Theme, theme), false)
		xAxisRange := calculateValueAxisRange(p, false, p.Width(),
			xValueAxis.Min, xValueAxis.Max, xValueAxis.RangeValuePaddingScale,
			xValueAxis.Labels,
			xValueAxis.LabelCount, xValueAxis.Unit, xValueAxis.LabelCountAdjustment,
			opt.seriesList, 0, opt.stackSeries,
			getPreferredValueFormatter(xValueAxis.ValueFormatter, opt.valueFormatter),
			xValueAxis.LabelRotation, xValueAxis.LabelFontStyle,
			xValueAxis.PreferNiceIntervals)
		xAxisOpts = xValueAxis.toAxisOption(xAxisRange)
	} else { // X is category axis (typical)
		xAxisRange := calculateCategoryAxisRange(p, p.Width(), false, flagIs(false, opt.categoryAxis.BoundaryGap),
			opt.categoryAxis.Labels,
			opt.categoryAxis.LabelCount, opt.categoryAxis.LabelCountAdjustment, opt.categoryAxis.Unit,
			opt.seriesList,
			opt.categoryAxis.LabelRotation, opt.categoryAxis.LabelFontStyle)
		xAxisOpts = opt.categoryAxis.toAxisOption(xAxisRange)
	}
	if xAxisOpts.position == "" {
		xAxisOpts.position = PositionBottom
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

	// prepare y-axis options, range data, and render
	yAxisCount := getSeriesYAxisCount(opt.seriesList)
	if yAxisCount < 0 {
		return nil, errors.New("series specified invalid y-axis index")
	}

	if opt.categoryY {
		// Y-slot renders the category axis directly
		ca := opt.categoryAxis
		catRange := calculateCategoryAxisRange(p, rangeHeight, true, false,
			ca.Labels,
			ca.LabelCount, ca.LabelCountAdjustment, ca.Unit,
			opt.seriesList,
			ca.LabelRotation, ca.LabelFontStyle)
		result.categoryAxisRange = catRange
		axisOpt := ca.toAxisOption(catRange)
		if axisOpt.position == "" {
			axisOpt.position = PositionLeft
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
	} else {
		// Y-slot renders value axis(es)
		type yAxisEntry struct {
			option ValueAxisOption
			prep   *valueAxisPrep
			r      axisRange
		}
		var entries []yAxisEntry
		var valuePreps []*valueAxisPrep
		var valuePrepIndices []int
		// reserve fixed pixel headroom above the data max so mark point pins are not clipped;
		// the pin head outer edge sits 5*symbolSize/4 above the data point
		var markPointClearance int
		if symSize := opt.seriesList.markPointSize(); symSize > 0 {
			markPointClearance = 5 * symSize / 4
		}
		if yAxisCount > 0 {
			entries = make([]yAxisEntry, yAxisCount)
			for yIndex := 0; yIndex < yAxisCount; yIndex++ {
				var yAxisOption ValueAxisOption
				if len(opt.valueAxis) > yIndex {
					yAxisOption = opt.valueAxis[yIndex]
				}
				yAxisOption = *yAxisOption.prep(getPreferredTheme(yAxisOption.Theme, theme), true)
				entries[yIndex].option = yAxisOption
				valueFormatter := getPreferredValueFormatter(yAxisOption.ValueFormatter, opt.valueFormatter)
				prep := prepareValueAxisRange(p, true, rangeHeight,
					yAxisOption.Min, yAxisOption.Max, yAxisOption.RangeValuePaddingScale,
					yAxisOption.Labels,
					yAxisOption.LabelCount, yAxisOption.Unit, yAxisOption.LabelCountAdjustment,
					opt.seriesList, yIndex, opt.stackSeries,
					valueFormatter, yAxisOption.LabelRotation, yAxisOption.LabelFontStyle,
					yAxisOption.PreferNiceIntervals)
				prep.maxClearancePx = markPointClearance
				entries[yIndex].prep = &prep
				valuePreps = append(valuePreps, entries[yIndex].prep)
				valuePrepIndices = append(valuePrepIndices, yIndex)
			}
		}

		// coordinate and resolve value axis ranges
		if len(valuePreps) > 0 {
			ranges := coordinateValueAxisRanges(p, valuePreps)
			for i, yIndex := range valuePrepIndices {
				entries[yIndex].r = ranges[i]
			}
		}

		// vertical overlay legend: if data (or a mark point pin) under the legend's horizontal span
		// would collide with it, reserve top headroom to push the chart below the legend.
		// single value axis only; dual axes share a label count, so re-resolving one shifts the other.
		if legendTopOverlay && len(valuePreps) == 1 {
			legendBottomRel := legendResult.Bottom - (p.box.Top - legendFrameTop)
			leftOffset := p.box.Left - legendFrameLeft
			n := getSeriesMaxDataCount(opt.seriesList)
			if lo, hi, ok := legendIndexSpan(legendResult.Left-leftOffset, legendResult.Right-leftOffset, p.Width(), n); ok {
				if localMax, ok := localMaxOverIndices(opt.seriesList, 0, lo, hi, opt.stackSeries); ok {
					if entries[0].r.getRestHeight(localMax)-markPointClearance < legendBottomRel {
						if needed := legendBottomRel + markPointClearance; needed > entries[0].prep.maxClearancePx {
							entries[0].prep.maxClearancePx = needed
							ranges := coordinateValueAxisRanges(p, valuePreps)
							for i, yIndex := range valuePrepIndices {
								entries[yIndex].r = ranges[i]
							}
						}
					}
				}
			}
		}

		// render y-axes (reverse order so mark lines from left axis don't extend into right axis)
		for yIndex := yAxisCount - 1; yIndex >= 0; yIndex-- {
			entry := entries[yIndex]
			result.valueAxisRanges[yIndex] = entry.r

			axisOpt := entry.option.toAxisOption(entry.r)
			if yIndex != 0 {
				axisOpt.splitLineShow = Ptr(false) // only show split lines on primary index axis
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
	}

	xAxisPadding := Box{
		Left:  rangeWidthLeft,
		Right: rangeWidthRight,
		IsSet: true,
	}
	if opt.categoryY {
		xAxisOpts.aRange.size = p.Width() - rangeWidthLeft - rangeWidthRight // adjust size to match new painter dimensions
		// right-positioned category axis mirrors the value axis so bars grow from the category baseline
		xAxisOpts.aRange.reversed = opt.categoryAxis.Position == PositionRight
		xAxisOpts = xValueAxis.toAxisOption(xAxisOpts.aRange) // regenerate axis options after value changes above
		if xAxisOpts.position == "" {
			xAxisOpts.position = PositionBottom
		}
	} else {
		xAxisOpts.aRange.size -= rangeWidthLeft + rangeWidthRight // adjust size to match new painter dimensions
		xAxisPadding.Top = p.Height() - xAxisHeight
		xAxisOpts.painterPrePositioned = true // we must provide the exact painter position which will meet with the y-axis exactly
	}

	_, err = newAxisPainter(p.Child(PainterPaddingOption(xAxisPadding)), xAxisOpts).Render()
	if err != nil {
		return nil, err
	}

	if opt.categoryY {
		result.valueAxisRanges[0] = xAxisOpts.aRange
	} else {
		result.categoryAxisRange = xAxisOpts.aRange
	}
	result.seriesPainter = p.Child(PainterPaddingOption(Box{
		Left:   rangeWidthLeft,
		Right:  rangeWidthRight,
		Bottom: xAxisHeight,
		IsSet:  true,
	}))
	return &result, nil
}

// legendIndexSpan maps a legend's horizontal pixel span to an inclusive data index range [lo, hi]
// over n points across plotWidth. ok is false when there are no points or no overlap.
func legendIndexSpan(legendLeft, legendRight, plotWidth, n int) (int, int, bool) {
	if n <= 0 || plotWidth <= 0 || legendRight < 0 || legendLeft > plotWidth {
		return 0, 0, false
	} else if n == 1 {
		return 0, 0, true
	}
	if legendLeft < 0 {
		legendLeft = 0
	}
	if legendRight > plotWidth {
		legendRight = plotWidth
	}
	last := n - 1
	lo := legendLeft * last / plotWidth
	hi := (legendRight*last + plotWidth - 1) / plotWidth // ceil toward the right edge
	if hi > last {
		hi = last
	}
	if lo > hi {
		lo = hi
	}
	return lo, hi, true
}

// localMaxOverIndices returns the max plotted value among series on yaxisIndex over the inclusive
// index range [lo, hi]. When stackSeries, per-index values are summed across series on that axis.
// Invalid (NaN/null) values are skipped; the bool is false when no valid value falls in the range.
func localMaxOverIndices(sl seriesList, yaxisIndex, lo, hi int, stackSeries bool) (float64, bool) {
	if lo > hi {
		return 0, false
	}
	max := -math.MaxFloat64
	var found bool
	if stackSeries {
		sums := make([]float64, hi-lo+1)
		valid := make([]bool, hi-lo+1)
		for i := 0; i < sl.len(); i++ {
			series := sl.getSeries(i)
			if series.getYAxisIndex() != yaxisIndex {
				continue
			}
			values := series.getValues()
			for idx := lo; idx <= hi && idx < len(values); idx++ {
				if !isValidExtent(values[idx]) {
					continue
				}
				sums[idx-lo] += values[idx]
				valid[idx-lo] = true
			}
		}
		for j, ok := range valid {
			if ok && sums[j] > max {
				max = sums[j]
				found = true
			}
		}
	} else {
		for i := 0; i < sl.len(); i++ {
			series := sl.getSeries(i)
			if series.getYAxisIndex() != yaxisIndex {
				continue
			}
			values := series.getValues()
			for idx := lo; idx <= hi && idx < len(values); idx++ {
				if !isValidExtent(values[idx]) {
					continue
				}
				if values[idx] > max {
					max = values[idx]
					found = true
				}
			}
		}
	}
	if !found {
		return 0, false
	}
	return max, true
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
	horizontalBarSeriesList := filterSeriesList[BarSeriesList](opt.SeriesList, ChartTypeHorizontalBar)
	pieSeriesList := filterSeriesList[PieSeriesList](opt.SeriesList, ChartTypePie)
	doughnutSeriesList := filterSeriesList[DoughnutSeriesList](opt.SeriesList, ChartTypeDoughnut)
	radarSeriesList := filterSeriesList[RadarSeriesList](opt.SeriesList, ChartTypeRadar)
	funnelSeriesList := filterSeriesList[FunnelSeriesList](opt.SeriesList, ChartTypeFunnel)
	violinSeriesList := filterSeriesList[ViolinSeriesList](opt.SeriesList, ChartTypeViolin)
	horizontalViolinSeriesList := filterSeriesList[ViolinSeriesList](opt.SeriesList, ChartTypeHorizontalViolin)

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
	} else if len(violinSeriesList) != 0 && len(violinSeriesList) != seriesCount {
		return nil, errors.New("violin can not mix other charts")
	} else if len(horizontalViolinSeriesList) != 0 && len(horizontalViolinSeriesList) != seriesCount {
		return nil, errors.New("horizontal violin can not mix other charts")
	}

	// boundary gap must be resolved here as it's shared between the axis and the chart handlers
	// TODO - BoundaryGap behavior may not be accurate for chart types which conditionally select the default behavior
	if opt.XAxis.BoundaryGap == nil {
		if len(barSeriesList) != 0 || len(candlestickSeries) != 0 {
			if len(lineSeriesList) != 0 || len(scatterSeriesList) != 0 {
				opt.XAxis.BoundaryGap = Ptr(true) // align points to the bar / candle slot centers
			}
		} else if len(scatterSeriesList) != 0 {
			opt.XAxis.BoundaryGap = Ptr(false)
		}
	}

	categoryY := len(horizontalBarSeriesList) != 0 || len(violinSeriesList) != 0
	renderOpt := defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		categoryAxis:   &opt.XAxis,
		valueAxis:      opt.YAxis,
		stackSeries:    flagIs(true, opt.StackSeries),
		title:          opt.Title,
		legend:         &opt.Legend,
		categoryY:      categoryY,
		valueFormatter: opt.ValueFormatter,
		// the background color has been set
		backgroundIsFilled: true,
	}
	if len(pieSeriesList) != 0 ||
		len(doughnutSeriesList) != 0 ||
		len(radarSeriesList) != 0 ||
		len(funnelSeriesList) != 0 {
		renderOpt.categoryAxis.Show = Ptr(false)
		renderOpt.valueAxis = []ValueAxisOption{
			{
				Show: Ptr(false),
			},
		}
	}
	if len(horizontalBarSeriesList) != 0 {
		// horizontal bar: XAxis carries value data, YAxis carries category data
		var catAxis CategoryAxisOption
		if len(opt.YAxis) > 0 {
			ya := opt.YAxis[0]
			catAxis = CategoryAxisOption{
				Show:                 ya.Show,
				Theme:                ya.Theme,
				Title:                ya.Title,
				TitleFontStyle:       ya.TitleFontStyle,
				Labels:               ya.Labels,
				Position:             ya.Position,
				LabelFontStyle:       ya.LabelFontStyle,
				LabelRotation:        ya.LabelRotation,
				ValueFormatter:       ya.ValueFormatter,
				Unit:                 ya.Unit,
				LabelCount:           ya.LabelCount,
				LabelCountAdjustment: ya.LabelCountAdjustment,
			}
		}
		catAxis.Unit = 1 // each category bar should have a label; axis fitting skips only when labels don't physically fit
		valAxis := ValueAxisOption{
			Show:                 opt.XAxis.Show,
			Theme:                opt.XAxis.Theme,
			Title:                opt.XAxis.Title,
			TitleFontStyle:       opt.XAxis.TitleFontStyle,
			Labels:               opt.XAxis.Labels,
			LabelFontStyle:       opt.XAxis.LabelFontStyle,
			LabelRotation:        opt.XAxis.LabelRotation,
			ValueFormatter:       opt.XAxis.ValueFormatter,
			Unit:                 opt.XAxis.Unit,
			LabelCount:           opt.XAxis.LabelCount,
			LabelCountAdjustment: opt.XAxis.LabelCountAdjustment,
		}
		renderOpt.categoryAxis = &catAxis
		renderOpt.valueAxis = []ValueAxisOption{valAxis}
	}
	if len(violinSeriesList) != 0 {
		var catAxis CategoryAxisOption
		var valAxis ValueAxisOption
		if len(renderOpt.valueAxis) > 0 {
			valAxis = renderOpt.valueAxis[0]
		}
		catAxis, valAxis = violinConfigureRenderOption(violinSeriesList, false, nil,
			getPreferredValueFormatter(opt.XAxis.ValueFormatter, opt.ValueFormatter),
			catAxis, valAxis)
		renderOpt.categoryAxis = &catAxis
		renderOpt.valueAxis = []ValueAxisOption{valAxis}
		renderOpt.categoryY = true
	} else if len(horizontalViolinSeriesList) != 0 {
		var catAxis CategoryAxisOption
		var valAxis ValueAxisOption
		if len(renderOpt.valueAxis) > 0 {
			valAxis = renderOpt.valueAxis[0]
		}
		var yAxisValueFormatter ValueFormatter
		if len(opt.YAxis) > 0 {
			yAxisValueFormatter = opt.YAxis[0].ValueFormatter
		}
		catAxis, valAxis = violinConfigureRenderOption(horizontalViolinSeriesList, true, nil,
			getPreferredValueFormatter(yAxisValueFormatter, opt.ValueFormatter),
			catAxis, valAxis)
		renderOpt.categoryAxis = &catAxis
		renderOpt.valueAxis = []ValueAxisOption{valAxis}
		renderOpt.categoryY = false
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
				Theme:        opt.Theme,
				CategoryAxis: opt.XAxis,
				SeriesList:   barSeriesList,
				StackSeries:  opt.StackSeries,
				BarSize:      opt.BarSize,
				BarMargin:    opt.BarMargin,
			}).renderChart(renderResult)
			return err
		})
	}

	// horizontal bar chart
	if len(horizontalBarSeriesList) != 0 {
		if len(opt.YAxis) > 1 {
			return nil, errors.New("horizontal bar chart only accepts a single Y-Axis")
		}

		handler.Add(func() error {
			_, err := newBarChart(p, BarChartOption{
				Horizontal:     true,
				Theme:          opt.Theme,
				BarSize:        opt.BarSize,
				BarMargin:      opt.BarMargin,
				SeriesList:     horizontalBarSeriesList,
				StackSeries:    opt.StackSeries,
				ValueFormatter: opt.ValueFormatter,
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
				CandleWidth:    opt.BarSize,
				WickWidth:      1.0,
				ValueFormatter: opt.ValueFormatter,
			}).renderChart(renderResult)
			return err
		})
	}

	// violin chart
	if len(violinSeriesList) != 0 || len(horizontalViolinSeriesList) != 0 {
		// TODO - ChartOption does not expose options: ShowSpine, SpineWidth, ViolinWidth
		vOpt := ViolinChartOption{
			Theme:          opt.Theme,
			ValueFormatter: opt.ValueFormatter,
		}
		if len(violinSeriesList) != 0 {
			vOpt.SeriesList = violinSeriesList
		} else {
			vOpt.Horizontal = true
			vOpt.SeriesList = horizontalViolinSeriesList
		}
		handler.Add(func() error {
			_, err := newViolinChart(p, vOpt).renderChart(renderResult)
			return err
		})
	}

	// line chart
	if len(lineSeriesList) != 0 {
		handler.Add(func() error {
			_, err := newLineChart(p, LineChartOption{
				Theme:           opt.Theme,
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
		if _, err = Render(item); err != nil {
			return nil, err
		}
	}

	return p, nil
}

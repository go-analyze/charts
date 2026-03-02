package charts

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/go-analyze/bulk"
)

const (
	// ViolinNormalizationPerSeries scales each violin to its own density max.
	ViolinNormalizationPerSeries = "series"
	// ViolinNormalizationGlobal scales all violins to a shared density max.
	ViolinNormalizationGlobal = "global"
	// Use an odd default so 0 always lands on the center tick.
	violinDefaultValueAxisLabelCount = 5
)

// ViolinAxis contains axis configuration options for violin charts.
type ViolinAxis struct {
	// Show specifies if the value axis should be rendered. Set to *false (via Ptr(false)) to hide the axis.
	Show *bool
	// Title specifies a name for the value axis.
	Title string
	// TitleFontStyle specifies the font, size, and color for the value axis title.
	TitleFontStyle FontStyle
	// LabelFontStyle specifies the font configuration for value axis labels.
	LabelFontStyle FontStyle
	// LabelRotation is the rotation angle in radians for value axis labels.
	LabelRotation float64
	// Limit forces the absolute value axis extent when set (Use Ptr(float64)).
	// Violin value axes are symmetric around zero, so this produces [-Limit, +Limit].
	Limit *float64
	// Unit suggests the value axis step size (recommendation only). Larger values result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the value axis.
	LabelCount int
	// LabelCountAdjustment specifies relative influence on label count.
	LabelCountAdjustment int
	// PreferNiceIntervals allows the label count to flex slightly to produce rounder axis intervals.
	PreferNiceIntervals *bool
}

// ViolinChartOption defines the options for rendering a violin chart. Render the chart using Painter.ViolinChart.
type ViolinChartOption struct {
	// Theme specifies the colors used for the violin chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// SeriesList provides the data population for the chart. Typically constructed using NewSeriesListViolin.
	SeriesList ViolinSeriesList
	// ValueAxis contains configuration options for the numeric value axis.
	ValueAxis ViolinAxis
	// Title contains options for rendering the chart title.
	Title TitleOption
	// Legend contains options for the data legend.
	Legend LegendOption
	// Horizontal when true renders horizontal violins (Y is numeric, X is category).
	Horizontal bool
	// ViolinWidth specifies the max width of each violin. Accepts a percentage (e.g. "40%") or
	// pixel value (e.g. "64"). May be reduced to fit.
	ViolinWidth string
	// ShowSpine when non-nil controls spine visibility. Default true.
	ShowSpine *bool
	// SpineWidth is the stroke width of the spine. Default 1.0.
	SpineWidth float64
	// ValueFormatter defines how float values are rendered to strings for axis labels.
	ValueFormatter ValueFormatter

	// TODO - allow an outline around the violin to be configured
}

// TODO - inner statistics overlay (box plot or quartile markers inside each violin)

type violinChart struct {
	p   *Painter
	opt *ViolinChartOption
}

// newViolinChart returns a violin chart renderer.
func newViolinChart(p *Painter, opt ViolinChartOption) *violinChart {
	return &violinChart{
		p:   p,
		opt: &opt,
	}
}

// NewViolinChartOptionWithData returns an initialized ViolinChartOption with data.
func NewViolinChartOptionWithData(data [][][2]float64) ViolinChartOption {
	return NewViolinChartOptionWithSeries(NewSeriesListViolin(data))
}

// NewViolinChartOptionWithSeries returns an initialized ViolinChartOption with the provided SeriesList.
func NewViolinChartOptionWithSeries(sl ViolinSeriesList) ViolinChartOption {
	return ViolinChartOption{
		SeriesList:     sl,
		Padding:        defaultPadding,
		Theme:          GetDefaultTheme(),
		ValueFormatter: defaultValueFormatter,
	}
}

// NewViolinChartOptionWithSamples builds a violin chart option using Gaussian KDE from sample values.
// Optional string arguments may include normalization mode ("series" or "global")
// and/or KDE bandwidth override ("bandwidth=<value>" or "bw=<value>").
func NewViolinChartOptionWithSamples(samples [][]float64, pointCount int, normalization ...string) (ViolinChartOption, error) {
	if pointCount <= 0 {
		return ViolinChartOption{}, errors.New("point count must be greater than 0")
	}

	mode, bandwidthOverride, err := parseViolinSampleOptions(normalization)
	if err != nil {
		return ViolinChartOption{}, err
	}

	type densitySeries struct {
		values []float64
		max    float64
	}

	densities := make([]densitySeries, len(samples))
	seriesData := make([][][2]float64, len(samples))
	seriesStats := make([]*PopulationSummary, len(samples))
	var globalMax float64
	var validSeries int
	for i := range samples {
		filtered := bulk.SliceFilter(isValidExtent, samples[i])
		if len(filtered) == 0 {
			continue
		}
		validSeries++
		seriesStats[i] = Ptr(summarizePopulationData(filtered))

		density := gaussianKDE(filtered, pointCount, bandwidthOverride)
		if len(density) == 0 {
			continue
		}

		maxDensity := maxFloatSlice(density)
		densities[i] = densitySeries{values: density, max: maxDensity}
		if maxDensity > globalMax {
			globalMax = maxDensity
		}
	}

	if validSeries == 0 {
		return ViolinChartOption{}, errors.New("series with no valid data")
	}

	for i, density := range densities {
		if len(density.values) == 0 {
			continue
		}

		normMax := density.max
		if mode == ViolinNormalizationGlobal {
			normMax = globalMax
		}
		if normMax <= 0 {
			continue
		}

		pairs := make([][2]float64, len(density.values))
		for j, d := range density.values {
			extent := d / normMax
			pairs[j] = [2]float64{extent, extent}
		}
		seriesData[i] = pairs
	}

	opt := NewViolinChartOptionWithData(seriesData)
	for i := range opt.SeriesList {
		opt.SeriesList[i].Stats = seriesStats[i]
	}
	return opt, nil
}

func parseViolinSampleOptions(options []string) (string, *float64, error) {
	mode := ViolinNormalizationPerSeries
	var modeSet bool
	var bandwidthOverride *float64

	for _, rawOption := range options {
		option := strings.TrimSpace(rawOption)
		if option == "" {
			continue
		}

		switch option {
		case ViolinNormalizationPerSeries, ViolinNormalizationGlobal:
			if modeSet && mode != option {
				return "", nil, errors.New("conflicting normalization mode")
			}
			mode = option
			modeSet = true
			continue
		}

		// TODO - v0.6 - Use strings.CutPrefix with go update
		var valueText string
		var ok bool
		if strings.HasPrefix(option, "bandwidth=") {
			valueText = option[len("bandwidth="):]
			ok = true
		} else if strings.HasPrefix(option, "bw=") {
			valueText = option[len("bw="):]
			ok = true
		}
		if ok {
			bandwidth, err := strconv.ParseFloat(strings.TrimSpace(valueText), 64)
			if err != nil || bandwidth <= 0 || math.IsNaN(bandwidth) || math.IsInf(bandwidth, 0) {
				return "", nil, errors.New("invalid bandwidth override")
			}
			bandwidthOverride = &bandwidth
			continue
		}

		return "", nil, errors.New("invalid normalization mode")
	}

	return mode, bandwidthOverride, nil
}

func gaussianKDE(samples []float64, pointCount int, bandwidthOverride *float64) []float64 {
	if len(samples) == 0 || pointCount <= 0 {
		return nil
	}

	summary := summarizePopulationData(samples)
	if summary.StandardDeviation == 0 || summary.Max <= summary.Min {
		return nil
	}

	n := float64(len(samples))
	bandwidth := 1.06 * summary.StandardDeviation * math.Pow(n, -0.2)
	if bandwidthOverride != nil {
		bandwidth = *bandwidthOverride
	}
	if bandwidth <= 0 || math.IsNaN(bandwidth) || math.IsInf(bandwidth, 0) {
		return nil
	}

	values := make([]float64, pointCount)
	norm := 1.0 / (n * bandwidth * math.Sqrt(2*math.Pi))
	for i := 0; i < pointCount; i++ {
		var x float64
		if pointCount > 1 {
			x = summary.Min + (summary.Max-summary.Min)*(float64(i)/float64(pointCount-1))
		} else {
			x = (summary.Min + summary.Max) / 2.0
		}

		var sum float64
		for _, sample := range samples {
			u := (x - sample) / bandwidth
			sum += math.Exp(-0.5 * u * u)
		}
		values[i] = norm * sum
	}
	return values
}

func maxFloatSlice(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, value := range values[1:] {
		if value > maxVal {
			maxVal = value
		}
	}
	return maxVal
}

func (v *violinChart) renderChart(result *defaultRenderResult) (Box, error) {
	p := v.p
	opt := v.opt
	seriesCount := len(opt.SeriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	for i := range opt.SeriesList {
		if opt.SeriesList[i].YAxisIndex != 0 {
			return BoxZero, errors.New("violin series specified invalid y-axis index")
		}
	}
	seriesPainter := result.seriesPainter

	// Vertical: xaxisRange=value, yaxisRanges[0]=category
	// Horizontal: xaxisRange=category, yaxisRanges[0]=value
	var valueRange, categoryRange axisRange
	if opt.Horizontal {
		categoryRange = result.xaxisRange
		valueRange = result.yaxisRanges[0]
	} else {
		valueRange = result.xaxisRange
		categoryRange = result.yaxisRanges[0]
	}
	if categoryRange.divideCount == 0 {
		return BoxZero, errors.New("violin category axis produced no slots")
	}

	c0, c1 := categoryRange.getRange(0)
	slotSize := int(c1 - c0)
	var parsedViolinWidth int
	if opt.ViolinWidth != "" {
		if w, err := parseFlexibleValue(opt.ViolinWidth, float64(slotSize)); err == nil && w > 0 {
			parsedViolinWidth = int(w)
		}
	}
	margin, _, violinWidth := calculateBarMarginsAndSize(1, slotSize, parsedViolinWidth, nil)
	divideValues := categoryRange.autoDivide()

	// Resolve style defaults
	showSpine := opt.ShowSpine == nil || *opt.ShowSpine
	spineWidth := opt.SpineWidth
	if spineWidth == 0 {
		spineWidth = 1.0
	}
	spineColor := opt.Theme.GetAxisSplitLineColor()

	plotH := seriesPainter.Height()
	var markLineRenderers []renderer

	for index, series := range opt.SeriesList {
		if index >= categoryRange.divideCount {
			break
		}
		seriesMarks := series.MarkLine.Lines.filterGlobal(false)
		if len(seriesMarks) > 0 && series.Stats == nil {
			return BoxZero, errors.New("violin mark line requires series stats")
		}
		seriesThemeIndex := index
		if series.absThemeIndex != nil {
			seriesThemeIndex = *series.absThemeIndex
		}
		seriesColor := opt.Theme.GetSeriesColor(seriesThemeIndex)

		// Category slot position for this series
		slotStart := divideValues[index]

		bandCount := len(series.Data)
		if bandCount != 0 {
			bandSize := float64(violinWidth) / float64(bandCount)

			for j, pair := range series.Data {
				a, b := pair[0], pair[1]
				// Sanitize: null/NaN/Inf => 0
				if !isValidExtent(a) {
					a = 0
				}
				if !isValidExtent(b) {
					b = 0
				}
				// Use absolute magnitude
				a, b = math.Abs(a), math.Abs(b)

				if a == 0 && b == 0 {
					continue
				}

				bandStart := slotStart + margin + int(float64(j)*bandSize)
				bandEnd := slotStart + margin + int(float64(j+1)*bandSize)

				if opt.Horizontal { // Horizontal: bands go left-right in slot, value goes up-down
					top := plotH - valueRange.getHeight(b)
					bottom := plotH - valueRange.getHeight(-a)
					if top > bottom { // Normalize
						top, bottom = bottom, top
					}
					seriesPainter.FilledRect(bandStart, top, bandEnd, bottom, seriesColor, ColorTransparent, 0)
				} else { // Vertical: bands go top-bottom in slot, value goes left-right
					left, right := valueRange.getHeight(-a), valueRange.getHeight(b)
					if left > right { // Normalize
						left, right = right, left
					}
					seriesPainter.FilledRect(left, bandStart, right, bandEnd, seriesColor, ColorTransparent, 0)
				}
			}
		}

		// Draw spine at the slot center above filled bands
		if showSpine {
			if opt.Horizontal {
				spineY := plotH - valueRange.getHeight(0)
				seriesPainter.LineStroke([]Point{
					{X: slotStart + margin, Y: spineY},
					{X: slotStart + margin + violinWidth, Y: spineY},
				}, spineColor, spineWidth)
			} else {
				spineX := valueRange.getHeight(0)
				seriesPainter.LineStroke([]Point{
					{X: spineX, Y: slotStart + margin},
					{X: spineX, Y: slotStart + margin + violinWidth},
				}, spineColor, spineWidth)
			}
		}

		// Queue mark lines so they render after all bands/spines
		markLine := buildViolinMarkLineRenderer(seriesPainter, opt, series, seriesMarks,
			seriesColor, slotStart, margin, violinWidth)
		if markLine != nil {
			markLineRenderers = append(markLineRenderers, markLine)
		}
	}

	if err := doRender(markLineRenderers...); err != nil {
		return BoxZero, err
	}
	return p.box, nil
}

func buildViolinMarkLineRenderer(seriesPainter *Painter, opt *ViolinChartOption, series ViolinSeries,
	seriesMarks []SeriesMark, seriesColor Color, slotStart, margin, violinWidth int) renderer {
	if len(seriesMarks) == 0 || series.Stats == nil {
		return nil
	}

	markValueFormatter := getPreferredValueFormatter(series.MarkLine.ValueFormatter, opt.ValueFormatter)
	markColor := seriesColor.WithAdjustHSL(0.0, 0.1, -0.2)
	violinStart := slotStart + margin
	violinEnd := violinStart + violinWidth

	var markLineBox Box
	if opt.Horizontal {
		markLineBox = NewBox(
			seriesPainter.box.Left+violinStart,
			seriesPainter.box.Top,
			seriesPainter.box.Left+violinEnd,
			seriesPainter.box.Bottom,
		)
	} else {
		markLineBox = NewBox(
			seriesPainter.box.Left,
			seriesPainter.box.Top+violinStart,
			seriesPainter.box.Right,
			seriesPainter.box.Top+violinEnd,
		)
	}
	markLinePainter := seriesPainter.Child(PainterBoxOption(markLineBox))

	// Map data values to positions along the violin body
	axis := axisRange{
		min:  series.Stats.Min,
		max:  series.Stats.Max,
		size: violinWidth,
	}

	markLine := newMarkLinePainter(markLinePainter)
	markLine.add(markLineRenderOption{
		fillColor:      markColor,
		fontColor:      opt.Theme.GetMarkTextColor(),
		strokeColor:    markColor,
		font:           seriesPainter.font,
		seriesSummary:  series.Stats,
		marklines:      seriesMarks,
		axisRange:      axis,
		valueFormatter: markValueFormatter,
		verticalLine:   opt.Horizontal, // Draw perpendicular to the spine
		// Vertical violins use top-to-bottom bucket mapping
		horizontalUsesHeight: !opt.Horizontal,
	})
	return markLine
}

// violinConfigureRenderOption sets up the axis configuration on a defaultRenderOption for violin charts.
// The limit parameter is only available through the Painter API's ViolinAxis.Limit.
func violinConfigureRenderOption(renderOpt *defaultRenderOption, sl ViolinSeriesList, horizontal bool, limit *float64, valueFormatter ValueFormatter) {
	if horizontal {
		for i := range sl {
			sl[i].horizontal = true
		}
	}
	categoryLabels := sl.names()
	absMax := violinMaxAbsExtent(sl)
	baseFormatter := getPreferredValueFormatter(valueFormatter)
	absoluteValueFormatter := func(val float64) string {
		return baseFormatter(math.Abs(val))
	}

	if horizontal {
		// Horizontal: X is category (hidden), Y is numeric value axis
		renderOpt.xAxis.Show = Ptr(false)
		renderOpt.xAxis.Labels = categoryLabels
		valueAxisLabelCount :=
			violinResolveValueAxisLabelCount(renderOpt.yAxis[0].LabelCount, renderOpt.yAxis[0].LabelCountAdjustment)
		axisMin, axisMax := violinResolveAxisBounds(absMax, limit, valueAxisLabelCount, renderOpt.yAxis[0].Unit)
		renderOpt.yAxis[0].Min = &axisMin
		renderOpt.yAxis[0].Max = &axisMax
		renderOpt.yAxis[0].LabelCount = valueAxisLabelCount
		renderOpt.yAxis[0].LabelCountAdjustment = 0
		renderOpt.yAxis[0].ValueFormatter = absoluteValueFormatter
	} else {
		// Vertical: Y is category (hidden), X is numeric value axis
		renderOpt.yAxis[0].Show = Ptr(false)
		renderOpt.yAxis[0].Labels = categoryLabels
		renderOpt.xAxis.Labels = nil
		renderOpt.xAxis.LabelCount =
			violinResolveValueAxisLabelCount(renderOpt.xAxis.LabelCount, renderOpt.xAxis.LabelCountAdjustment)
		renderOpt.xAxis.LabelCountAdjustment = 0
		renderOpt.xAxis.ValueFormatter = absoluteValueFormatter
		axisMin, axisMax := violinResolveAxisBounds(absMax, limit, renderOpt.xAxis.LabelCount, renderOpt.xAxis.Unit)
		renderOpt.axisRangeOverride = &[2]float64{axisMin, axisMax}
	}
}

func violinMaxAbsExtent(sl ViolinSeriesList) float64 {
	var absMax float64
	for _, series := range sl {
		for _, pair := range series.Data {
			if isValidExtent(pair[0]) {
				if abs := math.Abs(pair[0]); abs > absMax {
					absMax = abs
				}
			}
			if isValidExtent(pair[1]) {
				if abs := math.Abs(pair[1]); abs > absMax {
					absMax = abs
				}
			}
		}
	}
	return absMax
}

func violinResolveAxisBounds(absMax float64, limitOpt *float64, labelCount int, unit float64) (float64, float64) {
	extent := absMax
	if extent <= 0 {
		extent = float64(zeroSpanAdjustment)
	}
	if limitOpt != nil {
		extent = math.Max(extent, math.Abs(*limitOpt))
		return -extent, extent
	}
	stepsFromZero := (violinEnsureZeroLabelCount(labelCount) - 1) / 2
	extent = roundSymmetricAxisExtent(extent, stepsFromZero, unit)
	return -extent, extent
}

func violinEnsureZeroLabelCount(labelCount int) int {
	if labelCount <= 0 {
		return violinDefaultValueAxisLabelCount
	} else if labelCount < 3 {
		return 3
	} else if labelCount%2 == 0 {
		return labelCount + 1
	}
	return labelCount
}

func violinResolveValueAxisLabelCount(labelCount, labelCountAdjustment int) int {
	baseCount := labelCount
	if baseCount <= 0 {
		baseCount = violinDefaultValueAxisLabelCount
	}
	return violinEnsureZeroLabelCount(baseCount + labelCountAdjustment)
}

func (v *violinChart) Render() (Box, error) {
	p := v.p
	opt := v.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	if opt.Legend.Symbol == "" {
		opt.Legend.Symbol = SymbolSquare
	}

	// Vertical (default): X is numeric, Y is hidden category
	// Horizontal: Y is numeric, X is hidden category
	axisOpt := opt.ValueAxis
	xAxisOpt := XAxisOption{}
	yAxisOpts := []YAxisOption{{}}

	// Map ViolinAxis fields onto the active value axis
	if opt.Horizontal {
		yAxisOpts[0].Show = axisOpt.Show
		yAxisOpts[0].Title = axisOpt.Title
		yAxisOpts[0].TitleFontStyle = axisOpt.TitleFontStyle
		yAxisOpts[0].LabelFontStyle = axisOpt.LabelFontStyle
		yAxisOpts[0].LabelRotation = axisOpt.LabelRotation
		yAxisOpts[0].Unit = axisOpt.Unit
		yAxisOpts[0].LabelCount = axisOpt.LabelCount
		yAxisOpts[0].LabelCountAdjustment = axisOpt.LabelCountAdjustment
		yAxisOpts[0].PreferNiceIntervals = axisOpt.PreferNiceIntervals
	} else {
		xAxisOpt.Show = axisOpt.Show
		xAxisOpt.Title = axisOpt.Title
		xAxisOpt.TitleFontStyle = axisOpt.TitleFontStyle
		xAxisOpt.LabelFontStyle = axisOpt.LabelFontStyle
		xAxisOpt.LabelRotation = axisOpt.LabelRotation
		xAxisOpt.Unit = axisOpt.Unit
		xAxisOpt.LabelCount = axisOpt.LabelCount
		xAxisOpt.LabelCountAdjustment = axisOpt.LabelCountAdjustment
	}

	renderOpt := defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		xAxis:          &xAxisOpt,
		yAxis:          yAxisOpts,
		title:          opt.Title,
		legend:         &opt.Legend,
		valueFormatter: opt.ValueFormatter,
		axisReversed:   !opt.Horizontal,
	}
	violinConfigureRenderOption(&renderOpt, opt.SeriesList, opt.Horizontal, axisOpt.Limit, opt.ValueFormatter)

	renderResult, err := defaultRender(p, renderOpt)
	if err != nil {
		return BoxZero, err
	}
	return v.renderChart(renderResult)
}

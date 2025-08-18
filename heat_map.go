package charts

import (
	"errors"
	"math"
	"strconv"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/matrix"
)

// HeatMapOption contains configuration options for a heat map chart. Render the chart using Painter.HeatMapChart.
type HeatMapOption struct {
	// Theme specifies the color palette used for rendering the heat map.
	Theme ColorPalette
	// BaseColorIndex specifies which color from the theme palette to use as the base for gradients.
	BaseColorIndex int
	// Padding specifies the padding around the heat map chart.
	Padding Box
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// Title contains options for rendering the chart title.
	Title TitleOption
	// Values provides the 2D data for the heat map.
	// The outer slice represents the rows (y-axis) and the inner slice represents the columns (x-axis).
	Values [][]float64
	// XAxis contains configuration options for the X-axis.
	XAxis HeatMapAxis
	// YAxis contains configuration options for the Y-axis.
	YAxis HeatMapAxis
	// ScaleMinValue overrides the minimum value for color gradient calculation. If nil, calculated from the data.
	ScaleMinValue *float64
	// ScaleMaxValue overrides the maximum value for color gradient calculation. If nil, calculated from the data.
	ScaleMaxValue *float64
	// ValuesLabel contains configuration for displaying numeric values on heat map cells.
	ValuesLabel SeriesLabel
}

// HeatMapAxis contains configuration options for an axis on a heat map chart.
type HeatMapAxis struct {
	// Title specifies the title to display next to the axis.
	Title string
	// TitleFontStyle specifies the font style for the axis title.
	TitleFontStyle FontStyle
	// Labels specifies custom labels for the axis. If empty, numeric indices are used.
	// Must match the data size for the corresponding axis.
	Labels []string
	// LabelFontStyle specifies the font style for the axis labels.
	LabelFontStyle FontStyle
	// LabelRotation is the rotation angle in radians for labels. Use DegreesToRadians(float64) to convert from degrees.
	LabelRotation float64
	// LabelCount is the number of labels to show on the axis. Use a smaller count to reduce text collisions.
	LabelCount int
	// LabelCountAdjustment specifies relative influence on label count.
	// Negative values result in cleaner graphs; positive values may cause text collisions.
	LabelCountAdjustment int
}

type heatMap struct {
	p   *Painter
	opt *HeatMapOption
}

// newHeatMapChart returns a heat map chart renderer.
func newHeatMapChart(p *Painter, opt HeatMapOption) *heatMap {
	return &heatMap{
		p:   p,
		opt: &opt,
	}
}

// NewHeatMapOptionWithData returns an initialized HeatMapOption with the provided data.
func NewHeatMapOptionWithData(data [][]float64) HeatMapOption {
	return HeatMapOption{
		Padding: defaultPadding,
		Values:  data,
	}
}

func (h *heatMap) renderChart(result *defaultRenderResult) (Box, error) {
	opt := h.opt
	if len(opt.Values) == 0 {
		return BoxZero, errors.New("empty values")
	}
	numRows := len(opt.Values)
	numCols := sliceMaxLen(opt.Values...)
	if numCols == 0 {
		return BoxZero, errors.New("heat map has no columns")
	}
	seriesPainter := result.seriesPainter.Child(PainterPaddingOption(NewBoxEqual(1)))

	// determine scale for map colors
	minVal, maxVal := computeMinMax(opt.Values, numCols)
	if opt.ScaleMinValue != nil {
		minVal = *opt.ScaleMinValue
	}
	if opt.ScaleMaxValue != nil {
		maxVal = *opt.ScaleMaxValue
	}
	valueRange := maxVal - minVal
	if math.Abs(valueRange) <= matrix.DefaultEpsilon {
		minVal = 0
		maxVal = 1
		valueRange = maxVal - minVal
	}

	baseColor := opt.Theme.GetSeriesColor(opt.BaseColorIndex)
	cellWidth := seriesPainter.Width() / numCols
	cellHeight := seriesPainter.Height() / numRows
	if cellWidth < 2 || cellHeight < 2 {
		return BoxZero, errors.New("insufficient space for heat map cells")
	}

	// Draw each cell, using the ratio to adjust the lightness of the base color.
	for y := range opt.Values {
		for x := 0; x < numCols; x++ {
			var value float64
			if x < len(opt.Values[y]) {
				value = opt.Values[y][x]
			}
			ratio := (value - minVal) / valueRange
			lightDelta := (1 - ratio) * 0.4
			satDelta := (1 - ratio) * 0.1
			if opt.Theme.IsDark() {
				lightDelta *= -1
			}
			cellColor := baseColor.WithAdjustHSL(0, satDelta, lightDelta)

			x1 := x * cellWidth
			y1 := y * cellHeight
			x2 := x1 + cellWidth
			y2 := y1 + cellHeight

			seriesPainter.FilledRect(x1, y1, x2, y2, cellColor, cellColor, 0)
		}
	}

	if flagIs(true, opt.ValuesLabel.Show) {
		opt.ValuesLabel.FontStyle =
			fillFontStyleDefaults(opt.ValuesLabel.FontStyle, defaultLabelFontSize, opt.Theme.GetLabelTextColor(), opt.Font)

		labelPainter := newSeriesLabelPainter(seriesPainter, []string{""}, opt.ValuesLabel, opt.Theme)
		for y := range opt.Values {
			for x := 0; x < numCols; x++ {
				var value float64
				if x < len(opt.Values[y]) {
					value = opt.Values[y][x]
				}
				xCenter := x*cellWidth + cellWidth/2
				yCenter := y*cellHeight + cellHeight/2
				labelPainter.Add(labelValue{
					index:     0,
					value:     value,
					x:         xCenter,
					y:         yCenter,
					fontStyle: opt.ValuesLabel.FontStyle,
				})
			}
		}
		if _, err := labelPainter.Render(); err != nil {
			return BoxZero, err
		}
	}

	return seriesPainter.box, nil
}

func computeMinMax(values [][]float64, numCol int) (float64, float64) {
	if len(values) == 0 || numCol == 0 {
		return 0, 0
	}

	var min, max float64
	if len(values[0]) != 0 {
		min = values[0][0]
		max = values[0][0]
	}
	for _, row := range values {
		rowMin, rowMax := chartdraw.MinMax(row...)
		if rowMin < min {
			min = rowMin
		}
		if rowMax > max {
			max = rowMax
		}
		if len(row) < numCol { // ensure range considers potential default values
			if min < 0 {
				min = 0
			}
			if max < 0 {
				max = 0
			}
		}
	}
	return min, max
}

func (h *heatMap) Render() (Box, error) {
	p := h.p
	opt := h.opt

	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}

	numRows := len(opt.Values)
	numCols := sliceMaxLen(opt.Values...)

	// Ensure X-axis labels cover all columns.
	for len(opt.XAxis.Labels) < numCols {
		opt.XAxis.Labels = append(opt.XAxis.Labels, strconv.Itoa(len(opt.XAxis.Labels)))
	}

	xAxisOption := XAxisOption{
		Title:                opt.XAxis.Title,
		TitleFontStyle:       opt.XAxis.TitleFontStyle,
		Labels:               opt.XAxis.Labels,
		LabelFontStyle:       opt.XAxis.LabelFontStyle,
		LabelRotation:        opt.XAxis.LabelRotation,
		LabelCount:           opt.XAxis.LabelCount,
		LabelCountAdjustment: opt.XAxis.LabelCountAdjustment,
	}

	// Ensure y-axis labels cover all rows.
	for len(opt.YAxis.Labels) < numRows {
		opt.YAxis.Labels = append(opt.YAxis.Labels, strconv.Itoa(len(opt.YAxis.Labels)))
	}
	yAxisOption := []YAxisOption{{
		Title:                  opt.YAxis.Title,
		TitleFontStyle:         opt.YAxis.TitleFontStyle,
		Labels:                 opt.YAxis.Labels,
		LabelFontStyle:         opt.YAxis.LabelFontStyle,
		LabelRotation:          opt.YAxis.LabelRotation,
		LabelCountAdjustment:   opt.YAxis.LabelCountAdjustment,
		LabelCount:             opt.YAxis.LabelCount,
		Min:                    Ptr(0.0),
		Max:                    Ptr(float64(numRows - 1)),
		RangeValuePaddingScale: Ptr(0.0),
		isCategoryAxis:         true,
	}}

	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:   opt.Theme,
		padding: opt.Padding,
		seriesList: heatMapFakeSeries{
			rows: numRows,
		},
		stackSeries: false,
		xAxis:       &xAxisOption,
		yAxis:       yAxisOption,
		title:       opt.Title,
		legend:      &LegendOption{Show: Ptr(false)},
	})
	if err != nil {
		return BoxZero, err
	}

	return h.renderChart(renderResult)
}

// heatMapFakeSeries is a dummy series type used solely to satisfy defaultRender's needs and notably drive axis rendering.
type heatMapFakeSeries struct {
	rows int
}

func (h heatMapFakeSeries) len() int {
	return 1
}

func (h heatMapFakeSeries) getSeries(_ int) series {
	return h
}

func (h heatMapFakeSeries) getSeriesName(_ int) string {
	return h.names()[0]
}

func (h heatMapFakeSeries) getSeriesValues(_ int) []float64 {
	return nil // not used, current usage is just in sumSeries, not used by defaultRender
}

func (h heatMapFakeSeries) getSeriesLen(_ int) int {
	return 0 // not used, current usage in getSeriesMaxDataCount, which is only used when axisReverse is true
}

func (h heatMapFakeSeries) names() []string {
	return []string{"Heat Map"}
}

func (h heatMapFakeSeries) hasMarkPoint() bool {
	return false
}

func (h heatMapFakeSeries) setSeriesName(_ int, _ string) {
	// ignored
}

func (h heatMapFakeSeries) sortByNameIndex(_ map[string]int) {
	// no-op
}

func (h heatMapFakeSeries) getSeriesSymbol(_ int) Symbol {
	return ""
}

func (h heatMapFakeSeries) getType() string {
	return ChartTypeHeatMap
}

func (h heatMapFakeSeries) getYAxisIndex() int {
	return 0
}

func (h heatMapFakeSeries) getValues() []float64 {
	return []float64{0, float64(h.rows)} // fake series data to get y-axis values set correctly
}

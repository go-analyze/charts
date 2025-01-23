package charts

import (
	"errors"
	"math"

	"github.com/golang/freetype/truetype"
)

type barChart struct {
	p   *Painter
	opt *BarChartOption
}

// newBarChart returns a bar chart renderer
func newBarChart(p *Painter, opt BarChartOption) *barChart {
	return &barChart{
		p:   p,
		opt: &opt,
	}
}

// NewBarChartOptionWithData returns an initialized BarChartOption with the SeriesList set for the provided data slice.
func NewBarChartOptionWithData(data [][]float64) BarChartOption {
	sl := NewSeriesListBar(data)
	return BarChartOption{
		SeriesList:     sl,
		Padding:        defaultPadding,
		Theme:          GetDefaultTheme(),
		Font:           GetDefaultFont(),
		YAxis:          make([]YAxisOption, sl.getYAxisCount()),
		ValueFormatter: defaultValueFormatter,
	}
}

type BarChartOption struct {
	// Theme specifies the colors used for the bar chart.
	Theme ColorPalette
	// Padding specifies the padding of bar chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
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
	// BarWidth specifies the width of each bar. Width may be reduced to ensure all series fit on the chart.
	BarWidth int
	// BarMargin specifies the margin between bars grouped together. BarWidth takes priority over the margin.
	BarMargin *float64
	// RoundedBarCaps set to `true` to produce a bar graph where the bars have rounded tops.
	RoundedBarCaps *bool
	// ValueFormatter defines how float values should be rendered to strings, notably for numeric axis labels.
	ValueFormatter ValueFormatter
}

func calculateBarMarginsAndSize(seriesCount, space int, configuredBarSize int, configuredBarMargin *float64) (int, int, int) {
	// default margins, adjusted below with config and series count
	margin := 10   // margin between each series block
	barMargin := 5 // margin between each bar
	if space < 20 {
		margin = 2
		barMargin = 2
	} else if space < 50 {
		margin = 5
		barMargin = 3
	}
	// check margin configuration if bar size allows margin
	if configuredBarSize+barMargin < space/seriesCount {
		// BarWidth is in range that we should also consider an optional margin configuration
		if configuredBarMargin != nil {
			barMargin = int(math.Round(*configuredBarMargin))
			if barMargin+configuredBarSize > space/seriesCount {
				barMargin = (space / seriesCount) - configuredBarSize
			}
		}
	} // else, bar width is out of range.  Ignore margin config

	barSize := (space - 2*margin - barMargin*(seriesCount-1)) / seriesCount
	// check bar size configuration, limited by the series count and space available
	if configuredBarSize > 0 && configuredBarSize < barSize {
		barSize = configuredBarSize
		// recalculate margin
		margin = (space - seriesCount*barSize - barMargin*(seriesCount-1)) / 2
	}

	return margin, barMargin, barSize
}

func (b *barChart) render(result *defaultRenderResult, seriesList SeriesList) (Box, error) {
	p := b.p
	opt := b.opt
	seriesPainter := result.seriesPainter

	xRange := newRange(b.p, getPreferredValueFormatter(opt.XAxis.ValueFormatter, opt.ValueFormatter),
		seriesPainter.Width(), len(opt.XAxis.Data), 0.0, 0.0, 0.0, 0.0)
	x0, x1 := xRange.GetRange(0)
	width := int(x1 - x0)
	seriesCount := len(seriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}
	margin, barMargin, barWidth := calculateBarMarginsAndSize(seriesCount, width, opt.BarWidth, opt.BarMargin)
	barMaxHeight := seriesPainter.Height()
	seriesNames := seriesList.Names()

	markPointPainter := newMarkPointPainter(seriesPainter)
	markLinePainter := newMarkLinePainter(seriesPainter)
	rendererList := []renderer{
		markPointPainter,
		markLinePainter,
	}
	for index := range seriesList {
		series := seriesList[index]
		yRange := result.axisRanges[series.YAxisIndex]
		seriesColor := opt.Theme.GetSeriesColor(series.index)

		divideValues := xRange.AutoDivide()
		points := make([]Point, len(series.Data))
		var labelPainter *seriesLabelPainter
		if flagIs(true, series.Label.Show) {
			labelPainter = newSeriesLabelPainter(seriesPainter, seriesNames, series.Label, opt.Theme, opt.Font)
			rendererList = append(rendererList, labelPainter)
		}

		for j, item := range series.Data {
			if j >= xRange.divideCount {
				continue
			}
			x := divideValues[j]
			x += margin
			if index != 0 {
				x += index * (barWidth + barMargin)
			}

			h := yRange.getHeight(item)
			top := barMaxHeight - h

			if flagIs(true, opt.RoundedBarCaps) {
				seriesPainter.roundedRect(Box{
					Top:    top,
					Left:   x,
					Right:  x + barWidth,
					Bottom: barMaxHeight - 1,
					IsSet:  true,
				}, barWidth, true, false,
					seriesColor, seriesColor, 0.0)
			} else {
				seriesPainter.FilledRect(x, top, x+barWidth, barMaxHeight-1,
					seriesColor, seriesColor, 0.0)
			}
			// generate marker point by hand
			points[j] = Point{
				X: x + (barWidth >> 1), // centered position
				Y: top,
			}
			// return if the label does not need to be displayed
			if labelPainter == nil {
				continue
			}
			y := barMaxHeight - h
			radians := float64(0)
			fontStyle := series.Label.FontStyle
			if series.Label.Position == PositionBottom {
				y = barMaxHeight
				radians = -math.Pi / 2
				if fontStyle.FontColor.IsZero() {
					if isLightColor(seriesColor) {
						fontStyle.FontColor = defaultLightFontColor
					} else {
						fontStyle.FontColor = defaultDarkFontColor
					}
				}
			}
			labelPainter.Add(labelValue{
				vertical:  true, // label is above bar
				index:     index,
				value:     item,
				fontStyle: fontStyle,
				x:         x + (barWidth >> 1),
				y:         y,
				radians:   radians,
				offset:    series.Label.Offset,
			})
		}

		markPointPainter.Add(markPointRenderOption{
			FillColor: seriesColor,
			Font:      opt.Font,
			Series:    series,
			Points:    points,
		})
		markLinePainter.Add(markLineRenderOption{
			FillColor:   seriesColor,
			FontColor:   opt.Theme.GetTextColor(),
			StrokeColor: seriesColor,
			Font:        opt.Font,
			Series:      series,
			Range:       yRange,
		})
	}
	// the largest and smallest mark point
	if err := doRender(rendererList...); err != nil {
		return BoxZero, err
	}

	return p.box, nil
}

func (b *barChart) Render() (Box, error) {
	p := b.p
	opt := b.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
	renderResult, err := defaultRender(p, defaultRenderOption{
		theme:          opt.Theme,
		padding:        opt.Padding,
		seriesList:     opt.SeriesList,
		xAxis:          &b.opt.XAxis,
		yAxis:          opt.YAxis,
		title:          opt.Title,
		legend:         &b.opt.Legend,
		valueFormatter: opt.ValueFormatter,
	})
	if err != nil {
		return BoxZero, err
	}
	seriesList := opt.SeriesList.Filter(ChartTypeBar)
	return b.render(renderResult, seriesList)
}

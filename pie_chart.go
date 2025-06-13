package charts

import (
	"errors"
	"fmt"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

const defaultPieRadiusFactor = 0.4

type pieChart struct {
	p   *Painter
	opt *PieChartOption
}

// NewPieChartOptionWithData returns an initialized PieChartOption with the SeriesList set from the provided data slice.
func NewPieChartOptionWithData(data []float64) PieChartOption {
	return PieChartOption{
		SeriesList: NewSeriesListPie(data),
		Padding:    defaultPadding,
		Theme:      GetDefaultTheme(),
		Font:       GetDefaultFont(),
	}
}

// PieChartOption defines the options for rendering a pie chart. Render the chart using Painter.PieChart.
type PieChartOption struct {
	// Theme specifies the colors used for the pie chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// Deprecated: Font is deprecated, instead the font needs to be set on the SeriesLabel, or other specific elements.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListPie.
	SeriesList PieSeriesList
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// Radius sets the default pie radius, for example "40%".
	// Default is "40%".
	Radius string
	// SegmentGap provides a gap between each pie slice.
	SegmentGap float64
	// ValueFormatter defines how float values should be rendered to strings, notably for series labels.
	ValueFormatter ValueFormatter
}

// newPieChart returns a pie chart renderer.
func newPieChart(p *Painter, opt PieChartOption) *pieChart {
	return &pieChart{
		p:   p,
		opt: &opt,
	}
}

type sector struct {
	value       float64
	percent     float64
	radius      float64
	startAngle  float64 // starting angle (radians)
	delta       float64 // sweep angle (radians)
	midAngle    float64
	quadrant    int  // 1: top-right, 2: top-left, 3: bottom-left, 4: bottom-right
	yCenter     bool // set to true if close to center in the y-axis
	label       string
	seriesLabel SeriesLabel
	color       Color
}

func newSector(radius float64, value, currentValue, totalValue float64,
	label string, seriesLabel SeriesLabel, altFormatter ValueFormatter, color Color) sector {
	s := sector{
		value:       value,
		percent:     value / totalValue,
		radius:      radius,
		startAngle:  chartdraw.PercentToRadians(currentValue/totalValue) - math.Pi/2,
		delta:       chartdraw.PercentToRadians(value / totalValue),
		seriesLabel: seriesLabel,
		color:       color,
	}
	s.midAngle = s.startAngle + s.delta/2

	// determine quadrant based on the mid-percentage of this sector
	p := (currentValue + value/2) / totalValue
	if p < 0.25 {
		s.quadrant = 1
	} else if p < 0.5 {
		s.quadrant = 4
	} else if p < 0.75 {
		s.quadrant = 3
	} else {
		s.quadrant = 2
	}
	s.yCenter = (p > .15 && p < .35) || (p > .65 && p < .85)

	if !flagIs(false, seriesLabel.Show) { // only set the label if it's being rendered
		valueFormatter := seriesLabel.ValueFormatter
		if valueFormatter == nil {
			valueFormatter = altFormatter
		}
		if valueFormatter != nil && seriesLabel.FormatTemplate == "" {
			s.label = valueFormatter(s.value)
		} else {
			s.label = labelFormatPie([]string{label}, seriesLabel.FormatTemplate, seriesLabel.ValueFormatter,
				0, s.value, s.percent)
		}
	}
	return s
}

// calculateOuterLabelLines computes the basic line positions for an outer label.
func (s *sector) calculateOuterLabelLines(cx, cy int, outerRadius, labelRadius float64, labelLineLength int) (lineStartX, lineStartY, lineBranchX, lineBranchY, lineEndX, lineEndY int) {
	lineStartX = cx + int(outerRadius*math.Cos(s.midAngle))
	lineStartY = cy + int(outerRadius*math.Sin(s.midAngle))
	lineBranchX = cx + int(labelRadius*math.Cos(s.midAngle))
	lineBranchY = cy + int(labelRadius*math.Sin(s.midAngle))
	offset := labelLineLength
	if lineBranchX <= cx {
		offset = -offset
	}
	lineEndX = lineBranchX + offset
	lineEndY = lineBranchY
	return
}

// calculateAdjustedOuterLabelPosition adds collision avoidance to the outer label positions.
func (s *sector) calculateAdjustedOuterLabelPosition(cx, cy int, outerRadius, labelRadius float64, labelLineLength, prevY int,
	labelFontSize float64, textBox Box) (lineStartX, lineStartY, lineBranchX, lineBranchY, lineEndX, lineEndY, textX, textY int) {
	lsX, lsY, lbX, lbY, leX, _ := s.calculateOuterLabelLines(cx, cy, outerRadius, labelRadius, labelLineLength)
	// adjust Y to avoid collisions
	threshold := ceilFloatToInt(labelFontSize) + 5
	adjustedBranchY := lbY
	if s.quadrant <= 2 {
		// quadrants in the top half
		for {
			if (prevY - adjustedBranchY) > threshold {
				break
			}
			adjustedBranchY--
		}
	} else {
		// quadrants in the bottom half
		for {
			if (adjustedBranchY - prevY) > threshold {
				break
			}
			adjustedBranchY++
		}
	}
	leY := adjustedBranchY

	// compute text position
	const textMargin = 3
	textX = leX + textMargin
	textY = leY + (textBox.Height() >> 1) - 1
	if leX <= cx {
		textX = leX - textBox.Width() - textMargin
	}
	return lsX, lsY, lbX, adjustedBranchY, leX, leY, textX, textY
}

func circleChartPosition(painter *Painter) (int, int, float64) {
	cx := painter.Width() >> 1
	cy := painter.Height() >> 1
	diameter := chartdraw.MinInt(painter.Width(), painter.Height())
	return cx, cy, float64(diameter)
}

func getFlexibleRadius(diameter, defaultRadiusFactor float64, radiusValue string) float64 {
	var radius float64
	if radiusValue != "" {
		radius, _ = parseFlexibleValue(radiusValue, diameter)
	}
	if radius <= 0 {
		radius = diameter * defaultRadiusFactor
	}
	return radius
}

func (p *pieChart) renderChart(result *defaultRenderResult) (Box, error) {
	opt := p.opt
	seriesPainter := result.seriesPainter
	cx, cy, diameter := circleChartPosition(seriesPainter)
	radius := getFlexibleRadius(diameter, defaultPieRadiusFactor, opt.Radius)
	var total float64
	for index, series := range opt.SeriesList {
		if opt.Radius == "" && series.Radius != "" {
			if seriesRadius := getFlexibleRadius(diameter, defaultPieRadiusFactor, series.Radius); seriesRadius > radius {
				radius = seriesRadius
			}
		}
		if series.Value < 0 {
			return BoxZero, fmt.Errorf("unsupported negative value for series index %d", index)
		}
		total += series.Value
	}

	_, err := renderPie(seriesPainter, cx, cy, diameter, radius, total, true, opt.SeriesList,
		opt.Theme, opt.SegmentGap, defaultPieRadiusFactor, opt.ValueFormatter, opt.Font)
	return p.p.box, err
}

func renderPie(p *Painter, cx, cy int, space, radius, total float64, renderLabels bool, seriesList PieSeriesList,
	theme ColorPalette, sliceGap, defaultRadiusFactor float64,
	valueFormatter ValueFormatter, fallbackFont *truetype.Font) ([]sector, error) {
	if len(seriesList) == 0 {
		return nil, errors.New("empty series list")
	} else if total <= 0 {
		return nil, errors.New("the sum value of pie chart should be greater than 0")
	}

	labelLineWidth := 15
	if radius < 50 {
		labelLineWidth = 5
	}
	// compute labelRadius for outer labels
	labelRadius := radius + float64(labelLineWidth)
	seriesNames := seriesList.names()

	var currentSum float64
	// organize sectors by quadrant
	var quadrant1, quadrant2, quadrant3, quadrant4 []sector
	for index, series := range seriesList {
		seriesRadius := radius
		if series.Radius != "" {
			seriesRadius = getFlexibleRadius(space, defaultRadiusFactor, series.Radius)
		}
		color := theme.GetSeriesColor(index)
		s := newSector(seriesRadius, series.Value, currentSum, total,
			seriesNames[index], series.Label, valueFormatter, color)

		switch s.quadrant {
		case 1:
			quadrant1 = append([]sector{s}, quadrant1...)
		case 2:
			quadrant2 = append(quadrant2, s)
		case 3:
			quadrant3 = append([]sector{s}, quadrant3...)
		case 4:
			quadrant4 = append(quadrant4, s)
		}
		currentSum += series.Value
	}
	sectors := append(append(append(quadrant1, quadrant4...), quadrant3...), quadrant2...)

	var currentQuadrant int
	var prevY, maxY, minY int
	for _, s := range sectors {
		// draw the pie slice
		p.moveTo(cx, cy)
		p.arcTo(cx, cy, s.radius, s.radius, s.startAngle, s.delta)
		p.lineTo(cx, cy)
		p.close()
		if sliceGap > 0 {
			p.fillStroke(s.color, theme.GetBackgroundColor(), sliceGap)
		} else {
			p.fill(s.color)
		}

		if !renderLabels || s.label == "" {
			continue
		}

		// initialize prevY for collision avoidance per quadrant
		if currentQuadrant != s.quadrant {
			currentQuadrant = s.quadrant
			if s.quadrant == 1 {
				minY = cy * 2
				maxY = 0
				prevY = cy * 2
			}
			if s.quadrant == 2 {
				prevY = minY
			}
			if s.quadrant == 3 {
				minY = cy * 2
				maxY = 0
				prevY = 0
			}
			if s.quadrant == 4 {
				prevY = maxY
			}
		}
		fontStyle := fillFontStyleDefaults(s.seriesLabel.FontStyle,
			defaultLabelFontSize, theme.GetLabelTextColor(), fallbackFont)
		textBox := p.MeasureText(s.label, 0, fontStyle)
		// for outer labels use the adjusted positions
		lsX, lsY, lbX, lbY, leX, leY, textX, textY :=
			s.calculateAdjustedOuterLabelPosition(cx, cy, s.radius, labelRadius, labelLineWidth, prevY, fontStyle.FontSize, textBox)
		prevY = leY
		if prevY > maxY {
			maxY = prevY
		}
		if prevY < minY {
			minY = prevY
		}
		p.moveTo(lsX, lsY)
		p.lineTo(lbX, lbY)
		p.moveTo(lbX, lbY)
		p.lineTo(leX, leY)
		p.stroke(s.color, 1)
		p.Text(s.label, textX, textY, 0, fontStyle)
	}
	return sectors, nil
}

func (p *pieChart) Render() (Box, error) {
	opt := p.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.p.theme)
	}
	if opt.Legend.Symbol == "" {
		opt.Legend.Symbol = SymbolSquare // default to square symbol for pie charts
	}

	renderResult, err := defaultRender(p.p, defaultRenderOption{
		theme:      opt.Theme,
		padding:    opt.Padding,
		seriesList: opt.SeriesList,
		xAxis: &XAxisOption{
			Show: Ptr(false),
		},
		yAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		title:  opt.Title,
		legend: &p.opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	return p.renderChart(renderResult)
}

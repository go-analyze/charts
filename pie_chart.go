package charts

import (
	"errors"
	"fmt"
	"math"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
)

type pieChart struct {
	p   *Painter
	opt *PieChartOption
}

// NewPieChartOptionWithData returns an initialized PieChartOption with the SeriesList set for the provided data slice.
func NewPieChartOptionWithData(data []float64) PieChartOption {
	return PieChartOption{
		SeriesList: NewSeriesListPie(data),
		Padding:    defaultPadding,
		Theme:      GetDefaultTheme(),
		Font:       GetDefaultFont(),
	}
}

type PieChartOption struct {
	// Theme specifies the colors used for the pie chart.
	Theme ColorPalette
	// Padding specifies the padding of pie chart.
	Padding Box
	// Font is the font used to render the chart.
	Font *truetype.Font
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListPie.
	SeriesList PieSeriesList
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// Radius default radius for pie e.g.: 40%, default is "40%"
	Radius string
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
	cx          int
	cy          int
	rx          float64
	ry          float64
	start       float64
	delta       float64
	offset      int
	quadrant    int
	lineStartX  int
	lineStartY  int
	lineBranchX int
	lineBranchY int
	lineEndX    int
	lineEndY    int
	label       string
	seriesLabel SeriesLabel
	color       Color
}

func newSector(cx int, cy int, radius float64, labelRadius float64,
	value float64, currentValue float64, totalValue float64,
	labelLineLength int, label string, seriesLabel SeriesLabel, altFormatter ValueFormatter, color Color) sector {
	s := sector{
		value:       value,
		percent:     value / totalValue,
		cx:          cx,
		cy:          cy,
		rx:          radius,
		ry:          radius,
		start:       chartdraw.PercentToRadians(currentValue/totalValue) - math.Pi/2,
		delta:       chartdraw.PercentToRadians(value / totalValue),
		offset:      labelLineLength,
		seriesLabel: seriesLabel,
		color:       color,
	}
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
	angle := s.start + s.delta/2
	s.lineStartX = cx + int(radius*math.Cos(angle))
	s.lineStartY = cy + int(radius*math.Sin(angle))
	s.lineBranchX = cx + int(labelRadius*math.Cos(angle))
	s.lineBranchY = cy + int(labelRadius*math.Sin(angle))
	if s.lineBranchX <= cx {
		s.offset *= -1
	}
	s.lineEndX = s.lineBranchX + s.offset
	s.lineEndY = s.lineBranchY
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

func (s *sector) calculateY(prevY int) int {
	for i := 0; i <= s.cy; i++ {
		if s.quadrant <= 2 {
			if (prevY - s.lineBranchY) > labelFontSize+5 {
				break
			}
			s.lineBranchY -= 1
		} else {
			if (s.lineBranchY - prevY) > labelFontSize+5 {
				break
			}
			s.lineBranchY += 1
		}
	}
	s.lineEndY = s.lineBranchY
	return s.lineBranchY
}

func (s *sector) calculateTextXY(textBox Box) (x int, y int) {
	textMargin := 3
	x = s.lineEndX + textMargin
	y = s.lineEndY + textBox.Height()>>1 - 1
	if s.offset < 0 {
		textWidth := textBox.Width()
		x = s.lineEndX - textWidth - textMargin
	}
	return
}

func (p *pieChart) render(result *defaultRenderResult, seriesList PieSeriesList) (Box, error) {
	opt := p.opt
	seriesCount := len(seriesList)
	if seriesCount == 0 {
		return BoxZero, errors.New("empty series list")
	}

	seriesPainter := result.seriesPainter
	cx := seriesPainter.Width() >> 1
	cy := seriesPainter.Height() >> 1
	diameter := chartdraw.MinInt(seriesPainter.Width(), seriesPainter.Height())
	radius := getRadius(float64(diameter), opt.Radius)
	var total float64
	for index, series := range seriesList {
		if opt.Radius == "" && series.Radius != "" {
			seriesRadius := getRadius(float64(diameter), series.Radius)
			if index == 0 || seriesRadius > radius {
				radius = seriesRadius
			}
		}
		if series.Value < 0 {
			return BoxZero, fmt.Errorf("unsupported negative value for series index %d", index)
		}
		total += series.Value
	}
	if total <= 0 {
		return BoxZero, errors.New("the sum value of pie chart should greater than 0")
	}

	labelLineWidth := 15
	if radius < 50 {
		labelLineWidth = 10
	}
	labelRadius := radius + float64(labelLineWidth)
	seriesNames := opt.Legend.SeriesNames
	if len(seriesNames) == 0 {
		seriesNames = seriesList.names()
	}
	theme := opt.Theme

	var currentValue float64
	var quadrant1, quadrant2, quadrant3, quadrant4 []sector
	seriesLen := len(seriesList)
	for index, series := range seriesList {
		seriesRadius := radius
		if series.Radius != "" {
			seriesRadius = getRadius(float64(diameter), series.Radius)
		}
		color := theme.GetSeriesColor(index)
		if index == seriesLen-1 {
			if color == theme.GetSeriesColor(0) {
				color = theme.GetSeriesColor(1)
			}
		}
		s := newSector(cx, cy, seriesRadius, labelRadius, series.Value, currentValue, total, labelLineWidth,
			seriesNames[index], series.Label, opt.ValueFormatter, color)
		switch quadrant := s.quadrant; quadrant {
		case 1:
			quadrant1 = append([]sector{s}, quadrant1...)
		case 2:
			quadrant2 = append(quadrant2, s)
		case 3:
			quadrant3 = append([]sector{s}, quadrant3...)
		case 4:
			quadrant4 = append(quadrant4, s)
		}
		currentValue += series.Value
	}
	sectors := append(quadrant1, quadrant4...)
	sectors = append(sectors, quadrant3...)
	sectors = append(sectors, quadrant2...)

	var currentQuadrant int
	var prevY, maxY, minY int
	for _, s := range sectors {
		seriesPainter.moveTo(s.cx, s.cy)
		seriesPainter.arcTo(s.cx, s.cy, s.rx, s.ry, s.start, s.delta)
		seriesPainter.lineTo(s.cx, s.cy)
		seriesPainter.close()
		seriesPainter.fillStroke(s.color, s.color, 1)
		if s.label == "" {
			continue
		}
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
		prevY = s.calculateY(prevY)
		if prevY > maxY {
			maxY = prevY
		}
		if prevY < minY {
			minY = prevY
		}
		seriesPainter.moveTo(s.lineStartX, s.lineStartY)
		seriesPainter.lineTo(s.lineBranchX, s.lineBranchY)
		seriesPainter.moveTo(s.lineBranchX, s.lineBranchY)
		seriesPainter.lineTo(s.lineEndX, s.lineEndY)
		seriesPainter.stroke(s.color, 1)
		textStyle := FontStyle{
			FontColor: theme.GetLabelTextColor(),
			FontSize:  labelFontSize,
			Font:      opt.Font,
		}
		if !s.seriesLabel.FontStyle.FontColor.IsZero() {
			textStyle.FontColor = s.seriesLabel.FontStyle.FontColor
		}
		if s.seriesLabel.FontStyle.FontSize > 0 {
			textStyle.FontSize = s.seriesLabel.FontStyle.FontSize
		}
		x, y := s.calculateTextXY(seriesPainter.MeasureText(s.label, 0, textStyle))
		seriesPainter.Text(s.label, x, y, 0, textStyle)
	}
	return p.p.box, nil
}

func (p *pieChart) Render() (Box, error) {
	opt := p.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.p.theme)
	}
	if opt.Legend.Symbol == "" {
		// default to square symbol for this chart type
		opt.Legend.Symbol = SymbolSquare
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
	return p.render(renderResult, opt.SeriesList)
}

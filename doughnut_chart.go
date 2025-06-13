package charts

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

type doughnutChart struct {
	p   *Painter
	opt *DoughnutChartOption
}

// NewDoughnutChartOptionWithData returns an initialized DoughnutChartOption with the SeriesList set from the provided data slice.
func NewDoughnutChartOptionWithData(data []float64) DoughnutChartOption {
	return DoughnutChartOption{
		SeriesList: NewSeriesListDoughnut(data),
		Padding:    defaultPadding,
		Theme:      GetDefaultTheme(),
	}
}

// DoughnutChartOption defines the options for rendering a doughnut chart. Render the chart using Painter.DoughnutChart.
type DoughnutChartOption struct {
	// Theme specifies the colors used for the doughnut chart.
	Theme ColorPalette
	// Padding specifies the padding around the chart.
	Padding Box
	// SeriesList provides the data population for the chart, typically constructed using NewSeriesListDoughnut.
	SeriesList DoughnutSeriesList
	// Title are options for rendering the title.
	Title TitleOption
	// Legend are options for the data legend.
	Legend LegendOption
	// RadiusRing sets the outer radius of the ring, for example "40%".
	// Default is "40%".
	RadiusRing string
	// RadiusCenter is the radius for the center hole of the doughnut and must be smaller than RadiusRing.
	RadiusCenter string
	// CenterValues specifies what should be rendered in the center of the doughnut,
	// current options are "none" (default), "labels", "sum".
	// * labels - Will render the labels on the inside of the circle instead of the outside (more risk of collision).
	// * sum - Will put the sum count of all the series (formatted using ValueFormatter).
	CenterValues string
	// CenterValuesFontStyle provides the styling for center values (series labels prefer their specific series styling).
	CenterValuesFontStyle FontStyle
	// SegmentGap provides a margin between each series section.
	SegmentGap float64
	// ValueFormatter defines how float values should be rendered to strings, notably for series labels.
	ValueFormatter ValueFormatter
}

// newDoughnutChart creates a new doughnut chart renderer.
func newDoughnutChart(p *Painter, opt DoughnutChartOption) *doughnutChart {
	return &doughnutChart{
		p:   p,
		opt: &opt,
	}
}

func (d *doughnutChart) Render() (Box, error) {
	opt := d.opt
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(d.p.theme)
	}
	if opt.Legend.Symbol == "" {
		opt.Legend.Symbol = SymbolSquare // default symbol for doughnut charts
	}

	renderResult, err := defaultRender(d.p, defaultRenderOption{
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
		legend: &d.opt.Legend,
	})
	if err != nil {
		return BoxZero, err
	}
	return d.renderChart(renderResult)
}

func (d *doughnutChart) renderChart(result *defaultRenderResult) (Box, error) {
	opt := d.opt
	seriesPainter := result.seriesPainter
	centerLabels := strings.EqualFold(opt.CenterValues, "labels")
	cx, cy, diameter := circleChartPosition(seriesPainter)
	radiusFactorDefault := defaultPieRadiusFactor
	if centerLabels {
		radiusFactorDefault = 0.5 // larger radius to help fit center labels
	}
	radiusRing := getFlexibleRadius(diameter, radiusFactorDefault, opt.RadiusRing)
	minRadius := radiusRing
	var total float64
	for index, series := range opt.SeriesList {
		if opt.RadiusRing == "" && series.Radius != "" {
			seriesRadius := getFlexibleRadius(diameter, radiusFactorDefault, series.Radius)
			if seriesRadius > radiusRing {
				radiusRing = seriesRadius
			} else if seriesRadius < minRadius {
				minRadius = seriesRadius
			}
		}
		if series.Value < 0 {
			return BoxZero, fmt.Errorf("unsupported negative value at series index %d", index)
		}
		total += series.Value
	}
	radiusCenter := minRadius * 0.6
	if opt.RadiusCenter != "" {
		var err error
		radiusCenter, err = parseFlexibleValue(opt.RadiusCenter, diameter)
		if err != nil {
			return BoxZero, fmt.Errorf("invalid RadiusCenter: %w", err)
		}
		if radiusCenter > radiusRing-10 {
			radiusCenter = radiusRing - 10 // set to maximum value
		}
	}

	sectors, err := renderPie(seriesPainter, cx, cy, diameter, radiusRing, total, !centerLabels,
		opt.SeriesList.toPieSeriesList(), opt.Theme, opt.SegmentGap, radiusFactorDefault, opt.ValueFormatter, nil)
	if err != nil {
		msg := strings.ReplaceAll(err.Error(), "pie", "doughnut")
		msg = strings.ReplaceAll(msg, "Pie", "Doughnut")
		return BoxZero, errors.New(msg)
	}

	// Draw doughnut center / hole
	circleColor := opt.Theme.GetBackgroundColor()
	if circleColor.IsZero() {
		circleColor = ColorWhite
	} else if circleColor.A != 255 {
		circleColor = circleColor.WithAlpha(255)
	}
	seriesPainter.Circle(radiusCenter, cx, cy, circleColor, circleColor, 0.0)

	if centerLabels {
		placements := placeCenterLabelsWithCollisionResolution(seriesPainter, opt, cx, cy, radiusCenter, sectors)

		for _, lp := range placements {
			s := lp.sector

			// choose label connection point based on quadrant for a shorter line
			var connectX, connectY int
			if s.quadrant == 2 || s.quadrant == 3 {
				connectX = lp.box.Left
			} else {
				connectX = lp.box.Right
			}
			if s.yCenter {
				connectY = lp.box.CenterY()
			} else if s.quadrant == 1 || s.quadrant == 2 {
				connectY = lp.box.Top
			} else {
				connectY = lp.box.Bottom
			}

			// compute the donut boundary connection point
			lineStartX, lineStartY := s.connectionPoint(cx, cy, radiusCenter, connectX, connectY)

			// draw a line from the donut boundary to the label connection point
			seriesPainter.moveTo(lineStartX, lineStartY)
			seriesPainter.lineTo(connectX, connectY)
			seriesPainter.stroke(s.color, defaultStrokeWidth)

			// finally, render the label text at its resolved position
			fontStyle := fillFontStyleDefaults(mergeFontStyles(s.seriesLabel.FontStyle, opt.CenterValuesFontStyle),
				defaultLabelFontSize, opt.Theme.GetLabelTextColor())
			seriesPainter.Text(s.label, lp.box.Left, lp.box.Bottom, 0, fontStyle)
		}
	} else if strings.EqualFold(opt.CenterValues, "sum") {
		opt.CenterValuesFontStyle = fillFontStyleDefaults(opt.CenterValuesFontStyle,
			defaultLabelFontSize, opt.Theme.GetLabelTextColor())
		valueFormatter := getPreferredValueFormatter(opt.ValueFormatter)
		sumStr := valueFormatter(total)
		centerTextBox := d.p.MeasureText(sumStr, 0, opt.CenterValuesFontStyle)

		seriesPainter.Text(sumStr, cx-(centerTextBox.Width()>>1), cy+(centerTextBox.Height()>>1),
			0, opt.CenterValuesFontStyle)
	}
	return d.p.box, nil
}

const (
	maxNudgeIterations     = 20
	nudgeAngleRange        = 0.24
	nudgeAngleStep         = 0.02
	innerLabelRadiusFactor = 0.8
)

type labelPlacement struct {
	sector       *sector
	labelMeasure Box
	box          Box // final coordinates for label rendering
}

// placeCenterLabelsWithCollisionResolution places labels inside doughnut, avoiding overlaps.
func placeCenterLabelsWithCollisionResolution(p *Painter, opt *DoughnutChartOption,
	cx, cy int, radiusCenter float64, sectors []sector) []*labelPlacement {
	// ensure sorted by ascending midAngle
	sort.Slice(sectors, func(i, j int) bool {
		return sectors[i].midAngle < sectors[j].midAngle
	})
	anchorRad := innerLabelRadiusFactor * radiusCenter

	placements := make([]*labelPlacement, 0, len(sectors))
	for i, s := range sectors {
		if s.label == "" {
			continue
		}
		fontStyle := fillFontStyleDefaults(mergeFontStyles(s.seriesLabel.FontStyle, opt.CenterValuesFontStyle),
			defaultLabelFontSize, opt.Theme.GetLabelTextColor())
		measured := p.MeasureText(s.label, 0, fontStyle)
		midAngle := s.midAngle

		// build candidate angles (midAngle ± small increments)
		candidateAngles := []float64{midAngle}
		for off := nudgeAngleStep; off <= nudgeAngleRange; off += nudgeAngleStep {
			candidateAngles = append(candidateAngles, midAngle+off, midAngle-off)
		}

		var chosenBox Box
		var placedOK bool
		for _, a := range candidateAngles {
			candidateBox := computeLabelBox(cx, cy, anchorRad, a, measured)
			if isInsideCircle(candidateBox, cx, cy, radiusCenter) &&
				!anyLabelCollision(candidateBox, placements) {
				chosenBox = candidateBox
				placedOK = true
				break
			}
		}
		if !placedOK { // fallback if no valid placement found
			chosenBox = computeLabelBox(cx, cy, anchorRad, midAngle, measured)
		}

		lp := &labelPlacement{
			sector:       &sectors[i],
			labelMeasure: measured,
			box:          chosenBox,
		}

		// a proactive shift towards the series sector keeps the labels more distinctly on their side
		shiftLabelHorizontallyTowardSector(lp, cx, lp.sector.midAngle)

		placements = append(placements, lp)
	}

	// iteratively nudge colliding labels apart while keeping them inside the center circle.
	for iteration := 0; iteration < maxNudgeIterations; iteration++ {
		collided := false
		for i := 0; i < len(placements); i++ {
			for j := i + 1; j < len(placements); j++ {
				if placements[i].box.Overlaps(placements[j].box) {
					collided = true
					// choose which label to move
					moveLP := placements[i]
					stayLP := placements[j]
					if placements[j].sector.value > placements[i].sector.value {
						moveLP, stayLP = placements[j], placements[i]
					}

					dx, dy := minimalRadialPush(moveLP, stayLP, cx, cy)
					moveLP.box = moveLP.box.Shift(dx, dy)

					// clamp the moved label inside the circle.
					clampInsideCircle(moveLP, cx, cy, radiusCenter)
				}
			}
		}
		if !collided {
			break
		}
	}

	return placements
}

// computeLabelBox computes and returns the bounding box for a label.
// Uses center (cx, cy), an anchor radius, angle (in radians), and text size.
// Aligns the box left or right based on the anchor's position relative to the center.
func computeLabelBox(cx, cy int, anchorRad, angle float64, textSize Box) Box {
	// compute anchor point on the circle
	ax := float64(cx) + anchorRad*math.Cos(angle)
	ay := float64(cy) + anchorRad*math.Sin(angle)
	w, h := textSize.Width(), textSize.Height()

	// align box left or right based on anchor position
	var left int
	if ax > float64(cx) {
		left = int(ax) - w
	} else {
		left = int(ax)
	}
	top := int(ay) - h
	return NewBox(left, top, left+w, top+h)
}

// isInsideCircle returns true if the provided box (set with absolute coordinates) is within
// the circle radius at the provided position.
func isInsideCircle(b Box, cx, cy int, r float64) bool {
	for _, pt := range b.Corners().ToPoints() {
		dx := float64(pt.X - cx)
		dy := float64(pt.Y - cy)
		if (dx*dx + dy*dy) > r*r {
			return false
		}
	}
	return true
}

// anyLabelCollision returns true if box b overlaps any already placed label boxes.
func anyLabelCollision(b Box, placed []*labelPlacement) bool {
	for _, p := range placed {
		if b.Overlaps(p.box) {
			return true
		}
	}
	return false
}

// shiftLabelHorizontallyTowardSector moves the label horizontally so that it’s on the same side as the sector.
func shiftLabelHorizontallyTowardSector(lp *labelPlacement, cx int, angle float64) {
	// the sectors “horizontal side” is determined by cos(angle)
	sectorSide := math.Cos(angle)

	// compute label center relative to chart center
	lcx := float64(lp.box.CenterX() - cx)

	// shift label if its center is on the wrong side
	if sectorSide > 0 && lcx < 0 {
		// shift center to the right
		shift := int(math.Abs(lcx)) + 1 // small offset to cross over
		lp.box = lp.box.Shift(shift, 0)
	} else if sectorSide < 0 && lcx > 0 {
		// shift center to the left
		shift := -int(math.Abs(lcx)) - 1
		lp.box = lp.box.Shift(shift, 0)
	}
}

// minimalRadialPush returns a small displacement vector that nudges 'p' away from 'q'
// along the radial direction (ideally outward, if p is farther along the radial axis).
func minimalRadialPush(p, q *labelPlacement, cx, cy int) (dx, dy int) {
	pcx := float64(p.box.CenterX())
	pcy := float64(p.box.CenterY())
	rx := pcx - float64(cx)
	ry := pcy - float64(cy)
	rLen := math.Hypot(rx, ry)
	if rLen == 0 {
		return 1, 0 // fallback nudge
	}
	ux, uy := rx/rLen, ry/rLen

	// measure overlap along that radial vector vs q
	pMin, pMax := projectBoxRadially(p.box, ux, uy)
	qMin, qMax := projectBoxRadially(q.box, ux, uy)
	overlapDist := math.Min(pMax, qMax) - math.Max(pMin, qMin)
	if overlapDist <= 0 {
		return int(math.Copysign(1, ux)), int(math.Copysign(1, uy)) // fallback nudge to ensure we don't stall nudging
	}
	// nudge p outward by overlap distance plus a small margin
	dist := overlapDist + 2
	dxF := ux * dist
	dyF := uy * dist
	return int(dxF + math.Copysign(0.5, dxF)), int(dyF + math.Copysign(0.5, dyF))
}

// projectBoxRadially projects the corners of the box onto (ux, uy), returning min and max.
func projectBoxRadially(b Box, ux, uy float64) (float64, float64) {
	minV, maxV := math.Inf(1), math.Inf(-1)
	for _, pt := range b.Corners().ToPoints() {
		dot := float64(pt.X)*ux + float64(pt.Y)*uy
		if dot < minV {
			minV = dot
		}
		if dot > maxV {
			maxV = dot
		}
	}
	return minV, maxV
}

// clampInsideCircle adjusts a label's position to keep its box inside the circle centered at (cx,cy)
// with radius r. It first attempts a horizontal shift if near the top/bottom, then nudges radially.
func clampInsideCircle(lp *labelPlacement, cx, cy int, r float64) {
	for i := 0; i < maxNudgeIterations; i++ {
		// if label center is near top or bottom, adjust horizontally towards center first
		if midY := lp.box.CenterY(); math.Abs(float64(midY-cy)) > 0.6*r {
			// Compute allowed horizontal offset based on the circle equation:
			// at a given y, x must be within [cx - sqrt(r² - (y-cy)²), cx + sqrt(r² - (y-cy)²)]
			dy := float64(midY - cy)
			if math.Abs(dy) > r {
				return // cannot clamp horizontally if outside vertical bounds
			}
			allowedXOffset := math.Sqrt(r*r - dy*dy)
			allowedLeft := cx - int(allowedXOffset)
			allowedRight := cx + int(allowedXOffset)

			// shift label if its edges exceed allowed horizontal bounds
			if lp.box.Left < allowedLeft {
				dx := allowedLeft - lp.box.Left
				lp.box = lp.box.Shift(dx, 0)
			} else if lp.box.Right > allowedRight {
				dx := allowedRight - lp.box.Right
				lp.box = lp.box.Shift(dx, 0)
			}
		}

		if isInsideCircle(lp.box, cx, cy, r) {
			break
		}
		// nudge the label radially toward the center
		dx := cx - lp.box.CenterX()
		dy := cy - lp.box.CenterY()
		stepLen := math.Hypot(float64(dx), float64(dy))
		if stepLen < 1 {
			break
		}
		ux, uy := float64(dx)/stepLen, float64(dy)/stepLen
		shiftX, shiftY := int(ux+math.Copysign(0.5, ux)), int(uy+math.Copysign(0.5, uy))
		if shiftX == 0 && shiftY == 0 {
			break
		}
		lp.box = lp.box.Shift(shiftX, shiftY)
	}
}

// connectionPoint returns a point along the sector arc (at radius r) that is closest to the labels center.
// If the angle to the label falls outside the sector, it clamps to the nearest boundary.
func (s *sector) connectionPoint(cx, cy int, r float64, labelX, labelY int) (int, int) {
	angleToLabel := math.Atan2(float64(labelY-cy), float64(labelX-cx))
	chosenAngle := clampAngleToSector(angleToLabel, s.startAngle, s.startAngle+s.delta)
	x := cx + int(r*math.Cos(chosenAngle))
	y := cy + int(r*math.Sin(chosenAngle))
	return x, y
}

// clampAngleToSector clamps an `angle` to the interval [sectorStart, sectorEnd].
func clampAngleToSector(angle, sectorStart, sectorEnd float64) float64 {
	angle = normalizeAngle(angle)
	sectorStart = normalizeAngle(sectorStart)
	sectorEnd = normalizeAngle(sectorEnd)

	// if sector does not wrap around 0, perform a simple clamp
	if sectorStart <= sectorEnd {
		if angle < sectorStart {
			return sectorStart
		} else if angle > sectorEnd {
			return sectorEnd
		}
		return angle
	}

	// for sectors that cross 0, treat as [sectorStart, 2π) U [0, sectorEnd]
	if angle > sectorStart || angle < sectorEnd {
		return angle
	}
	// clamp to the closer boundary
	if math.Abs(angle-sectorStart) > math.Abs(angle-sectorEnd) {
		return sectorEnd // start is greater, end is closer
	}
	return sectorStart
}

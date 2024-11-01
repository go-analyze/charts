package charts

import (
	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

type Box = chartdraw.Box
type Color = drawing.Color

var BoxZero = chartdraw.BoxZero

// Offset provides an ability to configure a shift from the top or left alignments.
type Offset struct {
	// Left indicates a vertical spacing adjustment from the top.
	Top int
	// Left indicates a horizontal spacing adjustment from the left.
	Left int
}

type Point struct {
	X int
	Y int
}

const (
	ChartTypeLine          = "line"
	ChartTypeBar           = "bar"
	ChartTypePie           = "pie"
	ChartTypeRadar         = "radar"
	ChartTypeFunnel        = "funnel"
	ChartTypeHorizontalBar = "horizontalBar"
)

const (
	ChartOutputSVG           = "svg"
	ChartOutputPNG           = "png"
	chartDefaultOutputFormat = ChartOutputPNG
)

const (
	PositionLeft   = "left"
	PositionRight  = "right"
	PositionCenter = "center"
	PositionTop    = "top"
	PositionBottom = "bottom"
)

const (
	AlignLeft   = "left"
	AlignRight  = "right"
	AlignCenter = "center"
)

const (
	OrientHorizontal = "horizontal"
	OrientVertical   = "vertical"
)

package charts

import (
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type Box = chart.Box
type Style = chart.Style
type Color = drawing.Color

var BoxZero = chart.BoxZero

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

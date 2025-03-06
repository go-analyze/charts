package charts

import (
	"strconv"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

type Box = chartdraw.Box
type Point = chartdraw.Point
type Color = drawing.Color
type FontStyle = chartdraw.FontStyle

var BoxZero = chartdraw.BoxZero

// NewBox returns a new box with the provided left, top, right, and bottom sizes.
func NewBox(left, top, right, bottom int) Box {
	return Box{
		IsSet:  true,
		Top:    top,
		Bottom: bottom,
		Left:   left,
		Right:  right,
	}
}

// NewBoxEqual returns a new box with equal sizes to each side.
func NewBoxEqual(size int) Box {
	return NewBox(size, size, size, size)
}

// NewFontStyleWithSize constructs a new FontStyle with the specified font size. If you want to avoid directly
// constructing the FontStyle struct, you can use this followed by additional `WithX` function calls on the returned
// FontStyle.
func NewFontStyleWithSize(size float64) FontStyle {
	return FontStyle{
		FontSize: size,
	}
}

func fillFontStyleDefaults(fs FontStyle, defaultSize float64, defaultColor Color, fontOptions ...*truetype.Font) FontStyle {
	if fs.FontSize == 0 {
		fs.FontSize = defaultSize
	}
	if fs.FontColor.IsZero() {
		fs.FontColor = defaultColor
	}
	if fs.Font == nil {
		fs.Font = getPreferredFont(fontOptions...)
	}
	return fs
}

// OffsetInt provides an ability to configure a shift from the top or left alignments.
type OffsetInt struct {
	// Left indicates a vertical spacing adjustment from the top.
	Top int
	// Left indicates a horizontal spacing adjustment from the left.
	Left int
}

func (o OffsetInt) WithTop(val int) OffsetInt {
	return OffsetInt{
		Left: o.Left,
		Top:  val,
	}
}

func (o OffsetInt) WithLeft(val int) OffsetInt {
	return OffsetInt{
		Left: val,
		Top:  o.Top,
	}
}

// OffsetStr provides an ability to configure a shift from the top or left alignments using flexible string inputs.
type OffsetStr struct {
	// Left is the distance between the component and the left side of the container.
	// It can be pixel value (20), percentage value (20%), or position description: 'left', 'right', 'center'.
	Left string
	// Top is the distance between the component and the top side of the container.
	// It can be pixel value (20), or percentage value (20%), or position description: 'top', 'bottom'.
	Top string
}

var OffsetLeft = OffsetStr{Left: PositionLeft}
var OffsetRight = OffsetStr{Left: PositionRight}
var OffsetCenter = OffsetStr{Left: PositionCenter}

func (o OffsetStr) WithTop(val string) OffsetStr {
	return OffsetStr{
		Left: o.Left,
		Top:  val,
	}
}

func (o OffsetStr) WithTopI(val int) OffsetStr {
	return OffsetStr{
		Left: o.Left,
		Top:  strconv.Itoa(val),
	}
}

func (o OffsetStr) WithLeft(val string) OffsetStr {
	return OffsetStr{
		Left: val,
		Top:  o.Top,
	}
}

func (o OffsetStr) WithLeftI(val int) OffsetStr {
	return OffsetStr{
		Left: strconv.Itoa(val),
		Top:  o.Top,
	}
}

const (
	ChartTypeLine          = "line"
	ChartTypeScatter       = "scatter"
	ChartTypeBar           = "bar"
	ChartTypePie           = "pie"
	ChartTypeRadar         = "radar"
	ChartTypeFunnel        = "funnel"
	ChartTypeHorizontalBar = "horizontalBar"
	ChartTypeHeatMap       = "heatMap"
)

const (
	ChartOutputSVG           = "svg"
	ChartOutputPNG           = "png"
	ChartOutputJPG           = "jpg"
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

type Symbol string

const (
	SymbolNone    = "none"
	SymbolCircle  = "circle"
	SymbolDot     = "dot"
	SymbolSquare  = "square"
	SymbolDiamond = "diamond"
)

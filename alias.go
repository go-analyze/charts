package charts

import (
	"strconv"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

// Box defines spacing boundaries around a component.
type Box = chartdraw.Box

// Point represents an X,Y coordinate pair.
type Point = chartdraw.Point

// Color represents an RGBA color.
type Color = drawing.Color

// FontStyle configures font properties including size, color, and family.
type FontStyle = chartdraw.FontStyle

// BoxZero is an unset Box with no dimensions.
var BoxZero = chartdraw.BoxZero

// NewBox returns a new Box with the specified left, top, right, and bottom values to define the position and dimensions.
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

// fillFontStyleDefaults returns a FontStyle with the given size, color, and a font (if not already set).
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

// mergeFontStyles sets from the default FontStyles the size, color, and font as
// provided by the default styles (in order).
func mergeFontStyles(primary FontStyle, defaultFs ...FontStyle) FontStyle {
	if primary.FontSize == 0 {
		for _, fs := range defaultFs {
			if fs.FontSize != 0 {
				primary.FontSize = fs.FontSize
				break
			}
		}
	}
	if primary.FontColor.IsZero() {
		for _, fs := range defaultFs {
			if !fs.FontColor.IsZero() {
				primary.FontColor = fs.FontColor
				break
			}
		}
	}
	if primary.Font == nil {
		for _, fs := range defaultFs {
			if fs.Font != nil {
				primary.Font = fs.Font
				break
			}
		}
	}
	return primary
}

// OffsetInt provides an ability to configure a shift from the top or left alignments.
type OffsetInt struct {
	// Top indicates a vertical spacing adjustment from the top.
	Top int
	// Left indicates a horizontal spacing adjustment from the left.
	Left int
}

// WithTop returns a copy of the offset with the Top value set.
func (o OffsetInt) WithTop(val int) OffsetInt {
	return OffsetInt{
		Left: o.Left,
		Top:  val,
	}
}

// WithLeft returns a copy of the offset with the Left value set.
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

// OffsetLeft positions a component on the left.
var OffsetLeft = OffsetStr{Left: PositionLeft}

// OffsetRight positions a component on the right.
var OffsetRight = OffsetStr{Left: PositionRight}

// OffsetCenter positions a component in the center.
var OffsetCenter = OffsetStr{Left: PositionCenter}

// WithTop returns a copy of the offset with the Top value set.
func (o OffsetStr) WithTop(val string) OffsetStr {
	return OffsetStr{
		Left: o.Left,
		Top:  val,
	}
}

// WithTopI sets Top using an integer value.
func (o OffsetStr) WithTopI(val int) OffsetStr {
	return OffsetStr{
		Left: o.Left,
		Top:  strconv.Itoa(val),
	}
}

// WithLeft returns a copy of the offset with the Left value set.
func (o OffsetStr) WithLeft(val string) OffsetStr {
	return OffsetStr{
		Left: val,
		Top:  o.Top,
	}
}

// WithLeftI sets Left using an integer value.
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
	ChartTypeDoughnut      = "doughnut"
	ChartTypeRadar         = "radar"
	ChartTypeFunnel        = "funnel"
	ChartTypeHorizontalBar = "horizontalBar"
	ChartTypeHeatMap       = "heatMap"
	ChartTypeCandlestick   = "candlestick"
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

// Symbol defines the shape used for data points and legends.
type Symbol string

const (
	SymbolNone        = "none"
	SymbolCircle      = "circle"
	SymbolDot         = "dot"
	SymbolSquare      = "square"
	SymbolDiamond     = "diamond"
	symbolCandlestick = "candlestick" // internal only, set automatically
)

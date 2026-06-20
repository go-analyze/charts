package charts

// SymbolShape defines the shape used for data points and legend icons.
type SymbolShape string

const (
	SymbolNone        SymbolShape = "none"
	SymbolCircle      SymbolShape = "circle"
	SymbolDot         SymbolShape = "dot"
	SymbolSquare      SymbolShape = "square"
	SymbolDiamond     SymbolShape = "diamond"
	symbolCandlestick SymbolShape = "candlestick" // internal only, set automatically
)

// Symbol configures the shape and size drawn at data points and legend icons.
type Symbol struct {
	// Shape is the symbol shape. Empty selects a chart-type default.
	Shape SymbolShape
	// Size is the symbol radius in pixels. 0 selects a chart-type default
	// (scatter: 2.0; line: derived from LineStrokeWidth).
	Size float64
}

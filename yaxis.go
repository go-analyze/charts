package charts

import (
	"github.com/golang/freetype/truetype"
)

type YAxisOption struct {
	// The minimun value of axis.
	Min *float64
	// The maximum value of axis.
	Max *float64
	// RangeValuePaddingScale suggest a scale of padding added to the max and min values
	RangeValuePaddingScale *float64
	// The font of y-axis
	Font *truetype.Font
	// The data value of y-axis
	Data []string
	// The theme of chart
	Theme ColorPalette
	// The font size of y-axis label
	FontSize float64
	// The position of axis, it can be 'left' or 'right'
	Position string
	// The color of label
	FontColor Color
	// Formatter for y-axis text value
	Formatter string
	// Color for y-axis
	Color Color
	// The flag for show axis, set this to *false will hide axis
	Show *bool
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis.  Specify a smaller number to reduce writing collisions.
	LabelCount int
	// LabelSkipCount specifies a number of lines between labels where there will be no label and instead just a horizontal line.
	LabelSkipCount int
	isCategoryAxis bool
}

// NewYAxisOptions returns a y-axis option
func NewYAxisOptions(data []string, others ...[]string) []YAxisOption {
	arr := [][]string{
		data,
	}
	arr = append(arr, others...)
	opts := make([]YAxisOption, 0)
	for _, data := range arr {
		opts = append(opts, YAxisOption{
			Data: data,
		})
	}
	return opts
}

func (opt *YAxisOption) ToAxisOption(p *Painter) AxisOption {
	position := PositionLeft
	if opt.Position == PositionRight {
		position = PositionRight
	}
	theme := opt.Theme
	if theme == nil {
		theme = p.theme
	}
	axisOpt := AxisOption{
		Formatter:      opt.Formatter,
		Theme:          theme,
		Data:           opt.Data,
		Position:       position,
		FontSize:       opt.FontSize,
		StrokeWidth:    -1,
		Font:           opt.Font,
		FontColor:      opt.FontColor,
		BoundaryGap:    FalseFlag(),
		Unit:           opt.Unit,
		LabelCount:     opt.LabelCount,
		LabelSkipCount: opt.LabelSkipCount,
		SplitLineShow:  true,
		SplitLineColor: theme.GetAxisSplitLineColor(),
		Show:           opt.Show,
	}
	if !opt.Color.IsZero() {
		axisOpt.FontColor = opt.Color
		axisOpt.StrokeColor = opt.Color
	}
	if opt.isCategoryAxis {
		axisOpt.BoundaryGap = TrueFlag()
		axisOpt.StrokeWidth = 1
		axisOpt.SplitLineShow = false
	}
	return axisOpt
}

// NewLeftYAxis returns a left y axis renderer
func NewLeftYAxis(p *Painter, opt YAxisOption) *axisPainter {
	p = p.Child(PainterPaddingOption(Box{
		Bottom: defaultXAxisHeight,
	}))
	return NewAxisPainter(p, opt.ToAxisOption(p))
}

// NewRightYAxis returns a right y axis renderer
func NewRightYAxis(p *Painter, opt YAxisOption) *axisPainter {
	p = p.Child(PainterPaddingOption(Box{
		Bottom: defaultXAxisHeight,
	}))
	axisOpt := opt.ToAxisOption(p)
	axisOpt.Position = PositionRight
	axisOpt.SplitLineShow = false
	return NewAxisPainter(p, axisOpt)
}

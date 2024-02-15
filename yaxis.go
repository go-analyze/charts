package charts

import (
	"github.com/golang/freetype/truetype"
)

type YAxisOption struct {
	// Show specifies if the y-axis should be rendered, set this to *false (through False()) to hide the axis.
	Show *bool
	// Theme specifies the colors used for the x-axis.
	Theme ColorPalette
	// Color for y-axis.
	AxisColor Color
	// Min, if set this will force the minimum value of y-axis.
	Min *float64
	// Max, if set this will force the maximum value of y-axis.
	Max *float64
	// RangeValuePaddingScale suggest a scale of padding added to the max and min values.
	RangeValuePaddingScale *float64
	// Data provides labels for the y-axis.
	Data []string
	// Position describes the position of y-axis, it can be 'left' or 'right'.
	Position string
	// FontSize specifies the font size of each label.
	FontSize float64
	// Font is the font used to render each label.
	Font *truetype.Font
	// FontColor is the color used for text rendered.
	FontColor Color
	// Formatter for replacing y-axis text values.
	Formatter string
	// Unit is a suggestion for how large the axis step is, this is a recommendation only. Larger numbers result in fewer labels.
	Unit float64
	// LabelCount is the number of labels to show on the axis.  Specify a smaller number to reduce writing collisions.
	LabelCount int
	// LabelCountAdjustment specifies a relative influence on how many labels should be rendered.
	// Typically, this is negative to result in cleaner graphs, positive values may result in text collisions.
	LabelCountAdjustment int
	// LabelSkipCount specifies a number of lines between labels where there will be no label and instead just a horizontal line.
	LabelSkipCount int
	isCategoryAxis bool
	// The flag for show axis split line, set this to true will show axis split line
	SplitLineShow *bool
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
		Formatter:            opt.Formatter,
		Theme:                theme,
		Data:                 opt.Data,
		Position:             position,
		FontSize:             opt.FontSize,
		StrokeWidth:          -1,
		Font:                 opt.Font,
		FontColor:            opt.FontColor,
		BoundaryGap:          False(),
		Unit:                 opt.Unit,
		LabelCount:           opt.LabelCount,
		LabelCountAdjustment: opt.LabelCountAdjustment,
		LabelSkipCount:       opt.LabelSkipCount,
		SplitLineShow:        true,
		Show:                 opt.Show,
	}
	if !opt.AxisColor.IsZero() {
		axisOpt.FontColor = opt.AxisColor
		axisOpt.Theme = theme.WithAxisColor(opt.AxisColor)
	}
	if opt.isCategoryAxis {
		axisOpt.BoundaryGap = True()
		axisOpt.StrokeWidth = 1
		axisOpt.SplitLineShow = false
	}
	if opt.SplitLineShow != nil {
		axisOpt.SplitLineShow = *opt.SplitLineShow
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

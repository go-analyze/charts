package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

type tableChart struct {
	p   *Painter
	opt *TableChartOption
}

// NewTableChart returns a table chart render
func NewTableChart(p *Painter, opt TableChartOption) *tableChart {
	return &tableChart{
		p:   p,
		opt: &opt,
	}
}

type TableCell struct {
	// Text the text of table cell
	Text string
	// Style the current style of table cell
	Style Style
	// Row the row index of table cell
	Row int
	// Column the column index of table cell
	Column int
}

type TableChartOption struct {
	// OutputFormat specifies the output type, "svg" or "png".
	OutputFormat string
	// Theme specifies the colors used for the table.
	Theme ColorPalette
	// Padding specifies the padding of table.
	Padding Box
	// Width specifies the width of the table.
	Width int
	// Header provides header data for the top of the table.
	Header []string
	// Data provides the row and column data for the table.
	Data [][]string
	// Spans provide the span for each column on the table.
	Spans []int
	// TextAligns specifies the text alignment for each cell on the table.
	TextAligns []string
	// The font size of table contents.
	FontSize float64
	// Font is the font used to render the table.
	Font *truetype.Font
	// FontColor is the color used for text on the table.
	FontColor Color
	// HeaderBackgroundColor provides a background color of header row.
	HeaderBackgroundColor Color
	// HeaderFontColor specifies a text color for the header text.
	HeaderFontColor Color
	// RowBackgroundColors specifies an array of colors for each row.
	RowBackgroundColors []Color
	// BackgroundColor specifies a general background color.
	BackgroundColor Color
	// CellTextStyle customize text style of table cell
	CellTextStyle func(TableCell) *Style
	// CellStyle customize drawing style of table cell
	CellStyle func(TableCell) *Style
}

type TableSetting struct {
	// HeaderColor specifies the color of the header.
	HeaderColor Color
	// HeaderFontColor specifies the color of header text.
	HeaderFontColor Color
	// FontColor specifies the color of table text.
	FontColor Color
	// RowColors specifies an array of colors for each row.
	RowColors []Color
	// Padding specifies the padding of each cell.
	CellPadding Box
}

var TableLightThemeSetting = TableSetting{
	HeaderColor:     Color{R: 220, G: 220, B: 220, A: 255},
	HeaderFontColor: Color{R: 80, G: 80, B: 80, A: 255},
	FontColor:       Color{R: 50, G: 50, B: 50, A: 255},
	RowColors: []Color{
		drawing.ColorWhite,
		{R: 245, G: 245, B: 245, A: 255},
	},
	CellPadding: Box{
		Left:   10,
		Top:    10,
		Right:  10,
		Bottom: 10,
	},
}

var TableDarkThemeSetting = TableSetting{
	HeaderColor:     Color{R: 38, G: 38, B: 42, A: 255},
	HeaderFontColor: Color{R: 216, G: 217, B: 218, A: 255},
	FontColor:       Color{R: 216, G: 217, B: 218, A: 255},
	RowColors: []Color{
		{R: 24, G: 24, B: 28, A: 255},
		{R: 38, G: 38, B: 42, A: 255},
	},
	CellPadding: Box{
		Left:   10,
		Top:    10,
		Right:  10,
		Bottom: 10,
	},
}

type renderInfo struct {
	Width        int
	Height       int
	HeaderHeight int
	RowHeights   []int
	ColumnWidths []int
}

func (t *tableChart) render() (*renderInfo, error) {
	info := renderInfo{
		RowHeights: make([]int, 0),
	}
	p := t.p
	if t.opt.Theme == nil {
		t.opt.Theme = getPreferredTheme(p.theme)
	}
	opt := t.opt
	if len(opt.Header) == 0 {
		return nil, errors.New("header can not be nil")
	}
	if opt.FontSize <= 0 {
		opt.FontSize = 12
	}
	if opt.FontColor.IsZero() {
		if opt.Theme.IsDark() {
			opt.FontColor = TableDarkThemeSetting.FontColor
		} else {
			opt.FontColor = TableLightThemeSetting.FontColor
		}
	}
	if opt.Font == nil {
		opt.Font = GetDefaultFont()
	}
	if opt.HeaderFontColor.IsZero() {
		if opt.Theme.IsDark() {
			opt.HeaderFontColor = TableDarkThemeSetting.HeaderFontColor
		} else {
			opt.HeaderFontColor = TableLightThemeSetting.HeaderFontColor
		}
	}

	spans := opt.Spans
	if len(spans) != len(opt.Header) {
		newSpans := make([]int, len(opt.Header))
		for index := range opt.Header {
			if index >= len(spans) {
				newSpans[index] = 1
			} else {
				newSpans[index] = spans[index]
			}
		}
		spans = newSpans
	}

	sum := sumInt(spans)
	values := autoDivideSpans(p.Width(), sum, spans)
	columnWidths := make([]int, 0)
	for index, v := range values {
		if index == len(values)-1 {
			break
		}
		columnWidths = append(columnWidths, values[index+1]-v)
	}
	info.ColumnWidths = columnWidths

	height := 0
	textStyle := Style{
		FontSize:  opt.FontSize,
		FontColor: opt.HeaderFontColor,
		FillColor: opt.HeaderFontColor,
		Font:      opt.Font,
	}

	headerHeight := 0
	if opt.Padding.IsZero() {
		if opt.Theme.IsDark() {
			opt.Padding = TableDarkThemeSetting.CellPadding
		} else {
			opt.Padding = TableLightThemeSetting.CellPadding
		}
	}
	getCellTextStyle := opt.CellTextStyle
	if getCellTextStyle == nil {
		getCellTextStyle = func(_ TableCell) *Style {
			return nil
		}
	}
	// textAligns := opt.TextAligns
	getTextAlign := func(index int) string {
		if len(opt.TextAligns) <= index {
			return ""
		}
		return opt.TextAligns[index]
	}

	// processing of the table cells
	renderTableCells := func(
		currentStyle Style,
		rowIndex int,
		textList []string,
		currentHeight int,
		cellPadding Box,
	) int {
		cellMaxHeight := 0
		paddingHeight := cellPadding.Top + cellPadding.Bottom
		paddingWidth := cellPadding.Left + cellPadding.Right
		for index, text := range textList {
			cellStyle := getCellTextStyle(TableCell{
				Text:   text,
				Row:    rowIndex,
				Column: index,
				Style:  currentStyle,
			})
			if cellStyle == nil {
				cellStyle = &currentStyle
			}
			p.SetStyle(*cellStyle)
			x := values[index]
			y := currentHeight + cellPadding.Top
			width := values[index+1] - x
			x += cellPadding.Left
			width -= paddingWidth
			box := p.TextFit(text, x, y+int(opt.FontSize), width, getTextAlign(index))
			// calculate the highest height
			if box.Height()+paddingHeight > cellMaxHeight {
				cellMaxHeight = box.Height() + paddingHeight
			}
		}
		return cellMaxHeight
	}

	// processing of the table headers
	headerHeight = renderTableCells(textStyle, 0, opt.Header, height, opt.Padding)
	height += headerHeight
	info.HeaderHeight = headerHeight

	// processing of the table contents
	textStyle.FontColor = opt.FontColor
	textStyle.FillColor = opt.FontColor
	for index, textList := range opt.Data {
		cellHeight := renderTableCells(textStyle, index+1, textList, height, opt.Padding)
		info.RowHeights = append(info.RowHeights, cellHeight)
		height += cellHeight
	}

	info.Width = p.Width()
	info.Height = height
	return &info, nil
}

func (t *tableChart) renderWithInfo(info *renderInfo) (Box, error) {
	p := t.p
	if t.opt.Theme == nil {
		t.opt.Theme = getPreferredTheme(p.theme)
	}
	opt := t.opt
	if !opt.BackgroundColor.IsZero() {
		p.SetBackground(p.Width(), p.Height(), opt.BackgroundColor)
	}
	if opt.HeaderBackgroundColor.IsZero() {
		if opt.Theme.IsDark() {
			opt.HeaderBackgroundColor = TableDarkThemeSetting.HeaderColor
		} else {
			opt.HeaderBackgroundColor = TableLightThemeSetting.HeaderColor
		}
	}

	// if the header background colors is set
	p.SetBackground(info.Width, info.HeaderHeight, opt.HeaderBackgroundColor, true)
	currentHeight := info.HeaderHeight
	if opt.RowBackgroundColors == nil {
		if opt.Theme.IsDark() {
			opt.RowBackgroundColors = TableDarkThemeSetting.RowColors
		} else {
			opt.RowBackgroundColors = TableLightThemeSetting.RowColors
		}
	}
	for index, h := range info.RowHeights {
		color := opt.RowBackgroundColors[index%len(opt.RowBackgroundColors)]
		child := p.Child(PainterPaddingOption(Box{
			Top: currentHeight,
		}))
		child.SetBackground(p.Width(), h, color, true)
		currentHeight += h
	}
	// adjust the background color according to the set table style
	getCellStyle := opt.CellStyle
	if getCellStyle != nil {
		arr := [][]string{
			opt.Header,
		}
		arr = append(arr, opt.Data...)
		top := 0
		heights := []int{info.HeaderHeight}
		heights = append(heights, info.RowHeights...)
		// loop through all table cells to generate background color
		for i, textList := range arr {
			left := 0
			for j, v := range textList {
				style := getCellStyle(TableCell{
					Text:   v,
					Row:    i,
					Column: j,
				})
				if style != nil && !style.FillColor.IsZero() {
					padding := style.Padding
					child := p.Child(PainterPaddingOption(Box{
						Top:  top + padding.Top,
						Left: left + padding.Left,
					}))
					w := info.ColumnWidths[j] - padding.Left - padding.Top
					h := heights[i] - padding.Top - padding.Bottom
					child.SetBackground(w, h, style.FillColor, true)
				}
				left += info.ColumnWidths[j]
			}
			top += heights[i]
		}
	}
	_, err := t.render()
	if err != nil {
		return BoxZero, err
	}

	return Box{
		Right:  info.Width,
		Bottom: info.Height,
	}, nil
}

func (t *tableChart) Render() (Box, error) {
	p := t.p
	if !t.opt.BackgroundColor.IsZero() {
		p.SetBackground(p.Width(), p.Height(), t.opt.BackgroundColor)
	}

	r := p.render
	fn := chartdraw.PNG
	if p.outputFormat == ChartOutputSVG {
		fn = chartdraw.SVG
	}
	newRender, err := fn(p.Width(), 100)
	if err != nil {
		return BoxZero, err
	}
	p.render = newRender
	info, err := t.render()
	if err != nil {
		return BoxZero, err
	}
	p.render = r
	return t.renderWithInfo(info)
}

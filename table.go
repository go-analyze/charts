package charts

import (
	"errors"

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
	// FontStyle contains the configuration for the cell text font.
	FontStyle FontStyle
	// FillColor sets a color for this table cell.
	FillColor drawing.Color
	// Row the row index of table cell.
	Row int
	// Column the column index of table cell.
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
	// FontStyle contains the configuration for the table text font.
	FontStyle FontStyle
	// HeaderBackgroundColor provides a background color of header row.
	HeaderBackgroundColor Color
	// HeaderFontColor specifies a text color for the header text.
	HeaderFontColor Color
	// RowBackgroundColors specifies an array of colors for each row.
	RowBackgroundColors []Color
	// BackgroundColor specifies a general background color.
	BackgroundColor Color
	// CellModifier is an optional function to modify the style or content of a specific TableCell before they are rendered.
	CellModifier func(TableCell) TableCell
}

type tableSetting struct {
	// headerColor specifies the color of the header.
	headerColor Color
	// headerFontColor specifies the color of header text.
	headerFontColor Color
	// fontColor specifies the color of table text.
	fontColor Color
	// rowColors specifies an array of colors for each row.
	rowColors []Color
}

var tableLightThemeSetting = tableSetting{
	headerColor:     Color{R: 220, G: 220, B: 220, A: 255},
	headerFontColor: Color{R: 80, G: 80, B: 80, A: 255},
	fontColor:       Color{R: 50, G: 50, B: 50, A: 255},
	rowColors: []Color{
		drawing.ColorWhite,
		{R: 245, G: 245, B: 245, A: 255},
	},
}

var tableDarkThemeSetting = tableSetting{
	headerColor:     Color{R: 38, G: 38, B: 42, A: 255},
	headerFontColor: Color{R: 216, G: 217, B: 218, A: 255},
	fontColor:       Color{R: 216, G: 217, B: 218, A: 255},
	rowColors: []Color{
		{R: 24, G: 24, B: 28, A: 255},
		{R: 38, G: 38, B: 42, A: 255},
	},
}

type renderInfo struct {
	width        int
	height       int
	headerHeight int
	rowHeights   []int
	columnWidths []int
	tableCells   [][]TableCell
}

func (t *tableChart) render() (*renderInfo, error) {
	info := renderInfo{
		rowHeights: make([]int, 0),
	}
	p := t.p
	if t.opt.Theme == nil {
		t.opt.Theme = getPreferredTheme(p.theme)
	}
	opt := t.opt
	if len(opt.Header) == 0 {
		return nil, errors.New("header can not be nil")
	}
	fontStyle := opt.FontStyle
	if fontStyle.FontSize <= 0 {
		fontStyle.FontSize = 12
	}
	if fontStyle.FontColor.IsZero() {
		if opt.Theme.IsDark() {
			fontStyle.FontColor = tableDarkThemeSetting.fontColor
		} else {
			fontStyle.FontColor = tableLightThemeSetting.fontColor
		}
	}
	if fontStyle.Font == nil {
		fontStyle.Font = GetDefaultFont()
	}
	if opt.HeaderFontColor.IsZero() {
		if opt.Theme.IsDark() {
			opt.HeaderFontColor = tableDarkThemeSetting.headerFontColor
		} else {
			opt.HeaderFontColor = tableLightThemeSetting.headerFontColor
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
	info.columnWidths = columnWidths

	height := 0
	style := chartdraw.Style{
		FontStyle: FontStyle{
			FontSize:  fontStyle.FontSize,
			FontColor: opt.HeaderFontColor,
			Font:      fontStyle.Font,
		},
		FillColor: opt.HeaderFontColor,
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
		style chartdraw.Style,
		rowIndex int,
		textList []string,
		currentHeight int,
		cellPadding Box,
	) ([]TableCell, int) {
		cellMaxHeight := 0
		paddingHeight := cellPadding.Top + cellPadding.Bottom
		paddingWidth := cellPadding.Left + cellPadding.Right
		cells := make([]TableCell, len(textList))
		for index, text := range textList {
			tc := TableCell{
				Text:      text,
				Row:       rowIndex,
				Column:    index,
				FontStyle: style.FontStyle,
				FillColor: style.FillColor,
			}
			if opt.CellModifier != nil {
				tc = opt.CellModifier(tc)
				// Update style values to capture any changes
				style.FontStyle = tc.FontStyle
				style.FillColor = tc.FillColor
			}
			cells[index] = tc
			p.SetStyle(style)
			x := values[index]
			y := currentHeight + cellPadding.Top
			width := values[index+1] - x
			x += cellPadding.Left
			width -= paddingWidth
			box := p.TextFit(text, x, y+int(fontStyle.FontSize), width, getTextAlign(index))
			// calculate the highest height
			if box.Height()+paddingHeight > cellMaxHeight {
				cellMaxHeight = box.Height() + paddingHeight
			}
		}
		return cells, cellMaxHeight
	}

	info.tableCells = make([][]TableCell, len(opt.Data)+1)

	// processing of the table headers
	headerCells, headerHeight := renderTableCells(style, 0, opt.Header, height, opt.Padding)
	info.tableCells[0] = headerCells
	height += headerHeight
	info.headerHeight = headerHeight

	// processing of the table contents
	style.FontColor = fontStyle.FontColor
	style.FillColor = fontStyle.FontColor
	for index, textList := range opt.Data {
		newCells, cellHeight := renderTableCells(style, index+1, textList, height, opt.Padding)
		info.tableCells[index+1] = newCells
		info.rowHeights = append(info.rowHeights, cellHeight)
		height += cellHeight
	}

	info.width = p.Width()
	info.height = height
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
			opt.HeaderBackgroundColor = tableDarkThemeSetting.headerColor
		} else {
			opt.HeaderBackgroundColor = tableLightThemeSetting.headerColor
		}
	}
	p.SetBackground(info.width, info.headerHeight, opt.HeaderBackgroundColor, true)

	if opt.RowBackgroundColors == nil {
		if opt.Theme.IsDark() {
			opt.RowBackgroundColors = tableDarkThemeSetting.rowColors
		} else {
			opt.RowBackgroundColors = tableLightThemeSetting.rowColors
		}
	}

	currentHeight := info.headerHeight
	for index, h := range info.rowHeights {
		color := opt.RowBackgroundColors[index%len(opt.RowBackgroundColors)]
		child := p.Child(PainterPaddingOption(Box{
			Top: currentHeight,
		}))
		child.SetBackground(p.Width(), h, color, true)
		currentHeight += h
	}
	// adjust the background color according to the set table style
	if opt.CellModifier != nil {
		arr := [][]string{
			opt.Header,
		}
		arr = append(arr, opt.Data...)
		top := 0
		heights := []int{info.headerHeight}
		heights = append(heights, info.rowHeights...)
		// loop through all table cells to generate background color
		for i, _ := range arr {
			left := 0
			for j, tc := range info.tableCells[i] {
				if !tc.FillColor.IsZero() {
					padding := opt.Padding
					child := p.Child(PainterPaddingOption(Box{
						Top:  top + padding.Top,
						Left: left + padding.Left,
					}))
					w := info.columnWidths[j] - padding.Left - padding.Top
					h := heights[i] - padding.Top - padding.Bottom
					child.SetBackground(w, h, tc.FillColor, true)
				}
				left += info.columnWidths[j]
			}
			top += heights[i]
		}
	}
	_, err := t.render()
	if err != nil {
		return BoxZero, err
	}

	return Box{
		Right:  info.width,
		Bottom: info.height,
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

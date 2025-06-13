package charts

import (
	"errors"

	"github.com/go-analyze/charts/chartdraw"
)

// TableOptionRenderDirect table render with the provided options directly to an image. Table options are different
// from other charts as they include the state for initializing the Painter, where other charts accept the Painter. If
// you want to write a Table on an existing Painter use TableOptionRender
func TableOptionRenderDirect(opt TableChartOption) (*Painter, error) {
	if opt.OutputFormat == "" {
		opt.OutputFormat = chartDefaultOutputFormat
	}
	if opt.Width <= 0 {
		opt.Width = defaultChartWidth
	}

	p := NewPainter(PainterOptions{
		OutputFormat: opt.OutputFormat,
		Width:        opt.Width,
		Height:       100, // is only used to calculate the height of the table
		Font:         opt.FontStyle.Font,
	})
	info, err := newTableChart(p, opt).render()
	if err != nil {
		return nil, err
	}

	p = NewPainter(PainterOptions{
		OutputFormat: opt.OutputFormat,
		Width:        info.width,
		Height:       info.height,
		Font:         opt.FontStyle.Font,
	})
	if _, err = newTableChart(p, opt).renderWithInfo(info); err != nil {
		return nil, err
	}
	return p, nil
}

// TableRenderValues renders a table chart with the simple header and data values provided.
func TableRenderValues(header []string, data [][]string, spanMaps ...map[int]int) (*Painter, error) {
	opt := TableChartOption{
		Header: header,
		Data:   data,
	}
	if len(spanMaps) != 0 {
		spanMap := spanMaps[0]
		spans := make([]int, len(opt.Header))
		for index := range spans {
			v, ok := spanMap[index]
			if !ok {
				v = 1
			}
			spans[index] = v
		}
		opt.Spans = spans
	}
	return TableOptionRenderDirect(opt)
}

type tableChart struct {
	p   *Painter
	opt *TableChartOption
}

// newTableChart returns a table chart renderer.
func newTableChart(p *Painter, opt TableChartOption) *tableChart {
	return &tableChart{
		p:   p,
		opt: &opt,
	}
}

// TableCell represents a single cell in a table.
type TableCell struct {
	// Text the text of table cell
	Text string
	// FontStyle contains the configuration for the cell text font.
	FontStyle FontStyle
	// FillColor sets a color for this table cell.
	FillColor Color
	// Row the row index of table cell.
	Row int
	// Column the column index of table cell.
	Column int
}

// TableChartOption defines options for rendering a table chart.
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
		ColorWhite,
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
	info := renderInfo{}
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
		fontStyle.FontSize = defaultFontSize
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

	sum := chartdraw.SumInt(spans...)
	values := autoDivideSpans(p.Width(), sum, spans)
	columnWidths := make([]int, 0, len(values))
	for index, v := range values {
		if index == len(values)-1 {
			break
		}
		columnWidths = append(columnWidths, values[index+1]-v)
	}
	info.columnWidths = columnWidths

	var height int
	headerFontStyle := FontStyle{
		FontSize:  fontStyle.FontSize,
		FontColor: opt.HeaderFontColor,
		Font:      fontStyle.Font,
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
		fontStyle FontStyle,
		fillColor Color,
		rowIndex int,
		textList []string,
		currentHeight int,
		cellPadding Box,
	) ([]TableCell, int) {
		var cellMaxHeight int
		paddingHeight := cellPadding.Top + cellPadding.Bottom
		paddingWidth := cellPadding.Left + cellPadding.Right
		cells := make([]TableCell, len(textList))
		for index, text := range textList {
			tc := TableCell{
				Text:      text,
				Row:       rowIndex,
				Column:    index,
				FontStyle: fontStyle,
				FillColor: fillColor,
			}
			if opt.CellModifier != nil {
				tc = opt.CellModifier(tc)
				// TODO - deprecate this behavior
				// Update style values to capture any changes
				fontStyle = tc.FontStyle
				fillColor = tc.FillColor
			}
			cells[index] = tc
			x := values[index]
			y := currentHeight + cellPadding.Top
			width := values[index+1] - x
			x += cellPadding.Left
			width -= paddingWidth
			box := p.TextFit(text, x, y+int(fontStyle.FontSize), width, fontStyle, getTextAlign(index))
			// calculate the highest height
			if box.Height()+paddingHeight > cellMaxHeight {
				cellMaxHeight = box.Height() + paddingHeight
			}
		}
		return cells, cellMaxHeight
	}

	info.tableCells = make([][]TableCell, len(opt.Data)+1)

	// processing of the table headers
	headerCells, headerHeight := renderTableCells(headerFontStyle, ColorTransparent,
		0, opt.Header, height, opt.Padding)
	info.tableCells[0] = headerCells
	height += headerHeight
	info.headerHeight = headerHeight

	// processing of the table contents
	for index, textList := range opt.Data {
		newCells, cellHeight := renderTableCells(fontStyle, ColorTransparent, index+1, textList, height, opt.Padding)
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
		p.drawBackground(opt.BackgroundColor)
	}

	if opt.HeaderBackgroundColor.IsZero() {
		if opt.Theme.IsDark() {
			opt.HeaderBackgroundColor = tableDarkThemeSetting.headerColor
		} else {
			opt.HeaderBackgroundColor = tableLightThemeSetting.headerColor
		}
	}
	p.FilledRect(0, 0, info.width, info.headerHeight, opt.HeaderBackgroundColor, ColorTransparent, 0.0)

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
			Top:   currentHeight,
			IsSet: true,
		}))
		child.FilledRect(0, 0, p.Width(), h, color, color, 0.0)
		currentHeight += h
	}
	// adjust the background color according to the set table style
	if opt.CellModifier != nil {
		arr := [][]string{
			opt.Header,
		}
		arr = append(arr, opt.Data...)
		var top int
		heights := []int{info.headerHeight}
		heights = append(heights, info.rowHeights...)
		// loop through all table cells to generate background color
		for i := range arr {
			var left int
			for j, tc := range info.tableCells[i] {
				if !tc.FillColor.IsZero() {
					padding := opt.Padding
					child := p.Child(PainterPaddingOption(Box{
						Top:   top + padding.Top,
						Left:  left + padding.Left,
						IsSet: true,
					}))
					w := info.columnWidths[j] - padding.Left - padding.Top
					h := heights[i] - padding.Top - padding.Bottom
					child.FilledRect(0, 0, w, h, tc.FillColor, tc.FillColor, 0.0)
				}
				left += info.columnWidths[j]
			}
			top += heights[i]
		}
	}
	if _, err := t.render(); err != nil {
		return BoxZero, err
	}

	return Box{
		Right:  info.width,
		Bottom: info.height,
		IsSet:  true,
	}, nil
}

func (t *tableChart) Render() (Box, error) {
	p := t.p
	if !t.opt.BackgroundColor.IsZero() {
		p.drawBackground(t.opt.BackgroundColor)
	}

	r := p.render
	fn := chartdraw.PNG
	if p.outputFormat == ChartOutputSVG {
		fn = chartdraw.SVG
	}
	p.render = fn(p.Width(), 100)
	info, err := t.render()
	if err != nil {
		return BoxZero, err
	}
	p.render = r
	return t.renderWithInfo(info)
}

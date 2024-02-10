package charts

import (
	"errors"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type tableChart struct {
	p   *Painter
	opt *TableChartOption
}

// NewTableChart returns a table chart render
func NewTableChart(p *Painter, opt TableChartOption) *tableChart {
	if opt.Theme == nil {
		opt.Theme = getPreferredTheme(p.theme)
	}
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
	// The output type
	Type string
	// The width of table
	Width int
	// The theme
	Theme ColorPalette
	// The padding of table cell
	Padding Box
	// The header data of table
	Header []string
	// The data of table
	Data [][]string
	// The span list of table column
	Spans []int
	// The text align list of table cell
	TextAligns []string
	// The font size of table
	FontSize float64
	// The font to render the table with
	Font *truetype.Font
	// The font color of table
	FontColor Color
	// The background color of header
	HeaderBackgroundColor Color
	// The header font color
	HeaderFontColor Color
	// The background color of row
	RowBackgroundColors []Color
	// The background color
	BackgroundColor Color
	// CellTextStyle customize text style of table cell
	CellTextStyle func(TableCell) *Style
	// CellStyle customize drawing style of table cell
	CellStyle func(TableCell) *Style
}

type TableSetting struct {
	// The color of header
	HeaderColor Color
	// The color of header text
	HeaderFontColor Color
	// The color of table text
	FontColor Color
	// The color list of row
	RowColors []Color
	// The padding of cell
	Padding Box
}

var TableLightThemeSetting = TableSetting{
	HeaderColor: Color{
		R: 240,
		G: 240,
		B: 240,
		A: 255,
	},
	HeaderFontColor: Color{
		R: 98,
		G: 105,
		B: 118,
		A: 255,
	},
	FontColor: Color{
		R: 70,
		G: 70,
		B: 70,
		A: 255,
	},
	RowColors: []Color{
		drawing.ColorWhite,
		{
			R: 247,
			G: 247,
			B: 247,
			A: 255,
		},
	},
	Padding: Box{
		Left:   10,
		Top:    10,
		Right:  10,
		Bottom: 10,
	},
}

var TableDarkThemeSetting = TableSetting{
	HeaderColor: Color{
		R: 38,
		G: 38,
		B: 42,
		A: 255,
	},
	HeaderFontColor: Color{
		R: 216,
		G: 217,
		B: 218,
		A: 255,
	},
	FontColor: Color{
		R: 216,
		G: 217,
		B: 218,
		A: 255,
	},
	RowColors: []Color{
		{
			R: 24,
			G: 24,
			B: 28,
			A: 255,
		},
		{
			R: 38,
			G: 38,
			B: 42,
			A: 255,
		},
	},
	Padding: Box{
		Left:   10,
		Top:    10,
		Right:  10,
		Bottom: 10,
	},
}

var tableDefaultSetting = TableLightThemeSetting

// SetDefaultTableSetting sets the default setting for table
func SetDefaultTableSetting(setting TableSetting) {
	tableDefaultSetting = setting
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
	opt := t.opt
	if len(opt.Header) == 0 {
		return nil, errors.New("header can not be nil")
	}
	theme := opt.Theme
	if theme == nil {
		theme = p.theme
	}
	fontSize := opt.FontSize
	if fontSize == 0 {
		fontSize = 12
	}
	fontColor := opt.FontColor
	if fontColor.IsZero() {
		fontColor = tableDefaultSetting.FontColor
	}
	font := opt.Font
	if font == nil {
		font = GetDefaultFont()
	}
	headerFontColor := opt.HeaderFontColor
	if opt.HeaderFontColor.IsZero() {
		headerFontColor = tableDefaultSetting.HeaderFontColor
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
		FontSize:  fontSize,
		FontColor: headerFontColor,
		FillColor: headerFontColor,
		Font:      font,
	}

	headerHeight := 0
	padding := opt.Padding
	if padding.IsZero() {
		padding = tableDefaultSetting.Padding
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
			box := p.TextFit(text, x, y+int(fontSize), width, getTextAlign(index))
			// calculate the highest height
			if box.Height()+paddingHeight > cellMaxHeight {
				cellMaxHeight = box.Height() + paddingHeight
			}
		}
		return cellMaxHeight
	}

	// processing of the table headers
	headerHeight = renderTableCells(textStyle, 0, opt.Header, height, padding)
	height += headerHeight
	info.HeaderHeight = headerHeight

	// processing of the table contents
	textStyle.FontColor = fontColor
	textStyle.FillColor = fontColor
	for index, textList := range opt.Data {
		cellHeight := renderTableCells(textStyle, index+1, textList, height, padding)
		info.RowHeights = append(info.RowHeights, cellHeight)
		height += cellHeight
	}

	info.Width = p.Width()
	info.Height = height
	return &info, nil
}

func (t *tableChart) renderWithInfo(info *renderInfo) (Box, error) {
	p := t.p
	opt := t.opt
	if !opt.BackgroundColor.IsZero() {
		p.SetBackground(p.Width(), p.Height(), opt.BackgroundColor)
	}
	headerBGColor := opt.HeaderBackgroundColor
	if headerBGColor.IsZero() {
		headerBGColor = tableDefaultSetting.HeaderColor
	}

	// if the header background colors is set
	p.SetBackground(info.Width, info.HeaderHeight, headerBGColor, true)
	currentHeight := info.HeaderHeight
	rowColors := opt.RowBackgroundColors
	if rowColors == nil {
		rowColors = tableDefaultSetting.RowColors
	}
	for index, h := range info.RowHeights {
		color := rowColors[index%len(rowColors)]
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
		heights := []int{
			info.HeaderHeight,
		}
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
	opt := t.opt
	if !opt.BackgroundColor.IsZero() {
		p.SetBackground(p.Width(), p.Height(), opt.BackgroundColor)
	}

	r := p.render
	fn := chart.PNG
	if p.outputType == ChartOutputSVG {
		fn = chart.SVG
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

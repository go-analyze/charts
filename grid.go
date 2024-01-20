package charts

type gridPainter struct {
	p   *Painter
	opt *GridPainterOption
}

type GridPainterOption struct {
	// The stroke width
	StrokeWidth float64
	// The stroke color
	StrokeColor Color
	// The spans of column
	ColumnSpans []int
	// The column of grid
	Column int
	// The row of grid
	Row int
	// Ignore first row
	IgnoreFirstRow bool
	// Ignore last row
	IgnoreLastRow bool
	// Ignore first column
	IgnoreFirstColumn bool
	// Ignore last column
	IgnoreLastColumn bool
}

// NewGridPainter returns new a grid renderer
func NewGridPainter(p *Painter, opt GridPainterOption) *gridPainter {
	return &gridPainter{
		p:   p,
		opt: &opt,
	}
}

func (g *gridPainter) Render() (Box, error) {
	opt := g.opt
	ignoreColumnLines := make([]int, 0)
	if opt.IgnoreFirstColumn {
		ignoreColumnLines = append(ignoreColumnLines, 0)
	}
	if opt.IgnoreLastColumn {
		ignoreColumnLines = append(ignoreColumnLines, opt.Column)
	}
	ignoreRowLines := make([]int, 0)
	if opt.IgnoreFirstRow {
		ignoreRowLines = append(ignoreRowLines, 0)
	}
	if opt.IgnoreLastRow {
		ignoreRowLines = append(ignoreRowLines, opt.Row)
	}
	strokeWidth := opt.StrokeWidth
	if strokeWidth <= 0 {
		strokeWidth = 1
	}

	g.p.SetDrawingStyle(Style{
		StrokeWidth: strokeWidth,
		StrokeColor: opt.StrokeColor,
	})
	g.p.Grid(GridOption{
		Column:            opt.Column,
		ColumnSpans:       opt.ColumnSpans,
		Row:               opt.Row,
		IgnoreColumnLines: ignoreColumnLines,
		IgnoreRowLines:    ignoreRowLines,
	})
	return g.p.box, nil
}

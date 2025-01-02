package charts

import (
	"github.com/go-analyze/charts/chartdraw"
)

type gridPainter struct {
	p   *Painter
	opt *gridPainterOption
}

type gridPainterOption struct {
	// StrokeWidth is the grid line width.
	StrokeWidth float64
	// StrokeColor is the grid line color.
	StrokeColor Color
	// ColumnSpans specifies the span for each column.
	ColumnSpans []int
	// Columns is the count of columns in the grid.
	Columns int
	// Rows are the count of rows in the grid.
	Rows int
	// IgnoreFirstRow can be set to true to ignore the first row.
	IgnoreFirstRow bool
	// IgnoreLastRow can be set to true to ignore the last row.
	IgnoreLastRow bool
	// IgnoreFirstColumn can be set to true to ignore the first colum.
	IgnoreFirstColumn bool
	// IgnoreLastColumn can be set to true to ignore the last columns.
	IgnoreLastColumn bool
}

// newGridPainter returns new a grid renderer
func newGridPainter(p *Painter, opt gridPainterOption) *gridPainter {
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
		ignoreColumnLines = append(ignoreColumnLines, opt.Columns)
	}
	ignoreRowLines := make([]int, 0)
	if opt.IgnoreFirstRow {
		ignoreRowLines = append(ignoreRowLines, 0)
	}
	if opt.IgnoreLastRow {
		ignoreRowLines = append(ignoreRowLines, opt.Rows)
	}
	strokeWidth := opt.StrokeWidth
	if strokeWidth <= 0 {
		strokeWidth = 1
	}

	g.p.SetDrawingStyle(chartdraw.Style{
		StrokeWidth: strokeWidth,
		StrokeColor: opt.StrokeColor,
	})
	g.p.grid(gridOption{
		Columns:           opt.Columns,
		ColumnSpans:       opt.ColumnSpans,
		Rows:              opt.Rows,
		IgnoreColumnLines: ignoreColumnLines,
		IgnoreRowLines:    ignoreRowLines,
	})
	return g.p.box, nil
}

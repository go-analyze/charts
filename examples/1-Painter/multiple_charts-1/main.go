package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example of building a painter to write multiple charts on the same image.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "multiple-charts-1.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        800,
		Height:       600,
	})
	p.FilledRect(0, 0, 800, 600, charts.ColorWhite, charts.ColorWhite, 0.0)

	// APPROACH 1: Grid-based layout - Easiest and best for aligned, structured layouts
	// Creates a 2x2 grid where the top row is merged into a single cell
	painters, err := p.LayoutByGrid(2, 2).
		CellAt("topCenter", 0, 0).Span(2, 1).
		CellAt("bottomLeft", 0, 1).
		CellAt("bottomRight", 1, 1).
		Build()
	if err != nil {
		panic(err)
	}

	// APPROACH 2: Row-based layout - Best for progressive, flexible layouts
	// (only included to demonstrate the same layout using the row-based approach)
	painters, err = p.LayoutByRows().
		Row().Height("300").Columns("topCenter").   // First row: single column, fixed height
		Row().Columns("bottomLeft", "bottomRight"). // Second row: two equal columns, remaining height
		Build()
	if err != nil {
		panic(err)
	}

	// APPROACH 3: Direct Child painters - Most control but requires manual calculations
	/*
		topCenterPainter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 800, 300)),
			charts.PainterThemeOption(charts.GetTheme(charts.ThemeVividLight)))
		bottomLeftPainter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 300, 400, 600)),
			charts.PainterThemeOption(charts.GetTheme(charts.ThemeAnt)))
		bottomRightPainter := p.Child(charts.PainterBoxOption(charts.NewBox(400, 300, 800, 600)),
	*/

	// Prepare chart data and options
	dataValues := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	lineOpt := charts.LineChartOption{
		Padding: charts.NewBoxEqual(10),
		XAxis: charts.XAxisOption{
			Labels: []string{
				"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
			},
			LabelCount: 7,
		},
		Legend: charts.LegendOption{
			Show: charts.Ptr(false),
		},
		SeriesList: charts.NewSeriesListLine(dataValues),
	}
	barOpt := charts.NewBarChartOptionWithData(dataValues)
	barOpt.XAxis = lineOpt.XAxis
	barOpt.Legend = lineOpt.Legend
	scatterOpt := charts.NewScatterChartOptionWithData(dataValues)
	scatterOpt.XAxis = lineOpt.XAxis
	scatterOpt.Legend = lineOpt.Legend

	// Render charts on each child painter
	if err := painters["topCenter"].ScatterChart(scatterOpt); err != nil {
		panic(err)
	} else if err := painters["bottomLeft"].BarChart(barOpt); err != nil {
		panic(err)
	} else if err := painters["bottomRight"].LineChart(lineOpt); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

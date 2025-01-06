package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
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
	p.SetBackground(800, 600, drawing.ColorWhite)
	// set the space and theme for each chart
	topCenterPainter := p.Child(charts.PainterBoxOption(chartdraw.NewBox(0, 0, 800, 300)),
		charts.PainterThemeOption(charts.GetTheme(charts.ThemeVividLight)))
	bottomLeftPainter := p.Child(charts.PainterBoxOption(chartdraw.NewBox(300, 0, 400, 600)),
		charts.PainterThemeOption(charts.GetTheme(charts.ThemeAnt)))
	bottomRightPainter := p.Child(charts.PainterBoxOption(chartdraw.NewBox(300, 400, 800, 600)),
		charts.PainterThemeOption(charts.GetTheme(charts.ThemeLight)))

	lineOpt := charts.LineChartOption{
		Padding: charts.Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		XAxis: charts.XAxisOption{
			Data: []string{
				"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
			},
			LabelCount: 7,
		},
		Legend: charts.LegendOption{
			Data: []string{
				"Email", "Union Ads", "Video Ads", "Direct", "Search Engine",
			},
		},
		SeriesList: charts.NewSeriesListLine([][]float64{
			{120, 132, 101, 134, 90, 230, 210},
			{220, 182, 191, 234, 290, 330, 310},
			{150, 232, 201, 154, 190, 330, 410},
			{320, 332, 301, 334, 390, 330, 320},
			{820, 932, 901, 934, 1290, 1330, 1320},
		}),
	}

	// render the same chart in each spot for the demo
	if err := bottomLeftPainter.LineChart(lineOpt); err != nil {
		panic(err)
	}
	if err := bottomRightPainter.LineChart(lineOpt); err != nil {
		panic(err)
	}
	if err := topCenterPainter.LineChart(lineOpt); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

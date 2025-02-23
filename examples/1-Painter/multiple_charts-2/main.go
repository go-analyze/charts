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

	file := filepath.Join(tmpPath, "multiple-charts-2.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        800,
		Height:       600,
	})
	p.FilledRect(0, 0, 800, 600, charts.ColorWhite, charts.ColorWhite, 0.0)
	// set the space and theme for each chart
	topCenterPainter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 300, 800)))
	bottomLeftPainter := p.Child(charts.PainterBoxOption(charts.NewBox(300, 0, 600, 400)))
	bottomRightPainter := p.Child(charts.PainterBoxOption(charts.NewBox(300, 400, 600, 800)))

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
	barOpt := charts.BarChartOption{
		Padding: charts.NewBoxEqual(10),
		XAxis:   lineOpt.XAxis,
		Legend: charts.LegendOption{
			Show: charts.Ptr(false),
		},
		SeriesList: charts.NewSeriesListBar(dataValues),
	}
	pieOpt := charts.PieChartOption{
		Padding: charts.NewBoxEqual(10),
		Legend: charts.LegendOption{
			SeriesNames: []string{
				"Email", "Union Ads", "Video Ads", "Direct", "Search Engine",
			},
		},
		SeriesList: charts.NewSeriesListPie(lineOpt.SeriesList.SumSeries()), // utilize SumSeries() to easily get a pie chart representation
	}

	// render the same chart in each spot for the demo
	if err := bottomLeftPainter.BarChart(barOpt); err != nil {
		panic(err)
	}
	if err := bottomRightPainter.LineChart(lineOpt); err != nil {
		panic(err)
	}
	if err := topCenterPainter.PieChart(pieOpt); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

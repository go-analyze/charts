package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example line chart demonstrating gradient label coloring.
This shows how to:
- Use LabelFormatterGradientGreenRed to color labels from green (min) to red (max)
- Show labels with color-coding based on their values
- Create visually appealing data presentations with automatic color scaling
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-10-gradient-labels.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{20, 15, 35, 40, 10, 55, 25, 45, 30, 50},
	}

	opt := charts.NewLineChartOptionWithData(values)
	opt.XAxis.Labels = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct"}
	opt.SeriesList = charts.NewSeriesListLine(values, charts.LineSeriesOption{
		Names: []string{"Sales Performance"},
		Label: charts.SeriesLabel{
			Show:           charts.Ptr(true),
			LabelFormatter: charts.LabelFormatterGradientGreenRed(values[0]),
		},
	})
	opt.Padding = charts.NewBoxEqual(20)
	opt.Title = charts.TitleOption{
		Text:    "Sales Performance with Gradient Label Colors",
		Subtext: "(Green = Low Values, Red = High Values)",
	}
	opt.Legend.Show = charts.Ptr(false)

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        800,
		Height:       500,
	})

	if err := p.LineChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

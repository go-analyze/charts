package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example of a "Stacked" horizontal bar chart. Stacked charts are a good way to represent data where the sum is important,
and you want to show what components produce that sum.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "horizontal-bar-chart-3.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{10, 30, 50, 70, 90, 110, 130},
		{20, 40, 60, 80, 100, 120, 140},
	}

	opt := charts.NewHorizontalBarChartOptionWithData(values)
	opt.Title.Text = "Some Numbers"
	opt.Padding = charts.Box{
		Top:    20,
		Right:  20,
		Bottom: 0,
		Left:   20,
	}
	opt.StackSeries = charts.True()
	for i := range opt.SeriesList {
		opt.SeriesList[i].Label.Show = charts.True()
	}
	opt.Legend.SeriesNames = []string{
		"2011", "2012",
	}
	opt.XAxis.Show = charts.False()
	opt.YAxis = charts.YAxisOption{
		Labels: []string{
			"UN", "Brazil", "Indonesia", "USA", "India", "China", "World",
		},
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.HorizontalBarChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

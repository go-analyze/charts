package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example bar chart with a variety of basic configuration options shown using the Painter API.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "bar-chart-1-basic.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3},
	}

	opt := charts.NewBarChartOptionWithData(values)
	opt.Title.Text = "Bar Chart"
	opt.XAxis.Labels = []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}
	opt.Legend = charts.LegendOption{
		SeriesNames: []string{
			"Rainfall", "Evaporation",
		},
		Offset:       charts.OffsetRight,
		OverlayChart: charts.Ptr(true),
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.BarChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

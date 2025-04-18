package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example bar chart with custom bar sizes and margins.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "bar-chart-2-size_margin.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7},
	}

	opt := charts.NewBarChartOptionWithData(values)
	opt.XAxis.Labels = []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	}
	opt.Legend.Show = charts.Ptr(false)

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        1200,
		Height:       400,
	})
	defaultPainter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 400, 400)))
	opt.Title.Text = "Default"
	if err := defaultPainter.BarChart(opt); err != nil {
		panic(err)
	}
	barSizePainter := p.Child(charts.PainterBoxOption(charts.NewBox(400, 0, 800, 400)))
	opt.Title.Text = "Small Bar"
	opt.BarWidth = 4
	if err := barSizePainter.BarChart(opt); err != nil {
		panic(err)
	}
	marginPainter := p.Child(charts.PainterBoxOption(charts.NewBox(800, 0, 1200, 400)))
	opt.Title.Text = "No Margin"
	opt.BarMargin = charts.Ptr(0.0)
	opt.BarWidth = 0 // reset to default size
	if err := marginPainter.BarChart(opt); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

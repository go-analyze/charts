package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example heat map chart using the Painter API.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "heat-map-chart-1-basic.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{4.4, 4.9, 7.0, 7.5, 4.3},
		{2.6, 5.9, 9.0, 6.4, 2.3},
		{3.3, 6.4, 7.0, 4.9, 3.2},
		{1.9, 6.0, 9.0, 5.9, 2.6},
		{4.4, 5.9, 7.0, 6.4, 4.6},
	}

	opt := charts.NewHeatMapOptionWithData(values)
	opt.Title.Text = "Heat Map Chart"
	opt.Title.Offset = charts.OffsetCenter
	opt.XAxis.Title = "X-Axis"
	opt.YAxis.Title = "Y-Axis"

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.HeatMapChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

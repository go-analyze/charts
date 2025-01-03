package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example funnel chart with a variety of basic configuration options shown using the Painter API.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "funnel-chart-2.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{100, 80, 60, 40, 20, 10, 2}

	opt := charts.FunnelChartOption{}
	opt.SeriesList = charts.NewFunnelSeriesList(values)
	opt.Title.Text = "Funnel"
	opt.Legend.Data = []string{
		"Show", "Click", "Visit", "Inquiry", "Order", "Pay", "Cancel",
	}
	opt.Legend.Padding = charts.Box{Left: 100}

	p, err := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err != nil {
		panic(err)
	}
	if _, err = charts.NewFunnelChart(p, opt).Render(); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

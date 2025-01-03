package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example funnel chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "funnel-chart-1.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{100, 80, 60, 40, 20, 10, 2}

	p, err := charts.FunnelRender(
		values,
		charts.TitleTextOptionFunc("Funnel"),
		charts.LegendLabelsOptionFunc([]string{
			"Show", "Click", "Visit", "Inquiry", "Order", "Pay", "Cancel",
		}),
		func(opt *charts.ChartOption) {
			opt.Legend.Padding = charts.Box{Left: 100}
		},
	)
	if err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

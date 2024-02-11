package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "funnel-chart.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{
		100,
		80,
		60,
		40,
		20,
		10,
		0,
	}
	p, err := charts.FunnelRender(
		values,
		charts.TitleTextOptionFunc("Funnel"),
		charts.LegendLabelsOptionFunc([]string{
			"Show",
			"Click",
			"Visit",
			"Inquiry",
			"Order",
			"Pay",
			"Cancel",
		}),
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

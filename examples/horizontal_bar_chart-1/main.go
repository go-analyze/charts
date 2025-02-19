package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example horizontal bar chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "horizontal-bar-chart-1.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{10, 30, 50, 70, 90, 110, 130},
		{20, 40, 60, 80, 100, 120, 140},
	}
	p, err := charts.HorizontalBarRender(
		values,
		charts.TitleTextOptionFunc("World Population"),
		charts.PaddingOptionFunc(charts.Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		charts.LegendLabelsOptionFunc([]string{
			"2011", "2012",
		}),
		charts.YAxisLabelsOptionFunc([]string{
			"UN", "Brazil", "Indonesia", "USA", "India", "China", "World",
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

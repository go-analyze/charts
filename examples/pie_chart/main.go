package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example pie chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "pie-chart.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{
		1048,
		735,
		580,
		484,
		300,
	}
	p, err := charts.PieRender(
		values,
		charts.TitleOptionFunc(charts.TitleOption{
			Text:            "Rainfall vs Evaporation",
			Subtext:         "(Fake Data)",
			Left:            charts.PositionCenter,
			FontSize:        16,
			SubtextFontSize: 10,
		}),
		charts.PaddingOptionFunc(charts.Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			Orient: charts.OrientVertical,
			Data: []string{
				"Search Engine",
				"Direct",
				"Email",
				"Union Ads",
				"Video Ads",
			},
			Left:     charts.PositionLeft,
			FontSize: 12,
		}),
		charts.PieSeriesShowLabel(),
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

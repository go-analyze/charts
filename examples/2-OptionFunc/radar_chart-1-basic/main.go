package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example radar chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "radar-chart-1-basic.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{4200, 3000, 20000, 35000, 50000, 18000},
		{5000, 14000, 28000, 26000, 42000, 21000},
	}
	p, err := charts.RadarRender(
		values,
		charts.TitleOptionFunc(charts.TitleOption{
			Text:      "Basic Radar Chart",
			FontStyle: charts.NewFontStyleWithSize(16),
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{
				"Allocated Budget", "Actual Spending",
			},
			Offset: charts.OffsetRight,
		}),
		charts.RadarIndicatorOptionFunc([]string{
			"Sales",
			"Administration",
			"Information Technology",
			"Customer Support",
			"Development",
			"Marketing",
		}, []float64{
			6500, 16000, 30000, 38000, 52000, 25000,
		}),
	)
	if err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

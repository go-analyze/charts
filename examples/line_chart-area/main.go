package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example line chart with the area below the line shaded.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-area.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
	}
	p, err := charts.LineRender(
		values,
		charts.TitleTextOptionFunc("Line"),
		charts.XAxisLabelsOptionFunc([]string{
			"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{"Email"},
			Padding: charts.Box{
				Top:    5,
				Bottom: 10,
			},
		}),
		charts.YAxisOptionFunc(charts.YAxisOption{
			Min: charts.FloatPointer(0.0), // ensure y-axis starts at 0
		}),
		// setup fill styling below
		func(opt *charts.ChartOption) {
			opt.FillArea = true                    // shade the area under the line
			opt.FillOpacity = 150                  // set the fill opacity a little lighter than default
			opt.XAxis.BoundaryGap = charts.False() // BoundaryGap is less appealing when enabling FillArea
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

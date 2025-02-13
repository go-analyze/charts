package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example of putting two charts together.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "multiple-charts-3.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
	}
	p, err := charts.LineRender(
		values,
		charts.XAxisLabelsOptionFunc([]string{
			"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{
				"Email", "Video Ads", "Direct",
			},
			OverlayChart: charts.Ptr(false),
			Offset: charts.OffsetStr{
				Top:  charts.PositionBottom,
				Left: "20%",
			},
		}),
		func(opt *charts.ChartOption) {
			opt.YAxis = []charts.YAxisOption{
				{
					Max: charts.Ptr(2000.0),
				},
			}
			opt.Symbol = charts.SymbolCircle
			opt.LineStrokeWidth = 1.2
			opt.ValueFormatter = func(f float64) string {
				return charts.FormatValueHumanize(f, 1, true)
			}

			opt.Children = []charts.ChartOption{
				{
					Box: charts.NewBox(10, 200, 200, 500),
					SeriesList: charts.NewSeriesListHorizontalBar([][]float64{
						{70, 90, 110, 130},
						{80, 100, 120, 140},
					}).ToGenericSeriesList(),
					Legend: charts.LegendOption{
						SeriesNames: []string{
							"2011", "2012",
						},
					},
					YAxis: []charts.YAxisOption{
						{
							Labels: []string{
								"USA", "India", "China", "World",
							},
						},
					},
				},
			}
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

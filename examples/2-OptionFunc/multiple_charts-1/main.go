package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example of overlaying two charts on each other.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "multiple-charts-1.png")
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
				return fmt.Sprintf("%.0f", f)
			}
		},
	)
	if err != nil {
		panic(err)
	}

	hBarOpt := charts.HorizontalBarChartOption{
		Padding: charts.NewBoxEqual(10),
		// set a transparent background theme to not overwrite the chart below
		Theme: charts.GetDefaultTheme().WithBackgroundColor(charts.ColorTransparent),
		Legend: charts.LegendOption{
			SeriesNames: []string{
				"2011", "2012",
			},
		},
		YAxis: charts.YAxisOption{
			Labels: []string{
				"USA", "India", "China", "World",
			},
		},
		SeriesList: charts.NewSeriesListHorizontalBar([][]float64{
			{70, 90, 110, 130},
			{80, 100, 120, 140},
		}),
	}
	p = p.Child(charts.PainterBoxOption(charts.NewBox(200, 0, 600, 200)))

	if err = p.HorizontalBarChart(hBarOpt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

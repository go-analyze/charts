package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example scatter chart demonstrating top-N label selection.
This shows how to:
- Use LabelFormatterTopN to show labels only for the highest N values
- Reduce visual clutter by highlighting only the most important data points
- Track website traffic over time with emphasis on peak traffic days
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "scatter-chart-4-top-n-labels.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	// Sample data representing daily website visitors over 30 days (in thousands)
	values := [][]float64{
		{
			15.2, 18.5, 22.1, 19.8, 25.4, 21.3, 17.9, 32.6, 28.1, 24.7,
			31.5, 29.3, 26.8, 35.2, 41.7, 38.9, 33.1, 29.6, 27.4, 30.8,
			36.3, 42.1, 39.5, 44.8, 48.3, 45.6, 40.2, 37.9, 34.5, 26.1,
		},
	}

	opt := charts.NewScatterChartOptionWithData(values)
	opt.XAxis.Labels = []string{"Day 1", "Day 2", "Day 3", "Day 4", "Day 5",
		"Day 6", "Day 7", "Day 8", "Day 9", "Day 10",
		"Day 11", "Day 12", "Day 13", "Day 14", "Day 15",
		"Day 16", "Day 17", "Day 18", "Day 19", "Day 20",
		"Day 21", "Day 22", "Day 23", "Day 24", "Day 25",
		"Day 26", "Day 27", "Day 28", "Day 29", "Day 30"}
	opt.YAxis[0].Min = charts.Ptr(0.0)
	opt.YAxis[0].Max = charts.Ptr(50.0)
	opt.SeriesList = charts.NewSeriesListScatter(values, charts.ScatterSeriesOption{
		Names: []string{"Daily Visitors (k)"},
		Label: charts.SeriesLabel{
			Show: charts.Ptr(true),
			// Show labels only for the top 5 traffic days
			LabelFormatter: charts.LabelFormatterTopN(values[0], 5),
			FontStyle: charts.FontStyle{
				FontSize:  16,
				FontColor: charts.ColorRedAlt1,
			},
		},
	})
	opt.Legend.Show = charts.Ptr(false)
	opt.Padding = charts.NewBoxEqual(20)
	opt.Title = charts.TitleOption{
		Text:    "Website Traffic Over 30 Days - Peak Days Highlighted",
		Subtext: "(Only top 5 traffic days show labels)",
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        800,
		Height:       500,
	})

	if err := p.ScatterChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

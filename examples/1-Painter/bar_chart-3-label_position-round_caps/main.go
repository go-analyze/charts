package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example bar chart showing different series label positions and with rounded caps.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "bar-chart-3-label_position.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{23.2, 25.6, 76.7, 135.6, 162.2, 32.6},
		{26.4, 28.7, 70.7, 175.6, 182.2, 48.7},
	}

	opt := charts.NewBarChartOptionWithData(values)
	opt.XAxis.Labels = []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	}
	opt.BarMargin = charts.Ptr(1.0) // TODO - v0.6 - Update to percentage
	for i := range opt.SeriesList {
		opt.SeriesList[i].Label.Show = charts.Ptr(true)
		opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
			return charts.FormatValueHumanizeShort(f, 0, false)
		}
	}
	opt.RoundedBarCaps = charts.Ptr(true)

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        1000,
		Height:       400,
	})
	defaultPainter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 500, 400)))
	opt.Title.Text = "Bar Chart Top Label"
	if err := defaultPainter.BarChart(opt); err != nil {
		panic(err)
	}
	bottomLabelPainter := p.Child(charts.PainterBoxOption(charts.NewBox(500, 0, 1000, 400)))
	opt.Title.Text = "Bar Chart Bottom Label"
	opt.SeriesLabelPosition = charts.PositionBottom
	if err := bottomLabelPainter.BarChart(opt); err != nil {
		panic(err)
	}
	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

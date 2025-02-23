package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example basic line chart with a variety of basic configuration options shown using the Painter API.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-8-dual_y_axis.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	opt := charts.NewLineChartOptionWithData(values)
	opt.Title.Text = "Dual Axis Line"
	opt.XAxis = charts.XAxisOption{
		Labels: []string{"A", "B", "C", "D", "E", "F", "G"},
	}
	opt.Legend = charts.LegendOption{
		SeriesNames: []string{"Left Series", "Right Series"},
	}
	opt.SeriesList[1].YAxisIndex = 1
	opt.YAxis = append(opt.YAxis, opt.YAxis[0])
	opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
	opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.LineChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

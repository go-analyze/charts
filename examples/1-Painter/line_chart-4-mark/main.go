package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example basic line chart with mark points and mark lines configured.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-4-mark.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 95, 90, 230, 210},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}

	opt := charts.NewLineChartOptionWithData(values)
	opt.Padding = charts.NewBox(20, 20, 20, 48)
	opt.Title.FontStyle.FontSize = 16
	opt.XAxis.Labels = []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	opt.Legend.SeriesNames = []string{
		"Email", "Direct", "Search Engine",
	}
	opt.Symbol = charts.SymbolCircle
	opt.LineStrokeWidth = 1.2
	for i := range opt.SeriesList {
		opt.SeriesList[i].MarkPoint.AddPoints(charts.SeriesMarkTypeMax)
		opt.SeriesList[i].MarkLine.AddLines(charts.SeriesMarkTypeAverage)
		opt.SeriesList[i].MarkLine.ValueFormatter = func(v float64) string {
			return charts.FormatValueHumanizeShort(v, 1, false)
		}
	}

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

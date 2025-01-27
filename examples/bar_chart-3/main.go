package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example of a "Stacked" bar chart. Stacked charts are a good way to represent data where the sum is important,
and you want to show what components produce that sum.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "bar-chart-3.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3},
	}

	opt := charts.NewBarChartOptionWithData(values)
	opt.Padding = charts.Box{
		Top:    20,
		Right:  45,
		Bottom: 20,
		Left:   20,
	}
	opt.StackSeries = charts.True()
	opt.XAxis.Data = []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}
	opt.XAxis.LabelCount = 12 // force label count due to the labels being very close
	opt.Legend = charts.LegendOption{
		Data:         []string{"Rainfall", "Evaporation"},
		Offset:       charts.OffsetRight,
		OverlayChart: charts.True(),
	}
	// Markline to show the max for the first series, as well as the average for the first series
	opt.SeriesList[0].MarkLine = charts.NewMarkLine(charts.SeriesMarkDataTypeAverage, charts.SeriesMarkDataTypeMax)
	opt.SeriesList[0].MarkLine.ValueFormatter = func(f float64) string {
		return charts.FormatValueHumanizeShort(f, 0, false)
	}
	// Mark point on the lass series to show the maximum value of this series
	opt.SeriesList[1].MarkPoint = charts.NewMarkPoint(charts.SeriesMarkDataTypeMax)
	opt.SeriesList[1].MarkPoint.ValueFormatter = func(f float64) string {
		return "Max:" + charts.FormatValueHumanizeShort(f, 0, false)
	}
	opt.SeriesList[1].MarkPoint.SymbolSize = 32
	//opt.SeriesList[1].MarkPoint.GlobalPoint = true	// would change the mark point to put the total of the stack

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.BarChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

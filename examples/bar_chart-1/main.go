package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example bar chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "bar-chart-1.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3},
	}
	p, err := charts.BarRender(
		values,
		charts.XAxisLabelsOptionFunc([]string{
			"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{
				"Rainfall", "Evaporation",
			},
			Offset:       charts.OffsetRight,
			OverlayChart: charts.Ptr(true),
		}),
		charts.MarkLineOptionFunc(0, charts.SeriesMarkTypeAverage),
		charts.MarkPointOptionFunc(0, charts.SeriesMarkTypeMax,
			charts.SeriesMarkTypeMin),
		func(opt *charts.ChartOption) {
			opt.SeriesList[1].MarkPoint = charts.NewMarkPoint(
				charts.SeriesMarkTypeMax,
				charts.SeriesMarkTypeMin,
			)
			opt.SeriesList[1].MarkLine = charts.NewMarkLine(
				charts.SeriesMarkTypeAverage,
			)
			opt.XAxis.LabelCount = 12 // force label count due to the labels being very close
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

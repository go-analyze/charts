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

	file := filepath.Join(tmpPath, "bar-chart.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{
			2.0,
			4.9,
			7.0,
			23.2,
			25.6,
			76.7,
			135.6,
			162.2,
			32.6,
			20.0,
			6.4,
			3.3,
		},
		{
			2.6,
			5.9,
			9.0,
			26.4,
			28.7,
			70.7,
			175.6,
			182.2,
			48.7,
			18.8,
			6.0,
			2.3,
		},
	}
	p, err := charts.BarRender(
		values,
		charts.XAxisDataOptionFunc([]string{
			"Jan",
			"Feb",
			"Mar",
			"Apr",
			"May",
			"Jun",
			"Jul",
			"Aug",
			"Sep",
			"Oct",
			"Nov",
			"Dec",
		}),
		charts.LegendLabelsOptionFunc([]string{
			"Rainfall",
			"Evaporation",
		}, charts.PositionRight),
		charts.MarkLineOptionFunc(0, charts.SeriesMarkDataTypeAverage),
		charts.MarkPointOptionFunc(0, charts.SeriesMarkDataTypeMax,
			charts.SeriesMarkDataTypeMin),
		// custom option func
		func(opt *charts.ChartOption) {
			opt.SeriesList[1].MarkPoint = charts.NewMarkPoint(
				charts.SeriesMarkDataTypeMax,
				charts.SeriesMarkDataTypeMin,
			)
			opt.SeriesList[1].MarkLine = charts.NewMarkLine(
				charts.SeriesMarkDataTypeAverage,
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

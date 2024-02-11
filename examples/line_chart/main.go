package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{
			120,
			132,
			101,
			// 134,
			charts.GetNullValue(),
			90,
			230,
			210,
		},
		{
			220,
			182,
			191,
			234,
			290,
			330,
			310,
		},
		{
			150,
			232,
			201,
			154,
			190,
			330,
			410,
		},
		{
			320,
			332,
			301,
			334,
			390,
			330,
			320,
		},
		{
			820,
			932,
			901,
			934,
			1290,
			1330,
			1320,
		},
	}
	p, err := charts.LineRender(
		values,
		charts.TitleTextOptionFunc("Line"),
		charts.XAxisDataOptionFunc([]string{
			"Mon",
			"Tue",
			"Wed",
			"Thu",
			"Fri",
			"Sat",
			"Sun",
		}),
		charts.LegendLabelsOptionFunc([]string{
			"Email",
			"Union Ads",
			"Video Ads",
			"Direct",
			"Search Engine",
		}, "50"),
		func(opt *charts.ChartOption) {
			opt.Legend.Padding = charts.Box{
				Top:    5,
				Bottom: 10,
			}
			opt.SymbolShow = charts.False()
			opt.LineStrokeWidth = 1
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
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

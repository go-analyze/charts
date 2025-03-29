package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example doughnut chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "doughnut-chart-1-basic.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{
		1048, 735, 580, 484, 300,
	}
	p, err := charts.DoughnutRender(
		values,
		charts.TitleOptionFunc(charts.TitleOption{
			Text:             "Doughnut Chart",
			Subtext:          "(Fake Data)",
			Offset:           charts.OffsetCenter,
			FontStyle:        charts.NewFontStyleWithSize(16),
			SubtextFontStyle: charts.NewFontStyleWithSize(10),
		}),
		charts.PaddingOptionFunc(charts.NewBoxEqual(20)),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{
				"Search Engine", "Direct", "Email", "Union Ads", "Video Ads",
			},
			Vertical: charts.Ptr(true),
			Offset: charts.OffsetStr{
				Left: "80%",
				Top:  charts.PositionBottom,
			},
			FontStyle: charts.NewFontStyleWithSize(10),
		}),
	)
	if err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

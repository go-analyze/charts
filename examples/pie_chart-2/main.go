package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example pie chart with a variety of basic configuration options shown using the Painter API.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "pie-chart-2.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{
		1048, 735, 580, 484, 300,
	}

	opt := charts.NewPieChartOptionWithData(values)
	opt.Title = charts.TitleOption{
		Text:             "Rainfall vs Evaporation",
		Subtext:          "(Fake Data)",
		Offset:           charts.OffsetCenter,
		FontStyle:        charts.NewFontStyleWithSize(16),
		SubtextFontStyle: charts.NewFontStyleWithSize(10),
	}
	opt.Padding = charts.NewBoxEqual(20)
	opt.Legend = charts.LegendOption{
		SeriesNames: []string{
			"Search Engine", "Direct", "Email", "Union Ads", "Video Ads",
		},
		Vertical: charts.True(),
		Offset: charts.OffsetStr{
			Left: "80%",
			Top:  charts.PositionBottom,
		},
		FontStyle: charts.NewFontStyleWithSize(10),
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.PieChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

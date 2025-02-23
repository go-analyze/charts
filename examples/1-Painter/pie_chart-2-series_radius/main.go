package main

import (
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-analyze/charts"
)

/*
Example pie chart which varies the series radius by the percentage of the series.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "pie-chart-2-series_radius.png")
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
		Vertical: charts.Ptr(true),
		Offset: charts.OffsetStr{
			Left: "80%",
			Top:  charts.PositionBottom,
		},
		FontStyle: charts.NewFontStyleWithSize(10),
	}
	largestRadius := 120.0
	seriesMax := opt.SeriesList.MaxValue()
	for i := range opt.SeriesList {
		radiusValue := largestRadius * math.Sqrt(opt.SeriesList[i].Value/seriesMax)
		opt.SeriesList[i].Radius = strconv.Itoa(int(math.Ceil(radiusValue)))
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

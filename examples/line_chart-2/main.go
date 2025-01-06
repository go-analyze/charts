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

	file := filepath.Join(tmpPath, "line-chart-2.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, charts.GetNullValue(), 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}

	opt := charts.LineChartOption{}
	opt.SeriesList = charts.NewSeriesListLine(values)
	opt.Title.Text = "Line"
	opt.Title.FontStyle.FontSize = 16
	opt.XAxis.Data = []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	opt.Legend.Data = []string{
		"Email", "Union Ads", "Video Ads", "Direct", "Search Engine",
	}
	opt.Legend.Padding = charts.Box{
		Left: 100,
	}
	opt.SymbolShow = charts.True()
	opt.StrokeWidth = 1.2

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if _, err := charts.NewLineChart(p, opt).Render(); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example basic line chart demonstrating the spacing between boundary gaps.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-7-boundary_gap.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 90, 230, 210},
		{220, 182, 191, 290, 330, 310},
		{150, 232, 201, 190, 330, 410},
		{320, 332, 301, 390, 330, 320},
		{820, 932, 901, 1290, 1330, 1320},
	}

	opt := charts.NewLineChartOptionWithData(values)
	opt.Padding = charts.NewBoxEqual(10)
	opt.Title.FontStyle.FontSize = 16
	opt.XAxis.Labels = []string{
		"A", "B", "C", "D", "E", "F",
	}
	opt.Legend.Show = charts.Ptr(false)
	opt.Symbol = charts.SymbolCircle

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        1200,
		Height:       400,
	})
	boundaryGapPainter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 600, 400)))
	opt.Title.Text = "Boundary Gap"
	if err := boundaryGapPainter.LineChart(opt); err != nil {
		panic(err)
	}
	boundaryGapDisabledPainter := p.Child(charts.PainterBoxOption(charts.NewBox(600, 0, 1200, 400)))
	opt.XAxis.BoundaryGap = charts.Ptr(false)
	opt.Title.Text = "Boundary Gap Disabled"
	if err := boundaryGapDisabledPainter.LineChart(opt); err != nil {
		panic(err)
	}
	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

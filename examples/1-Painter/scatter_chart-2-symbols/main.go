package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example basic scatter chart showing per-series symbols.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "scatter-chart-2-symbols.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 95, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
	}

	opt := charts.NewScatterChartOptionWithData(values)
	opt.Title.FontStyle.FontSize = 16
	opt.XAxis.Labels = []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	opt.Legend.SeriesNames = []string{
		"Email", "Union Ads", "Video Ads", "Direct",
	}
	opt.SymbolSize = 4
	opt.SeriesList[0].Symbol = charts.SymbolCircle
	opt.SeriesList[1].Symbol = charts.SymbolDiamond
	opt.SeriesList[2].Symbol = charts.SymbolSquare
	opt.SeriesList[3].Symbol = charts.SymbolDot

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.ScatterChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

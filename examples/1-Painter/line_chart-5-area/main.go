package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example line chart with the area below the line shaded.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}
	file := filepath.Join(tmpPath, "line-chart-5-area.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
	}

	opt := charts.NewLineChartOptionWithData(values)
	opt.Title.Text = "Line"
	opt.XAxis.Labels = []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	opt.Legend.SeriesNames = []string{"Email"}
	opt.Legend.Padding = charts.Box{
		Top:    5,
		Bottom: 10,
	}
	opt.YAxis[0].Min = charts.Ptr(0.0) // Ensure y-axis starts at 0

	// Setup fill styling below
	opt.FillArea = charts.Ptr(true)           // Enable fill area
	opt.FillOpacity = 150                     // Set fill opacity
	opt.XAxis.BoundaryGap = charts.Ptr(false) // Disable boundary gap

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.LineChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err := writeFile(buf); err != nil {
		panic(err)
	}
}

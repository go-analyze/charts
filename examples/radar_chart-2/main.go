package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example radar chart with a variety of basic configuration options shown using the Painter API.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "radar-chart-2.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := [][]float64{
		{4200, 3000, 20000, 35000, 50000, 18000},
		{5000, 14000, 28000, 26000, 42000, 21000},
	}

	opt := charts.RadarChartOption{}
	opt.SeriesList = charts.NewSeriesListRadar(values)
	opt.Title = charts.TitleOption{
		Text: "Basic Radar Chart",
		FontStyle: charts.FontStyle{
			FontSize: 16,
		},
	}
	opt.Legend = charts.LegendOption{
		Data: []string{
			"Allocated Budget", "Actual Spending",
		},
		Offset: charts.OffsetRight,
	}
	opt.RadarIndicators = charts.NewRadarIndicators([]string{
		"Sales",
		"Administration",
		"Information Technology",
		"Customer Support",
		"Development",
		"Marketing",
	}, []float64{
		6500,
		16000,
		30000,
		38000,
		52000,
		25000,
	})

	p, err := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err != nil {
		panic(err)
	}
	if _, err = charts.NewRadarChart(p, opt).Render(); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

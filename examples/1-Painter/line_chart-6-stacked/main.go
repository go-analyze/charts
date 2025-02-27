package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
A "Stacked" line chart example. Stacked charts are a good way to represent data where the sum is important,
and you want to show what components produce that sum. When the line chart is stacked each series represents a layer in
the chart, with the last series line being drawn as the sum of all the values at any given point on the x-axis.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-6-stacked.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	markPoint := charts.NewMarkPoint("max")
	markPoint.ValueFormatter = func(f float64) string {
		return "Max:" + charts.FormatValueHumanizeShort(f, 0, false)
	}
	markPoint.SymbolSize = 30
	seriesList := charts.NewSeriesListLine([][]float64{
		{1.9, 23.2, 25.6, 102.6, 142.2, 32.6, 20.0, 2.3},
		{12.0, 26.4, 28.7, 144.6, 122.2, 48.7, 18.8, 13.3},
		{80.0, 40.4, 28.4, 48.8, 24.4, 24.2, 40.8, 80.8},
	}, charts.LineSeriesOption{
		Label: charts.SeriesLabel{
			Show: charts.Ptr(true),
			ValueFormatter: func(f float64) string {
				return charts.FormatValueHumanizeShort(f, 1, true)
			},
		},
		MarkPoint: markPoint,
	})
	dataLabels := []string{"A", "B", "C"}
	opt := charts.LineChartOption{
		Padding: charts.Box{
			Top:    10,
			Right:  40,
			Left:   10,
			Bottom: 10,
		},
		SeriesList:  seriesList,
		StackSeries: charts.Ptr(true),
		XAxis: charts.XAxisOption{
			Labels: []string{
				"1", "2", "3", "4", "5", "6", "7", "8",
			},
			BoundaryGap: charts.Ptr(false),
		},
		Legend: charts.LegendOption{
			SeriesNames: dataLabels,
		},
		YAxis: []charts.YAxisOption{
			{
				Title: "A+B+C Sum",
				TitleFontStyle: charts.FontStyle{
					FontSize:  12,
					FontColor: charts.ColorBlack,
				},
				Labels: dataLabels,
				LabelFontStyle: charts.FontStyle{
					FontSize:  8,
					FontColor: charts.ColorBlack,
					Font:      charts.GetDefaultFont(),
				},
			},
		},
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.LineChart(opt); err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

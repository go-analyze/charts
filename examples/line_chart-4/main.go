package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
	"github.com/go-analyze/charts/chartdraw"
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

	file := filepath.Join(tmpPath, "line-chart-4.png")
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
			Show: charts.True(),
			ValueFormatter: func(f float64) string {
				return charts.FormatValueHumanizeShort(f, 1, true)
			},
		},
		MarkPoint: markPoint,
	})
	dataLabels := []string{"A", "B", "C"}
	opt := charts.LineChartOption{
		Padding: charts.Box{
			Top:    20,
			Right:  40,
			Left:   25,
			Bottom: 20,
		},
		SeriesList:  seriesList,
		StackSeries: charts.True(),
		XAxis: charts.XAxisOption{
			Data: []string{
				"1", "2", "3", "4", "5", "6", "7", "8",
			},
			BoundaryGap: charts.False(),
		},
		Legend: charts.LegendOption{
			Data: dataLabels,
		},
		YAxis: []charts.YAxisOption{
			{
				Data: dataLabels,
				FontStyle: charts.FontStyle{
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
	p.Text("A+B+C Sum", 20, 240, chartdraw.DegreesToRadians(270), charts.FontStyle{
		FontSize:  12,
		FontColor: charts.ColorBlack,
		Font:      charts.GetDefaultFont(),
	})

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

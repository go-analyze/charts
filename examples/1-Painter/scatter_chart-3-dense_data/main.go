package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-analyze/charts"
)

/*
Another scatter chart example with large data point counts, and more significant theming and other customization.
*/

const dataPointCount = 1000

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "scatter-chart-3-dense_data.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := generateRandomData(3, dataPointCount, 10)
	xAxisLabels := generateLabels(dataPointCount, "foo ")
	axisFont := charts.NewFontStyleWithSize(6.0)
	trend := charts.NewTrendLine(charts.SeriesTrendTypeSMA)
	trend[0].Period = 100
	series := charts.NewSeriesListScatterMultiValue(values, charts.ScatterSeriesOption{
		TrendLine: trend,
		Label: charts.SeriesLabel{
			ValueFormatter: func(f float64) string {
				return charts.FormatValueHumanizeShort(f, 0, false)
			},
		},
	})
	series[0].MarkLine.AddLines(charts.SeriesMarkTypeMax)
	series[1].MarkLine.AddLines(charts.SeriesMarkTypeMax)

	opt := charts.ScatterChartOption{
		SeriesList: series,
		Padding:    charts.NewBoxEqual(16).WithRight(32),
		Theme:      charts.GetTheme(charts.ThemeVividLight),
		SymbolSize: 0.5,
		Title: charts.TitleOption{
			Text:   "Dense Scatter Chart Demo",
			Offset: charts.OffsetCenter,
		},
		Legend: charts.LegendOption{
			SeriesNames: []string{
				"One", "Two", "Three",
			},
			// Legend Vertical, on the right, and with smaller font to give more space for data
			Vertical:  charts.Ptr(true),
			Offset:    charts.OffsetRight,
			Align:     charts.AlignRight,
			FontStyle: charts.NewFontStyleWithSize(6.0),
		},
		YAxis: []charts.YAxisOption{
			{
				Min:            charts.Ptr(0.0), // force min to be zero
				Max:            charts.Ptr(280.0),
				FontStyle:      axisFont,
				Unit:           10,
				LabelSkipCount: 1,
			},
		},
		XAxis: charts.XAxisOption{
			Labels:        xAxisLabels,
			FontStyle:     axisFont,
			BoundaryGap:   charts.Ptr(false),
			LabelCount:    10,
			LabelRotation: charts.DegreesToRadians(45),
		},
	}

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

func generateRandomData(seriesCount int, dataPointCount int, maxVariationPercentage float64) [][][]float64 {
	data := make([][][]float64, seriesCount)
	for i := 0; i < seriesCount; i++ {
		data[i] = make([][]float64, dataPointCount)
	}

	for i := 0; i < seriesCount; i++ {
		for j := 0; j < dataPointCount; j++ {
			if j == 0 {
				// Set the initial value for the line
				data[i][j] = []float64{rand.Float64() * 100}
			} else {
				// Calculate the allowed variation range
				variationRange := data[i][j-1][0] * maxVariationPercentage / 100
				min := data[i][j-1][0] - variationRange
				max := data[i][j-1][0] + variationRange

				// Generate a random value within the allowed range
				values := []float64{min + rand.Float64()*(max-min)}
				if j%2 == 0 {
					values = append(values, min+rand.Float64()*(max-min))
				}
				if j%10 == 0 {
					values = append(values, min+rand.Float64()*(max-min))
				}
				data[i][j] = values
			}
		}
	}

	return data
}

func generateLabels(dataPointCount int, prefix string) []string {
	labels := make([]string, dataPointCount)
	for i := 0; i < dataPointCount; i++ {
		labels[i] = prefix + strconv.Itoa(i)
	}
	return labels
}

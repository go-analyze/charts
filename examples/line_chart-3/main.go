package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-analyze/charts"
)

/*
Another line chart example with large data point counts, and more significant theming and other customization.
*/

const dataPointCount = 100

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-3.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := generateRandomData(4, dataPointCount, 10)
	xAxisLabels := generateLabels(dataPointCount, "foo ")
	axisFont := charts.NewFontStyleWithSize(6.0)

	p, err := charts.LineRender(
		values,
		charts.ThemeNameOptionFunc(charts.ThemeVividLight), // custom color theme
		charts.DimensionsOptionFunc(800, 600),
		charts.TitleOptionFunc(charts.TitleOption{
			Text:   "Line Chart Demo",
			Offset: charts.OffsetCenter,
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{
				"Critical", "High", "Medium", "Low",
			},
			// Legend Vertical, on the right, and with smaller font to give more space for data
			Vertical:  charts.Ptr(true),
			Offset:    charts.OffsetRight,
			Align:     charts.AlignRight,
			FontStyle: charts.NewFontStyleWithSize(6.0),
		}),
		charts.PaddingOptionFunc(charts.NewBoxEqual(12)),
		charts.YAxisOptionFunc(charts.YAxisOption{
			Min:       charts.Ptr(0.0), // force min to be zero
			FontStyle: axisFont,
			// y-axis labels well spaced to keep a clean look
			Unit:           10,
			LabelSkipCount: 1,
		}),
		charts.XAxisOptionFunc(charts.XAxisOption{
			Labels:       xAxisLabels,
			FontStyle:    axisFont,
			BoundaryGap:  charts.Ptr(true),
			LabelCount:   10,
			TextRotation: charts.DegreesToRadians(45),
		}),
		func(opt *charts.ChartOption) {
			// disable the symbols and reduce the stroke width to give more fidelity on the line
			opt.SymbolShow = charts.Ptr(false)
			opt.LineStrokeWidth = 1.6
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
		},
	)
	if err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

func generateRandomData(lineCount int, dataPointCount int, maxVariationPercentage float64) [][]float64 {
	data := make([][]float64, lineCount)
	for i := 0; i < lineCount; i++ {
		data[i] = make([]float64, dataPointCount)
	}

	for i := 0; i < lineCount; i++ {
		for j := 0; j < dataPointCount; j++ {
			if j == 0 {
				// Set the initial value for the line
				data[i][j] = rand.Float64() * 100
			} else {
				// Calculate the allowed variation range
				variationRange := data[i][j-1] * maxVariationPercentage / 100
				min := data[i][j-1] - variationRange
				max := data[i][j-1] + variationRange

				// Generate a random value within the allowed range
				data[i][j] = min + rand.Float64()*(max-min)
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

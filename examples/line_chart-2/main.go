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
Another line chart example with more significant theming and other configuration changes.
*/

const dataPointCount = 100

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-2.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := generateRandomData(4, dataPointCount, 10)
	axisFontSize := 6.0

	p, err := charts.LineRender(
		values,
		charts.ThemeNameOptionFunc(charts.ThemeVividLight),
		charts.WidthOptionFunc(800),
		charts.HeightOptionFunc(600),
		charts.TitleTextOptionFunc("Line Chart Demo"),
		charts.LegendLabelsOptionFunc([]string{"Critical", "High", "Medium", "Low"}),
		charts.PaddingOptionFunc(charts.Box{
			Top:    12,
			Bottom: 12,
			Left:   12,
			Right:  12,
		}),
		charts.YAxisOptionFunc(charts.YAxisOption{
			Min: charts.FloatPointer(0.0), // force min to be zero
			FontStyle: charts.FontStyle{
				FontSize: axisFontSize,
			},
			Unit:           10,
			LabelSkipCount: 1,
		}),
		charts.XAxisOptionFunc(charts.XAxisOption{
			Data: generateLabels(dataPointCount, "foo "),
			FontStyle: charts.FontStyle{
				FontSize: axisFontSize,
			},
			BoundaryGap: charts.True(),
			LabelCount:  10,
		}),
		func(opt *charts.ChartOption) {
			opt.Legend.Left = charts.PositionRight
			opt.Legend.Align = charts.AlignRight
			opt.Legend.Orient = charts.OrientVertical
			opt.Legend.FontStyle.FontSize = 6
			opt.Title.Left = charts.PositionCenter
			opt.SymbolShow = charts.False()
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

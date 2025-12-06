package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-analyze/charts"
)

/*
Another line chart example with large data point counts, and use of gaps. This helps demonstrate additional
configuration and styling, and notably includes custom drawing directly on the Painter.

The data shown here is an example of comparing the f-stop across multiple Canon lens options.
*/

const startMM = 60
const endMM = 510

var lensDefinitions = map[string]string{
	"70-200mm f/2.8":           "70f2.8,201f-", // use a string encoding to define the f/stop point changes
	"70-200mm f/2.8 + 1.4x TC": "98f4,281f-",
	"70-200mm f/2.8 + 2x TC":   "140f5.6,401f-",
	"100-500mm f/4.5-7.1":      "100f4.5,151f5,254f5.6,363f6.3,472f7.1,501f-",
	//"200-800mm f/6.3-9":        "200f6.3,268f7.1,455f8,637f9,801f-",
}

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-9-custom.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values, xAxisLabels, labels, err := populateData()
	if err != nil {
		panic(err)
	}

	opt := charts.NewLineChartOptionWithData(values)
	opt.Theme = charts.GetTheme(charts.ThemeAnt)
	opt.Padding = charts.Box{
		Top:    20,
		Left:   20,
		Right:  20,
		Bottom: 10,
	}
	opt.Title.Text = "Canon RF Zoom Lenses"
	opt.Title.Offset = charts.OffsetCenter
	opt.Title.FontStyle.FontSize = 16

	opt.XAxis.Labels = xAxisLabels
	opt.XAxis.Unit = 40
	opt.XAxis.LabelCount = 10
	opt.XAxis.LabelRotation = charts.DegreesToRadians(45)
	opt.XAxis.BoundaryGap = charts.Ptr(true)
	opt.XAxis.LabelFontStyle = charts.NewFontStyleWithSize(6.0)
	opt.YAxis = []charts.YAxisOption{
		{
			Show:          charts.Ptr(false), // disabling in favor of manually printed y-values
			Min:           charts.Ptr(1.4),
			Max:           charts.Ptr(8.0),
			LabelCount:    4,
			SpineLineShow: charts.Ptr(true),
			FontStyle:     charts.NewFontStyleWithSize(8.0),
		},
	}
	opt.Legend.Show = charts.Ptr(false)
	opt.Symbol = charts.SymbolNone
	opt.LineStrokeWidth = 1.5

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		// positions drawn below depend on the canvas size set here
		Width:  600,
		Height: 400,
	})
	if err = p.LineChart(opt); err != nil {
		panic(err)
	}

	// Custom drawing directly on the Painter
	fontStyle := charts.FontStyle{
		FontSize:  12,
		FontColor: charts.ColorBlack,
		Font:      charts.GetDefaultFont(),
	}
	//p.Text("f/stop", 10, 170, chartdraw.DegreesToRadians(90), fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(0)
	p.Text(labels[0], 420, 84, 0, fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(1)
	p.Text(labels[1], 45, 284, 0, fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(2)
	p.Text(labels[2], 140, 230, 0, fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(3)
	p.Text(labels[3], 160, 155, 0, fontStyle)

	fontStyle.FontSize = 8
	fontStyle.FontColor = opt.Theme.GetSeriesColor(0)
	p.Text("f/4.5", 42, 220, 0, fontStyle)
	p.Text("f/5.0", 105, 196, 0, fontStyle)
	p.Text("f/6.3", 370, 137, 0, fontStyle)
	p.Text("f/7.1", 570, 100, 0, fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(1)
	p.Text("f/2.8", 5, 298, 0, fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(2)
	p.Text("f/4.0", 40, 244, 0, fontStyle)

	fontStyle.FontColor = opt.Theme.GetSeriesColor(3)
	p.Text("f/5.6", 92, 168, 0, fontStyle)

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

func populateData() (values [][]float64, xAxisLabels []string, labels []string, err error) {
	for k := range lensDefinitions {
		labels = append(labels, k)
	}
	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	for i := startMM; i <= endMM; i++ {
		xAxisLabels = append(xAxisLabels, fmt.Sprintf("%vmm", i))
	}

	for _, lens := range labels {
		parts := strings.Split(lensDefinitions[lens], ",")
		count := (endMM - startMM) + 1
		lensValues := make([]float64, count)
		currentFValue := charts.GetNullValue()
		// for code simplicity we assume startMM is strictly BEFORE the first lens, this allows us to set null
		// values until the start point (which will be loaded on the first run of the loop)
		nextPartIndex := 0
		nextMM := startMM
		nextFStop := currentFValue
		for i := 0; i < count; i++ {
			if i+startMM == nextMM {
				currentFValue = nextFStop
				if nextPartIndex < len(parts) {
					nextFStop, nextMM, err = parseFStopMM(parts[nextPartIndex])
					nextPartIndex++
					if err != nil {
						return
					}
				} else {
					nextFStop = charts.GetNullValue()
				}
			}
			lensValues[i] = currentFValue
		}
		values = append(values, lensValues)
	}

	return
}

func parseFStopMM(str string) (float64, int, error) {
	parts := strings.Split(str, "f")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid lens spec str: '%s'", str)
	}
	mm, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if parts[1] == "-" {
		return charts.GetNullValue(), mm, nil
	}
	fstop, err := strconv.ParseFloat(parts[1], 64)
	return fstop, mm, err
}

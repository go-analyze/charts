package main

//go:generate go run main.go

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-analyze/charts/chartdraw"
)

func random(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func main() {
	numValues := 1024
	numSeries := 100
	series := make([]chartdraw.Series, numSeries)

	for i := 0; i < numSeries; i++ {
		xValues := make([]time.Time, numValues)
		yValues := make([]float64, numValues)

		for j := 0; j < numValues; j++ {
			xValues[j] = time.Now().AddDate(0, 0, (numValues-j)*-1)
			yValues[j] = random(float64(-500), float64(500))
		}

		series[i] = chartdraw.TimeSeries{
			Name:    fmt.Sprintf("aaa.bbb.hostname-%v.ccc.ddd.eee.fff.ggg.hhh.iii.jjj.kkk.lll.mmm.nnn.value", i),
			XValues: xValues,
			YValues: yValues,
		}
	}

	graph := chartdraw.Chart{
		XAxis: chartdraw.XAxis{
			Name: "Time",
		},
		YAxis: chartdraw.YAxis{
			Name: "Value",
		},
		Series: series,
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

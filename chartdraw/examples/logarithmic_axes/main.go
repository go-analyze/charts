package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	/*
	   In this example we set the primary YAxis to have logarithmic range.
	*/

	graph := chartdraw.Chart{
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top:  20,
				Left: 20,
			},
		},
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				Name:    "A test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1, 10, 100, 1000, 10000},
			},
		},
		YAxis: chartdraw.YAxis{
			Style:     chartdraw.Shown(),
			NameStyle: chartdraw.Shown(),
			Range:     &chartdraw.LogarithmicRange{},
		},
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

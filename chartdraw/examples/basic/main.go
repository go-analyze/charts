package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	graph := chartdraw.Chart{
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}
	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

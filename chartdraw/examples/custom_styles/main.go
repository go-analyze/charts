package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func main() {
	/*
	   In this example we set some custom colors for the series and the chart background and canvas.
	*/
	graph := chartdraw.Chart{
		Background: chartdraw.Style{
			FillColor: drawing.ColorBlue,
		},
		Canvas: chartdraw.Style{
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				Style: chartdraw.Style{
					StrokeColor: drawing.ColorRed,               // will supercede defaults
					FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
				},
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

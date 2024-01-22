package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func main() {
	graph := chartdraw.Chart{
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				XValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
				YValues: chartdraw.Seq{Sequence: chartdraw.NewRandomSequence().WithLen(100).WithMin(100).WithMax(512)}.Values(),
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chartdraw.PNG, f)
}

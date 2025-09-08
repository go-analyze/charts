package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	/*
		In this example we add a new type of series, a `SimpleMovingAverageSeries` that takes another series as a required argument.
		InnerSeries only needs to implement `ValuesProvider`, so really you could chain `SimpleMovingAverageSeries` together if you wanted.
	*/

	mainSeries := chartdraw.ContinuousSeries{
		Name:    "A test series",
		XValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),        //generates a []float64 from 1.0 to 100.0 in 1.0 step increments, or 100 elements.
		YValues: chartdraw.Seq{Sequence: chartdraw.NewRandomSequence().WithLen(100).WithMin(0).WithMax(100)}.Values(), //generates a []float64 randomly from 0 to 100 with 100 elements.
	}

	// note we create a LinearRegressionSeries series by assignin the inner series.
	// we need to use a reference because `.Render()` needs to modify state within the series.
	linRegSeries := &chartdraw.LinearRegressionSeries{
		InnerSeries: mainSeries,
	} // we can optionally set the `WindowSize` property which alters how the moving average is calculated.

	graph := chartdraw.Chart{
		Series: []chartdraw.Series{
			mainSeries,
			linRegSeries,
		},
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

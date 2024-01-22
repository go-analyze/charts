package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {

	/*
		In this example we add a new type of series, a `PolynomialRegressionSeries` that takes another series as a required argument.
		InnerSeries only needs to implement `ValuesProvider`, so really you could chain `PolynomialRegressionSeries` together if you wanted.
	*/

	mainSeries := chartdraw.ContinuousSeries{
		Name:    "A test series",
		XValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),        //generates a []float64 from 1.0 to 100.0 in 1.0 step increments, or 100 elements.
		YValues: chartdraw.Seq{Sequence: chartdraw.NewRandomSequence().WithLen(100).WithMin(0).WithMax(100)}.Values(), //generates a []float64 randomly from 0 to 100 with 100 elements.
	}

	polyRegSeries := &chartdraw.PolynomialRegressionSeries{
		Degree:      3,
		InnerSeries: mainSeries,
	}

	graph := chartdraw.Chart{
		Series: []chartdraw.Series{
			mainSeries,
			polyRegSeries,
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chartdraw.PNG, f)
}

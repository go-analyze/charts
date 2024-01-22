package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	/*
	   In this example we add a `Renderable` or a custom component to the `Elements` array.
	   In this specific case it is a pre-built renderable (`CreateLegend`) that draws a legend for the chart's series.
	   If you like, you can use `CreateLegend` as a template for writing your own renderable, or even your own legend.
	*/

	graph := chartdraw.Chart{
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top:  20,
				Left: 260,
			},
		},
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				Name:    "A test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Another test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Yet Another test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Even More series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Foo Bar",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Bar Baz",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Moo Bar",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Zoo Bar Baz",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "Fast and the Furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "2 Fast 2 Furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},

			chartdraw.ContinuousSeries{
				Name:    "They only get more fast and more furious",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{5.0, 4.0, 3.0, 2.0, 1.0},
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []chartdraw.Renderable{
		chartdraw.LegendLeft(&graph),
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chartdraw.PNG, f)
}

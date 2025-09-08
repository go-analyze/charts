package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	/*
	   In this example we set a custom set of ticks to use for the y-axis. It can be (almost) whatever you want, including some custom labels for ticks.
	   Custom ticks will supercede a custom range, which will supercede automatic generation based on series values.
	*/

	graph := chartdraw.Chart{
		YAxis: chartdraw.YAxis{
			Range: &chartdraw.ContinuousRange{
				Min: 0.0,
				Max: 4.0,
			},
			Ticks: []chartdraw.Tick{
				{Value: 0.0, Label: "0.00"},
				{Value: 2.0, Label: "2.00"},
				{Value: 4.0, Label: "4.00"},
				{Value: 6.0, Label: "6.00"},
				{Value: 8.0, Label: "Eight"},
				{Value: 10.0, Label: "Ten"},
			},
		},
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

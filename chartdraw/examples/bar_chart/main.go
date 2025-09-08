package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	graph := chartdraw.BarChart{
		Title: "Test Bar Chart",
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars: []chartdraw.Value{
			{Value: 5.25, Label: "Blue"},
			{Value: 4.88, Label: "Green"},
			{Value: 4.74, Label: "Gray"},
			{Value: 3.22, Label: "Orange"},
			{Value: 3, Label: "Test"},
			{Value: 2.27, Label: "??"},
			{Value: 1, Label: "!!"},
		},
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

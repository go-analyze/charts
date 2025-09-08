package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	mainSeries := chartdraw.ContinuousSeries{
		Name:    "A test series",
		XValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values(),
		YValues: chartdraw.Seq{Sequence: chartdraw.NewRandomSequence().WithLen(100).WithMin(50).WithMax(150)}.Values(),
	}

	minSeries := &chartdraw.MinSeries{
		Style: chartdraw.Style{
			StrokeColor:     chartdraw.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	maxSeries := &chartdraw.MaxSeries{
		Style: chartdraw.Style{
			StrokeColor:     chartdraw.ColorAlternateGray,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: mainSeries,
	}

	graph := chartdraw.Chart{
		Width:  1920,
		Height: 1080,
		YAxis: chartdraw.YAxis{
			Name: "Random Values",
			Range: &chartdraw.ContinuousRange{
				Min: 25,
				Max: 175,
			},
		},
		XAxis: chartdraw.XAxis{
			Name: "Random Other Values",
		},
		Series: []chartdraw.Series{
			mainSeries,
			minSeries,
			maxSeries,
			chartdraw.LastValueAnnotationSeries(minSeries),
			chartdraw.LastValueAnnotationSeries(maxSeries),
		},
	}

	graph.Elements = []chartdraw.Renderable{chartdraw.Legend(&graph)}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = graph.Render(chartdraw.PNG, f)
}

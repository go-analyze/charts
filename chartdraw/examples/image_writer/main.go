package main

import (
	"fmt"
	"log"

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
	collector := &chartdraw.ImageWriter{}
	graph.Render(chartdraw.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)
}

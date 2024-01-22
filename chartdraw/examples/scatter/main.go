package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func drawChart(res http.ResponseWriter, req *http.Request) {
	viridisByY := func(xr, yr chartdraw.Range, index int, x, y float64) drawing.Color {
		return chartdraw.Viridis(y, yr.GetMin(), yr.GetMax())
	}

	graph := chartdraw.Chart{
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				Style: chartdraw.Style{
					StrokeWidth:      chartdraw.Disabled,
					DotWidth:         5,
					DotColorProvider: viridisByY,
				},
				XValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(0).WithEnd(127)}.Values(),
				YValues: chartdraw.Seq{Sequence: chartdraw.NewRandomSequence().WithLen(128).WithMin(0).WithMax(1024)}.Values(),
			},
		},
	}

	res.Header().Set("Content-Type", chartdraw.ContentTypePNG)
	err := graph.Render(chartdraw.PNG, res)
	if err != nil {
		log.Println(err.Error())
	}
}

func unit(res http.ResponseWriter, req *http.Request) {
	graph := chartdraw.Chart{
		Height: 50,
		Width:  50,
		Canvas: chartdraw.Style{
			Padding: chartdraw.BoxZero,
		},
		Background: chartdraw.Style{
			Padding: chartdraw.BoxZero,
		},
		Series: []chartdraw.Series{
			chartdraw.ContinuousSeries{
				XValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(0).WithEnd(4)}.Values(),
				YValues: chartdraw.Seq{Sequence: chartdraw.NewLinearSequence().WithStart(0).WithEnd(4)}.Values(),
			},
		},
	}

	res.Header().Set("Content-Type", chartdraw.ContentTypePNG)
	err := graph.Render(chartdraw.PNG, res)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	http.HandleFunc("/", drawChart)
	http.HandleFunc("/unit", unit)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

//go:generate go run main.go

import (
	"bytes"
	"log"
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	var b float64
	b = 1000

	ts1 := chartdraw.ContinuousSeries{ //TimeSeries{
		Name:    "Time Series",
		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{1.0, 2.0, 30.0, 4.0, 50.0, 6.0, 7.0, 88.0},
	}

	ts2 := chartdraw.ContinuousSeries{ //TimeSeries{
		Style: chartdraw.Style{
			StrokeColor: chartdraw.GetDefaultColor(1),
		},

		XValues: []float64{10 * b, 20 * b, 30 * b, 40 * b, 50 * b, 60 * b, 70 * b, 80 * b},
		YValues: []float64{15.0, 52.0, 30.0, 42.0, 50.0, 26.0, 77.0, 38.0},
	}

	graph := chartdraw.Chart{

		XAxis: chartdraw.XAxis{
			Name:           "The XAxis",
			ValueFormatter: chartdraw.TimeMinuteValueFormatter, //TimeHourValueFormatter,
		},

		YAxis: chartdraw.YAxis{
			Name: "The YAxis",
		},

		Series: []chartdraw.Series{
			ts1,
			ts2,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chartdraw.PNG, buffer)
	if err != nil {
		log.Fatal(err)
	}

	fo, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	if _, err := fo.Write(buffer.Bytes()); err != nil {
		panic(err)
	}
}

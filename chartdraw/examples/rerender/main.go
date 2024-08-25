package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-analyze/charts/chartdraw"
)

var lock sync.Mutex
var graph *chartdraw.Chart
var ts *chartdraw.TimeSeries

func addData(t time.Time, e time.Duration) {
	lock.Lock()
	ts.XValues = append(ts.XValues, t)
	ts.YValues = append(ts.YValues, float64(e.Milliseconds()))
	lock.Unlock()
}

func drawChart(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		addData(start, time.Since(start))
	}()
	if len(ts.XValues) == 0 {
		http.Error(res, "no data (yet)", http.StatusBadRequest)
		return
	}
	res.Header().Set("Content-Type", "image/png")
	if err := graph.Render(chartdraw.PNG, res); err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	ts = &chartdraw.TimeSeries{
		XValues: []time.Time{},
		YValues: []float64{},
	}
	graph = &chartdraw.Chart{
		Series: []chartdraw.Series{ts},
	}
	http.HandleFunc("/", drawChart)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

//go:generate go run main.go

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	log := chartdraw.NewLogger()
	drawChart(log)
}

func parseInt(str string) int {
	v, _ := strconv.Atoi(str)
	return v
}

func parseFloat64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

func readData() ([]time.Time, []float64) {
	var xvalues []time.Time
	var yvalues []float64
	err := chartdraw.ReadLines("requests.csv", func(line string) error {
		parts := chartdraw.SplitCSV(line)
		year := parseInt(parts[0])
		month := parseInt(parts[1])
		day := parseInt(parts[2])
		hour := parseInt(parts[3])
		elapsedMillis := parseFloat64(parts[4])
		xvalues = append(xvalues, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.UTC))
		yvalues = append(yvalues, elapsedMillis)
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return xvalues, yvalues
}

func releases() []chartdraw.GridLine {
	return []chartdraw.GridLine{
		{Value: chartdraw.TimeToFloat64(time.Date(2016, 8, 1, 9, 30, 0, 0, time.UTC))},
		{Value: chartdraw.TimeToFloat64(time.Date(2016, 8, 2, 9, 30, 0, 0, time.UTC))},
		{Value: chartdraw.TimeToFloat64(time.Date(2016, 8, 2, 15, 30, 0, 0, time.UTC))},
		{Value: chartdraw.TimeToFloat64(time.Date(2016, 8, 4, 9, 30, 0, 0, time.UTC))},
		{Value: chartdraw.TimeToFloat64(time.Date(2016, 8, 5, 9, 30, 0, 0, time.UTC))},
		{Value: chartdraw.TimeToFloat64(time.Date(2016, 8, 6, 9, 30, 0, 0, time.UTC))},
	}
}

func drawChart(log chartdraw.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		xvalues, yvalues := readData()
		mainSeries := chartdraw.TimeSeries{
			Name: "Prod Request Timings",
			Style: chartdraw.Style{
				StrokeColor: chartdraw.ColorBlue,
				FillColor:   chartdraw.ColorBlue.WithAlpha(100),
			},
			XValues: xvalues,
			YValues: yvalues,
		}

		linreg := &chartdraw.LinearRegressionSeries{
			Name: "Linear Regression",
			Style: chartdraw.Style{
				StrokeColor:     chartdraw.ColorAlternateBlue,
				StrokeDashArray: []float64{5.0, 5.0},
			},
			InnerSeries: mainSeries,
		}

		sma := &chartdraw.SMASeries{
			Name: "SMA",
			Style: chartdraw.Style{
				StrokeColor:     chartdraw.ColorRed,
				StrokeDashArray: []float64{5.0, 5.0},
			},
			InnerSeries: mainSeries,
		}

		graph := chartdraw.Chart{
			Log:    log,
			Width:  1280,
			Height: 720,
			Background: chartdraw.Style{
				Padding: chartdraw.Box{
					Top: 50,
				},
			},
			YAxis: chartdraw.YAxis{
				Name: "Elapsed Millis",
				TickStyle: chartdraw.Style{
					TextRotationDegrees: 45.0,
				},
				ValueFormatter: func(v interface{}) string {
					return fmt.Sprintf("%d ms", int(v.(float64)))
				},
			},
			XAxis: chartdraw.XAxis{
				ValueFormatter: chartdraw.TimeHourValueFormatter,
				GridMajorStyle: chartdraw.Style{
					StrokeColor: chartdraw.ColorAlternateGray,
					StrokeWidth: 1.0,
				},
				GridLines: releases(),
			},
			Series: []chartdraw.Series{
				mainSeries,
				linreg,
				chartdraw.LastValueAnnotationSeries(linreg),
				sma,
				chartdraw.LastValueAnnotationSeries(sma),
			},
		}

		graph.Elements = []chartdraw.Renderable{chartdraw.LegendThin(&graph)}

		f, _ := os.Create("output.png")
		defer f.Close()
		graph.Render(chartdraw.PNG, f)
	}
}

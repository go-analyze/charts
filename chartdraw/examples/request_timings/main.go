package main

//go:generate go run main.go

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	drawChart()
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
	err := readLines("requests.csv", func(line string) error {
		parts, err := splitCSV(line)
		if err != nil {
			return err
		}
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

// readLines reads a file and calls the handler for each line.
func readLines(filePath string, handler func(string) error) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		err = handler(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func splitCSV(line string) ([]string, error) {
	if len(line) == 0 {
		return []string{}, nil
	}

	r := csv.NewReader(strings.NewReader(line))
	return r.Read()
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

func drawChart() http.HandlerFunc {
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
		defer func() { _ = f.Close() }()
		_ = graph.Render(chartdraw.PNG, f)
	}
}

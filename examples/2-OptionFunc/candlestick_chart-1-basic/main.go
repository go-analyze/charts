package main

import (
	"os"

	"github.com/go-analyze/charts"
)

func main() {
	// Create sample OHLC data
	ohlcData := []charts.OHLCData{
		{Open: 100.0, High: 110.0, Low: 95.0, Close: 105.0},
		{Open: 105.0, High: 115.0, Low: 100.0, Close: 112.0},
		{Open: 112.0, High: 118.0, Low: 108.0, Close: 115.0},
		{Open: 115.0, High: 120.0, Low: 110.0, Close: 108.0}, // bearish candle
		{Open: 108.0, High: 113.0, Low: 105.0, Close: 109.0},
		{Open: 109.0, High: 116.0, Low: 106.0, Close: 114.0},
		{Open: 114.0, High: 121.0, Low: 111.0, Close: 119.0},
		{Open: 119.0, High: 125.0, Low: 116.0, Close: 122.0},
	}

	// Convert to generic series list for use with ChartOption
	candlestickSeriesList := charts.CandlestickSeriesList{{Data: ohlcData}}
	seriesList := candlestickSeriesList.ToGenericSeriesList()

	// Create chart using ChartOption and OptionFunc pattern
	painter, err := charts.Render(charts.ChartOption{
		SeriesList: seriesList,
		Title: charts.TitleOption{
			Text: "Basic Candlestick Chart (OptionFunc)",
			FontStyle: charts.FontStyle{
				FontSize: 18,
			},
		},
		XAxis: charts.XAxisOption{
			Labels: []string{"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7", "Day 8"},
		},
		YAxis: []charts.YAxisOption{
			{
				Unit: 1,
			},
		},
		Legend: charts.LegendOption{
			SeriesNames: []string{"Stock Price"},
			Show:        charts.Ptr(true),
		},
		Padding:      charts.NewBoxEqual(20),
		Width:        800,
		Height:       600,
		OutputFormat: charts.ChartOutputPNG,
	})

	if err != nil {
		panic(err)
	}

	// Save the chart to file
	buf, err := painter.Bytes()
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("candlestick_basic_optionfunc.png", buf, 0644); err != nil {
		panic(err)
	}
}

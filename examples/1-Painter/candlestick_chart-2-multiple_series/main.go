package main

import (
	"os"

	"github.com/go-analyze/charts"
)

func main() {
	stockAData := []charts.OHLCData{
		{Open: 100.0, High: 110.0, Low: 95.0, Close: 105.0},
		{Open: 105.0, High: 115.0, Low: 100.0, Close: 112.0},
		{Open: 112.0, High: 118.0, Low: 108.0, Close: 115.0},
		{Open: 115.0, High: 120.0, Low: 110.0, Close: 108.0},
		{Open: 108.0, High: 113.0, Low: 105.0, Close: 109.0},
	}
	stockBData := []charts.OHLCData{
		{Open: 150.0, High: 160.0, Low: 145.0, Close: 155.0},
		{Open: 155.0, High: 165.0, Low: 150.0, Close: 162.0},
		{Open: 162.0, High: 168.0, Low: 158.0, Close: 165.0},
		{Open: 165.0, High: 170.0, Low: 160.0, Close: 158.0},
		{Open: 158.0, High: 163.0, Low: 155.0, Close: 159.0},
	}
	stockCData := []charts.OHLCData{
		{Open: 200.0, High: 210.0, Low: 195.0, Close: 205.0},
		{Open: 205.0, High: 215.0, Low: 200.0, Close: 212.0},
		{Open: 212.0, High: 218.0, Low: 208.0, Close: 215.0},
		{Open: 215.0, High: 220.0, Low: 210.0, Close: 208.0},
		{Open: 208.0, High: 213.0, Low: 205.0, Close: 209.0},
	}

	// Create candlestick chart option with multiple series
	opt := charts.CandlestickChartOption{
		Theme: charts.GetTheme(charts.ThemeLight),
		SeriesList: charts.CandlestickSeriesList{
			charts.CandlestickSeries{Data: stockAData, Name: "Stock A", CandleStyle: charts.CandleStyleFilled},
			charts.CandlestickSeries{Data: stockBData, Name: "Stock B", CandleStyle: charts.CandleStyleTraditional},
			charts.CandlestickSeries{Data: stockCData, Name: "Stock C", CandleStyle: charts.CandleStyleOutline},
		},
		XAxis: charts.XAxisOption{
			Labels: []string{"Day 1", "Day 2", "Day 3", "Day 4", "Day 5"},
		},
		YAxis: []charts.YAxisOption{
			{Unit: 1},
		},
		Legend: charts.LegendOption{
			SeriesNames: []string{"Stock A", "Stock B", "Stock C"},
			Show:        charts.Ptr(true),
		},
		Padding: charts.NewBoxEqual(20),
	}

	// Create painter with specific dimensions and output format
	painterOptions := charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        1000,
		Height:       700,
	}
	p := charts.NewPainter(painterOptions)

	// Render and save the candlestick chart
	if err := p.CandlestickChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err := os.WriteFile("candlestick_multiple_series.png", buf, 0644); err != nil {
		panic(err)
	}
}

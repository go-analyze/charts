package main

import (
	"os"

	"github.com/go-analyze/charts"
)

func main() {
	ohlcData := []charts.OHLCData{
		{Open: 100.0, High: 110.0, Low: 95.0, Close: 105.0},
		{Open: 105.0, High: 115.0, Low: 100.0, Close: 112.0},
		{Open: 112.0, High: 118.0, Low: 108.0, Close: 115.0},
		{Open: 115.0, High: 120.0, Low: 104.0, Close: 108.0}, // bearish candle
		{Open: 108.0, High: 113.0, Low: 105.0, Close: 109.0},
		{Open: 109.0, High: 116.0, Low: 106.0, Close: 114.0},
		{Open: 114.0, High: 121.0, Low: 111.0, Close: 119.0},
	}

	opt := charts.NewCandlestickOptionWithData(ohlcData)
	opt.Title = charts.TitleOption{
		Text: "Candlestick Chart",
		FontStyle: charts.FontStyle{
			FontSize: 16,
		},
	}
	opt.XAxis = charts.XAxisOption{
		Labels: []string{"Day 1", "Day 2", "Day 3", "Day 4", "Day 5", "Day 6", "Day 7"},
	}
	opt.Legend = charts.LegendOption{
		SeriesNames: []string{"Stock Price"},
		Show:        charts.Ptr(true),
	}

	// Create painter with specific dimensions and output format
	painterOptions := charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	}
	p := charts.NewPainter(painterOptions)

	// Render and save the candlestick chart
	if err := p.CandlestickChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err := os.WriteFile("candlestick_basic.png", buf, 0644); err != nil {
		panic(err)
	}
}

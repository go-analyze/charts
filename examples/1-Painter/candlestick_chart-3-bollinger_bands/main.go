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
		{Open: 115.0, High: 125.0, Low: 110.0, Close: 120.0},
		{Open: 120.0, High: 130.0, Low: 115.0, Close: 125.0},
		{Open: 125.0, High: 135.0, Low: 120.0, Close: 130.0},
		{Open: 130.0, High: 140.0, Low: 125.0, Close: 135.0},
		{Open: 135.0, High: 145.0, Low: 130.0, Close: 140.0},
		{Open: 140.0, High: 150.0, Low: 135.0, Close: 145.0},
		{Open: 145.0, High: 155.0, Low: 140.0, Close: 150.0},
		{Open: 150.0, High: 160.0, Low: 145.0, Close: 148.0}, // Price pullback
		{Open: 148.0, High: 153.0, Low: 143.0, Close: 146.0},
		{Open: 146.0, High: 151.0, Low: 141.0, Close: 144.0},
		{Open: 144.0, High: 149.0, Low: 139.0, Close: 142.0},
		{Open: 142.0, High: 147.0, Low: 137.0, Close: 140.0},
		{Open: 140.0, High: 145.0, Low: 135.0, Close: 138.0},
		{Open: 138.0, High: 143.0, Low: 133.0, Close: 136.0},
		{Open: 136.0, High: 141.0, Low: 131.0, Close: 134.0},
		{Open: 134.0, High: 139.0, Low: 129.0, Close: 132.0},
		{Open: 132.0, High: 137.0, Low: 127.0, Close: 130.0},
	}

	chartOpt := charts.CandlestickChartOption{
		SeriesList: []charts.CandlestickSeries{{
			Data: ohlcData,
			CloseTrendLine: []charts.SeriesTrendLine{
				{Type: charts.SeriesTrendTypeBollingerUpper, Period: 10},
				{Type: charts.SeriesTrendTypeSMA, Period: 10}, // Middle band
				{Type: charts.SeriesTrendTypeBollingerLower, Period: 10},
			},
		}},
		Title: charts.TitleOption{
			Text: "Candlestick Chart with Bollinger Bands",
			FontStyle: charts.FontStyle{
				FontSize: 18,
			},
		},
		XAxis: charts.XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
				"11", "12", "13", "14", "15", "16", "17", "18", "19", "20"},
		},
		YAxis: []charts.YAxisOption{
			{
				Unit: 1,
			},
		},
		Legend: charts.LegendOption{
			Show: charts.Ptr(true),
		},
		Padding: charts.NewBoxEqual(20),
	}

	// Create painter and render the chart
	painter := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        800,
		Height:       600,
	}, charts.PainterThemeOption(charts.GetDefaultTheme()))

	if err := painter.CandlestickChart(chartOpt); err != nil {
		panic(err)
	} else if buf, err := painter.Bytes(); err != nil {
		panic(err)
	} else if err := os.WriteFile("candlestick_bollinger_bands.png", buf, 0644); err != nil {
		panic(err)
	}
}

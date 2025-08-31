package main

import (
    "os"

    "github.com/go-analyze/charts"
)

// This example demonstrates aggregating 1-minute OHLC candles into 5-minute
// candles, rendering two charts on the same painter: the original data on top
// and the aggregated data below.
func main() {
    // Create 1-minute OHLC data (simulated)
    minuteData := []charts.OHLCData{
        // First 5-minute period
        {Open: 100.0, High: 102.0, Low: 99.0, Close: 101.0},  // Minute 1
        {Open: 101.0, High: 103.0, Low: 100.0, Close: 102.0}, // Minute 2
        {Open: 102.0, High: 105.0, Low: 101.0, Close: 104.0}, // Minute 3
        {Open: 104.0, High: 106.0, Low: 103.0, Close: 105.0}, // Minute 4
        {Open: 105.0, High: 107.0, Low: 104.0, Close: 106.0}, // Minute 5

        // Second 5-minute period
        {Open: 106.0, High: 108.0, Low: 105.0, Close: 107.0}, // Minute 6
        {Open: 107.0, High: 109.0, Low: 106.0, Close: 108.0}, // Minute 7
        {Open: 108.0, High: 110.0, Low: 107.0, Close: 109.0}, // Minute 8
        {Open: 109.0, High: 111.0, Low: 108.0, Close: 110.0}, // Minute 9
        {Open: 110.0, High: 112.0, Low: 109.0, Close: 111.0}, // Minute 10

        // Third 5-minute period
        {Open: 111.0, High: 113.0, Low: 110.0, Close: 112.0}, // Minute 11
        {Open: 112.0, High: 114.0, Low: 111.0, Close: 113.0}, // Minute 12
        {Open: 113.0, High: 115.0, Low: 112.0, Close: 114.0}, // Minute 13
        {Open: 114.0, High: 116.0, Low: 113.0, Close: 115.0}, // Minute 14
        {Open: 115.0, High: 117.0, Low: 114.0, Close: 116.0}, // Minute 15
    }

    // Aggregate to 5-minute candles
    minuteSeries := charts.CandlestickSeries{Data: minuteData, Name: "1-Minute"}
    fiveMinuteSeries := charts.AggregateCandlestick(minuteSeries, 5)

    // Build painter and create two child regions (top/bottom)
    p := charts.NewPainter(charts.PainterOptions{
        OutputFormat: charts.ChartOutputPNG,
        Width:        1200,
        Height:       800,
    })

    // Optional: white background
    p.FilledRect(0, 0, 1200, 800, charts.ColorWhite, charts.ColorWhite, 0)

    top := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 1200, 400)))
    bottom := p.Child(charts.PainterBoxOption(charts.NewBox(0, 400, 1200, 800)))

    // Top chart: original 1-minute data
    topOpt := charts.CandlestickChartOption{
        Title: charts.TitleOption{
            Text: "1-Minute Candles (Before Aggregation)",
            FontStyle: charts.FontStyle{FontSize: 16},
        },
        XAxis: charts.XAxisOption{
            Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"},
        },
        YAxis: []charts.YAxisOption{{Unit: 1}},
        Legend: charts.LegendOption{SeriesNames: []string{"1-Minute"}, Show: charts.Ptr(true)},
        Padding:   charts.NewBoxEqual(20),
        SeriesList: charts.CandlestickSeriesList{{Data: minuteData, Name: "1-Minute"}},
    }

    if err := top.CandlestickChart(topOpt); err != nil {
        panic(err)
    }

    // Bottom chart: aggregated 5-minute data
    bottomOpt := charts.CandlestickChartOption{
        Title: charts.TitleOption{
            Text: "5-Minute Aggregated Candles",
            FontStyle: charts.FontStyle{FontSize: 16},
        },
        XAxis: charts.XAxisOption{
            Labels: []string{"1-5", "6-10", "11-15"},
        },
        YAxis: []charts.YAxisOption{{Unit: 1}},
        Legend: charts.LegendOption{SeriesNames: []string{"5-Minute"}, Show: charts.Ptr(true)},
        Padding:   charts.NewBoxEqual(20),
        SeriesList: charts.CandlestickSeriesList{{Data: fiveMinuteSeries.Data, Name: "5-Minute"}},
    }

    if err := bottom.CandlestickChart(bottomOpt); err != nil {
        panic(err)
    }

    // Save the combined image
    if buf, err := p.Bytes(); err != nil {
        panic(err)
    } else if err := os.WriteFile("candlestick_aggregation.png", buf, 0644); err != nil {
        panic(err)
    }
}

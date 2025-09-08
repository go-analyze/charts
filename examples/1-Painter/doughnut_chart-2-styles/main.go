package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Several doughnut charts showing different possible styles using the Painter API.
*/

func writeChart(p *charts.Painter, filename string) error {
	buf, err := p.Bytes()
	if err != nil {
		return err
	}
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, filename)
	return os.WriteFile(file, buf, 0600)
}

func main() {
	values := []float64{
		1048, 735, 580, 484, 300,
	}

	opt := charts.NewDoughnutChartOptionWithData(values)
	opt.Title.Offset = charts.OffsetCenter
	opt.Padding = charts.NewBoxEqual(10).WithBottom(15)
	opt.Legend = charts.LegendOption{
		Show: charts.Ptr(false),
		SeriesNames: []string{
			"Direct", "Search Engine", "Referral", "Email", "Video Ads",
		},
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        1400,
		Height:       400,
	})
	segmentGapPainter := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	labelsInsidePainter := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        400,
		Height:       400,
	})
	style1Painter := p.Child(charts.PainterBoxOption(charts.NewBox(0, 0, 600, 400)))
	opt.Title.Text = "Labels Outside"
	opt.SegmentGap = 24
	if err := segmentGapPainter.DoughnutChart(opt); err != nil {
		panic(err)
	}
	if err := style1Painter.DoughnutChart(opt); err != nil {
		panic(err)
	}
	style2Painter := p.Child(charts.PainterBoxOption(charts.NewBox(600, 0, 1000, 400)))
	opt.Title.Text = "Labels Inside"
	opt.SegmentGap = 0
	opt.CenterValues = "labels"
	opt.RadiusCenter = "40%"
	if err := labelsInsidePainter.DoughnutChart(opt); err != nil {
		panic(err)
	}
	if err := style2Painter.DoughnutChart(opt); err != nil {
		panic(err)
	}
	style3Painter := p.Child(charts.PainterBoxOption(charts.NewBox(1000, 0, 1400, 400)))
	opt.Title.Text = "Legend"
	opt.CenterValues = "sum"
	opt.CenterValuesFontStyle.FontSize = 12
	opt.CenterValuesFontStyle.FontColor = charts.ColorDarkGray
	opt.ValueFormatter = func(f float64) string {
		return "Total Response: " + charts.FormatValueHumanizeShort(f, 2, false)
	}
	opt.SegmentGap = 8
	opt.RadiusCenter = "32%"
	opt.Legend.Show = charts.Ptr(true)
	opt.Legend.Offset = charts.OffsetStr{Top: charts.PositionBottom}
	opt.Legend.OverlayChart = charts.Ptr(true)
	for i := range opt.SeriesList {
		opt.SeriesList[i].Label.Show = charts.Ptr(false)
	}
	if err := style3Painter.DoughnutChart(opt); err != nil {
		panic(err)
	}

	if err := writeChart(p, "doughnut-chart-2-styles.png"); err != nil {
		panic(err)
	}
	if err := writeChart(segmentGapPainter, "doughnut-chart-2-styles-gap.png"); err != nil {
		panic(err)
	}
	if err := writeChart(labelsInsidePainter, "doughnut-chart-2-styles-labels.png"); err != nil {
		panic(err)
	}
}

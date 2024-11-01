package main

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func main() {
	barWidth := 80

	var (
		colorWhite          = drawing.Color{R: 241, G: 241, B: 241, A: 255}
		colorMariner        = drawing.Color{R: 60, G: 100, B: 148, A: 255}
		colorLightSteelBlue = drawing.Color{R: 182, G: 195, B: 220, A: 255}
		colorPoloBlue       = drawing.Color{R: 126, G: 155, B: 200, A: 255}
		colorSteelBlue      = drawing.Color{R: 73, G: 120, B: 177, A: 255}
	)

	stackedBarChart := chartdraw.StackedBarChart{
		Title:      "Quarterly Sales",
		TitleStyle: chartdraw.Shown(),
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top: 75,
			},
		},
		Width:        800,
		Height:       600,
		XAxis:        chartdraw.Shown(),
		YAxis:        chartdraw.Shown(),
		BarSpacing:   40,
		IsHorizontal: true,
		Bars: []chartdraw.StackedBar{
			{
				Name:  "Q1",
				Width: barWidth,
				Values: []chartdraw.Value{
					{
						Label: "32K",
						Value: 32,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "46K",
						Value: 46,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "48K",
						Value: 48,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "42K",
						Value: 42,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
				},
			},
			{
				Name:  "Q2",
				Width: barWidth,
				Values: []chartdraw.Value{
					{
						Label: "45K",
						Value: 45,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "60K",
						Value: 60,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "62K",
						Value: 62,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "53K",
						Value: 53,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
				},
			},
			{
				Name:  "Q3",
				Width: barWidth,
				Values: []chartdraw.Value{
					{
						Label: "54K",
						Value: 54,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "58K",
						Value: 58,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "55K",
						Value: 55,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "47K",
						Value: 47,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
				},
			},
			{
				Name:  "Q4",
				Width: barWidth,
				Values: []chartdraw.Value{
					{
						Label: "46K",
						Value: 46,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorMariner,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "70K",
						Value: 70,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorLightSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "74K",
						Value: 74,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorPoloBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
					{
						Label: "60K",
						Value: 60,
						Style: chartdraw.Style{
							StrokeWidth: .01,
							FillColor:   colorSteelBlue,
							FontStyle: chartdraw.FontStyle{
								FontColor: colorWhite,
							},
						},
					},
				},
			},
		},
	}

	pngFile, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	if err := stackedBarChart.Render(chartdraw.PNG, pngFile); err != nil {
		panic(err)
	}

	if err := pngFile.Close(); err != nil {
		panic(err)
	}
}

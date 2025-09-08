package main

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	sbc := chartdraw.StackedBarChart{
		Title: "Test Stacked Bar Chart",
		Background: chartdraw.Style{
			Padding: chartdraw.Box{
				Top: 40,
			},
		},
		Height: 512,
		Bars: []chartdraw.StackedBar{
			{
				Name: "This is a very long string to test word break wrapping.",
				Values: []chartdraw.Value{
					{Value: 5, Label: "Blue"},
					{Value: 5, Label: "Green"},
					{Value: 4, Label: "Gray"},
					{Value: 3, Label: "Orange"},
					{Value: 3, Label: "Test"},
					{Value: 2, Label: "??"},
					{Value: 1, Label: "!!"},
				},
			},
			{
				Name: "Test",
				Values: []chartdraw.Value{
					{Value: 10, Label: "Blue"},
					{Value: 5, Label: "Green"},
					{Value: 1, Label: "Gray"},
				},
			},
			{
				Name: "Test 2",
				Values: []chartdraw.Value{
					{Value: 10, Label: "Blue"},
					{Value: 6, Label: "Green"},
					{Value: 4, Label: "Gray"},
				},
			},
		},
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = sbc.Render(chartdraw.PNG, f)
}

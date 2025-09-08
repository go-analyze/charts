package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
)

func main() {
	pie := chartdraw.DonutChart{
		Width:  512,
		Height: 512,
		Values: []chartdraw.Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 4, Label: "Gray"},
			{Value: 4, Label: "Orange"},
			{Value: 3, Label: "Deep Blue"},
			{Value: 3, Label: "test"},
		},
	}

	f, _ := os.Create("output.png")
	defer func() { _ = f.Close() }()
	_ = pie.Render(chartdraw.PNG, f)
}

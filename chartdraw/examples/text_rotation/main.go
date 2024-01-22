package main

//go:generate go run main.go

import (
	"os"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func main() {
	f, _ := chartdraw.GetDefaultFont()
	r, _ := chartdraw.PNG(1024, 1024)

	chartdraw.Draw.Text(r, "Test", 64, 64, chartdraw.Style{
		FontColor: drawing.ColorBlack,
		FontSize:  18,
		Font:      f,
	})

	chartdraw.Draw.Text(r, "Test", 64, 64, chartdraw.Style{
		FontColor:           drawing.ColorBlack,
		FontSize:            18,
		Font:                f,
		TextRotationDegrees: 45.0,
	})

	tb := chartdraw.Draw.MeasureText(r, "Test", chartdraw.Style{
		FontColor: drawing.ColorBlack,
		FontSize:  18,
		Font:      f,
	}).Shift(64, 64)

	tbc := tb.Corners().Rotate(45)

	chartdraw.Draw.BoxCorners(r, tbc, chartdraw.Style{
		StrokeColor: drawing.ColorRed,
		StrokeWidth: 2,
	})

	tbcb := tbc.Box()
	chartdraw.Draw.Box(r, tbcb, chartdraw.Style{
		StrokeColor: drawing.ColorBlue,
		StrokeWidth: 2,
	})

	file, _ := os.Create("output.png")
	defer file.Close()
	r.Save(file)
}

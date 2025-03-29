package chartdraw

import (
	"bytes"
	"testing"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func BenchmarkRaterCircle(b *testing.B) {
	testRadius := []float64{400, 200, 128, 64, 16, 8, 2}
	bb := &bytes.Buffer{}
	for i := 0; i < b.N; i++ {
		png := PNG(800, 800)
		jpg := JPG(800, 800)

		var flip bool
		for _, r := range testRadius {
			color := drawing.ColorNavy
			if flip {
				color = drawing.ColorThistle
				flip = false
			} else {
				flip = true
			}

			png.SetFillColor(color)
			png.Circle(r, 400, 400)
			png.Fill()

			jpg.SetFillColor(color)
			jpg.Circle(r, 400, 400)
			jpg.Fill()
		}

		bb.Reset()
		_ = png.Save(bb)
		bb.Reset()
		_ = jpg.Save(bb)
	}
}

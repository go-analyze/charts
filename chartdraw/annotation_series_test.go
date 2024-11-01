package chartdraw

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestAnnotationSeriesMeasure(t *testing.T) {
	as := AnnotationSeries{
		Annotations: []Value2{
			{XValue: 1.0, YValue: 1.0, Label: "1.0"},
			{XValue: 2.0, YValue: 2.0, Label: "2.0"},
			{XValue: 3.0, YValue: 3.0, Label: "3.0"},
			{XValue: 4.0, YValue: 4.0, Label: "4.0"},
		},
	}

	r, err := PNG(110, 110)
	require.NoError(t, err)

	f, err := GetDefaultFont()
	require.NoError(t, err)

	xrange := &ContinuousRange{
		Min:    1.0,
		Max:    4.0,
		Domain: 100,
	}
	yrange := &ContinuousRange{
		Min:    1.0,
		Max:    4.0,
		Domain: 100,
	}

	cb := Box{
		Top:    5,
		Left:   5,
		Right:  105,
		Bottom: 105,
	}
	sd := Style{
		FontStyle: FontStyle{
			FontSize: 10.0,
			Font:     f,
		},
	}

	box := as.Measure(r, cb, xrange, yrange, sd)
	assert.False(t, box.IsZero())
	assert.Equal(t, -5, box.Top)
	assert.Equal(t, 5, box.Left)
	assert.Equal(t, 146, box.Right) //the top,left annotation sticks up 5px and out ~44px.
	assert.Equal(t, 115, box.Bottom)
}

func TestAnnotationSeriesRender(t *testing.T) {
	as := AnnotationSeries{
		Style: Style{
			FillColor:   drawing.ColorWhite,
			StrokeColor: drawing.ColorBlack,
		},
		Annotations: []Value2{
			{XValue: 1.0, YValue: 1.0, Label: "1.0"},
			{XValue: 2.0, YValue: 2.0, Label: "2.0"},
			{XValue: 3.0, YValue: 3.0, Label: "3.0"},
			{XValue: 4.0, YValue: 4.0, Label: "4.0"},
		},
	}

	r, err := PNG(110, 110)
	require.NoError(t, err)

	f, err := GetDefaultFont()
	require.NoError(t, err)

	xrange := &ContinuousRange{
		Min:    1.0,
		Max:    4.0,
		Domain: 100,
	}
	yrange := &ContinuousRange{
		Min:    1.0,
		Max:    4.0,
		Domain: 100,
	}

	cb := Box{
		Top:    5,
		Left:   5,
		Right:  105,
		Bottom: 105,
	}
	sd := Style{
		FontStyle: FontStyle{
			FontSize: 10.0,
			Font:     f,
		},
	}

	as.Render(r, cb, xrange, yrange, sd)

	rr, isRaster := r.(*rasterRenderer)
	assert.True(t, isRaster)
	assert.NotNil(t, rr)

	c := rr.i.At(38, 70)
	converted, isRGBA := color.RGBAModel.Convert(c).(color.RGBA)
	assert.True(t, isRGBA)
	assert.Equal(t, uint8(0), converted.R)
	assert.Equal(t, uint8(0), converted.G)
	assert.Equal(t, uint8(0), converted.B)
}

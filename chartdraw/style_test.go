package chartdraw

import (
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestStyleIsZero(t *testing.T) {
	zero := Style{}
	assert.True(t, zero.IsZero())

	strokeColor := Style{StrokeColor: drawing.ColorWhite}
	assert.False(t, strokeColor.IsZero())

	fillColor := Style{FillColor: drawing.ColorWhite}
	assert.False(t, fillColor.IsZero())

	strokeWidth := Style{StrokeWidth: 5.0}
	assert.False(t, strokeWidth.IsZero())

	fontSize := Style{FontStyle: FontStyle{FontSize: 12.0}}
	assert.False(t, fontSize.IsZero())

	fontColor := Style{FontStyle: FontStyle{FontColor: drawing.ColorWhite}}
	assert.False(t, fontColor.IsZero())

	font := Style{FontStyle: FontStyle{Font: &truetype.Font{}}}
	assert.False(t, font.IsZero())
}

func TestStyleGetStrokeColor(t *testing.T) {
	unset := Style{}
	assert.Equal(t, drawing.ColorTransparent, unset.GetStrokeColor())
	assert.Equal(t, drawing.ColorWhite, unset.GetStrokeColor(drawing.ColorWhite))

	set := Style{StrokeColor: drawing.ColorWhite}
	assert.Equal(t, drawing.ColorWhite, set.GetStrokeColor())
	assert.Equal(t, drawing.ColorWhite, set.GetStrokeColor(drawing.ColorBlack))
}

func TestStyleGetFillColor(t *testing.T) {
	unset := Style{}
	assert.Equal(t, drawing.ColorTransparent, unset.GetFillColor())
	assert.Equal(t, drawing.ColorWhite, unset.GetFillColor(drawing.ColorWhite))

	set := Style{FillColor: drawing.ColorWhite}
	assert.Equal(t, drawing.ColorWhite, set.GetFillColor())
	assert.Equal(t, drawing.ColorWhite, set.GetFillColor(drawing.ColorBlack))
}

func TestStyleGetStrokeWidth(t *testing.T) {
	unset := Style{}
	assert.Equal(t, DefaultStrokeWidth, unset.GetStrokeWidth())
	assert.Equal(t, DefaultStrokeWidth+1, unset.GetStrokeWidth(DefaultStrokeWidth+1))

	set := Style{StrokeWidth: DefaultStrokeWidth + 2}
	assert.Equal(t, DefaultStrokeWidth+2, set.GetStrokeWidth())
	assert.Equal(t, DefaultStrokeWidth+2, set.GetStrokeWidth(DefaultStrokeWidth+1))
}

func TestStyleGetFontSize(t *testing.T) {
	unset := Style{}
	assert.Equal(t, DefaultFontSize, unset.GetFontSize())
	assert.Equal(t, DefaultFontSize+1, unset.GetFontSize(DefaultFontSize+1))

	set := Style{FontStyle: FontStyle{FontSize: DefaultFontSize + 2}}
	assert.Equal(t, DefaultFontSize+2, set.GetFontSize())
	assert.Equal(t, DefaultFontSize+2, set.GetFontSize(DefaultFontSize+1))
}

func TestStyleGetFontColor(t *testing.T) {
	unset := Style{}
	assert.Equal(t, drawing.ColorTransparent, unset.GetFontColor())
	assert.Equal(t, drawing.ColorWhite, unset.GetFontColor(drawing.ColorWhite))

	set := Style{FontStyle: FontStyle{FontColor: drawing.ColorWhite}}
	assert.Equal(t, drawing.ColorWhite, set.GetFontColor())
	assert.Equal(t, drawing.ColorWhite, set.GetFontColor(drawing.ColorBlack))
}

func TestStyleGetFont(t *testing.T) {
	f, err := GetDefaultFont()
	require.NoError(t, err)

	unset := Style{}
	require.Nil(t, unset.GetFont())
	assert.Equal(t, f, unset.GetFont(f))

	set := Style{FontStyle: FontStyle{Font: f}}
	require.NotNil(t, set.GetFont())
}

func TestStyleGetPadding(t *testing.T) {
	unset := Style{}
	assert.True(t, unset.GetPadding().IsZero())
	assert.False(t, unset.GetPadding(DefaultBackgroundPadding).IsZero())
	assert.Equal(t, DefaultBackgroundPadding, unset.GetPadding(DefaultBackgroundPadding))

	set := Style{Padding: DefaultBackgroundPadding}
	assert.False(t, set.GetPadding().IsZero())
	assert.Equal(t, DefaultBackgroundPadding, set.GetPadding())
	assert.Equal(t, DefaultBackgroundPadding, set.GetPadding(Box{
		Top:    DefaultBackgroundPadding.Top + 1,
		Left:   DefaultBackgroundPadding.Left + 1,
		Right:  DefaultBackgroundPadding.Right + 1,
		Bottom: DefaultBackgroundPadding.Bottom + 1,
	}))
}

func TestStyleWithDefaultsFrom(t *testing.T) {
	f, err := GetDefaultFont()
	require.NoError(t, err)

	unset := Style{}
	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontStyle: FontStyle{
			FontColor: drawing.ColorWhite,
			Font:      f,
		},
		Padding: DefaultBackgroundPadding,
	}

	coalesced := unset.InheritFrom(set)
	assert.Equal(t, set, coalesced)
}

func TestStyleGetStrokeOptions(t *testing.T) {
	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontStyle: FontStyle{
			FontColor: drawing.ColorWhite,
		},
		Padding: DefaultBackgroundPadding,
	}
	svgStroke := set.GetStrokeOptions()
	assert.False(t, svgStroke.StrokeColor.IsZero())
	assert.NotZero(t, svgStroke.StrokeWidth)
	assert.True(t, svgStroke.FillColor.IsZero())
	assert.True(t, svgStroke.FontColor.IsZero())
}

func TestStyleGetFillOptions(t *testing.T) {
	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontStyle: FontStyle{
			FontColor: drawing.ColorWhite,
		},
		Padding: DefaultBackgroundPadding,
	}
	svgFill := set.GetFillOptions()
	assert.False(t, svgFill.FillColor.IsZero())
	assert.Zero(t, svgFill.StrokeWidth)
	assert.True(t, svgFill.StrokeColor.IsZero())
	assert.True(t, svgFill.FontColor.IsZero())
}

func TestStyleGetFillAndStrokeOptions(t *testing.T) {
	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontStyle: FontStyle{
			FontColor: drawing.ColorWhite,
		},
		Padding: DefaultBackgroundPadding,
	}
	svgFillAndStroke := set.GetFillAndStrokeOptions()
	assert.False(t, svgFillAndStroke.FillColor.IsZero())
	assert.NotZero(t, svgFillAndStroke.StrokeWidth)
	assert.False(t, svgFillAndStroke.StrokeColor.IsZero())
	assert.True(t, svgFillAndStroke.FontColor.IsZero())
}

func TestStyleGetTextOptions(t *testing.T) {
	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontStyle: FontStyle{
			FontColor: drawing.ColorWhite,
		},
		Padding: DefaultBackgroundPadding,
	}
	svgStroke := set.GetTextOptions()
	assert.True(t, svgStroke.StrokeColor.IsZero())
	assert.Zero(t, svgStroke.StrokeWidth)
	assert.True(t, svgStroke.FillColor.IsZero())
	assert.False(t, svgStroke.FontColor.IsZero())
}

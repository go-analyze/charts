package chartdraw

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestVectorRendererPath(t *testing.T) {
	t.Parallel()

	vr := SVG(100, 100)

	typed, isTyped := vr.(*vectorRenderer)
	assert.True(t, isTyped)

	typed.MoveTo(0, 0)
	typed.LineTo(100, 100)
	typed.LineTo(0, 100)
	typed.Close()
	typed.FillStroke()

	buffer := bytes.NewBuffer([]byte{})
	require.NoError(t, typed.Save(buffer))

	raw := buffer.String()
	assert.True(t, strings.HasPrefix(raw, "<svg"))
	assert.True(t, strings.HasSuffix(raw, "</svg>"))
}

func TestVectorRendererMeasureText(t *testing.T) {
	t.Parallel()

	vr := SVG(100, 100)

	vr.SetDPI(DefaultDPI)
	vr.SetFont(GetDefaultFont())
	vr.SetFontSize(12.0)

	tb := vr.MeasureText("Ljp")
	assert.Equal(t, 21, tb.Width())
	assert.Equal(t, 15, tb.Height())
}

func TestCanvasStyleSVG(t *testing.T) {
	t.Parallel()

	set := Style{
		StrokeColor: drawing.ColorWhite,
		StrokeWidth: 5.0,
		FillColor:   drawing.ColorWhite,
		FontStyle: FontStyle{
			FontColor: drawing.ColorWhite,
			Font:      GetDefaultFont(),
		},
		Padding: DefaultBackgroundPadding,
	}

	canvas := &canvas{dpi: DefaultDPI}

	svgString := canvas.styleAsSVG(set)
	assert.NotEmpty(t, svgString)
	assert.True(t, strings.HasPrefix(svgString, "style=\""))
	assert.True(t, strings.Contains(svgString, "stroke:rgba(255,255,255,1.0)"))
	assert.True(t, strings.Contains(svgString, "stroke-width:5"))
	assert.True(t, strings.Contains(svgString, "fill:rgba(255,255,255,1.0)"))
	assert.True(t, strings.HasSuffix(svgString, "\""))
}

func TestCanvasClassSVG(t *testing.T) {
	t.Parallel()

	set := Style{
		ClassName: "test-class",
	}

	canvas := &canvas{dpi: DefaultDPI}

	assert.Equal(t, "class=\"test-class\"", canvas.styleAsSVG(set))
}

func TestCanvasCustomInlineStylesheet(t *testing.T) {
	t.Parallel()

	b := strings.Builder{}

	canvas := &canvas{
		w:   &b,
		css: ".background { fill: red }",
	}

	canvas.Start(200, 200)

	assert.Contains(t, b.String(), fmt.Sprintf(`<style type="text/css"><![CDATA[%s]]></style>`, canvas.css))
}

func TestCanvasCustomInlineStylesheetWithNonce(t *testing.T) {
	t.Parallel()

	b := strings.Builder{}

	canvas := &canvas{
		w:     &b,
		css:   ".background { fill: red }",
		nonce: "RAND0MSTRING",
	}

	canvas.Start(200, 200)

	assert.Contains(t, b.String(), fmt.Sprintf(`<style type="text/css" nonce="%s"><![CDATA[%s]]></style>`, canvas.nonce, canvas.css))
}

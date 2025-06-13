package chartdraw

import (
	"bytes"
	"fmt"
	"math"
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
	assert.Equal(t, 16, tb.Height())
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
			FontSize:  12,
		},
		Padding: DefaultBackgroundPadding,
	}

	var bb bytes.Buffer
	styleAsSVG(&bb, set, DefaultDPI, false)
	svgString := bb.String()
	assert.NotEmpty(t, svgString)
	assert.True(t, strings.HasPrefix(svgString, "style=\""))
	assert.Contains(t, svgString, "stroke:white")
	assert.Contains(t, svgString, "stroke-width:5")
	assert.Contains(t, svgString, "fill:white")
	assert.NotContains(t, svgString, "font-size")
	assert.NotContains(t, svgString, "font-family")
	assert.True(t, strings.HasSuffix(svgString, "\""))

	bb.Reset()
	styleAsSVG(&bb, set, DefaultDPI, true)
	svgString = bb.String()
	assert.True(t, strings.HasPrefix(svgString, "style=\""))
	assert.Contains(t, svgString, "stroke:white")
	assert.Contains(t, svgString, "stroke-width:5")
	assert.Contains(t, svgString, "fill:white")
	assert.Contains(t, svgString, "font-size")
	assert.Contains(t, svgString, "font-family")
	assert.True(t, strings.HasSuffix(svgString, "\""))
}

func TestCanvasClassSVG(t *testing.T) {
	t.Parallel()

	set := Style{
		ClassName: "test-class",
	}

	var bb bytes.Buffer
	styleAsSVG(&bb, set, DefaultDPI, false)
	assert.Equal(t, "class=\"test-class\"", bb.String())
}

func TestCanvasCustomInlineStylesheet(t *testing.T) {
	t.Parallel()

	b := strings.Builder{}

	canvas := &canvas{
		w:   &b,
		bb:  bytes.NewBuffer(make([]byte, 0, 80)),
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
		bb:    bytes.NewBuffer(make([]byte, 0, 80)),
		css:   ".background { fill: red }",
		nonce: "RAND0MSTRING",
	}

	canvas.Start(200, 200)

	assert.Contains(t, b.String(), fmt.Sprintf(`<style type="text/css" nonce="%s"><![CDATA[%s]]></style>`, canvas.nonce, canvas.css))
}

func TestSVGWithCSS(t *testing.T) {
	t.Parallel()

	maker := SVGWithCSS(".cls{fill:red}", "nonce")
	r := maker(10, 10)
	vr, ok := r.(*vectorRenderer)
	require.True(t, ok)

	b := bytes.Buffer{}
	require.NoError(t, vr.Save(&b))
	out := b.String()
	assert.Contains(t, out, "nonce=\"nonce\"")
	assert.Contains(t, out, ".cls{fill:red}")
}

func TestCanvasBasicElements(t *testing.T) {
	t.Parallel()

	b := strings.Builder{}
	c := &canvas{w: &b, bb: bytes.NewBuffer(make([]byte, 0, 80))}
	c.Start(50, 50)

	c.Path([]string{"M 0 0", "L 10 10"}, Style{StrokeDashArray: []float64{1, 2}, StrokeWidth: 2, StrokeColor: drawing.ColorBlack})
	c.Text(5, 5, "hi", Style{FontStyle: FontStyle{Font: GetDefaultFont(), FontSize: 10}})
	c.Circle(5, 5, 3, Style{FillColor: drawing.ColorRed})
	c.End()

	out := b.String()
	assert.Contains(t, out, "stroke-dasharray=\"1.0, 2.0\"")
	assert.Contains(t, out, "<text")
	assert.Contains(t, out, "<circle")
	assert.True(t, strings.HasSuffix(out, "</svg>"))
}

func TestFormatFloatMinimized(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1", formatFloatMinimized(1))
	assert.Equal(t, "1.2", formatFloatMinimized(1.20))
	assert.Equal(t, "2.5", formatFloatMinimized(2.50))
}

func TestVectorRendererTextRotation(t *testing.T) {
	t.Parallel()

	vr := SVG(20, 20).(*vectorRenderer)
	vr.SetClassName("cls")
	vr.SetStrokeColor(drawing.ColorBlack)
	vr.SetFillColor(drawing.ColorRed)
	vr.SetTextRotation(math.Pi / 2)
	vr.Text("A", 10, 10)
	vr.ClearTextRotation()
	vr.Text("B", 10, 15)

	buf := bytes.Buffer{}
	require.NoError(t, vr.Save(&buf))
	out := buf.String()
	assert.Contains(t, out, "class=\"cls\"")
	assert.Contains(t, out, "rotate(90.00")
	assert.Contains(t, out, "B</text>")
}

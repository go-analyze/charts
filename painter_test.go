package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestPainterOption(t *testing.T) {
	t.Parallel()

	font := &truetype.Font{}
	d, err := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        800,
		Height:       600,
	},
		PainterBoxOption(Box{Right: 400, Bottom: 300}),
		PainterPaddingOption(Box{Left: 1, Top: 2, Right: 3, Bottom: 4}),
		PainterFontOption(font),
		PainterStyleOption(chartdraw.Style{ClassName: "test"}),
	)
	require.NoError(t, err)
	assert.Equal(t, Box{
		Left:   1,
		Top:    2,
		Right:  397,
		Bottom: 296,
	}, d.box)
	assert.Equal(t, font, d.font)
	assert.Equal(t, "test", d.style.ClassName)
}

func TestPainter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		result string
	}{
		{
			name: "moveTo_lineTo",
			fn: func(p *Painter) {
				p.MoveTo(1, 1)
				p.LineTo(2, 2)
				p.Stroke()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 6 11\nL 7 12\" style=\"stroke-width:0;stroke:none;fill:none\"/></svg>",
		},
		{
			name: "circle",
			fn: func(p *Painter) {
				p.Circle(5, 2, 3)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"7\" cy=\"13\" r=\"5\" style=\"stroke-width:0;stroke:none;fill:none\"/></svg>",
		},
		{
			name: "text",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"8\" y=\"16\" style=\"stroke-width:0;stroke:none;fill:none;font-family:'Roboto Medium',sans-serif\">hello world!</text></svg>",
		},
		{
			name: "line",
			fn: func(p *Painter) {
				p.SetDrawingStyle(chartdraw.Style{
					StrokeColor: drawing.ColorBlack,
					StrokeWidth: 1,
				})
				p.LineStroke([]Point{
					{X: 1, Y: 2},
					{X: 3, Y: 4},
				})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 6 12\nL 8 14\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/></svg>",
		},
		{
			name: "background",
			fn: func(p *Painter) {
				p.SetBackground(400, 300, chartdraw.ColorWhite)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 0 0\nL 400 0\nL 400 300\nL 0 300\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "arc",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: drawing.ColorBlack,
					FillColor:   drawing.ColorBlue,
				})
				p.ArcTo(100, 100, 100, 100, 0, math.Pi/2)
				p.Close()
				p.FillStroke()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 205 110\nA 100 100 90.00 0 1 105 210\nZ\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,255,1.0)\"/></svg>",
		},
		{
			name: "pin",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: Color{R: 84, G: 112, B: 198, A: 255},
					FillColor:   Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.Pin(30, 30, 30)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 32 47\nA 15 15 330.00 1 1 38 47\nL 35 33\nZ\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"M 20 33\nQ35,70 50,33\nZ\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/></svg>",
		},
		{
			name: "arrow_left",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: Color{R: 84, G: 112, B: 198, A: 255},
					FillColor:   Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.ArrowLeft(30, 30, 16, 10)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 51 35\nL 35 40\nL 51 45\nL 46 40\nL 51 35\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/></svg>",
		},
		{
			name: "arrow_right",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: Color{R: 84, G: 112, B: 198, A: 255},
					FillColor:   Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.ArrowRight(30, 30, 16, 10)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 19 35\nL 35 40\nL 19 45\nL 24 40\nL 19 35\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/></svg>",
		},
		{
			name: "arrow_top",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: Color{R: 84, G: 112, B: 198, A: 255},
					FillColor:   Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.ArrowTop(30, 30, 10, 16)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 30 40\nL 35 24\nL 40 40\nL 35 35\nL 30 40\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/></svg>",
		},
		{
			name: "arrow_bottom",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: Color{R: 84, G: 112, B: 198, A: 255},
					FillColor:   Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.ArrowBottom(30, 30, 10, 16)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 30 24\nL 35 40\nL 40 24\nL 35 30\nL 30 24\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/></svg>",
		},
		{
			name: "mark_line",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth:     1,
					StrokeColor:     Color{R: 84, G: 112, B: 198, A: 255},
					FillColor:       Color{R: 84, G: 112, B: 198, A: 255},
					StrokeDashArray: []float64{4, 2},
				})
				p.MarkLine(0, 20, 300)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"8\" cy=\"30\" r=\"3\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 14 30\nL 289 30\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 289 25\nL 305 30\nL 289 35\nL 294 30\nL 289 25\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/></svg>",
		},
		{
			name: "polygon",
			fn: func(p *Painter) {
				p.SetStyle(chartdraw.Style{
					StrokeWidth: 1,
					StrokeColor: Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.Polygon(Point{X: 100, Y: 100}, 50, 6)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 105 60\nL 148 85\nL 148 134\nL 105 160\nL 62 135\nL 62 86\nL 105 60\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:none\"/></svg>",
		},
		{
			name: "fill_area",
			fn: func(p *Painter) {
				p.SetDrawingStyle(chartdraw.Style{
					FillColor: Color{R: 84, G: 112, B: 198, A: 255},
				})
				p.FillArea([]Point{
					{X: 0, Y: 0},
					{X: 0, Y: 100},
					{X: 100, Y: 100},
					{X: 0, Y: 0},
				})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 5 110\nL 105 110\nL 5 10\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/></svg>",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			d, err := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(chartdraw.Box{Left: 5, Top: 10}))
			require.NoError(t, err)
			tt.fn(d)
			data, err := d.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, string(data))
		})
	}
}

func TestRoundedRect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		result string
	}{
		{
			name: "round_fully",
			fn: func(p *Painter) {
				p.RoundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, true, true)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 10\nL 25 10\nL 25 10\nA 5 5 90.00 0 1 30 15\nL 30 145\nL 30 145\nA 5 5 90.00 0 1 25 150\nL 15 150\nL 15 150\nA 5 5 90.00 0 1 10 145\nL 10 15\nL 10 15\nA 5 5 90.00 0 1 15 10\nZ\" style=\"stroke-width:1;stroke:rgba(0,0,255,1.0);fill:rgba(0,0,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,255,1.0);fill:rgba(0,0,255,1.0)\"/></svg>",
		},
		{
			name: "square_top",
			fn: func(p *Painter) {
				p.RoundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, false, true)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 10 10\nL 30 10\nL 30 145\nL 30 145\nA 5 5 90.00 0 1 25 150\nL 15 150\nL 15 150\nA 5 5 90.00 0 1 10 145\nL 10 10\nZ\" style=\"stroke-width:1;stroke:rgba(0,0,255,1.0);fill:rgba(0,0,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,255,1.0);fill:rgba(0,0,255,1.0)\"/></svg>",
		},
		{
			name: "square_bottom",
			fn: func(p *Painter) {
				p.RoundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, true, false)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 10\nL 25 10\nL 25 10\nA 5 5 90.00 0 1 30 15\nL 30 150\nL 10 150\nL 10 15\nL 10 15\nA 5 5 90.00 0 1 15 10\nZ\" style=\"stroke-width:1;stroke:rgba(0,0,255,1.0);fill:rgba(0,0,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,255,1.0);fill:rgba(0,0,255,1.0)\"/></svg>",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				Width:        400,
				Height:       300,
				OutputFormat: ChartOutputSVG,
			})
			require.NoError(t, err)
			p.OverrideDrawingStyle(chartdraw.Style{
				FillColor:   drawing.ColorBlue,
				StrokeWidth: 1,
				StrokeColor: drawing.ColorBlue,
			})
			tc.fn(p)
			buf, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.result, string(buf))
		})
	}
}

func TestPainterTextFit(t *testing.T) {
	t.Parallel()

	p, err := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        400,
		Height:       300,
	})
	require.NoError(t, err)
	style := FontStyle{
		FontSize:  12,
		FontColor: chartdraw.ColorBlack,
		Font:      GetDefaultFont(),
	}
	p.SetStyle(chartdraw.Style{FontStyle: style})
	box := p.TextFit("Hello World!", 0, 20, 80)
	assert.Equal(t, chartdraw.Box{Right: 45, Bottom: 35}, box)

	box = p.TextFit("Hello World!", 0, 100, 200)
	assert.Equal(t, chartdraw.Box{Right: 84, Bottom: 15}, box)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"0\" y=\"20\" style=\"stroke-width:0;stroke:none;fill:rgba(51,51,51,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello</text><text x=\"0\" y=\"40\" style=\"stroke-width:0;stroke:none;fill:rgba(51,51,51,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World!</text><text x=\"0\" y=\"100\" style=\"stroke-width:0;stroke:none;fill:rgba(51,51,51,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello World!</text></svg>", string(buf))
}

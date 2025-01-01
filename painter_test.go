package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestPainterOption(t *testing.T) {
	t.Parallel()

	d, err := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        800,
		Height:       600,
	},
		PainterBoxOption(Box{Right: 400, Bottom: 300}),
		PainterPaddingOption(Box{Left: 1, Top: 2, Right: 3, Bottom: 4}),
	)
	require.NoError(t, err)
	assert.Equal(t, Box{
		Left:   1,
		Top:    2,
		Right:  397,
		Bottom: 296,
	}, d.box)
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
		{
			name: "child_chart",
			fn: func(p *Painter) {
				_, _ = NewLineChart(p, makeMinimalLineChartOption()).Render()
				p = p.Child(PainterBoxOption(chartdraw.NewBox(0, 200, 400, 200)))
				opt := makeMinimalLineChartOption()
				opt.Theme = GetDefaultTheme().WithBackgroundColor(drawing.ColorFromAlphaMixedRGBA(0, 0, 0, 0))
				_, _ = NewLineChart(p, opt).Render()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 0 0\nL 395 0\nL 395 290\nL 0 290\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"15\" y=\"27\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"15\" y=\"53\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"15\" y=\"80\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"27\" y=\"107\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"27\" y=\"133\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"27\" y=\"160\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"27\" y=\"187\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"27\" y=\"213\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"27\" y=\"240\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"45\" y=\"267\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 64 20\nL 390 20\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 46\nL 390 46\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 73\nL 390 73\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 100\nL 390 100\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 126\nL 390 126\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 153\nL 390 153\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 180\nL 390 180\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 206\nL 390 206\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 64 233\nL 390 233\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 87 240\nL 133 238\nL 180 244\nL 226 238\nL 273 245\nL 319 222\nL 366 225\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><path  d=\"M 87 124\nL 133 105\nL 180 110\nL 226 105\nL 273 45\nL 319 39\nL 366 40\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><path  d=\"M 0 0\nL 200 0\nL 200 200\nL 0 200\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:none\"/><text x=\"210\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"210\" y=\"33\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"210\" y=\"50\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"222\" y=\"67\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"222\" y=\"83\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"222\" y=\"100\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"222\" y=\"117\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"222\" y=\"133\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"222\" y=\"150\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"240\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 259 10\nL 390 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 26\nL 390 26\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 43\nL 390 43\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 60\nL 390 60\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 76\nL 390 76\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 93\nL 390 93\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 110\nL 390 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 126\nL 390 126\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 143\nL 390 143\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 259 148\nL 280 147\nL 302 150\nL 324 147\nL 346 151\nL 368 137\nL 390 139\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><path  d=\"M 259 75\nL 280 63\nL 302 67\nL 324 63\nL 346 26\nL 368 22\nL 390 23\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/></svg>",
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
			assertEqualSVG(t, tt.result, data)
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
			assertEqualSVG(t, tc.result, buf)
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"0\" y=\"20\" style=\"stroke-width:0;stroke:none;fill:rgba(51,51,51,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello</text><text x=\"0\" y=\"40\" style=\"stroke-width:0;stroke:none;fill:rgba(51,51,51,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World!</text><text x=\"0\" y=\"100\" style=\"stroke-width:0;stroke:none;fill:rgba(51,51,51,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello World!</text></svg>", buf)
}

func TestMultipleChartsOnPainter(t *testing.T) {
	t.Parallel()

	p, err := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        800,
		Height:       600,
	})
	require.NoError(t, err)
	p.SetBackground(800, 600, drawing.ColorWhite)
	// set the space and theme for each chart
	topCenterPainter := p.Child(PainterBoxOption(chartdraw.NewBox(0, 0, 800, 300)),
		PainterThemeOption(GetTheme(ThemeVividLight)))
	bottomLeftPainter := p.Child(PainterBoxOption(chartdraw.NewBox(300, 0, 400, 600)),
		PainterThemeOption(GetTheme(ThemeAnt)))
	bottomRightPainter := p.Child(PainterBoxOption(chartdraw.NewBox(300, 400, 800, 600)),
		PainterThemeOption(GetTheme(ThemeLight)))

	pieOpt := makeBasicPieChartOption()
	pieOpt.Legend.Show = False()
	_, err = NewPieChart(bottomLeftPainter, pieOpt).Render()
	require.NoError(t, err)
	_, err = NewBarChart(bottomRightPainter, makeBasicBarChartOption()).Render()
	require.NoError(t, err)
	_, err = NewLineChart(topCenterPainter, makeBasicLineChartOption()).Render()
	require.NoError(t, err)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 800 600\"><path  d=\"M 0 0\nL 800 0\nL 800 600\nL 0 600\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 0 0\nL 400 0\nL 400 300\nL 0 300\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"122\" y=\"335\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"166\" y=\"350\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path  d=\"M 200 457\nL 200 383\nA 74 74 119.89 0 1 264 493\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgba(91,143,249,1.0);fill:rgba(91,143,249,1.0)\"/><path  d=\"M 264 420\nL 277 413\nM 277 413\nL 292 413\" style=\"stroke-width:1;stroke:rgba(91,143,249,1.0);fill:rgba(91,143,249,1.0)\"/><text x=\"295\" y=\"418\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Search Engine: 33.3%</text><path  d=\"M 200 457\nL 264 493\nA 74 74 84.08 0 1 170 524\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgba(90,216,166,1.0);fill:rgba(90,216,166,1.0)\"/><path  d=\"M 222 527\nL 227 541\nM 227 541\nL 242 541\" style=\"stroke-width:1;stroke:rgba(90,216,166,1.0);fill:rgba(90,216,166,1.0)\"/><text x=\"245\" y=\"546\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Direct: 23.35%</text><path  d=\"M 200 457\nL 170 524\nA 74 74 66.35 0 1 127 457\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgba(93,112,146,1.0);fill:rgba(93,112,146,1.0)\"/><path  d=\"M 138 497\nL 126 505\nM 126 505\nL 111 505\" style=\"stroke-width:1;stroke:rgba(93,112,146,1.0);fill:rgba(93,112,146,1.0)\"/><text x=\"27\" y=\"510\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Email: 18.43%</text><path  d=\"M 200 457\nL 127 457\nA 74 74 55.37 0 1 159 396\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgba(246,189,22,1.0);fill:rgba(246,189,22,1.0)\"/><path  d=\"M 135 423\nL 122 416\nM 122 416\nL 107 416\" style=\"stroke-width:1;stroke:rgba(246,189,22,1.0);fill:rgba(246,189,22,1.0)\"/><text x=\"-4\" y=\"421\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Union Ads: 15.37%</text><path  d=\"M 200 457\nL 159 396\nA 74 74 34.32 0 1 200 383\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgba(111,94,249,1.0);fill:rgba(111,94,249,1.0)\"/><path  d=\"M 179 387\nL 174 372\nM 174 372\nL 159 372\" style=\"stroke-width:1;stroke:rgba(111,94,249,1.0);fill:rgba(111,94,249,1.0)\"/><text x=\"56\" y=\"377\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Video Ads: 9.53%</text><path  d=\"M 0 0\nL 400 0\nL 400 300\nL 0 300\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"410\" y=\"317\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">189</text><text x=\"410\" y=\"344\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">168</text><text x=\"410\" y=\"372\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">147</text><text x=\"410\" y=\"400\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">126</text><text x=\"410\" y=\"428\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">105</text><text x=\"419\" y=\"455\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84</text><text x=\"419\" y=\"483\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">63</text><text x=\"419\" y=\"511\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"419\" y=\"539\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">21</text><text x=\"428\" y=\"567\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 447 310\nL 790 310\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 337\nL 790 337\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 365\nL 790 365\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 393\nL 790 393\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 421\nL 790 421\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 448\nL 790 448\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 476\nL 790 476\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 504\nL 790 504\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 532\nL 790 532\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 447 565\nL 447 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 509 565\nL 509 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 571 565\nL 571 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 634 565\nL 634 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 665 565\nL 665 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 727 565\nL 727 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 790 565\nL 790 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 447 560\nL 790 560\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"446\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"508\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"570\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"633\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"664\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"726\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"763\" y=\"585\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 452 558\nL 459 558\nL 459 559\nL 452 559\nL 452 558\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 480 554\nL 487 554\nL 487 559\nL 480 559\nL 480 554\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 509 551\nL 516 551\nL 516 559\nL 509 559\nL 509 551\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 537 530\nL 544 530\nL 544 559\nL 537 559\nL 537 530\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 566 527\nL 573 527\nL 573 559\nL 566 559\nL 566 527\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 594 459\nL 601 459\nL 601 559\nL 594 559\nL 594 459\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 623 381\nL 630 381\nL 630 559\nL 623 559\nL 623 381\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 652 346\nL 659 346\nL 659 559\nL 652 559\nL 652 346\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 680 517\nL 687 517\nL 687 559\nL 680 559\nL 680 517\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 709 534\nL 716 534\nL 716 559\nL 709 559\nL 709 534\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 737 552\nL 744 552\nL 744 559\nL 737 559\nL 737 552\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 766 556\nL 773 556\nL 773 559\nL 766 559\nL 766 556\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 462 557\nL 469 557\nL 469 559\nL 462 559\nL 462 557\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 490 553\nL 497 553\nL 497 559\nL 490 559\nL 490 553\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 519 549\nL 526 549\nL 526 559\nL 519 559\nL 519 549\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 547 526\nL 554 526\nL 554 559\nL 547 559\nL 547 526\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 576 523\nL 583 523\nL 583 559\nL 576 559\nL 576 523\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 604 467\nL 611 467\nL 611 559\nL 604 559\nL 604 467\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 633 328\nL 640 328\nL 640 559\nL 633 559\nL 633 328\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 662 319\nL 669 319\nL 669 559\nL 662 559\nL 662 319\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 690 496\nL 697 496\nL 697 559\nL 690 559\nL 690 496\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 719 536\nL 726 536\nL 726 559\nL 719 559\nL 719 536\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 747 553\nL 754 553\nL 754 559\nL 747 559\nL 747 553\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 776 557\nL 783 557\nL 783 559\nL 776 559\nL 776 557\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"452\" y=\"553\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"472\" y=\"549\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4.9</text><text x=\"509\" y=\"546\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"526\" y=\"525\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23.2</text><text x=\"555\" y=\"522\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25.6</text><text x=\"583\" y=\"454\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">76.7</text><text x=\"606\" y=\"376\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">135.6</text><text x=\"635\" y=\"341\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">162.2</text><text x=\"669\" y=\"512\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32.6</text><text x=\"703\" y=\"529\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"729\" y=\"547\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6.4</text><text x=\"758\" y=\"551\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3.3</text><text x=\"454\" y=\"552\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.6</text><text x=\"482\" y=\"548\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">5.9</text><text x=\"519\" y=\"544\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">9</text><text x=\"536\" y=\"521\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26.4</text><text x=\"565\" y=\"518\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28.7</text><text x=\"593\" y=\"462\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">70.7</text><text x=\"616\" y=\"323\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">175.6</text><text x=\"645\" y=\"314\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">182.2</text><text x=\"679\" y=\"491\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.7</text><text x=\"708\" y=\"531\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18.8</text><text x=\"747\" y=\"548\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"768\" y=\"552\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><path  d=\"M 0 0\nL 800 0\nL 800 300\nL 0 300\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 350 19\nL 380 19\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"365\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><text x=\"382\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 411 19\nL 441 19\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"426\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><text x=\"443\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"10\" y=\"52\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"75\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"99\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"22\" y=\"123\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"171\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"195\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"219\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"243\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"267\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 790 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 68\nL 790 68\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 92\nL 790 92\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 116\nL 790 116\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 140\nL 790 140\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 164\nL 790 164\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 188\nL 790 188\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 212\nL 790 212\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 236\nL 790 236\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 265\nL 59 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 163 265\nL 163 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 267 265\nL 267 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 372 265\nL 372 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 476 265\nL 476 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 581 265\nL 581 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 685 265\nL 685 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 790 265\nL 790 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 59 260\nL 790 260\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"106\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"210\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"314\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"419\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"524\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"629\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"732\" y=\"285\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text><path  d=\"M 111 243\nL 215 241\nL 319 245\nL 424 240\nL 528 247\nL 633 226\nL 737 229\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:none\"/><circle cx=\"111\" cy=\"243\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"215\" cy=\"241\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"245\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"424\" cy=\"240\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"528\" cy=\"247\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"633\" cy=\"226\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"737\" cy=\"229\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 111 138\nL 215 121\nL 319 126\nL 424 121\nL 528 68\nL 633 62\nL 737 63\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:none\"/><circle cx=\"111\" cy=\"138\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"215\" cy=\"121\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"126\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"424\" cy=\"121\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"528\" cy=\"68\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"633\" cy=\"62\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"737\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,255,255,1.0)\"/></svg>", buf)
}

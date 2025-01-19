package charts

import (
	"fmt"
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

	t.Run("box", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        800,
			Height:       600,
		},
			PainterBoxOption(Box{Right: 400, Bottom: 300}),
			PainterPaddingOption(Box{Left: 1, Top: 2, Right: 3, Bottom: 4}),
		)

		assert.Equal(t, Box{
			Left:   1,
			Top:    2,
			Right:  397,
			Bottom: 296,
		}, p.box)
	})
	t.Run("theme", func(t *testing.T) {
		theme := GetTheme(ThemeAnt)
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        800,
			Height:       600,
		}, PainterThemeOption(theme))

		assert.Equal(t, theme, p.theme)
	})
	t.Run("font", func(t *testing.T) {
		font := GetDefaultFont()
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        800,
			Height:       600,
		})
		require.Nil(t, p.font)

		p = p.Child(PainterFontOption(font))

		assert.Equal(t, font, p.font)
	})
}

func TestPainterInternal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		result string
	}{
		{
			name: "circle",
			fn: func(p *Painter) {
				p.Circle(5, 2, 3, drawing.ColorTransparent, drawing.ColorTransparent, 1.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"7\" cy=\"13\" r=\"5\" style=\"stroke:none;fill:none\"/></svg>",
		},
		{
			name: "moveTo_lineTo",
			fn: func(p *Painter) {
				p.moveTo(1, 1)
				p.lineTo(2, 2)
				p.stroke(drawing.ColorTransparent, 1.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 6 11\nL 7 12\" style=\"stroke:none;fill:none\"/></svg>",
		},
		{
			name: "arc",
			fn: func(p *Painter) {
				p.arcTo(100, 100, 100, 100, 0, math.Pi/2)
				p.close()
				p.fillStroke(drawing.ColorBlue, drawing.ColorBlack, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 205 110\nA 100 100 90.00 0 1 105 210\nZ\" style=\"stroke-width:1;stroke:black;fill:blue\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(chartdraw.Box{Left: 5, Top: 10}))
			tt.fn(p)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

func TestPainterExternal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		result string
	}{
		{
			name: "text",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6, 0, FontStyle{})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"8\" y=\"16\" style=\"stroke:none;fill:none;font-family:'Roboto Medium',sans-serif\">hello world!</text></svg>",
		},
		{
			name: "text_rotated",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6, DegreesToRadians(90), FontStyle{})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"8\" y=\"16\" style=\"stroke:none;fill:none;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,8,16)\">hello world!</text></svg>",
		},
		{
			name: "line_stroke",
			fn: func(p *Painter) {
				p.LineStroke([]Point{
					{X: 10, Y: 20},
					{X: 30, Y: 40},
					{X: 50, Y: 20},
				}, drawing.ColorBlack, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 30\nL 35 50\nL 55 30\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "smooth_line_stroke",
			fn: func(p *Painter) {
				p.SmoothLineStroke([]Point{
					{X: 10, Y: 20},
					{X: 20, Y: 40},
					{X: 30, Y: 60},
					{X: 40, Y: 50},
					{X: 50, Y: 40},
					{X: 60, Y: 80},
				}, 0.5, drawing.ColorBlack, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 30\nQ25,50 27,55\nQ35,70 37,67\nQ45,60 47,57\nQ55,50 57,60\nQ55,50 65,90\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "background",
			fn: func(p *Painter) {
				p.SetBackground(400, 300, chartdraw.ColorWhite)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 405 10\nL 405 310\nL 5 310\nL 5 10\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "pin",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.Pin(30, 30, 30, c, c, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 32 47\nA 15 15 330.00 1 1 38 47\nL 35 33\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path  d=\"M 20 33\nQ35,70 50,33\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "arrow_left",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowLeft(30, 30, 16, 10, c, c, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 51 35\nL 35 40\nL 51 45\nL 46 40\nL 51 35\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "arrow_right",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowRight(30, 30, 16, 10, c, c, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 19 35\nL 35 40\nL 19 45\nL 24 40\nL 19 35\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "arrow_up",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowUp(30, 30, 10, 16, c, c, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 30 40\nL 35 24\nL 40 40\nL 35 35\nL 30 40\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "arrow_down",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowDown(30, 30, 10, 16, c, c, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 30 24\nL 35 40\nL 40 24\nL 35 30\nL 30 24\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.MarkLine(0, 20, 300, c, c, 1, []float64{4, 2})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"8\" cy=\"30\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 14 30\nL 289 30\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 289 25\nL 305 30\nL 289 35\nL 294 30\nL 289 25\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "polygon",
			fn: func(p *Painter) {
				p.Polygon(Point{X: 100, Y: 100}, 50, 6, Color{R: 84, G: 112, B: 198, A: 255}, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 105 60\nL 148 85\nL 148 134\nL 105 160\nL 62 135\nL 62 86\nL 105 60\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/></svg>",
		},
		{
			name: "fill_area",
			fn: func(p *Painter) {
				p.FillArea([]Point{
					{X: 0, Y: 0},
					{X: 0, Y: 100},
					{X: 100, Y: 100},
					{X: 0, Y: 0},
				}, Color{R: 84, G: 112, B: 198, A: 255})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 5 110\nL 105 110\nL 5 10\" style=\"stroke:none;fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "child_chart",
			fn: func(p *Painter) {
				opt := makeMinimalLineChartOption()
				opt.ValueFormatter = func(f float64) string {
					return fmt.Sprintf("%.0f", f)
				}
				_ = p.LineChart(opt)
				p = p.Child(PainterBoxOption(chartdraw.NewBox(0, 200, 400, 200)))
				opt = makeMinimalLineChartOption()
				opt.Theme = GetDefaultTheme().WithBackgroundColor(drawing.ColorFromAlphaMixedRGBA(0, 0, 0, 0))
				_ = p.LineChart(opt)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 400 10\nL 400 300\nL 5 300\nL 5 10\" style=\"stroke:none;fill:white\"/><text x=\"15\" y=\"27\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1440</text><text x=\"15\" y=\"53\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1280</text><text x=\"15\" y=\"80\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1120</text><text x=\"23\" y=\"107\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"23\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"23\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"23\" y=\"187\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"23\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"23\" y=\"240\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"41\" y=\"267\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 60 20\nL 390 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 46\nL 390 46\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 73\nL 390 73\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 100\nL 390 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 126\nL 390 126\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 153\nL 390 153\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 180\nL 390 180\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 206\nL 390 206\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 60 233\nL 390 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 83 240\nL 130 238\nL 177 244\nL 224 238\nL 271 245\nL 318 222\nL 366 225\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path  d=\"M 83 124\nL 130 105\nL 177 110\nL 224 105\nL 271 45\nL 318 39\nL 366 40\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><path  d=\"M 200 0\nL 400 0\nL 400 200\nL 200 200\nL 200 0\" style=\"stroke:none;fill:none\"/><text x=\"210\" y=\"17\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"210\" y=\"33\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"210\" y=\"50\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"222\" y=\"67\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"222\" y=\"83\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"222\" y=\"100\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"222\" y=\"117\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"222\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"222\" y=\"150\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"240\" y=\"167\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 259 10\nL 390 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 26\nL 390 26\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 43\nL 390 43\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 60\nL 390 60\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 76\nL 390 76\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 93\nL 390 93\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 110\nL 390 110\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 126\nL 390 126\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 143\nL 390 143\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 259 148\nL 280 147\nL 302 150\nL 324 147\nL 346 151\nL 368 137\nL 390 139\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path  d=\"M 259 75\nL 280 63\nL 302 67\nL 324 63\nL 346 26\nL 368 22\nL 390 23\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(chartdraw.Box{Left: 5, Top: 10}))
			tt.fn(p)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

func TestTextRotationHeightAdjustment(t *testing.T) {
	t.Parallel()

	text := "hello world "
	expectedTemplate := "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"200\" y=\"%d\" style=\"stroke:none;fill:black;font-size:40.9px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(%d.00,200,%d)\">hello world %s</text></svg>"
	fontStyle := FontStyle{
		Font:      GetDefaultFont(),
		FontSize:  32,
		FontColor: drawing.ColorBlack,
	}
	drawDebugBox := false

	tests := []struct {
		degrees   int
		expectedY int
	}{
		{
			degrees:   15,
			expectedY: 127,
		},
		{
			degrees:   30,
			expectedY: 60,
		},
		{
			degrees:   45,
			expectedY: 1,
		},
		{
			degrees:   60,
			expectedY: -43,
		},
		{
			degrees:   75,
			expectedY: -71,
		},
		{
			degrees:   90,
			expectedY: -81,
		},
		{
			degrees:   105,
			expectedY: -71,
		},
		{
			degrees:   120,
			expectedY: -43,
		},
		{
			degrees:   135,
			expectedY: 1,
		},
		{
			degrees:   150,
			expectedY: 60,
		},
		{
			degrees:   165,
			expectedY: 127,
		},
		{
			degrees:   180,
			expectedY: 160,
		},
		{
			degrees:   195,
			expectedY: 170,
		},
		{
			degrees:   210,
			expectedY: 180,
		},
		{
			degrees:   225,
			expectedY: 188,
		},
		{
			degrees:   240,
			expectedY: 195,
		},
		{
			degrees:   255,
			expectedY: 199,
		},
		{
			degrees:   270,
			expectedY: 200,
		},
		{
			degrees:   285,
			expectedY: 200,
		},
		{
			degrees:   300,
			expectedY: 200,
		},
		{
			degrees:   315,
			expectedY: 200,
		},
		{
			degrees:   330,
			expectedY: 200,
		},
		{
			degrees:   345,
			expectedY: 200,
		},
		{
			degrees:   360,
			expectedY: 200,
		},
	}

	for _, tt := range tests {
		name := strconv.Itoa(tt.degrees)
		for len(name) < 3 {
			name = "0" + name
		}
		t.Run(name, func(t *testing.T) {
			padding := 200
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterPaddingOption(chartdraw.Box{Left: padding, Top: padding}))

			radians := DegreesToRadians(float64(tt.degrees))
			testText := text + name
			textBox := p.MeasureText(testText, 0, fontStyle)

			if drawDebugBox {
				debugBox := []Point{
					{X: 0, Y: 0},
					{X: 0, Y: -textBox.Height()},
					{X: textBox.Width(), Y: -textBox.Height()},
					{X: textBox.Width(), Y: 0},
					{X: 0, Y: 0},
				}
				p.LineStroke(debugBox, drawing.ColorBlue, 1)
			}

			assert.Equal(t, tt.expectedY, padding-textRotationHeightAdjustment(textBox.Width(), textBox.Height(), radians))

			p.Text(testText, 0, -textRotationHeightAdjustment(textBox.Width(), textBox.Height(), radians), radians, fontStyle)

			data, err := p.Bytes()
			require.NoError(t, err)

			if drawDebugBox {
				assertEqualSVG(t, "", data)
			} else {
				expectedResult := fmt.Sprintf(expectedTemplate, tt.expectedY, tt.degrees%360, tt.expectedY, name)
				assertEqualSVG(t, expectedResult, data)
			}
		})
	}
}

func TestPainterRoundedRect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		result string
	}{
		{
			name: "round_fully",
			fn: func(p *Painter) {
				p.roundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, true, true, drawing.ColorBlue, drawing.ColorBlue, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 10\nL 25 10\nL 25 10\nA 5 5 90.00 0 1 30 15\nL 30 145\nL 30 145\nA 5 5 90.00 0 1 25 150\nL 15 150\nL 15 150\nA 5 5 90.00 0 1 10 145\nL 10 15\nL 10 15\nA 5 5 90.00 0 1 15 10\nZ\" style=\"stroke-width:1;stroke:blue;fill:blue\"/></svg>",
		},
		{
			name: "square_top",
			fn: func(p *Painter) {
				p.roundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, false, true, drawing.ColorBlue, drawing.ColorBlue, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 10 10\nL 30 10\nL 30 145\nL 30 145\nA 5 5 90.00 0 1 25 150\nL 15 150\nL 15 150\nA 5 5 90.00 0 1 10 145\nL 10 10\nZ\" style=\"stroke-width:1;stroke:blue;fill:blue\"/></svg>",
		},
		{
			name: "square_bottom",
			fn: func(p *Painter) {
				p.roundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, true, false, drawing.ColorBlue, drawing.ColorBlue, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 10\nL 25 10\nL 25 10\nA 5 5 90.00 0 1 30 15\nL 30 150\nL 10 150\nL 10 15\nL 10 15\nA 5 5 90.00 0 1 15 10\nZ\" style=\"stroke-width:1;stroke:blue;fill:blue\"/></svg>",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				Width:        400,
				Height:       300,
				OutputFormat: ChartOutputSVG,
			})
			tc.fn(p)
			buf, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.result, buf)
		})
	}
}

func TestPainterMeasureText(t *testing.T) {
	t.Parallel()

	svgP := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        400,
		Height:       300,
	})
	pngP := NewPainter(PainterOptions{
		OutputFormat: ChartOutputPNG,
		Width:        400,
		Height:       300,
	})
	style := FontStyle{
		FontSize:  12,
		FontColor: chartdraw.ColorBlack,
		Font:      GetDefaultFont(),
	}

	assert.Equal(t, chartdraw.Box{Right: 84, Bottom: 15, IsSet: true},
		svgP.MeasureText("Hello World!", 0, style))
	assert.Equal(t, chartdraw.Box{Right: 99, Bottom: 14, IsSet: true},
		pngP.MeasureText("Hello World!", 0, style))
}

func TestPainterTextFit(t *testing.T) {
	t.Parallel()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        400,
		Height:       300,
	})
	fontStyle := FontStyle{
		FontSize:  12,
		FontColor: chartdraw.ColorBlack,
		Font:      GetDefaultFont(),
	}

	box := p.TextFit("Hello World!", 0, 20, 80, fontStyle)
	assert.Equal(t, chartdraw.Box{Right: 45, Bottom: 35, IsSet: true}, box)

	box = p.TextFit("Hello World!", 0, 100, 200, fontStyle)
	assert.Equal(t, chartdraw.Box{Right: 84, Bottom: 15, IsSet: true}, box)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"0\" y=\"20\" style=\"stroke:none;fill:rgb(51,51,51);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello</text><text x=\"0\" y=\"40\" style=\"stroke:none;fill:rgb(51,51,51);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World!</text><text x=\"0\" y=\"100\" style=\"stroke:none;fill:rgb(51,51,51);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello World!</text></svg>", buf)
}

func TestMultipleChartsOnPainter(t *testing.T) {
	t.Parallel()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        800,
		Height:       600,
	})
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
	err := bottomLeftPainter.PieChart(pieOpt)
	require.NoError(t, err)
	err = bottomRightPainter.BarChart(makeBasicBarChartOption())
	require.NoError(t, err)
	err = topCenterPainter.LineChart(makeBasicLineChartOption())
	require.NoError(t, err)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 800 600\"><path  d=\"M 0 0\nL 800 0\nL 800 600\nL 0 600\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 0 300\nL 400 300\nL 400 600\nL 0 600\nL 0 300\" style=\"stroke:none;fill:white\"/><text x=\"122\" y=\"335\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"166\" y=\"350\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path  d=\"M 200 457\nL 200 383\nA 74 74 119.89 0 1 264 493\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:rgb(91,143,249)\"/><path  d=\"M 264 420\nL 277 413\nM 277 413\nL 292 413\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:none\"/><text x=\"295\" y=\"418\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Search Engine: 33.3%</text><path  d=\"M 200 457\nL 264 493\nA 74 74 84.08 0 1 170 524\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:rgb(90,216,166)\"/><path  d=\"M 222 527\nL 227 541\nM 227 541\nL 242 541\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:none\"/><text x=\"245\" y=\"546\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Direct: 23.35%</text><path  d=\"M 200 457\nL 170 524\nA 74 74 66.35 0 1 127 457\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:rgb(93,112,146)\"/><path  d=\"M 138 497\nL 126 505\nM 126 505\nL 111 505\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:none\"/><text x=\"27\" y=\"510\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Email: 18.43%</text><path  d=\"M 200 457\nL 127 457\nA 74 74 55.37 0 1 159 396\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:rgb(246,189,22)\"/><path  d=\"M 135 423\nL 122 416\nM 122 416\nL 107 416\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:none\"/><text x=\"-4\" y=\"421\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Union Ads: 15.37%</text><path  d=\"M 200 457\nL 159 396\nA 74 74 34.32 0 1 200 383\nL 200 457\nZ\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:rgb(111,94,249)\"/><path  d=\"M 179 387\nL 174 372\nM 174 372\nL 159 372\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:none\"/><text x=\"56\" y=\"377\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Video Ads: 9.53%</text><path  d=\"M 400 300\nL 800 300\nL 800 600\nL 400 600\nL 400 300\" style=\"stroke:none;fill:white\"/><text x=\"410\" y=\"317\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">189</text><text x=\"410\" y=\"344\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">168</text><text x=\"410\" y=\"372\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">147</text><text x=\"410\" y=\"400\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">126</text><text x=\"410\" y=\"428\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">105</text><text x=\"419\" y=\"455\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84</text><text x=\"419\" y=\"483\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">63</text><text x=\"419\" y=\"511\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"419\" y=\"539\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">21</text><text x=\"428\" y=\"567\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 447 310\nL 790 310\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 337\nL 790 337\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 365\nL 790 365\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 393\nL 790 393\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 421\nL 790 421\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 448\nL 790 448\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 476\nL 790 476\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 504\nL 790 504\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 532\nL 790 532\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 447 565\nL 447 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 509 565\nL 509 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 571 565\nL 571 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 634 565\nL 634 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 665 565\nL 665 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 727 565\nL 727 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 790 565\nL 790 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 447 560\nL 790 560\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"446\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"508\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"570\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"633\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"664\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"726\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"763\" y=\"585\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 452 558\nL 459 558\nL 459 559\nL 452 559\nL 452 558\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 480 554\nL 487 554\nL 487 559\nL 480 559\nL 480 554\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 509 551\nL 516 551\nL 516 559\nL 509 559\nL 509 551\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 537 530\nL 544 530\nL 544 559\nL 537 559\nL 537 530\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 566 527\nL 573 527\nL 573 559\nL 566 559\nL 566 527\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 594 459\nL 601 459\nL 601 559\nL 594 559\nL 594 459\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 623 381\nL 630 381\nL 630 559\nL 623 559\nL 623 381\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 652 346\nL 659 346\nL 659 559\nL 652 559\nL 652 346\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 680 517\nL 687 517\nL 687 559\nL 680 559\nL 680 517\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 709 534\nL 716 534\nL 716 559\nL 709 559\nL 709 534\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 737 552\nL 744 552\nL 744 559\nL 737 559\nL 737 552\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 766 556\nL 773 556\nL 773 559\nL 766 559\nL 766 556\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 462 557\nL 469 557\nL 469 559\nL 462 559\nL 462 557\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 490 553\nL 497 553\nL 497 559\nL 490 559\nL 490 553\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 519 549\nL 526 549\nL 526 559\nL 519 559\nL 519 549\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 547 526\nL 554 526\nL 554 559\nL 547 559\nL 547 526\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 576 523\nL 583 523\nL 583 559\nL 576 559\nL 576 523\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 604 467\nL 611 467\nL 611 559\nL 604 559\nL 604 467\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 633 328\nL 640 328\nL 640 559\nL 633 559\nL 633 328\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 662 319\nL 669 319\nL 669 559\nL 662 559\nL 662 319\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 690 496\nL 697 496\nL 697 559\nL 690 559\nL 690 496\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 719 536\nL 726 536\nL 726 559\nL 719 559\nL 719 536\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 747 553\nL 754 553\nL 754 559\nL 747 559\nL 747 553\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 776 557\nL 783 557\nL 783 559\nL 776 559\nL 776 557\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"451\" y=\"553\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"474\" y=\"549\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4.9</text><text x=\"508\" y=\"546\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"527\" y=\"525\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23.2</text><text x=\"556\" y=\"522\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25.6</text><text x=\"584\" y=\"454\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">76.7</text><text x=\"610\" y=\"376\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">135.6</text><text x=\"639\" y=\"341\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">162.2</text><text x=\"670\" y=\"512\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32.6</text><text x=\"705\" y=\"529\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"731\" y=\"547\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6.4</text><text x=\"760\" y=\"551\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3.3</text><text x=\"456\" y=\"552\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.6</text><text x=\"484\" y=\"548\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">5.9</text><text x=\"518\" y=\"544\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">9</text><text x=\"537\" y=\"521\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26.4</text><text x=\"566\" y=\"518\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28.7</text><text x=\"594\" y=\"462\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">70.7</text><text x=\"620\" y=\"323\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">175.6</text><text x=\"649\" y=\"314\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">182.2</text><text x=\"680\" y=\"491\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.7</text><text x=\"709\" y=\"531\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18.8</text><text x=\"746\" y=\"548\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"770\" y=\"552\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><path  d=\"M 0 0\nL 800 0\nL 800 300\nL 0 300\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 350 19\nL 380 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"365\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"382\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 411 19\nL 441 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"426\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"443\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"10\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"99\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"22\" y=\"123\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"147\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"195\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"219\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"243\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"267\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 790 45\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 68\nL 790 68\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 92\nL 790 92\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 116\nL 790 116\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 140\nL 790 140\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 164\nL 790 164\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 188\nL 790 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 212\nL 790 212\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 236\nL 790 236\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 265\nL 59 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 163 265\nL 163 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 267 265\nL 267 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 372 265\nL 372 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 476 265\nL 476 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 581 265\nL 581 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 685 265\nL 685 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 790 265\nL 790 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 59 260\nL 790 260\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"106\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"210\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"314\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"419\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"524\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"629\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"732\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text><path  d=\"M 111 243\nL 215 241\nL 319 245\nL 424 240\nL 528 247\nL 633 226\nL 737 229\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"111\" cy=\"243\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"215\" cy=\"241\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"319\" cy=\"245\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"424\" cy=\"240\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"528\" cy=\"247\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"633\" cy=\"226\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"737\" cy=\"229\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><path  d=\"M 111 138\nL 215 121\nL 319 126\nL 424 121\nL 528 68\nL 633 62\nL 737 63\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"111\" cy=\"138\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"215\" cy=\"121\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"319\" cy=\"126\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"424\" cy=\"121\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"528\" cy=\"68\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"633\" cy=\"62\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"737\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/></svg>", buf)
}

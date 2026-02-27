package charts

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestBytesFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		outputFormat string
		magicStart   []byte
	}{
		{
			outputFormat: ChartOutputPNG,
			magicStart:   []byte{0x89, 0x50, 0x4E, 0x47},
		},
		{
			outputFormat: ChartOutputJPG,
			magicStart:   []byte{0xFF, 0xD8},
		},
		{
			outputFormat: ChartOutputSVG,
			magicStart:   []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 "),
		},
	}
	for _, tc := range tests {
		t.Run(tc.outputFormat, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: tc.outputFormat,
				Width:        200,
				Height:       100,
			})
			err := p.LineChart(makeFullLineChartStackedOption())
			require.NoError(t, err)
			content, err := p.Bytes()
			require.NoError(t, err)

			assert.True(t, bytes.HasPrefix(content, tc.magicStart))
		})
	}
}

func TestBytesCompareRenderedOutputs(t *testing.T) {
	t.Parallel()

	pngP := NewPainter(PainterOptions{
		OutputFormat: ChartOutputPNG,
		Width:        200,
		Height:       100,
	})
	jpgP := NewPainter(PainterOptions{
		OutputFormat: ChartOutputJPG,
		Width:        200,
		Height:       100,
	})
	err := pngP.LineChart(makeFullLineChartStackedOption())
	require.NoError(t, err)
	err = jpgP.LineChart(makeFullLineChartStackedOption())
	require.NoError(t, err)

	pngData, err := pngP.Bytes()
	require.NoError(t, err)
	jpgData, err := jpgP.Bytes()
	require.NoError(t, err)

	pngImg, err := decodeImage(pngData)
	require.NoError(t, err)
	jpgImg, err := decodeImage(jpgData)
	require.NoError(t, err)

	assert.Equal(t, color.RGBAModel, pngImg.ColorModel())
	assert.Equal(t, color.YCbCrModel, jpgImg.ColorModel())

	assert.Equal(t, pngImg.Bounds().Size().X, jpgImg.Bounds().Size().X)
	assert.Equal(t, pngImg.Bounds().Size().Y, jpgImg.Bounds().Size().Y)
}

func decodeImage(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	return img, err
}

func TestPainterInternal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		svg    string
		pngCRC uint32
	}{
		{
			name: "circle",
			fn: func(p *Painter) {
				p.Circle(5, 2, 3, ColorTransparent, ColorTransparent, 1.0)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"7\" cy=\"13\" r=\"5\" style=\"stroke:none;fill:none\"/></svg>",
			pngCRC: 0x2083f7bd,
		},
		{
			name: "moveTo_lineTo",
			fn: func(p *Painter) {
				p.moveTo(1, 1)
				p.lineTo(2, 2)
				p.stroke(ColorTransparent, 1.0)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 6 11\nL 7 12\" style=\"stroke:none;fill:none\"/></svg>",
			pngCRC: 0x2083f7bd,
		},
		{
			name: "arc",
			fn: func(p *Painter) {
				p.arcTo(100, 100, 100, 100, 0, math.Pi/2)
				p.close()
				p.fillStroke(ColorBlue, ColorBlack, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 205 110\nA 100 100 90.00 0 1 105 210\nZ\" style=\"stroke-width:1;stroke:black;fill:blue\"/></svg>",
			pngCRC: 0x5445e3e7,
		},
		{
			name: "draw_background",
			fn: func(p *Painter) {
				p.drawBackground(ColorWhite)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 5 10\nL 400 10\nL 400 300\nL 5 300\nL 5 10\" style=\"stroke:none;fill:white\"/></svg>",
			pngCRC: 0x60f3dd98,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			svgP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(Box{Left: 5, Top: 10}))
			tt.fn(svgP)
			data, err := svgP.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.svg, data)

			pngP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(Box{Left: 5, Top: 10}))
			tt.fn(pngP)
			data, err = pngP.Bytes()
			require.NoError(t, err)
			assertEqualPNGCRC(t, tt.pngCRC, data)
		})
	}
}

func TestPainterExternal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		svg    string
		pngCRC uint32
	}{
		{
			name: "text",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6, 0, FontStyle{})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"8\" y=\"16\" style=\"stroke:none;fill:none;font-family:'Roboto Medium',sans-serif\">hello world!</text></svg>",
			pngCRC: 0x2083f7bd,
		},
		{
			name: "text_rotated",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6, DegreesToRadians(90), FontStyle{})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"8\" y=\"16\" style=\"stroke:none;fill:none;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,8,16)\">hello world!</text></svg>",
			pngCRC: 0x2083f7bd,
		},
		{
			name: "line_stroke",
			fn: func(p *Painter) {
				p.LineStroke([]Point{
					{X: 10, Y: 20},
					{X: 30, Y: 40},
					{X: 50, Y: 20},
				}, ColorBlack, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 15 30\nL 35 50\nL 55 30\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
			pngCRC: 0x8cfe7b4b,
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
				}, 0.5, ColorBlack, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 15 30\nQ25,50 27,55\nQ35,70 37,67\nQ45,60 47,57\nQ55,50 57,60\nQ55,50 65,90\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
			pngCRC: 0x3056a98,
		},
		{
			name: "filled_rect",
			fn: func(p *Painter) {
				p.FilledRect(0, 0, 400, 300, ColorWhite, ColorWhite, 0.0)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 5 10\nL 405 10\nL 405 310\nL 5 310\nL 5 10\" style=\"stroke:none;fill:white\"/></svg>",
			pngCRC: 0x60f3dd98,
		},
		{
			name: "filled_rect_center",
			fn: func(p *Painter) {
				p.FilledRect(100, 100, 200, 150, ColorWhite, ColorWhite, 0.0)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 105 110\nL 205 110\nL 205 160\nL 105 160\nL 105 110\" style=\"stroke:none;fill:white\"/></svg>",
			pngCRC: 0x540b2357,
		},
		{
			name: "filled_rect_center_border",
			fn: func(p *Painter) {
				p.FilledRect(100, 100, 200, 150, ColorWhite, ColorBlue, 1.0)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 105 110\nL 205 110\nL 205 160\nL 105 160\nL 105 110\" style=\"stroke-width:1;stroke:blue;fill:white\"/></svg>",
			pngCRC: 0xd6c2a417,
		},
		{
			name: "pin",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.Pin(30, 30, 30, c, c, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 31 47\nA 15 15 330.00 1 1 39 47\nL 35 33\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path d=\"M 20 33\nQ35,70 50,33\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0x981d8eb5,
		},
		{
			name: "arrow_left",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowLeft(30, 30, 16, 10, c, c, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 51 35\nL 35 40\nL 51 45\nL 46 40\nL 51 35\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0x3415dab,
		},
		{
			name: "arrow_right",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowRight(30, 30, 16, 10, c, c, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 19 35\nL 35 40\nL 19 45\nL 24 40\nL 19 35\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0x142dfb03,
		},
		{
			name: "arrow_up",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowUp(30, 30, 10, 16, c, c, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 30 40\nL 35 24\nL 40 40\nL 35 35\nL 30 40\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0xe17c9204,
		},
		{
			name: "arrow_down",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowDown(30, 30, 10, 16, c, c, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 30 24\nL 35 40\nL 40 24\nL 35 30\nL 30 24\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0xd56c309d,
		},
		{
			name: "horizontal_mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.HorizontalMarkLine(0, 20, 300, c, c, 1, []float64{4, 2})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"8\" cy=\"30\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 14 30\nL 289 30\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 289 25\nL 305 30\nL 289 35\nL 294 30\nL 289 25\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0xa4ca1cb8,
		},
		{
			name: "vertical_mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.VerticalMarkLine(200, 100, 100, c, c, 1, []float64{4, 2})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"205\" cy=\"207\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 205 110\nL 205 210\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 200 126\nL 205 110\nL 210 126\nL 205 121\nL 200 126\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0x681c0b4e,
		},
		{
			name: "polygon",
			fn: func(p *Painter) {
				p.Polygon(Point{X: 100, Y: 100}, 50, 6, Color{R: 84, G: 112, B: 198, A: 255}, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 105 60\nL 148 85\nL 148 134\nL 105 160\nL 62 135\nL 62 86\nL 105 60\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/></svg>",
			pngCRC: 0xdb739c98,
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
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 5 10\nL 5 110\nL 105 110\nL 5 10\" style=\"stroke:none;fill:rgb(84,112,198)\"/></svg>",
			pngCRC: 0xf2b066ae,
		},
		{
			name: "child_chart",
			fn: func(p *Painter) {
				opt := makeMinimalLineChartOption()
				opt.ValueFormatter = func(f float64) string {
					return fmt.Sprintf("%.0f", f)
				}
				_ = p.LineChart(opt)
				p = p.Child(PainterBoxOption(NewBox(200, 0, 400, 200)))
				opt = makeMinimalLineChartOption()
				opt.Theme = GetDefaultTheme().WithBackgroundColor(ColorTransparent)
				_ = p.LineChart(opt)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 5 10\nL 400 10\nL 400 300\nL 5 300\nL 5 10\" style=\"stroke:none;fill:white\"/><text x=\"14\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1440</text><text x=\"14\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1280</text><text x=\"14\" y=\"85\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1120</text><text x=\"22\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"145\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"174\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"204\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"234\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"294\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 55 20\nL 390 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 50\nL 390 50\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 80\nL 390 80\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 110\nL 390 110\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 140\nL 390 140\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 170\nL 390 170\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 200\nL 390 200\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 230\nL 390 230\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 55 260\nL 390 260\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 82 268\nL 129 266\nL 176 272\nL 224 265\nL 271 274\nL 318 247\nL 366 251\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path d=\"M 82 137\nL 129 116\nL 176 122\nL 224 115\nL 271 49\nL 318 41\nL 366 43\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 200 0\nL 400 0\nL 400 200\nL 200 200\nL 200 0\" style=\"stroke:none;fill:none\"/><text x=\"209\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"209\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"209\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"221\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"221\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"221\" y=\"114\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"221\" y=\"134\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"221\" y=\"154\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"221\" y=\"174\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"239\" y=\"194\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 254 10\nL 390 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 30\nL 390 30\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 50\nL 390 50\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 70\nL 390 70\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 90\nL 390 90\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 110\nL 390 110\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 130\nL 390 130\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 150\nL 390 150\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 254 170\nL 390 170\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 258 175\nL 280 174\nL 302 178\nL 324 174\nL 346 179\nL 368 162\nL 390 164\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path d=\"M 258 88\nL 280 74\nL 302 78\nL 324 74\nL 346 29\nL 368 24\nL 390 25\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/></svg>",
			pngCRC: 0xc3dc7021,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			svgP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(Box{Left: 5, Top: 10}))
			tt.fn(svgP)
			data, err := svgP.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.svg, data)

			pngP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(Box{Left: 5, Top: 10}))
			tt.fn(pngP)
			data, err = pngP.Bytes()
			require.NoError(t, err)
			assertEqualPNGCRC(t, tt.pngCRC, data)
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
		FontColor: ColorBlack,
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
			expectedY: 159,
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
			}, PainterPaddingOption(Box{Left: padding, Top: padding}))

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
				p.LineStroke(debugBox, ColorBlue, 1)
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
		svg    string
		pngCRC uint32
	}{
		{
			name: "round_fully",
			fn: func(p *Painter) {
				p.roundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, true, true, ColorBlue, ColorBlue, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 15 10\nL 25 10\nL 25 10\nA 5 5 90.00 0 1 30 15\nL 30 145\nL 30 145\nA 5 5 90.00 0 1 25 150\nL 15 150\nL 15 150\nA 5 5 90.00 0 1 10 145\nL 10 15\nL 10 15\nA 5 5 90.00 0 1 15 10\nZ\" style=\"stroke-width:1;stroke:blue;fill:blue\"/></svg>",
			pngCRC: 0x63ab7f9f,
		},
		{
			name: "square_top",
			fn: func(p *Painter) {
				p.roundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, false, true, ColorBlue, ColorBlue, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 10 10\nL 30 10\nL 30 145\nL 30 145\nA 5 5 90.00 0 1 25 150\nL 15 150\nL 15 150\nA 5 5 90.00 0 1 10 145\nL 10 10\nZ\" style=\"stroke-width:1;stroke:blue;fill:blue\"/></svg>",
			pngCRC: 0x47fb7794,
		},
		{
			name: "square_bottom",
			fn: func(p *Painter) {
				p.roundedRect(Box{
					Left:   10,
					Right:  30,
					Bottom: 150,
					Top:    10,
				}, 5, true, false, ColorBlue, ColorBlue, 1)
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path d=\"M 15 10\nL 25 10\nL 25 10\nA 5 5 90.00 0 1 30 15\nL 30 150\nL 10 150\nL 10 15\nL 10 15\nA 5 5 90.00 0 1 15 10\nZ\" style=\"stroke-width:1;stroke:blue;fill:blue\"/></svg>",
			pngCRC: 0xe8bb388c,
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			svgP := NewPainter(PainterOptions{
				Width:        400,
				Height:       300,
				OutputFormat: ChartOutputSVG,
			})
			tc.fn(svgP)
			buf, err := svgP.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.svg, buf)

			pngP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        400,
				Height:       300,
			})
			tc.fn(pngP)
			data, err := pngP.Bytes()
			require.NoError(t, err)
			assertEqualPNGCRC(t, tc.pngCRC, data)
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
	jpgP := NewPainter(PainterOptions{
		OutputFormat: ChartOutputJPG,
		Width:        400,
		Height:       300,
	})
	style := FontStyle{
		FontSize:  12,
		FontColor: ColorBlack,
		Font:      GetDefaultFont(),
	}
	styleLargeNoto := FontStyle{
		FontSize:  28,
		FontColor: ColorBlue,
		Font:      GetFont(FontFamilyNotoSans),
	}
	styleLargeRoboto := FontStyle{
		FontSize:  28,
		FontColor: ColorBlue,
		Font:      GetFont(FontFamilyRoboto),
	}

	t.Run("basic", func(t *testing.T) {
		assert.Equal(t, Box{Right: 84, Bottom: 16, IsSet: true},
			svgP.MeasureText("Hello World!", 0, style))
		assert.Equal(t, Box{Right: 99, Bottom: 14, IsSet: true},
			pngP.MeasureText("Hello World!", 0, style))
		assert.Equal(t, Box{Right: 99, Bottom: 14, IsSet: true},
			jpgP.MeasureText("Hello World!", 0, style))
	})

	t.Run("basic_large", func(t *testing.T) {
		assert.Equal(t, Box{Right: 200, Bottom: 36, IsSet: true},
			svgP.MeasureText("Hello World!", 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 235, Bottom: 33, IsSet: true},
			pngP.MeasureText("Hello World!", 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 235, Bottom: 33, IsSet: true},
			jpgP.MeasureText("Hello World!", 0, styleLargeNoto))
	})

	t.Run("rotated-90", func(t *testing.T) {
		radians := DegreesToRadians(90)

		box := svgP.MeasureText("Hello World!", radians, style)
		assert.Equal(t, 84, box.Height())
		assert.Equal(t, 16, box.Width())

		box = pngP.MeasureText("Hello World!", radians, style)
		assert.Equal(t, 99, box.Height())
		assert.Equal(t, 14, box.Width())

		box = jpgP.MeasureText("Hello World!", radians, style)
		assert.Equal(t, 99, box.Height())
		assert.Equal(t, 14, box.Width())
	})

	t.Run("rotated-270", func(t *testing.T) {
		radians := DegreesToRadians(270)

		box := svgP.MeasureText("Hello World!", radians, style)
		assert.Equal(t, 84, box.Height())
		assert.Equal(t, 14, box.Width())

		box = pngP.MeasureText("Hello World!", radians, style)
		assert.Equal(t, 99, box.Height())
		assert.Equal(t, 12, box.Width())

		box = jpgP.MeasureText("Hello World!", radians, style)
		assert.Equal(t, 99, box.Height())
		assert.Equal(t, 12, box.Width())
	})

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, Box{IsSet: true}, svgP.MeasureText("", 0, style))
		assert.Equal(t, Box{IsSet: true}, pngP.MeasureText("", 0, style))
		assert.Equal(t, Box{IsSet: true}, jpgP.MeasureText("", 0, style))
	})

	t.Run("green_circle_emoji", func(t *testing.T) {
		text := "🟢"

		assert.Equal(t, Box{Right: 12, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 28, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("star_emoji", func(t *testing.T) {
		text := "⭐"

		assert.Equal(t, Box{Right: 12, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 28, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("question_emoji", func(t *testing.T) {
		text := "❓"

		assert.Equal(t, Box{Right: 12, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 28, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("money_emoji", func(t *testing.T) {
		text := "💰"

		assert.Equal(t, Box{Right: 12, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 8, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 28, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 24, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("multiple_emojis", func(t *testing.T) {
		text := "🟢⭐❓💰"

		assert.Equal(t, Box{Right: 48, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 32, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 32, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 112, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 95, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 95, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("mixed_text_emoji", func(t *testing.T) {
		text := "Status: 🟢 OK"

		assert.Equal(t, Box{Right: 89, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 99, Bottom: 14, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 99, Bottom: 14, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 205, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 233, Bottom: 32, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 233, Bottom: 32, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("transport_symbols", func(t *testing.T) {
		text := "🚗🚕"

		assert.Equal(t, Box{Right: 24, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 56, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("geometric_shapes", func(t *testing.T) {
		text := "▪▫"

		assert.Equal(t, Box{Right: 24, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 56, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("wave_dash_and_part_alternation_mark", func(t *testing.T) {
		text := "〰〽"

		assert.Equal(t, Box{Right: 24, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 56, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("playing_cards", func(t *testing.T) {
		text := "🂡🂢"

		assert.Equal(t, Box{Right: 24, Bottom: 16, IsSet: true}, svgP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, pngP.MeasureText(text, 0, style))
		assert.Equal(t, Box{Right: 16, Bottom: 13, IsSet: true}, jpgP.MeasureText(text, 0, style))

		assert.Equal(t, Box{Right: 56, Bottom: 36, IsSet: true}, svgP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, pngP.MeasureText(text, 0, styleLargeNoto))
		assert.Equal(t, Box{Right: 48, Bottom: 30, IsSet: true}, jpgP.MeasureText(text, 0, styleLargeNoto))
	})

	t.Run("rendered", func(t *testing.T) {
		tests := []struct {
			name        string
			input       string
			font        FontStyle
			expectedSVG string
			expectedCRC uint32
		}{
			{
				name:        "basic",
				input:       "Hello World!",
				font:        styleLargeNoto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">Hell</text><path d=\"M 50 64\nL 114 64\nL 114 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">o Wo</text><path d=\"M 50 164\nL 133 164\nL 133 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">rld!</text><path d=\"M 50 264\nL 104 264\nL 104 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0x1bede00d,
			},
			{
				name:        "emojis",
				input:       "⭐❓💰🔥💯🎯🚀⚡🌟🎉🎊",
				font:        styleLargeNoto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">⭐❓💰</text><path d=\"M 50 64\nL 134 64\nL 134 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">🔥💯🎯</text><path d=\"M 50 164\nL 134 164\nL 134 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">🚀⚡🌟🎉🎊</text><path d=\"M 50 264\nL 190 264\nL 190 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0xaaf41b1b,
			},
			{
				name:        "shapes",
				input:       "▫●□▲▼◇★○△▪▴▾◆◯⬟⬠⬡⬢⬣⬤⬥",
				font:        styleLargeNoto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">▫●□▲▼◇★</text><path d=\"M 50 64\nL 246 64\nL 246 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">○△▪▴▾◆◯</text><path d=\"M 50 164\nL 246 164\nL 246 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">⬟⬠⬡⬢⬣⬤⬥</text><path d=\"M 50 264\nL 191 264\nL 191 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0x2a17094,
			},
			{
				name:        "playing_cards",
				input:       "🂡🂢🂫🃄🃍🃘🃞🃟",
				font:        styleLargeNoto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">🂡🂢</text><path d=\"M 50 64\nL 106 64\nL 106 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">🂫🃄</text><path d=\"M 50 164\nL 106 164\nL 106 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">🃍🃘🃞🃟</text><path d=\"M 50 264\nL 162 264\nL 162 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0x94611d5b,
			},
			{
				name:        "faces",
				input:       "😂😍🤣😊😭😘😎🤔😴😋😉😏😬😐😑😮😯",
				font:        styleLargeNoto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">😂😍🤣😊😭</text><path d=\"M 50 64\nL 190 64\nL 190 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">😘😎🤔😴😋</text><path d=\"M 50 164\nL 190 164\nL 190 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">😉😏😬😐😑😮😯</text><path d=\"M 50 264\nL 246 264\nL 246 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0x64607851,
			},
			{
				name:        "fallback_notosans_currency",
				input:       "₠₡₢₥₭₮₯₰₲₳₴₵₶₷₸₻₾₿",
				font:        styleLargeRoboto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">₠₡₢₥₭₮</text><path d=\"M 50 64\nL 181 64\nL 181 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">₯₰₲₳₴₵</text><path d=\"M 50 164\nL 189 164\nL 189 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">₶₷₸₻₾₿</text><path d=\"M 50 264\nL 183 264\nL 183 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0xe1e34ff1,
			},
			{
				name:        "fallback_notosans_letterlike",
				input:       "℀℁ℂ℃℄℆ℇ℈℉ℊℋℌℍℎℏℐℑℒ℔ℕ℗℘ℙℚℛℜℝ℞℟℣ℤ℥℧ℨ℩KÅℬℭℯ",
				font:        styleLargeRoboto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">℀℁ℂ℃℄℆ℇ℈℉ℊℋℌℍ</text><path d=\"M 50 64\nL 391 64\nL 391 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">ℎℏℐℑℒ℔ℕ℗℘ℙℚℛℜ</text><path d=\"M 50 164\nL 384 164\nL 384 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">ℝ℞℟℣ℤ℥℧ℨ℩KÅℬℭℯ</text><path d=\"M 50 264\nL 355 264\nL 355 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0xbe623813,
			},
			{
				name:        "fallback_notosans_subscripts",
				input:       "ⁱₐₑₒₓₔₕₖₗₘₙₚₛₜ",
				font:        styleLargeRoboto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">ⁱₐₑₒ</text><path d=\"M 50 64\nL 114 64\nL 114 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">ₓₔₕₖ</text><path d=\"M 50 164\nL 114 164\nL 114 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Roboto Medium',sans-serif\">ₗₘₙₚₛₜ</text><path d=\"M 50 264\nL 146 264\nL 146 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0x6a15590d,
			},
			{
				name:        "fallback_roboto_mathematical",
				input:       "∂∆∏∑-√∞∫≈≠≤≥◊",
				font:        styleLargeNoto,
				expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"50\" y=\"100\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">∂∆∏∑</text><path d=\"M 50 64\nL 143 64\nL 143 100\nL 50 100\nL 50 64\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"200\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">-√∞∫</text><path d=\"M 50 164\nL 139 164\nL 139 200\nL 50 200\nL 50 164\" style=\"stroke-width:1;stroke:red;fill:none\"/><text x=\"50\" y=\"300\" style=\"stroke:none;fill:blue;font-size:35.8px;font-family:'Noto Sans Display Medium',sans-serif\">≈≠≤≥◊</text><path d=\"M 50 264\nL 151 264\nL 151 300\nL 50 300\nL 50 264\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
				expectedCRC: 0x90d8dde6,
			},
		}

		renderTextBox := func(p *Painter, y int, text string, font FontStyle) {
			const xShift = 50
			textBox := p.MeasureText(text, 0, font)
			p.Text(text, xShift, y, 0, font)
			p.FilledRect(xShift, y-textBox.Height(), xShift+textBox.Width(), y,
				ColorTransparent, ColorRed, 1.0)
		}

		t.Run("svg", func(t *testing.T) {
			for _, tc := range tests {
				runes := []rune(tc.input)
				div := len(runes) / 3
				str1, str2, str3 := string(runes[:div]), string(runes[div:div*2]), string(runes[div*2:])
				t.Run(tc.name, func(t *testing.T) {
					svgP := NewPainter(PainterOptions{
						OutputFormat: ChartOutputSVG,
						Width:        600,
						Height:       400,
					})

					renderTextBox(svgP, 100, str1, tc.font)
					renderTextBox(svgP, 200, str2, tc.font)
					renderTextBox(svgP, 300, str3, tc.font)

					data, err := svgP.Bytes()
					require.NoError(t, err)
					assertEqualSVG(t, tc.expectedSVG, data)
				})
			}
		})
		t.Run("png", func(t *testing.T) {
			for _, tc := range tests {
				runes := []rune(tc.input)
				div := len(runes) / 3
				str1, str2, str3 := string(runes[:div]), string(runes[div:div*2]), string(runes[div*2:])
				t.Run(tc.name, func(t *testing.T) {
					pngP := NewPainter(PainterOptions{
						OutputFormat: ChartOutputPNG,
						Width:        600,
						Height:       400,
					})

					renderTextBox(pngP, 100, str1, tc.font)
					renderTextBox(pngP, 200, str2, tc.font)
					renderTextBox(pngP, 300, str3, tc.font)

					data, err := pngP.Bytes()
					require.NoError(t, err)
					assertEqualPNGCRC(t, tc.expectedCRC, data)
				})
			}
		})
	})
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
		FontColor: ColorBlackAlt1,
		Font:      GetDefaultFont(),
	}

	text := "Hello World!"
	box := p.TextFit(text, 0, 20, 80, fontStyle)
	assert.Equal(t, Box{Right: 45, Bottom: 37, IsSet: true}, box)

	box = p.TextFit(text, 0, 100, 200, fontStyle)
	assert.Equal(t, Box{Right: 84, Bottom: 16, IsSet: true}, box)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><text x=\"0\" y=\"20\" style=\"stroke:none;fill:rgb(51,51,51);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello</text><text x=\"0\" y=\"41\" style=\"stroke:none;fill:rgb(51,51,51);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World!</text><text x=\"0\" y=\"100\" style=\"stroke:none;fill:rgb(51,51,51);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Hello World!</text></svg>", buf)
}

func TestMultipleChartsOnPainter(t *testing.T) {
	t.Parallel()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        800,
		Height:       600,
	})
	p.FilledRect(0, 0, 800, 600, ColorWhite, ColorTransparent, 0.0)
	// set the space and theme for each chart
	topCenterPainter := p.Child(PainterBoxOption(NewBox(0, 0, 800, 300)),
		PainterThemeOption(GetTheme(ThemeVividLight)))
	bottomLeftPainter := p.Child(PainterBoxOption(NewBox(0, 300, 400, 600)),
		PainterThemeOption(GetTheme(ThemeAnt)))
	bottomRightPainter := p.Child(PainterBoxOption(NewBox(400, 300, 800, 600)),
		PainterThemeOption(GetTheme(ThemeLight)))

	pieOpt := makeBasicPieChartOption()
	pieOpt.Legend.Show = Ptr(false)
	pieOpt.Legend.Symbol = ""
	err := bottomLeftPainter.PieChart(pieOpt)
	require.NoError(t, err)
	err = bottomRightPainter.BarChart(makeBasicBarChartOption())
	require.NoError(t, err)
	lineOpt := makeBasicLineChartOption()
	lineOpt.YAxis[0].PreferNiceIntervals = Ptr(true)
	err = topCenterPainter.LineChart(lineOpt)
	require.NoError(t, err)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 800 600\"><path d=\"M 0 0\nL 800 0\nL 800 600\nL 0 600\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 0 300\nL 400 300\nL 400 600\nL 0 600\nL 0 300\" style=\"stroke:none;fill:white\"/><text x=\"185\" y=\"336\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><text x=\"187\" y=\"352\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sub</text><path d=\"M 200 473\nL 200 388\nA 85 85 119.89 0 1 274 515\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(91,143,249)\"/><path d=\"M 273 431\nL 286 423\nM 286 423\nL 301 423\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:none\"/><text x=\"304\" y=\"428\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path d=\"M 200 473\nL 274 515\nA 85 85 84.08 0 1 165 551\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(90,216,166)\"/><path d=\"M 226 553\nL 231 568\nM 231 568\nL 246 568\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:none\"/><text x=\"249\" y=\"573\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path d=\"M 200 473\nL 165 551\nA 85 85 66.35 0 1 115 473\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(93,112,146)\"/><path d=\"M 129 519\nL 116 527\nM 116 527\nL 101 527\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:none\"/><text x=\"0\" y=\"532\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path d=\"M 200 473\nL 115 473\nA 85 85 55.37 0 1 152 403\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(246,189,22)\"/><path d=\"M 125 434\nL 112 426\nM 112 426\nL 97 426\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:none\"/><text x=\"-4\" y=\"431\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path d=\"M 200 473\nL 152 403\nA 85 85 34.32 0 1 200 388\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(111,94,249)\"/><path d=\"M 175 392\nL 171 378\nM 171 378\nL 156 378\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:none\"/><text x=\"64\" y=\"383\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><path d=\"M 400 300\nL 800 300\nL 800 600\nL 400 600\nL 400 300\" style=\"stroke:none;fill:white\"/><text x=\"409\" y=\"316\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"409\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">175</text><text x=\"409\" y=\"379\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"409\" y=\"410\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">125</text><text x=\"409\" y=\"442\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">100</text><text x=\"418\" y=\"473\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">75</text><text x=\"418\" y=\"505\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">50</text><text x=\"418\" y=\"536\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">25</text><text x=\"427\" y=\"568\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 442 310\nL 790 310\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 341\nL 790 341\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 373\nL 790 373\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 405\nL 790 405\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 437\nL 790 437\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 468\nL 790 468\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 500\nL 790 500\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 442 532\nL 790 532\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 446 564\nL 790 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 446 569\nL 446 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 477 569\nL 477 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 508 569\nL 508 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 539 569\nL 539 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 571 569\nL 571 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 602 569\nL 602 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 633 569\nL 633 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 664 569\nL 664 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 696 569\nL 696 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 727 569\nL 727 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 758 569\nL 758 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 790 569\nL 790 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"445\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"507\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"570\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"632\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"663\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"726\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"763\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path d=\"M 451 562\nL 458 562\nL 458 563\nL 451 563\nL 451 562\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 479 558\nL 486 558\nL 486 563\nL 479 563\nL 479 558\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 508 556\nL 515 556\nL 515 563\nL 508 563\nL 508 556\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 537 535\nL 544 535\nL 544 563\nL 537 563\nL 537 535\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 565 532\nL 572 532\nL 572 563\nL 565 563\nL 565 532\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 594 467\nL 601 467\nL 601 563\nL 594 563\nL 594 467\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 623 392\nL 630 392\nL 630 563\nL 623 563\nL 623 392\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 651 359\nL 658 359\nL 658 563\nL 651 563\nL 651 359\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 680 523\nL 687 523\nL 687 563\nL 680 563\nL 680 523\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 709 539\nL 716 539\nL 716 563\nL 709 563\nL 709 539\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 737 556\nL 744 556\nL 744 563\nL 737 563\nL 737 556\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 766 560\nL 773 560\nL 773 563\nL 766 563\nL 766 560\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 461 561\nL 468 561\nL 468 563\nL 461 563\nL 461 561\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 489 557\nL 496 557\nL 496 563\nL 489 563\nL 489 557\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 518 553\nL 525 553\nL 525 563\nL 518 563\nL 518 553\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 547 531\nL 554 531\nL 554 563\nL 547 563\nL 547 531\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 575 528\nL 582 528\nL 582 563\nL 575 563\nL 575 528\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 604 475\nL 611 475\nL 611 563\nL 604 563\nL 604 475\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 633 341\nL 640 341\nL 640 563\nL 633 563\nL 633 341\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 661 333\nL 668 333\nL 668 563\nL 661 563\nL 661 333\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 690 503\nL 697 503\nL 697 563\nL 690 563\nL 690 503\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 719 541\nL 726 541\nL 726 563\nL 719 563\nL 719 541\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 747 557\nL 754 557\nL 754 563\nL 747 563\nL 747 557\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 776 562\nL 783 562\nL 783 563\nL 776 563\nL 776 562\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"450\" y=\"557\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"473\" y=\"553\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4.9</text><text x=\"507\" y=\"551\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"527\" y=\"530\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23.2</text><text x=\"555\" y=\"527\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25.6</text><text x=\"584\" y=\"462\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">76.7</text><text x=\"610\" y=\"387\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">135.6</text><text x=\"638\" y=\"354\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">162.2</text><text x=\"670\" y=\"518\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32.6</text><text x=\"705\" y=\"534\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"731\" y=\"551\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6.4</text><text x=\"760\" y=\"555\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3.3</text><text x=\"455\" y=\"556\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.6</text><text x=\"483\" y=\"552\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">5.9</text><text x=\"517\" y=\"548\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">9</text><text x=\"537\" y=\"526\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26.4</text><text x=\"565\" y=\"523\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28.7</text><text x=\"594\" y=\"470\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">70.7</text><text x=\"620\" y=\"336\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">175.6</text><text x=\"648\" y=\"328\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">182.2</text><text x=\"680\" y=\"498\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.7</text><text x=\"709\" y=\"536\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18.8</text><text x=\"746\" y=\"552\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"770\" y=\"557\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><path d=\"M 0 0\nL 800 0\nL 800 300\nL 0 300\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><path d=\"M 350 19\nL 380 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"365\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"382\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path d=\"M 411 19\nL 441 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"426\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"443\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"9\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.6k</text><text x=\"9\" y=\"79\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"9\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.2k</text><text x=\"22\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1k</text><text x=\"12\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"12\" y=\"187\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600</text><text x=\"12\" y=\"214\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">400</text><text x=\"12\" y=\"241\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"30\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 45 46\nL 790 46\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 73\nL 790 73\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 100\nL 790 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 127\nL 790 127\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 155\nL 790 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 182\nL 790 182\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 209\nL 790 209\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 45 236\nL 790 236\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 49 264\nL 790 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 49 269\nL 49 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 154 269\nL 154 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 260 269\nL 260 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 366 269\nL 366 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 472 269\nL 472 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 578 269\nL 578 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 684 269\nL 684 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 790 269\nL 790 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"96\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"202\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"308\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"414\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"521\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"627\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"732\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text><path d=\"M 101 248\nL 207 247\nL 313 251\nL 419 246\nL 525 252\nL 631 233\nL 737 236\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"101\" cy=\"248\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"207\" cy=\"247\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"313\" cy=\"251\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"419\" cy=\"246\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"525\" cy=\"252\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"631\" cy=\"233\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"737\" cy=\"236\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><path d=\"M 101 153\nL 207 138\nL 313 142\nL 419 137\nL 525 89\nL 631 83\nL 737 85\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"101\" cy=\"153\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"207\" cy=\"138\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"313\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"419\" cy=\"137\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"525\" cy=\"89\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"631\" cy=\"83\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"737\" cy=\"85\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/></svg>", buf)
}

func TestDashedLineStroke(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		svg    string
		pngCRC uint32
	}{
		{
			name: "simple_dashed",
			fn: func(p *Painter) {
				p.DashedLineStroke([]Point{
					{X: 10, Y: 20},
					{X: 30, Y: 40},
					{X: 50, Y: 20},
				}, ColorBlack, 1, []float64{5, 3})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"5.0, 3.0\" d=\"M 10 20\nL 30 40\nL 50 20\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
			pngCRC: 0xfb96dbae,
		},
		{
			name: "thick_dashed",
			fn: func(p *Painter) {
				p.DashedLineStroke([]Point{
					{X: 10, Y: 50},
					{X: 100, Y: 50},
					{X: 100, Y: 150},
					{X: 200, Y: 150},
				}, ColorRed, 3, []float64{8, 4})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"8.0, 4.0\" d=\"M 10 50\nL 100 50\nL 100 150\nL 200 150\" style=\"stroke-width:3;stroke:red;fill:none\"/></svg>",
			pngCRC: 0x5a3cb7f6,
		},
		{
			name: "short_dash_pattern",
			fn: func(p *Painter) {
				p.DashedLineStroke([]Point{
					{X: 50, Y: 100},
					{X: 350, Y: 100},
				}, ColorBlue, 2, []float64{2, 2})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"2.0, 2.0\" d=\"M 50 100\nL 350 100\" style=\"stroke-width:2;stroke:blue;fill:none\"/></svg>",
			pngCRC: 0xd1da5da7,
		},
		{
			name: "long_dash_pattern",
			fn: func(p *Painter) {
				p.DashedLineStroke([]Point{
					{X: 50, Y: 150},
					{X: 350, Y: 150},
				}, ColorGreen, 2, []float64{15, 10})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"15.0, 10.0\" d=\"M 50 150\nL 350 150\" style=\"stroke-width:2;stroke:green;fill:none\"/></svg>",
			pngCRC: 0xa031e6a2,
		},
		{
			name: "complex_path",
			fn: func(p *Painter) {
				p.DashedLineStroke([]Point{
					{X: 50, Y: 200},
					{X: 100, Y: 180},
					{X: 150, Y: 220},
					{X: 200, Y: 160},
					{X: 250, Y: 240},
					{X: 300, Y: 180},
				}, ColorFromHex("#FF6B35"), 2, []float64{6, 4})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"6.0, 4.0\" d=\"M 50 200\nL 100 180\nL 150 220\nL 200 160\nL 250 240\nL 300 180\" style=\"stroke-width:2;stroke:rgb(255,107,53);fill:none\"/></svg>",
			pngCRC: 0x193a9802,
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			svgP := NewPainter(PainterOptions{
				Width:        400,
				Height:       300,
				OutputFormat: ChartOutputSVG,
			})
			tc.fn(svgP)
			buf, err := svgP.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.svg, buf)

			pngP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        400,
				Height:       300,
			})
			tc.fn(pngP)
			data, err := pngP.Bytes()
			require.NoError(t, err)
			assertEqualPNGCRC(t, tc.pngCRC, data)
		})
	}
}

func TestSmoothDashedLineStroke(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
		svg    string
		pngCRC uint32
	}{
		{
			name: "smooth_dashed_basic",
			fn: func(p *Painter) {
				p.SmoothDashedLineStroke([]Point{
					{X: 10, Y: 20},
					{X: 20, Y: 40},
					{X: 30, Y: 60},
					{X: 40, Y: 50},
					{X: 50, Y: 40},
				}, 0.5, ColorBlack, 1, []float64{4, 2})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"4.0, 2.0\" d=\"M 10 20\nQ20,40 22,45\nQ30,60 32,57\nQ40,50 42,47\nQ40,50 50,40\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
			pngCRC: 0x13085dce,
		},
		{
			name: "smooth_dashed_high_tension",
			fn: func(p *Painter) {
				p.SmoothDashedLineStroke([]Point{
					{X: 50, Y: 50},
					{X: 100, Y: 100},
					{X: 150, Y: 60},
					{X: 200, Y: 120},
					{X: 250, Y: 80},
				}, 0.8, ColorRed, 2, []float64{8, 4})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"8.0, 4.0\" d=\"M 50 50\nQ100,100 120,84\nQ150,60 170,84\nQ200,120 220,104\nQ200,120 250,80\" style=\"stroke-width:2;stroke:red;fill:none\"/></svg>",
			pngCRC: 0x31c46648,
		},
		{
			name: "smooth_dashed_low_tension",
			fn: func(p *Painter) {
				p.SmoothDashedLineStroke([]Point{
					{X: 50, Y: 150},
					{X: 100, Y: 180},
					{X: 150, Y: 140},
					{X: 200, Y: 200},
					{X: 250, Y: 160},
				}, 0.2, ColorBlue, 1.5, []float64{6, 3})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"6.0, 3.0\" d=\"M 50 150\nQ100,180 105,176\nQ150,140 155,146\nQ200,200 205,196\nQ200,200 250,160\" style=\"stroke-width:1.5;stroke:blue;fill:none\"/></svg>",
			pngCRC: 0xb7a95b2c,
		},
		{
			name: "smooth_dashed_complex",
			fn: func(p *Painter) {
				points := []Point{
					{X: 30, Y: 220},
					{X: 60, Y: 200},
					{X: 90, Y: 240},
					{X: 120, Y: 180},
					{X: 150, Y: 260},
					{X: 180, Y: 200},
					{X: 210, Y: 220},
					{X: 240, Y: 180},
				}
				p.SmoothDashedLineStroke(points, 0.6, ColorFromHex("#9932CC"), 2, []float64{10, 5})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"10.0, 5.0\" d=\"M 30 220\nQ60,200 69,212\nQ90,240 99,222\nQ120,180 129,204\nQ150,260 159,242\nQ180,200 189,206\nQ210,220 219,208\nQ210,220 240,180\" style=\"stroke-width:2;stroke:rgb(153,50,204);fill:none\"/></svg>",
			pngCRC: 0x50be7e30,
		},
		{
			name: "smooth_dashed_dotted",
			fn: func(p *Painter) {
				p.SmoothDashedLineStroke([]Point{
					{X: 300, Y: 50},
					{X: 320, Y: 80},
					{X: 340, Y: 40},
					{X: 360, Y: 90},
					{X: 380, Y: 60},
				}, 0.7, ColorGreen, 3, []float64{1, 3})
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path stroke-dasharray=\"1.0, 3.0\" d=\"M 300 50\nQ320,80 327,66\nQ340,40 347,57\nQ360,90 367,79\nQ360,90 380,60\" style=\"stroke-width:3;stroke:green;fill:none\"/></svg>",
			pngCRC: 0xd1b5c8bf,
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			svgP := NewPainter(PainterOptions{
				Width:        400,
				Height:       300,
				OutputFormat: ChartOutputSVG,
			})
			tc.fn(svgP)
			buf, err := svgP.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.svg, buf)

			pngP := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        400,
				Height:       300,
			})
			tc.fn(pngP)
			data, err := pngP.Bytes()
			require.NoError(t, err)
			assertEqualPNGCRC(t, tc.pngCRC, data)
		})
	}
}

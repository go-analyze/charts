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
		result string
	}{
		{
			name: "circle",
			fn: func(p *Painter) {
				p.Circle(5, 2, 3, ColorTransparent, ColorTransparent, 1.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"7\" cy=\"13\" r=\"5\" style=\"stroke:none;fill:none\"/></svg>",
		},
		{
			name: "moveTo_lineTo",
			fn: func(p *Painter) {
				p.moveTo(1, 1)
				p.lineTo(2, 2)
				p.stroke(ColorTransparent, 1.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 6 11\nL 7 12\" style=\"stroke:none;fill:none\"/></svg>",
		},
		{
			name: "arc",
			fn: func(p *Painter) {
				p.arcTo(100, 100, 100, 100, 0, math.Pi/2)
				p.close()
				p.fillStroke(ColorBlue, ColorBlack, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 205 110\nA 100 100 90.00 0 1 105 210\nZ\" style=\"stroke-width:1;stroke:black;fill:blue\"/></svg>",
		},
		{
			name: "draw_background",
			fn: func(p *Painter) {
				p.drawBackground(ColorWhite)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 400 10\nL 400 300\nL 5 300\nL 5 10\" style=\"stroke:none;fill:white\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(Box{Left: 5, Top: 10}))
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
				}, ColorBlack, 1)
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
				}, 0.5, ColorBlack, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 15 30\nQ25,50 27,55\nQ35,70 37,67\nQ45,60 47,57\nQ55,50 57,60\nQ55,50 65,90\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "filled_rect",
			fn: func(p *Painter) {
				p.FilledRect(0, 0, 400, 300, ColorWhite, ColorWhite, 0.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 405 10\nL 405 310\nL 5 310\nL 5 10\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "filled_rect_center",
			fn: func(p *Painter) {
				p.FilledRect(100, 100, 200, 150, ColorWhite, ColorWhite, 0.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 105 110\nL 205 110\nL 205 160\nL 105 160\nL 105 110\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "filled_rect_center_border",
			fn: func(p *Painter) {
				p.FilledRect(100, 100, 200, 150, ColorWhite, ColorBlue, 1.0)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 105 110\nL 205 110\nL 205 160\nL 105 160\nL 105 110\" style=\"stroke-width:1;stroke:blue;fill:white\"/></svg>",
		},
		{
			name: "pin",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.Pin(30, 30, 30, c, c, 1)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 31 47\nA 15 15 330.00 1 1 39 47\nL 35 33\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path  d=\"M 20 33\nQ35,70 50,33\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
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
			name: "horizontal_mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.HorizontalMarkLine(0, 20, 300, c, c, 1, []float64{4, 2})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"8\" cy=\"30\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 14 30\nL 289 30\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 289 25\nL 305 30\nL 289 35\nL 294 30\nL 289 25\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "vertical_mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.VerticalMarkLine(200, 100, 100, c, c, 1, []float64{4, 2})
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><circle cx=\"205\" cy=\"207\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 205 110\nL 205 210\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 200 126\nL 205 110\nL 210 126\nL 205 121\nL 200 126\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
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
				p = p.Child(PainterBoxOption(NewBox(200, 0, 400, 200)))
				opt = makeMinimalLineChartOption()
				opt.Theme = GetDefaultTheme().WithBackgroundColor(ColorTransparent)
				_ = p.LineChart(opt)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 300\"><path  d=\"M 5 10\nL 400 10\nL 400 300\nL 5 300\nL 5 10\" style=\"stroke:none;fill:white\"/><text x=\"14\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1440</text><text x=\"14\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1280</text><text x=\"14\" y=\"85\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1120</text><text x=\"22\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"145\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"174\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"204\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"234\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"294\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 55 20\nL 390 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 50\nL 390 50\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 80\nL 390 80\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 110\nL 390 110\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 140\nL 390 140\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 170\nL 390 170\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 200\nL 390 200\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 230\nL 390 230\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 55 260\nL 390 260\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 82 268\nL 129 266\nL 176 272\nL 224 265\nL 271 274\nL 318 247\nL 366 251\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path  d=\"M 82 137\nL 129 116\nL 176 122\nL 224 115\nL 271 49\nL 318 41\nL 366 43\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><path  d=\"M 200 0\nL 400 0\nL 400 200\nL 200 200\nL 200 0\" style=\"stroke:none;fill:none\"/><text x=\"209\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"209\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"209\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"221\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"221\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"221\" y=\"114\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"221\" y=\"134\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"221\" y=\"154\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"221\" y=\"174\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"239\" y=\"194\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 254 10\nL 390 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 30\nL 390 30\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 50\nL 390 50\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 70\nL 390 70\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 90\nL 390 90\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 110\nL 390 110\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 130\nL 390 130\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 150\nL 390 150\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 254 170\nL 390 170\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 258 175\nL 280 174\nL 302 178\nL 324 174\nL 346 179\nL 368 162\nL 390 164\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path  d=\"M 258 88\nL 280 74\nL 302 78\nL 324 74\nL 346 29\nL 368 24\nL 390 25\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterPaddingOption(Box{Left: 5, Top: 10}))
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
				}, 5, true, true, ColorBlue, ColorBlue, 1)
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
				}, 5, false, true, ColorBlue, ColorBlue, 1)
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
				}, 5, true, false, ColorBlue, ColorBlue, 1)
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

	t.Run("basic", func(t *testing.T) {
		assert.Equal(t, Box{Right: 84, Bottom: 16, IsSet: true},
			svgP.MeasureText("Hello World!", 0, style))
		assert.Equal(t, Box{Right: 99, Bottom: 14, IsSet: true},
			pngP.MeasureText("Hello World!", 0, style))
		assert.Equal(t, Box{Right: 99, Bottom: 14, IsSet: true},
			jpgP.MeasureText("Hello World!", 0, style))
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
	err = topCenterPainter.LineChart(makeBasicLineChartOption())
	require.NoError(t, err)

	buf, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 800 600\"><path  d=\"M 0 0\nL 800 0\nL 800 600\nL 0 600\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 0 300\nL 400 300\nL 400 600\nL 0 600\nL 0 300\" style=\"stroke:none;fill:white\"/><text x=\"185\" y=\"336\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><text x=\"187\" y=\"352\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sub</text><path  d=\"M 200 473\nL 200 388\nA 85 85 119.89 0 1 274 515\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(91,143,249)\"/><path  d=\"M 273 431\nL 286 423\nM 286 423\nL 301 423\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:none\"/><text x=\"304\" y=\"428\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 200 473\nL 274 515\nA 85 85 84.08 0 1 165 551\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(90,216,166)\"/><path  d=\"M 226 553\nL 231 568\nM 231 568\nL 246 568\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:none\"/><text x=\"249\" y=\"573\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 200 473\nL 165 551\nA 85 85 66.35 0 1 115 473\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(93,112,146)\"/><path  d=\"M 129 519\nL 116 527\nM 116 527\nL 101 527\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:none\"/><text x=\"1\" y=\"532\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 200 473\nL 115 473\nA 85 85 55.37 0 1 152 403\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(246,189,22)\"/><path  d=\"M 125 434\nL 112 426\nM 112 426\nL 97 426\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:none\"/><text x=\"-3\" y=\"431\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 200 473\nL 152 403\nA 85 85 34.32 0 1 200 388\nL 200 473\nZ\" style=\"stroke:none;fill:rgb(111,94,249)\"/><path  d=\"M 175 392\nL 171 378\nM 171 378\nL 156 378\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:none\"/><text x=\"64\" y=\"383\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><path  d=\"M 400 300\nL 800 300\nL 800 600\nL 400 600\nL 400 300\" style=\"stroke:none;fill:white\"/><text x=\"409\" y=\"316\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">189</text><text x=\"409\" y=\"344\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">168</text><text x=\"409\" y=\"372\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">147</text><text x=\"409\" y=\"400\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">126</text><text x=\"409\" y=\"428\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">105</text><text x=\"418\" y=\"456\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84</text><text x=\"418\" y=\"484\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">63</text><text x=\"418\" y=\"512\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"418\" y=\"540\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">21</text><text x=\"427\" y=\"568\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 442 310\nL 790 310\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 338\nL 790 338\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 366\nL 790 366\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 394\nL 790 394\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 422\nL 790 422\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 451\nL 790 451\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 479\nL 790 479\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 507\nL 790 507\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 442 535\nL 790 535\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 446 564\nL 790 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 446 569\nL 446 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 477 569\nL 477 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 508 569\nL 508 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 539 569\nL 539 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 571 569\nL 571 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 602 569\nL 602 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 633 569\nL 633 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 664 569\nL 664 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 696 569\nL 696 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 727 569\nL 727 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 758 569\nL 758 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 790 569\nL 790 564\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"445\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"507\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"570\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"632\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"663\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"726\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"763\" y=\"590\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 451 562\nL 458 562\nL 458 563\nL 451 563\nL 451 562\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 479 558\nL 486 558\nL 486 563\nL 479 563\nL 479 558\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 508 555\nL 515 555\nL 515 563\nL 508 563\nL 508 555\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 537 533\nL 544 533\nL 544 563\nL 537 563\nL 537 533\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 565 530\nL 572 530\nL 572 563\nL 565 563\nL 565 530\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 594 461\nL 601 461\nL 601 563\nL 594 563\nL 594 461\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 623 382\nL 630 382\nL 630 563\nL 623 563\nL 623 382\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 651 347\nL 658 347\nL 658 563\nL 651 563\nL 651 347\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 680 521\nL 687 521\nL 687 563\nL 680 563\nL 680 521\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 709 538\nL 716 538\nL 716 563\nL 709 563\nL 709 538\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 737 556\nL 744 556\nL 744 563\nL 737 563\nL 737 556\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 766 560\nL 773 560\nL 773 563\nL 766 563\nL 766 560\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 461 561\nL 468 561\nL 468 563\nL 461 563\nL 461 561\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 489 557\nL 496 557\nL 496 563\nL 489 563\nL 489 557\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 518 552\nL 525 552\nL 525 563\nL 518 563\nL 518 552\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 547 529\nL 554 529\nL 554 563\nL 547 563\nL 547 529\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 575 526\nL 582 526\nL 582 563\nL 575 563\nL 575 526\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 604 469\nL 611 469\nL 611 563\nL 604 563\nL 604 469\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 633 329\nL 640 329\nL 640 563\nL 633 563\nL 633 329\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 661 320\nL 668 320\nL 668 563\nL 661 563\nL 661 320\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 690 499\nL 697 499\nL 697 563\nL 690 563\nL 690 499\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 719 539\nL 726 539\nL 726 563\nL 719 563\nL 719 539\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 747 556\nL 754 556\nL 754 563\nL 747 563\nL 747 556\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 776 561\nL 783 561\nL 783 563\nL 776 563\nL 776 561\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"450\" y=\"557\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"473\" y=\"553\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4.9</text><text x=\"507\" y=\"550\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"527\" y=\"528\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23.2</text><text x=\"555\" y=\"525\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25.6</text><text x=\"584\" y=\"456\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">76.7</text><text x=\"610\" y=\"377\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">135.6</text><text x=\"638\" y=\"342\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">162.2</text><text x=\"670\" y=\"516\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32.6</text><text x=\"705\" y=\"533\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"731\" y=\"551\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6.4</text><text x=\"760\" y=\"555\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3.3</text><text x=\"455\" y=\"556\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.6</text><text x=\"483\" y=\"552\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">5.9</text><text x=\"517\" y=\"547\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">9</text><text x=\"537\" y=\"524\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26.4</text><text x=\"565\" y=\"521\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28.7</text><text x=\"594\" y=\"464\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">70.7</text><text x=\"620\" y=\"324\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">175.6</text><text x=\"648\" y=\"315\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">182.2</text><text x=\"680\" y=\"494\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.7</text><text x=\"709\" y=\"534\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18.8</text><text x=\"746\" y=\"551\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"770\" y=\"556\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><path  d=\"M 0 0\nL 800 0\nL 800 300\nL 0 300\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 350 19\nL 380 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"365\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"382\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 411 19\nL 441 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"426\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"443\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"9\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"9\" y=\"76\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"9\" y=\"100\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"21\" y=\"124\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"21\" y=\"148\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"21\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"21\" y=\"196\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"21\" y=\"220\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"21\" y=\"244\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"39\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 54 46\nL 790 46\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 70\nL 790 70\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 94\nL 790 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 118\nL 790 118\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 142\nL 790 142\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 167\nL 790 167\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 191\nL 790 191\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 215\nL 790 215\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 54 239\nL 790 239\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 58 264\nL 790 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 58 269\nL 58 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 162 269\nL 162 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 267 269\nL 267 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 371 269\nL 371 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 476 269\nL 476 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 580 269\nL 580 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 685 269\nL 685 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 790 269\nL 790 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"105\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"209\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"314\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"418\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"524\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"628\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"732\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text><path  d=\"M 110 246\nL 214 245\nL 319 249\nL 423 244\nL 528 251\nL 632 230\nL 737 233\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"110\" cy=\"246\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"214\" cy=\"245\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"319\" cy=\"249\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"423\" cy=\"244\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"528\" cy=\"251\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"632\" cy=\"230\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"737\" cy=\"233\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><path  d=\"M 110 140\nL 214 123\nL 319 128\nL 423 123\nL 528 69\nL 632 63\nL 737 65\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"110\" cy=\"140\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"214\" cy=\"123\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"319\" cy=\"128\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"423\" cy=\"123\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"528\" cy=\"69\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"632\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"737\" cy=\"65\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/></svg>", buf)
}

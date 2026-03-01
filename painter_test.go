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
		pngCRC uint32
	}{
		{
			name: "circle",
			fn: func(p *Painter) {
				p.Circle(5, 2, 3, ColorTransparent, ColorTransparent, 1.0)
			},
			pngCRC: 0x2083f7bd,
		},
		{
			name: "moveTo_lineTo",
			fn: func(p *Painter) {
				p.moveTo(1, 1)
				p.lineTo(2, 2)
				p.stroke(ColorTransparent, 1.0)
			},
			pngCRC: 0x2083f7bd,
		},
		{
			name: "arc",
			fn: func(p *Painter) {
				p.arcTo(100, 100, 100, 100, 0, math.Pi/2)
				p.close()
				p.fillStroke(ColorBlue, ColorBlack, 1)
			},
			pngCRC: 0x5445e3e7,
		},
		{
			name: "draw_background",
			fn: func(p *Painter) {
				p.drawBackground(ColorWhite)
			},
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
			assertTestdataSVG(t, data)

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
		pngCRC uint32
	}{
		{
			name: "text",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6, 0, FontStyle{})
			},
			pngCRC: 0x2083f7bd,
		},
		{
			name: "text_rotated",
			fn: func(p *Painter) {
				p.Text("hello world!", 3, 6, DegreesToRadians(90), FontStyle{})
			},
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
			pngCRC: 0x3056a98,
		},
		{
			name: "filled_rect",
			fn: func(p *Painter) {
				p.FilledRect(0, 0, 400, 300, ColorWhite, ColorWhite, 0.0)
			},
			pngCRC: 0x60f3dd98,
		},
		{
			name: "filled_rect_center",
			fn: func(p *Painter) {
				p.FilledRect(100, 100, 200, 150, ColorWhite, ColorWhite, 0.0)
			},
			pngCRC: 0x540b2357,
		},
		{
			name: "filled_rect_center_border",
			fn: func(p *Painter) {
				p.FilledRect(100, 100, 200, 150, ColorWhite, ColorBlue, 1.0)
			},
			pngCRC: 0xd6c2a417,
		},
		{
			name: "pin",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.Pin(30, 30, 30, c, c, 1)
			},
			pngCRC: 0x981d8eb5,
		},
		{
			name: "arrow_left",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowLeft(30, 30, 16, 10, c, c, 1)
			},
			pngCRC: 0x3415dab,
		},
		{
			name: "arrow_right",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowRight(30, 30, 16, 10, c, c, 1)
			},
			pngCRC: 0x142dfb03,
		},
		{
			name: "arrow_up",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowUp(30, 30, 10, 16, c, c, 1)
			},
			pngCRC: 0xe17c9204,
		},
		{
			name: "arrow_down",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.ArrowDown(30, 30, 10, 16, c, c, 1)
			},
			pngCRC: 0xd56c309d,
		},
		{
			name: "horizontal_mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.HorizontalMarkLine(0, 20, 300, c, c, 1, []float64{4, 2})
			},
			pngCRC: 0xa4ca1cb8,
		},
		{
			name: "vertical_mark_line",
			fn: func(p *Painter) {
				c := Color{R: 84, G: 112, B: 198, A: 255}
				p.VerticalMarkLine(200, 100, 100, c, c, 1, []float64{4, 2})
			},
			pngCRC: 0x681c0b4e,
		},
		{
			name: "polygon",
			fn: func(p *Painter) {
				p.Polygon(Point{X: 100, Y: 100}, 50, 6, Color{R: 84, G: 112, B: 198, A: 255}, 1)
			},
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
			assertTestdataSVG(t, data)

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
	var drawDebugBox bool

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
				assertEqualSVG(t, nil, data)
			} else {
				expectedResult := fmt.Sprintf(expectedTemplate, tt.expectedY, tt.degrees%360, tt.expectedY, name)
				assertEqualSVG(t, []byte(expectedResult), data)
			}
		})
	}
}

func TestPainterRoundedRect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
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
			assertTestdataSVG(t, buf)

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
			expectedCRC uint32
		}{
			{
				name:        "basic",
				input:       "Hello World!",
				font:        styleLargeNoto,
				expectedCRC: 0x1bede00d,
			},
			{
				name:        "emojis",
				input:       "⭐❓💰🔥💯🎯🚀⚡🌟🎉🎊",
				font:        styleLargeNoto,
				expectedCRC: 0xaaf41b1b,
			},
			{
				name:        "shapes",
				input:       "▫●□▲▼◇★○△▪▴▾◆◯⬟⬠⬡⬢⬣⬤⬥",
				font:        styleLargeNoto,
				expectedCRC: 0x2a17094,
			},
			{
				name:        "playing_cards",
				input:       "🂡🂢🂫🃄🃍🃘🃞🃟",
				font:        styleLargeNoto,
				expectedCRC: 0x94611d5b,
			},
			{
				name:        "faces",
				input:       "😂😍🤣😊😭😘😎🤔😴😋😉😏😬😐😑😮😯",
				font:        styleLargeNoto,
				expectedCRC: 0x64607851,
			},
			{
				name:        "fallback_notosans_currency",
				input:       "₠₡₢₥₭₮₯₰₲₳₴₵₶₷₸₻₾₿",
				font:        styleLargeRoboto,
				expectedCRC: 0xe1e34ff1,
			},
			{
				name:        "fallback_notosans_letterlike",
				input:       "℀℁ℂ℃℄℆ℇ℈℉ℊℋℌℍℎℏℐℑℒ℔ℕ℗℘ℙℚℛℜℝ℞℟℣ℤ℥℧ℨ℩KÅℬℭℯ",
				font:        styleLargeRoboto,
				expectedCRC: 0xbe623813,
			},
			{
				name:        "fallback_notosans_subscripts",
				input:       "ⁱₐₑₒₓₔₕₖₗₘₙₚₛₜ",
				font:        styleLargeRoboto,
				expectedCRC: 0x6a15590d,
			},
			{
				name:        "fallback_roboto_mathematical",
				input:       "∂∆∏∑-√∞∫≈≠≤≥◊",
				font:        styleLargeNoto,
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
					assertTestdataSVG(t, data)
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
	assertTestdataSVG(t, buf)
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
	assertTestdataSVG(t, buf)
}

func TestDashedLineStroke(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fn     func(*Painter)
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
			assertTestdataSVG(t, buf)

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
			assertTestdataSVG(t, buf)

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

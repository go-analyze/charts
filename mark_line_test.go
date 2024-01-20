package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestMarkLine(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				markLine := NewMarkLinePainter(p)
				series := NewSeriesFromValues([]float64{
					1,
					2,
					3,
				})
				series.MarkLine = NewMarkLine(
					SeriesMarkDataTypeMax,
					SeriesMarkDataTypeAverage,
					SeriesMarkDataTypeMin,
				)
				markLine.Add(markLineRenderOption{
					FillColor:   drawing.ColorBlack,
					FontColor:   drawing.ColorBlack,
					StrokeColor: drawing.ColorBlack,
					Series:      series,
					Range: NewRange(AxisRangeOption{
						Painter:     p,
						Min:         0,
						Max:         5,
						Size:        p.Height(),
						DivideCount: 6,
					}),
				})
				_, err := markLine.Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<circle cx=\"23\" cy=\"272\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 272\nL 562 272\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 267\nL 578 272\nL 562 277\nL 567 272\nL 562 267\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"276\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><circle cx=\"23\" cy=\"308\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 308\nL 562 308\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 303\nL 578 308\nL 562 313\nL 567 308\nL 562 303\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"312\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"23\" cy=\"344\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 344\nL 562 344\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 339\nL 578 344\nL 562 349\nL 567 344\nL 562 339\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"348\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1</text></svg>",
		},
	}
	for _, tt := range tests {
		p, err := NewPainter(PainterOptions{
			Type:   ChartOutputSVG,
			Width:  600,
			Height: 400,
		}, PainterThemeOption(defaultTheme))
		assert.Nil(err)
		data, err := tt.render(p.Child(PainterPaddingOption(Box{
			Left:   20,
			Top:    20,
			Right:  20,
			Bottom: 20,
		})))
		assert.Nil(err)
		assert.Equal(tt.result, string(data))
	}
}

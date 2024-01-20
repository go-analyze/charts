package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestMarkPoint(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				series := NewSeriesFromValues([]float64{
					1,
					2,
					3,
				})
				series.MarkPoint = NewMarkPoint(SeriesMarkDataTypeMax)
				markPoint := NewMarkPointPainter(p)
				markPoint.Add(markPointRenderOption{
					FillColor: drawing.ColorBlack,
					Series:    series,
					Points: []Point{
						{
							X: 10,
							Y: 10,
						},
						{
							X: 30,
							Y: 30,
						},
						{
							X: 50,
							Y: 50,
						},
					},
				})
				_, err := markPoint.Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 67 62\nA 15 15 330.00 1 1 73 62\nL 70 48\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0)\"/><path  d=\"M 55 48\nQ70,85 85,48\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0)\"/><text x=\"66\" y=\"53\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text></svg>",
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

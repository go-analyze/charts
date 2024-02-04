package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestMarkLine(t *testing.T) {
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
					Range:       NewRange(p, p.Height(), 6, 0.0, 5.0, true),
				})
				if _, err := markLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<circle cx=\"23\" cy=\"175\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 175\nL 562 175\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 170\nL 578 175\nL 562 180\nL 567 175\nL 562 170\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"179\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><circle cx=\"23\" cy=\"243\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 243\nL 562 243\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 238\nL 578 243\nL 562 248\nL 567 243\nL 562 238\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"247\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"23\" cy=\"312\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 312\nL 562 312\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 307\nL 578 312\nL 562 317\nL 567 312\nL 562 307\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"316\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1</text></svg>",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				Type:   ChartOutputSVG,
				Width:  600,
				Height: 400,
			}, PainterThemeOption(defaultTheme))
			require.NoError(t, err)
			data, err := tt.render(p.Child(PainterPaddingOption(Box{
				Left:   20,
				Top:    20,
				Right:  20,
				Bottom: 20,
			})))
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, string(data))
		})
	}
}

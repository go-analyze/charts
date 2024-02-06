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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<circle cx=\"23\" cy=\"168\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 168\nL 562 168\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 163\nL 578 168\nL 562 173\nL 567 168\nL 562 163\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"172\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><circle cx=\"23\" cy=\"233\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 233\nL 562 233\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 228\nL 578 233\nL 562 238\nL 567 233\nL 562 228\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"237\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"23\" cy=\"299\" r=\"3\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 299\nL 562 299\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 294\nL 578 299\nL 562 304\nL 567 299\nL 562 294\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(0,0,0,1.0)\"/><text x=\"580\" y=\"303\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1</text></svg>",
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

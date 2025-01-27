package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				markLine := newMarkLinePainter(p)
				series := Series{
					Data: []float64{1, 2, 3},
				}
				series.MarkLine = NewMarkLine(
					SeriesMarkDataTypeMax,
					SeriesMarkDataTypeAverage,
					SeriesMarkDataTypeMin,
				)
				markLine.Add(markLineRenderOption{
					fillColor:   ColorBlack,
					fontColor:   ColorBlack,
					strokeColor: ColorBlack,
					series:      series,
					axisRange: newRange(p, nil,
						p.Height(), 6, 0.0, 5.0, 0.0, 0.0),
				})
				if _, err := markLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><circle cx=\"23\" cy=\"164\" r=\"3\" style=\"stroke-width:1;stroke:black;fill:black\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 164\nL 562 164\" style=\"stroke-width:1;stroke:black;fill:black\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 159\nL 578 164\nL 562 169\nL 567 164\nL 562 159\" style=\"stroke-width:1;stroke:black;fill:black\"/><text x=\"580\" y=\"168\" style=\"stroke:none;fill:black;font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><circle cx=\"23\" cy=\"236\" r=\"3\" style=\"stroke-width:1;stroke:black;fill:black\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 236\nL 562 236\" style=\"stroke-width:1;stroke:black;fill:black\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 231\nL 578 236\nL 562 241\nL 567 236\nL 562 231\" style=\"stroke-width:1;stroke:black;fill:black\"/><text x=\"580\" y=\"240\" style=\"stroke:none;fill:black;font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"23\" cy=\"308\" r=\"3\" style=\"stroke-width:1;stroke:black;fill:black\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 29 308\nL 562 308\" style=\"stroke-width:1;stroke:black;fill:black\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 303\nL 578 308\nL 562 313\nL 567 308\nL 562 303\" style=\"stroke-width:1;stroke:black;fill:black\"/><text x=\"580\" y=\"312\" style=\"stroke:none;fill:black;font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1</text></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)))
			data, err := tt.render(p.Child(PainterPaddingOption(Box{
				Left:   20,
				Top:    20,
				Right:  20,
				Bottom: 20,
			})))
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

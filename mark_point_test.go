package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkPoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				series := Series{
					Data: []float64{1, 2, 3},
				}
				series.MarkPoint = NewMarkPoint(SeriesMarkDataTypeMax)
				markPoint := newMarkPointPainter(p)
				markPoint.Add(markPointRenderOption{
					fillColor: ColorBlack,
					series:    series,
					points: []Point{
						{X: 10, Y: 10},
						{X: 30, Y: 30},
						{X: 50, Y: 50},
					},
				})
				if _, err := markPoint.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 67 62\nA 14 14 330.00 1 1 73 62\nL 70 49\nZ\" style=\"stroke:none;fill:black\"/><path  d=\"M 56 49\nQ70,84 84,49\nZ\" style=\"stroke:none;fill:black\"/><text x=\"66\" y=\"54\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text></svg>",
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

package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestGrid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				_, err := NewGridPainter(p, GridPainterOption{
					StrokeColor:       drawing.ColorBlack,
					Columns:           6,
					Rows:              6,
					IgnoreFirstRow:    true,
					IgnoreLastRow:     true,
					IgnoreFirstColumn: true,
					IgnoreLastColumn:  true,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 100 0\nL 100 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 200 0\nL 200 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 300 0\nL 300 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 400 0\nL 400 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 500 0\nL 500 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 66\nL 600 66\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 133\nL 600 133\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 200\nL 600 200\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 266\nL 600 266\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 333\nL 600 333\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/></svg>",
		},
		{
			render: func(p *Painter) ([]byte, error) {
				_, err := NewGridPainter(p, GridPainterOption{
					StrokeColor: drawing.ColorBlack,
					ColumnSpans: []int{2, 5, 3},
					Rows:        6,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 0 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 120 0\nL 120 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 420 0\nL 420 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 600 0\nL 600 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 0\nL 600 0\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 66\nL 600 66\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 133\nL 600 133\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 200\nL 600 200\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 266\nL 600 266\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 333\nL 600 333\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/><path  d=\"M 0 400\nL 600 400\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:rgba(255,255,255,0.0)\"/></svg>",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)))
			require.NoError(t, err)
			data, err := tt.render(p)
			require.NoError(t, err)
			assert.Equal(t, tt.result, string(data))
		})
	}
}

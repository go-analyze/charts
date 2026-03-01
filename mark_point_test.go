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
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				markPoint := newMarkPointPainter(p)
				markPoint.add(markPointRenderOption{
					fillColor:    ColorBlack,
					seriesValues: []float64{1, 2, 3},
					markpoints:   NewSeriesMarkList(SeriesMarkTypeMax),
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
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)))
			data, err := tt.render(p.Child(PainterPaddingOption(NewBoxEqual(20))))
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

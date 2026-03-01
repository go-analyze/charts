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
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				markLine := newMarkLinePainter(p)
				markLine.add(markLineRenderOption{
					fillColor:    ColorBlack,
					fontColor:    ColorBlack,
					strokeColor:  ColorBlack,
					seriesValues: []float64{1, 2, 3},
					marklines:    NewSeriesMarkList(SeriesMarkTypeMax, SeriesMarkTypeAverage, SeriesMarkTypeMin),
					axisRange:    newTestRange(p.Height(), 6, 0.0, 5.0, 0.0, 0.0),
				})
				if _, err := markLine.Render(); err != nil {
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

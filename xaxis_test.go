package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBottomXAxis(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		makeOption func() XAxisOption
		makeValue  bool
	}{
		{
			name: "basic",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "value",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			makeValue: true,
		},
		{
			name: "title",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Title:  "Title",
					Labels: []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "boundary_gap_disabled",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:      []string{"a", "b", "c", "d"},
					BoundaryGap: Ptr(false),
				}
			},
		},
		{
			name: "font_style",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:        []string{"abc", "def", "ghi"},
					FontStyle:     NewFontStyleWithSize(20),
					LabelRotation: DegreesToRadians(90),
				}
			},
		},
		{
			name: "label_start_offset",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:         []string{"a", "b", "c", "d", "e", "f"},
					DataStartIndex: 2,
				}
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			theme := GetTheme(ThemeLight)
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(theme), PainterPaddingOption(NewBoxEqual(100)))

			xAxisOpt := tt.makeOption()
			xAxisOpt = *xAxisOpt.prep(theme)
			aRange := newTestRangeForLabels(xAxisOpt.Labels, xAxisOpt.LabelRotation,
				fillFontStyleDefaults(xAxisOpt.LabelFontStyle, defaultFontSize, theme.GetXAxisTextColor()))
			aRange.isCategory = !tt.makeValue
			_, err := newAxisPainter(p, xAxisOpt.toAxisOption(aRange)).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

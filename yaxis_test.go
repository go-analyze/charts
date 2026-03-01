package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYAxis(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOption       func() *YAxisOption
		makeCategory     bool
		disableSplitLine bool
	}{
		{
			name: "basic",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionLeft, // typically defaulted in defaultRender
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "category",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position:       PositionLeft, // typically defaulted in defaultRender
					Labels:         []string{"a", "b", "c", "d"},
					isCategoryAxis: true,
				}
			},
			makeCategory: true,
		},
		{
			name: "font_style_with_rotation",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position:      PositionLeft, // typically defaulted in defaultRender
					Labels:        []string{"a", "b", "c"},
					FontStyle:     NewFontStyleWithSize(20),
					LabelRotation: DegreesToRadians(270),
				}
			},
		},
		{
			name: "lines",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position:      PositionLeft, // typically defaulted in defaultRender
					Labels:        []string{"a", "b", "c", "d"},
					SplitLineShow: Ptr(true),
					SpineLineShow: Ptr(true),
				}
			},
		},
		{
			name: "title_left",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionLeft, // typically defaulted in defaultRender
					Title:    "title",
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "title_right",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionRight,
					Title:    "title",
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "split_line_disable",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionRight,
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
			disableSplitLine: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			theme := GetTheme(ThemeLight)
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(theme),
				PainterPaddingOption(NewBoxEqual(100)))

			yAxisOpt := tt.makeOption()
			yAxisOpt = yAxisOpt.prep(theme)
			aRange := newTestRangeForLabels(yAxisOpt.Labels, yAxisOpt.LabelRotation,
				fillFontStyleDefaults(yAxisOpt.LabelFontStyle, defaultFontSize, theme.GetYAxisTextColor()))
			aRange.isCategory = tt.makeCategory
			axisOpt := yAxisOpt.toAxisOption(aRange)
			if tt.disableSplitLine {
				axisOpt.splitLineShow = false
			}
			_, err := newAxisPainter(p, axisOpt).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

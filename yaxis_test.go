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
		result           string
	}{
		{
			name: "basic",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionLeft, // typically defaulted in defaultRender
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"101\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"102\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"101\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"101\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><path d=\"M 116 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 116 166\nL 500 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 116 233\nL 500 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 120 100\nL 120 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 100\nL 120 100\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 150\nL 120 150\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 200\nL 120 200\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 250\nL 120 250\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 300\nL 120 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"101\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"102\" y=\"180\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"101\" y=\"229\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"101\" y=\"279\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"96\" y=\"104\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,96,104)\">c</text><text x=\"96\" y=\"203\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,96,203)\">b</text><text x=\"96\" y=\"302\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,96,302)\">a</text><path d=\"M 126 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 126 200\nL 500 200\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 120 100\nL 120 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 100\nL 120 100\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 166\nL 120 166\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 233\nL 120 233\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 115 300\nL 120 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"101\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"102\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"101\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"101\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><path d=\"M 116 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 116 166\nL 500 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 116 233\nL 500 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"114\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,114,213)\">title</text><text x=\"121\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"122\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"121\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"121\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><path d=\"M 136 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 136 166\nL 500 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 136 233\nL 500 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"486\" y=\"187\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,486,187)\">title</text><text x=\"470\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"470\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"470\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"470\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><path d=\"M 100 100\nL 460 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 100 166\nL 460 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 100 233\nL 460 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result:           "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"490\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"490\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"490\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"490\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text></svg>",
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
			assertEqualSVG(t, tt.result, data)
		})
	}
}

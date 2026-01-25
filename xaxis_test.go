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
		result     string
	}{
		{
			name: "basic",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 100 272\nL 500 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 100 277\nL 100 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 200 277\nL 200 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 300 277\nL 300 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 400 277\nL 400 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 500 277\nL 500 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"146\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"246\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"346\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"446\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>",
		},
		{
			name: "value",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			makeValue: true,
			result:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"99\" y=\"293\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"232\" y=\"293\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"365\" y=\"293\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"491\" y=\"293\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path d=\"M 233 100\nL 233 272\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 366 100\nL 366 272\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 500 100\nL 500 272\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "title",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Title:  "Title",
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"284\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><path d=\"M 100 253\nL 500 253\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 100 258\nL 100 253\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 200 258\nL 200 253\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 300 258\nL 300 253\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 400 258\nL 400 253\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 500 258\nL 500 253\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"146\" y=\"277\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"246\" y=\"277\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"346\" y=\"277\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"446\" y=\"277\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>",
		},
		{
			name: "boundary_gap_disabled",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:      []string{"a", "b", "c", "d"},
					BoundaryGap: Ptr(false),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 100 272\nL 500 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 100 277\nL 100 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 233 277\nL 233 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 366 277\nL 366 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 500 277\nL 500 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"99\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"232\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"365\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"491\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 100 237\nL 500 237\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 100 242\nL 100 237\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 233 242\nL 233 237\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 366 242\nL 366 237\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 500 242\nL 500 237\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"153\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,153,255)\">abc</text><text x=\"286\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,286,255)\">def</text><text x=\"420\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,420,255)\">ghi</text></svg>",
		},
		{
			name: "label_start_offset",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:         []string{"a", "b", "c", "d", "e", "f"},
					DataStartIndex: 2,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 100 272\nL 500 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 100 277\nL 100 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 166 277\nL 166 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 233 277\nL 233 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 300 277\nL 300 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 366 277\nL 366 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 433 277\nL 433 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 500 277\nL 500 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"129\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"195\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"262\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"329\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><text x=\"395\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">e</text><text x=\"463\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text></svg>",
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
			assertEqualSVG(t, tt.result, data)
		})
	}
}

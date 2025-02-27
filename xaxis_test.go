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
		result     string
	}{
		{
			name: "basic",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 100 275\nL 100 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 200 275\nL 200 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 300 275\nL 300 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 400 275\nL 400 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 500 275\nL 500 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 100 270\nL 500 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"146\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"246\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"346\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"446\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>",
		},
		{
			name: "title",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Title:  "Title",
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"285\" y=\"298\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><path  d=\"M 100 260\nL 100 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 200 260\nL 200 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 300 260\nL 300 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 400 260\nL 400 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 500 260\nL 500 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 100 255\nL 500 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"146\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"246\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"346\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"446\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>",
		},
		{
			name: "boundary_gap_disabled",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:      []string{"a", "b", "c", "d"},
					BoundaryGap: Ptr(false),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 100 275\nL 100 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 233 275\nL 233 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 366 275\nL 366 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 500 275\nL 500 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 100 270\nL 500 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"99\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"232\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"365\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"491\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 100 275\nL 100 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 233 275\nL 233 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 366 275\nL 366 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 500 275\nL 500 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 100 270\nL 500 270\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"154\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,154,280)\">abc</text><text x=\"287\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,287,280)\">def</text><text x=\"421\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,421,280)\">ghi</text></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)), PainterPaddingOption(NewBoxEqual(100)))

			_, err := newBottomXAxis(p, tt.makeOption()).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYAxis(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		makeOption func() *YAxisOption
		result     string
	}{
		{
			name: "basic",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"99\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"99\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"99\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"99\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 114 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 114 166\nL 500 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 114 233\nL 500 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "font_style_with_rotation",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Labels:        []string{"a", "b", "c"},
					FontStyle:     NewFontStyleWithSize(20),
					LabelRotation: DegreesToRadians(270),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"99\" y=\"104\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,99,104)\">a</text><text x=\"99\" y=\"203\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,99,203)\">b</text><text x=\"99\" y=\"302\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,99,302)\">c</text><path  d=\"M 129 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 129 200\nL 500 200\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "lines",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Labels:        []string{"a", "b", "c", "d"},
					SplitLineShow: Ptr(true),
					SpineLineShow: Ptr(true),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 118 100\nL 118 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 113 100\nL 118 100\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 113 166\nL 118 166\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 113 233\nL 118 233\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 113 300\nL 118 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"99\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"99\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"99\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"99\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 114 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 114 166\nL 500 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 114 233\nL 500 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "title_left",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Title:  "title",
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"114\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,114,213)\">title</text><text x=\"119\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"119\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"119\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"119\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 134 100\nL 500 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 134 166\nL 500 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 134 233\nL 500 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "title_right",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Title:    "title",
					Position: PositionRight,
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"486\" y=\"187\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,486,187)\">title</text><text x=\"472\" y=\"106\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"472\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"472\" y=\"238\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"472\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 100 100\nL 462 100\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 100 166\nL 462 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 100 233\nL 462 233\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)),
				PainterPaddingOption(NewBoxEqual(100)))

			_, err := newAxisPainter(p, tt.makeOption().toAxisOption(p.theme)).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

func TestYAxisSplitLineDisabled(t *testing.T) {
	t.Parallel()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	}, PainterPaddingOption(NewBoxEqual(10)))
	yaxisOpt := &YAxisOption{
		Position: PositionRight,
		Labels:   []string{"a", "b", "c", "d"},
	}

	opt := yaxisOpt.toAxisOption(GetTheme(ThemeLight))
	opt.splitLineShow = false
	_, err := newAxisPainter(p, opt).Render()
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"582\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"582\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"582\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"582\" y=\"394\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>", data)
}

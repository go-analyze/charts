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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"10\" y=\"17\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"10\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"10\" y=\"250\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"10\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 29 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 29 126\nL 590 126\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 29 243\nL 590 243\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "font_style",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Labels:    []string{"a", "b", "c"},
					FontStyle: NewFontStyleWithSize(20),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"11\" y=\"22\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"10\" y=\"197\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"11\" y=\"372\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\">c</text><path  d=\"M 35 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 35 185\nL 590 185\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 24 10\nL 29 10\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 24 126\nL 29 126\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 24 243\nL 29 243\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 24 360\nL 29 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 29 10\nL 29 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"10\" y=\"17\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"10\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"10\" y=\"250\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"10\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 29 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 29 126\nL 590 126\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 29 243\nL 590 243\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)),
				PainterPaddingOption(NewBoxEqual(10)), PainterPaddingOption(Box{Bottom: defaultXAxisHeight}))

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
	}, PainterPaddingOption(NewBoxEqual(10)), PainterPaddingOption(Box{Bottom: defaultXAxisHeight}))
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"581\" y=\"17\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"581\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"581\" y=\"250\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"581\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>", data)
}

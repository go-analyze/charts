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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"9\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"9\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"9\" y=\"394\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 24 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 24 136\nL 590 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 24 263\nL 590 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "font_style",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Labels:    []string{"a", "b", "c"},
					FontStyle: NewFontStyleWithSize(20),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"10\" y=\"21\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"9\" y=\"210\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"10\" y=\"399\" style=\"stroke:none;fill:rgb(70,70,70);font-size:25.6px;font-family:'Roboto Medium',sans-serif\">c</text><path  d=\"M 30 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 30 200\nL 590 200\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 28 10\nL 28 390\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 23 10\nL 28 10\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 23 136\nL 28 136\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 23 263\nL 28 263\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 23 390\nL 28 390\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"9\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"9\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"9\" y=\"394\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 24 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 24 136\nL 590 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 24 263\nL 590 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
		{
			name: "title_left",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Title:  "title",
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"22\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(270.00,22,213)\">title</text><text x=\"25\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"25\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"25\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"25\" y=\"394\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 40 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 40 136\nL 590 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 40 263\nL 590 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"578\" y=\"187\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,578,187)\">title</text><text x=\"562\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"562\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"562\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"562\" y=\"394\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text><path  d=\"M 10 10\nL 552 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 10 136\nL 552 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 10 263\nL 552 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)),
				PainterPaddingOption(NewBoxEqual(10)))

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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"578\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">a</text><text x=\"578\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">b</text><text x=\"578\" y=\"268\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">c</text><text x=\"578\" y=\"394\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">d</text></svg>", data)
}

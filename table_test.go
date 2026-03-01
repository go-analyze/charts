package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeDefaultTableChartOptions() TableChartOption {
	return TableChartOption{
		Padding: NewBoxEqual(10),
		Header: []string{
			"Name", "Age", "Address", "Tag", "Action",
		},
		Data: [][]string{
			{
				"John Brown",
				"32",
				"New York No. 1 Lake Park",
				"nice, developer",
				"Send Mail",
			},
			{
				"Jim Green	",
				"42",
				"London No. 1 Lake Park",
				"wow",
				"Send Mail",
			},
			{
				"Joe Black	",
				"32",
				"Sidney No. 1 Lake Park",
				"cool, teacher",
				"Send Mail",
			},
		},
	}
}

func TestTableChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		theme         ColorPalette
		makeOptions   func() TableChartOption
		errorExpected bool
	}{
		{
			name:  "default_theme_spans",
			theme: nil, // use default
			makeOptions: func() TableChartOption {
				opt := makeDefaultTableChartOptions()
				opt.Spans = []int{
					1,
					1,
					2,
					1,
					// span and header do not match, and are automatically set to 1 at the end
					// 1,
				}
				return opt
			},
		},
		{
			name:        "dark_theme",
			theme:       GetTheme(ThemeVividDark),
			makeOptions: makeDefaultTableChartOptions,
		},
		{
			name: "cell_modified",
			makeOptions: func() TableChartOption {
				opt := makeDefaultTableChartOptions()
				opt.CellModifier = func(tc TableCell) TableCell {
					if tc.Column%2 == 0 {
						tc.FillColor = ColorWhite
						tc.FontStyle.FontColor = ColorBlue
					} else if tc.Row%2 == 0 {
						tc.FillColor = ColorAqua
						tc.FontStyle.FontColor = ColorBlack
					} else {
						tc.FillColor = ColorYellow
						tc.FontStyle.FontColor = ColorPurple
					}
					return tc
				}
				return opt
			},
		},
		{
			name: "error_no_header",
			makeOptions: func() TableChartOption {
				opt := makeDefaultTableChartOptions()
				opt.Header = nil
				return opt
			},
			errorExpected: true,
		},
		{
			name: "oversized_data_row",
			makeOptions: func() TableChartOption {
				return TableChartOption{
					Padding: NewBoxEqual(10),
					Header:  []string{"A", "B"},
					Data: [][]string{
						{"1", "2", "extra1", "extra2"},
					},
				}
			},
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		runName := strconv.Itoa(i) + "-" + tt.name
		if tt.theme != nil {
			t.Run(runName+"-theme_painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(tt.theme))
				opt := tt.makeOptions()

				validateTableChartRender(t, p, opt, tt.errorExpected)
			})
			runName += "-theme_opt"
		}
		t.Run(runName, func(t *testing.T) {
			p := NewPainter(painterOptions)
			opt := tt.makeOptions()
			opt.Theme = tt.theme

			validateTableChartRender(t, p, opt, tt.errorExpected)
		})
	}
}

func validateTableChartRender(t *testing.T, p *Painter, opt TableChartOption, errorExpected bool) {
	t.Helper()

	err := p.TableChart(opt)
	if errorExpected {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

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
		result        string
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 36\nL 0 36\nL 0 0\" style=\"stroke:none;fill:rgb(220,220,220)\"/><path d=\"M 0 36\nL 600 36\nL 600 93\nL 0 93\nL 0 36\" style=\"stroke:none;fill:white\"/><path d=\"M 0 93\nL 600 93\nL 600 129\nL 0 129\nL 0 93\" style=\"stroke:none;fill:rgb(245,245,245)\"/><path d=\"M 0 129\nL 600 129\nL 600 186\nL 0 186\nL 0 129\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"22\" style=\"stroke:none;fill:rgb(80,80,80);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"110\" y=\"22\" style=\"stroke:none;fill:rgb(80,80,80);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"210\" y=\"22\" style=\"stroke:none;fill:rgb(80,80,80);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"410\" y=\"22\" style=\"stroke:none;fill:rgb(80,80,80);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"510\" y=\"22\" style=\"stroke:none;fill:rgb(80,80,80);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"58\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John</text><text x=\"10\" y=\"79\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brown</text><text x=\"110\" y=\"58\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"210\" y=\"58\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No. 1 Lake Park</text><text x=\"410\" y=\"58\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"410\" y=\"79\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"510\" y=\"58\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"115\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"110\" y=\"115\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"210\" y=\"115\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1 Lake Park</text><text x=\"410\" y=\"115\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"510\" y=\"115\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"151\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"110\" y=\"151\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"210\" y=\"151\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1 Lake Park</text><text x=\"410\" y=\"151\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool,</text><text x=\"410\" y=\"172\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">teacher</text><text x=\"510\" y=\"151\" style=\"stroke:none;fill:rgb(50,50,50);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
		},
		{
			name:        "dark_theme",
			theme:       GetTheme(ThemeVividDark),
			makeOptions: makeDefaultTableChartOptions,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 36\nL 0 36\nL 0 0\" style=\"stroke:none;fill:rgb(38,38,42)\"/><path d=\"M 0 36\nL 600 36\nL 600 93\nL 0 93\nL 0 36\" style=\"stroke:none;fill:rgb(24,24,28)\"/><path d=\"M 0 93\nL 600 93\nL 600 150\nL 0 150\nL 0 93\" style=\"stroke:none;fill:rgb(38,38,42)\"/><path d=\"M 0 150\nL 600 150\nL 600 207\nL 0 207\nL 0 150\" style=\"stroke:none;fill:rgb(24,24,28)\"/><text x=\"10\" y=\"22\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"130\" y=\"22\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"250\" y=\"22\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"370\" y=\"22\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"490\" y=\"22\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"58\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John Brown</text><text x=\"130\" y=\"58\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"58\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No.</text><text x=\"250\" y=\"79\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1 Lake Park</text><text x=\"370\" y=\"58\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"370\" y=\"79\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"490\" y=\"58\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"115\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"130\" y=\"115\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"250\" y=\"115\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1</text><text x=\"250\" y=\"136\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"115\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"490\" y=\"115\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"172\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"130\" y=\"172\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"172\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1</text><text x=\"250\" y=\"193\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"172\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool, teacher</text><text x=\"490\" y=\"172\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 36\nL 0 36\nL 0 0\" style=\"stroke:none;fill:rgb(220,220,220)\"/><path d=\"M 0 36\nL 600 36\nL 600 93\nL 0 93\nL 0 36\" style=\"stroke:none;fill:white\"/><path d=\"M 0 93\nL 600 93\nL 600 150\nL 0 150\nL 0 93\" style=\"stroke:none;fill:rgb(245,245,245)\"/><path d=\"M 0 150\nL 600 150\nL 600 207\nL 0 207\nL 0 150\" style=\"stroke:none;fill:white\"/><path d=\"M 10 10\nL 110 10\nL 110 26\nL 10 26\nL 10 10\" style=\"stroke:none;fill:white\"/><path d=\"M 130 10\nL 230 10\nL 230 26\nL 130 26\nL 130 10\" style=\"stroke:none;fill:aqua\"/><path d=\"M 250 10\nL 350 10\nL 350 26\nL 250 26\nL 250 10\" style=\"stroke:none;fill:white\"/><path d=\"M 370 10\nL 470 10\nL 470 26\nL 370 26\nL 370 10\" style=\"stroke:none;fill:aqua\"/><path d=\"M 490 10\nL 590 10\nL 590 26\nL 490 26\nL 490 10\" style=\"stroke:none;fill:white\"/><path d=\"M 10 46\nL 110 46\nL 110 83\nL 10 83\nL 10 46\" style=\"stroke:none;fill:white\"/><path d=\"M 130 46\nL 230 46\nL 230 83\nL 130 83\nL 130 46\" style=\"stroke:none;fill:yellow\"/><path d=\"M 250 46\nL 350 46\nL 350 83\nL 250 83\nL 250 46\" style=\"stroke:none;fill:white\"/><path d=\"M 370 46\nL 470 46\nL 470 83\nL 370 83\nL 370 46\" style=\"stroke:none;fill:yellow\"/><path d=\"M 490 46\nL 590 46\nL 590 83\nL 490 83\nL 490 46\" style=\"stroke:none;fill:white\"/><path d=\"M 10 103\nL 110 103\nL 110 140\nL 10 140\nL 10 103\" style=\"stroke:none;fill:white\"/><path d=\"M 130 103\nL 230 103\nL 230 140\nL 130 140\nL 130 103\" style=\"stroke:none;fill:aqua\"/><path d=\"M 250 103\nL 350 103\nL 350 140\nL 250 140\nL 250 103\" style=\"stroke:none;fill:white\"/><path d=\"M 370 103\nL 470 103\nL 470 140\nL 370 140\nL 370 103\" style=\"stroke:none;fill:aqua\"/><path d=\"M 490 103\nL 590 103\nL 590 140\nL 490 140\nL 490 103\" style=\"stroke:none;fill:white\"/><path d=\"M 10 160\nL 110 160\nL 110 197\nL 10 197\nL 10 160\" style=\"stroke:none;fill:white\"/><path d=\"M 130 160\nL 230 160\nL 230 197\nL 130 197\nL 130 160\" style=\"stroke:none;fill:yellow\"/><path d=\"M 250 160\nL 350 160\nL 350 197\nL 250 197\nL 250 160\" style=\"stroke:none;fill:white\"/><path d=\"M 370 160\nL 470 160\nL 470 197\nL 370 197\nL 370 160\" style=\"stroke:none;fill:yellow\"/><path d=\"M 490 160\nL 590 160\nL 590 197\nL 490 197\nL 490 160\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"22\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"130\" y=\"22\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"250\" y=\"22\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"370\" y=\"22\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"490\" y=\"22\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"58\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John Brown</text><text x=\"130\" y=\"58\" style=\"stroke:none;fill:purple;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"58\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No.</text><text x=\"250\" y=\"79\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1 Lake Park</text><text x=\"370\" y=\"58\" style=\"stroke:none;fill:purple;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"370\" y=\"79\" style=\"stroke:none;fill:purple;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"490\" y=\"58\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"115\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"130\" y=\"115\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"250\" y=\"115\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1</text><text x=\"250\" y=\"136\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"115\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"490\" y=\"115\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"172\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"130\" y=\"172\" style=\"stroke:none;fill:purple;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"172\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1</text><text x=\"250\" y=\"193\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"172\" style=\"stroke:none;fill:purple;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool, teacher</text><text x=\"490\" y=\"172\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
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

				validateTableChartRender(t, p, opt, tt.result, tt.errorExpected)
			})
			runName += "-theme_opt"
		}
		t.Run(runName, func(t *testing.T) {
			p := NewPainter(painterOptions)
			opt := tt.makeOptions()
			opt.Theme = tt.theme

			validateTableChartRender(t, p, opt, tt.result, tt.errorExpected)
		})
	}
}

func validateTableChartRender(t *testing.T, p *Painter, opt TableChartOption,
	expectedResult string, errorExpected bool) {
	t.Helper()

	err := p.TableChart(opt)
	if errorExpected {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
}

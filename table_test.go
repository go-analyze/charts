package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func makeDefaultTableChartOptions() TableChartOption {
	return TableChartOption{
		Padding: Box{
			Left:   10,
			Top:    10,
			Right:  10,
			Bottom: 10,
		},
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 35\nL 0 35\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(220,220,220,1.0)\"/><path  d=\"M 0 35\nL 600 35\nL 600 90\nL 0 90\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 0 90\nL 600 90\nL 600 125\nL 0 125\nL 0 90\" style=\"stroke-width:0;stroke:none;fill:rgba(245,245,245,1.0)\"/><path  d=\"M 0 125\nL 600 125\nL 600 180\nL 0 180\nL 0 125\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(80,80,80,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"110\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(80,80,80,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"210\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(80,80,80,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"410\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(80,80,80,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"510\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(80,80,80,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John</text><text x=\"10\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brown</text><text x=\"110\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"210\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No. 1 Lake Park</text><text x=\"410\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"410\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"510\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"110\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"210\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1 Lake Park</text><text x=\"410\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"510\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"110\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"210\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1 Lake Park</text><text x=\"410\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool,</text><text x=\"410\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">teacher</text><text x=\"510\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(50,50,50,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
		},
		{
			name:        "dark_theme",
			theme:       GetTheme(ThemeVividDark),
			makeOptions: makeDefaultTableChartOptions,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 35\nL 0 35\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(38,38,42,1.0)\"/><path  d=\"M 0 35\nL 600 35\nL 600 90\nL 0 90\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(24,24,28,1.0)\"/><path  d=\"M 0 90\nL 600 90\nL 600 145\nL 0 145\nL 0 90\" style=\"stroke-width:0;stroke:none;fill:rgba(38,38,42,1.0)\"/><path  d=\"M 0 145\nL 600 145\nL 600 200\nL 0 200\nL 0 145\" style=\"stroke-width:0;stroke:none;fill:rgba(24,24,28,1.0)\"/><text x=\"10\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"130\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"250\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"370\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"490\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John Brown</text><text x=\"130\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No.</text><text x=\"250\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1 Lake Park</text><text x=\"370\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"370\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"490\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"130\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"250\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1</text><text x=\"250\" y=\"132\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"490\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"130\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1</text><text x=\"250\" y=\"187\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool, teacher</text><text x=\"490\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(216,217,218,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
		},
		{
			name: "cell_modified",
			makeOptions: func() TableChartOption {
				opt := makeDefaultTableChartOptions()
				opt.CellModifier = func(tc TableCell) TableCell {
					if tc.Column%2 == 0 {
						tc.FillColor = drawing.ColorWhite
						tc.FontStyle.FontColor = drawing.ColorBlue
					} else if tc.Row%2 == 0 {
						tc.FillColor = drawing.ColorAqua
						tc.FontStyle.FontColor = drawing.ColorBlack
					} else {
						tc.FillColor = drawing.ColorYellow
						tc.FontStyle.FontColor = drawing.ColorPurple
					}
					return tc
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 35\nL 0 35\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(220,220,220,1.0)\"/><path  d=\"M 0 35\nL 600 35\nL 600 90\nL 0 90\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 0 90\nL 600 90\nL 600 145\nL 0 145\nL 0 90\" style=\"stroke-width:0;stroke:none;fill:rgba(245,245,245,1.0)\"/><path  d=\"M 0 145\nL 600 145\nL 600 200\nL 0 200\nL 0 145\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 10 10\nL 110 10\nL 110 25\nL 10 25\nL 10 10\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 130 10\nL 230 10\nL 230 25\nL 130 25\nL 130 10\" style=\"stroke-width:0;stroke:none;fill:rgba(0,255,255,1.0)\"/><path  d=\"M 250 10\nL 350 10\nL 350 25\nL 250 25\nL 250 10\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 370 10\nL 470 10\nL 470 25\nL 370 25\nL 370 10\" style=\"stroke-width:0;stroke:none;fill:rgba(0,255,255,1.0)\"/><path  d=\"M 490 10\nL 590 10\nL 590 25\nL 490 25\nL 490 10\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 10 45\nL 110 45\nL 110 80\nL 10 80\nL 10 45\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 130 45\nL 230 45\nL 230 80\nL 130 80\nL 130 45\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,0,1.0)\"/><path  d=\"M 250 45\nL 350 45\nL 350 80\nL 250 80\nL 250 45\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 370 45\nL 470 45\nL 470 80\nL 370 80\nL 370 45\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,0,1.0)\"/><path  d=\"M 490 45\nL 590 45\nL 590 80\nL 490 80\nL 490 45\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 10 100\nL 110 100\nL 110 135\nL 10 135\nL 10 100\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 130 100\nL 230 100\nL 230 135\nL 130 135\nL 130 100\" style=\"stroke-width:0;stroke:none;fill:rgba(0,255,255,1.0)\"/><path  d=\"M 250 100\nL 350 100\nL 350 135\nL 250 135\nL 250 100\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 370 100\nL 470 100\nL 470 135\nL 370 135\nL 370 100\" style=\"stroke-width:0;stroke:none;fill:rgba(0,255,255,1.0)\"/><path  d=\"M 490 100\nL 590 100\nL 590 135\nL 490 135\nL 490 100\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 10 155\nL 110 155\nL 110 190\nL 10 190\nL 10 155\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 130 155\nL 230 155\nL 230 190\nL 130 190\nL 130 155\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,0,1.0)\"/><path  d=\"M 250 155\nL 350 155\nL 350 190\nL 250 190\nL 250 155\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 370 155\nL 470 155\nL 470 190\nL 370 190\nL 370 155\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,0,1.0)\"/><path  d=\"M 490 155\nL 590 155\nL 590 190\nL 490 190\nL 490 155\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"130\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"250\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"370\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"490\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John Brown</text><text x=\"130\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(128,0,128,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No.</text><text x=\"250\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1 Lake Park</text><text x=\"370\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(128,0,128,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"370\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(128,0,128,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"490\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"130\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"250\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1</text><text x=\"250\" y=\"132\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,0,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"490\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"130\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(128,0,128,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1</text><text x=\"250\" y=\"187\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(128,0,128,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool, teacher</text><text x=\"490\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
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
		runName := strconv.Itoa(i)
		if tt.theme != nil {
			t.Run(runName+"-theme_painter", func(t *testing.T) {
				p, err := NewPainter(painterOptions, PainterThemeOption(tt.theme))
				require.NoError(t, err)
				opt := tt.makeOptions()

				validateTableChartRender(t, p, opt, tt.result, tt.errorExpected)
			})
			runName += "-theme_opt"
		}
		t.Run(runName, func(t *testing.T) {
			p, err := NewPainter(painterOptions)
			require.NoError(t, err)
			opt := tt.makeOptions()
			opt.Theme = tt.theme

			validateTableChartRender(t, p, opt, tt.result, tt.errorExpected)
		})
	}
}

func validateTableChartRender(t *testing.T, p *Painter, opt TableChartOption,
	expectedResult string, errorExpected bool) {
	t.Helper()

	_, err := NewTableChart(p, opt).Render()
	if errorExpected {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, string(data))
}

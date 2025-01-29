package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/roboto"
)

func TestInstallGetFont(t *testing.T) {
	t.Parallel()

	fontFamily := "test"
	err := InstallFont(fontFamily, roboto.Roboto)
	require.NoError(t, err)

	font := GetFont(fontFamily)
	assert.NotNil(t, font)
}

func TestGetPreferredFont(t *testing.T) {
	t.Parallel()

	t.Run("nill-default", func(t *testing.T) {
		require.Equal(t, GetDefaultFont(), getPreferredFont(nil))
	})
}

func TestCustomFontSizeRender(t *testing.T) {
	t.Parallel()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	}, PainterThemeOption(GetTheme(ThemeLight)))

	opt := makeBasicLineChartOption()
	opt.XAxis.FontStyle.FontSize = 4.0
	opt.YAxis = []YAxisOption{
		{
			FontStyle: NewFontStyleWithSize(4.0),
		},
	}
	opt.Title.FontStyle.FontSize = 4.0
	opt.Legend.FontStyle.FontSize = 4.0

	err := p.LineChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 256 19\nL 286 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"271\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"288\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"10\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"10\" y=\"37\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"73\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"109\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"14\" y=\"145\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"14\" y=\"181\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"14\" y=\"217\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"14\" y=\"253\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"14\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"14\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"20\" y=\"362\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 33 35\nL 590 35\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 71\nL 590 71\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 107\nL 590 107\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 143\nL 590 143\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 179\nL 590 179\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 215\nL 590 215\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 251\nL 590 251\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 287\nL 590 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 323\nL 590 323\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 33 365\nL 33 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 112 365\nL 112 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 192 365\nL 192 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 271 365\nL 271 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 351 365\nL 351 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 430 365\nL 430 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 510 365\nL 510 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"70\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"150\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"229\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"309\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"389\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"469\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"548\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">G</text><path  d=\"M 72 333\nL 152 331\nL 231 338\nL 311 330\nL 390 340\nL 470 309\nL 550 313\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"72\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"152\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"231\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"311\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"390\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"470\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"550\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path  d=\"M 72 175\nL 152 150\nL 231 157\nL 311 150\nL 390 69\nL 470 60\nL 550 63\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"72\" cy=\"175\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"152\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"231\" cy=\"157\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"311\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"390\" cy=\"69\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"470\" cy=\"60\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"550\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>", data)
}

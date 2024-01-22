package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/roboto"
)

func TestInstallGetFont(t *testing.T) {
	fontFamily := "test"
	err := InstallFont(fontFamily, roboto.Roboto)
	require.NoError(t, err)

	font := GetFont(fontFamily)
	assert.NotNil(t, font)
}

func TestGetPreferredFont(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		require.Equal(t, getPreferredFont(nil), GetDefaultFont())
	})
}

func TestCustomFontSizeRender(t *testing.T) {
	t.Parallel()

	p, err := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	}, PainterThemeOption(GetTheme(ThemeLight)))
	require.NoError(t, err)

	opt := makeBasicLineChartOption()
	opt.XAxis.FontSize = 4.0
	opt.YAxis = []YAxisOption{
		{
			FontSize: 4.0,
		},
	}
	opt.Title.FontSize = 4.0
	opt.Legend.FontSize = 4.0

	_, err = NewLineChart(p, opt).Render()
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 256 19\nL 286 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"271\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"288\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"343\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"10\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"10\" y=\"47\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"82\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"117\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"14\" y=\"152\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"14\" y=\"187\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"14\" y=\"222\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"14\" y=\"257\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"14\" y=\"292\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"14\" y=\"327\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"20\" y=\"362\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 33 45\nL 590 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 80\nL 590 80\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 115\nL 590 115\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 150\nL 590 150\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 220\nL 590 220\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 255\nL 590 255\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 290\nL 590 290\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 325\nL 590 325\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 33 365\nL 33 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 112 365\nL 112 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 192 365\nL 192 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 271 365\nL 271 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 351 365\nL 351 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 430 365\nL 430 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 510 365\nL 510 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 33 360\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"70\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"150\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"229\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"309\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"389\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"469\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"548\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">G</text><path  d=\"M 72 334\nL 152 332\nL 231 338\nL 311 331\nL 390 341\nL 470 310\nL 550 315\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"72\" cy=\"334\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"152\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"231\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"311\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"390\" cy=\"341\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"550\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 72 181\nL 152 157\nL 231 163\nL 311 156\nL 390 78\nL 470 70\nL 550 72\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"72\" cy=\"181\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"152\" cy=\"157\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"231\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"311\" cy=\"156\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"390\" cy=\"78\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"550\" cy=\"72\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>", string(data))
}

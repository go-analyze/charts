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

	t.Run("nil_default", func(t *testing.T) {
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
	opt.Legend.Symbol = SymbolDot

	err := p.LineChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 256 19\nL 286 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"271\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"288\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"10\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"9\" y=\"37\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"9\" y=\"73\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"9\" y=\"109\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"13\" y=\"146\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"13\" y=\"182\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"13\" y=\"218\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"13\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"13\" y=\"291\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"13\" y=\"327\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"19\" y=\"364\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 28 36\nL 590 36\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 72\nL 590 72\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 109\nL 590 109\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 145\nL 590 145\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 182\nL 590 182\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 218\nL 590 218\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 255\nL 590 255\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 291\nL 590 291\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 28 328\nL 590 328\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 32 365\nL 590 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 32 370\nL 32 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 111 370\nL 111 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 191 370\nL 191 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 271 370\nL 271 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 350 370\nL 350 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 430 370\nL 430 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 510 370\nL 510 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 370\nL 590 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"69\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"149\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"229\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"308\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"389\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"469\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"548\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">G</text><path  d=\"M 71 338\nL 151 335\nL 231 342\nL 310 335\nL 390 345\nL 470 313\nL 550 318\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"71\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"151\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"231\" cy=\"342\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"310\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"390\" cy=\"345\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"470\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"550\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path  d=\"M 71 178\nL 151 153\nL 231 160\nL 310 152\nL 390 71\nL 470 62\nL 550 64\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"71\" cy=\"178\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"151\" cy=\"153\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"231\" cy=\"160\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"310\" cy=\"152\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"390\" cy=\"71\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"470\" cy=\"62\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"550\" cy=\"64\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>", data)
}

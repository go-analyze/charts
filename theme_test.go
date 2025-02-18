package charts

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallGetTheme(t *testing.T) {
	t.Parallel()

	name := "TestInstallGetTheme"
	InstallTheme(name, ThemeOption{IsDarkMode: true})
	getThemeResult := GetTheme(name)
	assert.NotNil(t, getThemeResult)
	assert.NotEqual(t, GetDefaultTheme(), getThemeResult)
	assert.Equal(t, name, fmt.Sprintf("%s", getThemeResult))
}

func TestGetPreferredTheme(t *testing.T) {
	t.Parallel()

	t.Run("default", func(t *testing.T) {
		assert.Equal(t, getPreferredTheme(nil), GetDefaultTheme())
	})
	t.Run("provided", func(t *testing.T) {
		theme := GetTheme(ThemeGrafana)
		assert.Equal(t, getPreferredTheme(theme), theme)
	})
}

func TestDefaultTheme(t *testing.T) {
	t.Parallel()

	assert.Equal(t, GetTheme(ThemeLight), GetDefaultTheme())
	assert.Equal(t, GetTheme("Unknown Theme"), GetDefaultTheme())
}

func TestSetDefaultThemeError(t *testing.T) {
	t.Parallel()

	assert.Error(t, SetDefaultFont("not a theme"))
}

func renderTestLineChartWithThemeName(t *testing.T, fullChart bool, themeName string) []byte {
	t.Helper()

	return renderTestLineChartWithTheme(t, fullChart, GetTheme(themeName))
}

func renderTestLineChartWithTheme(t *testing.T, fullChart bool, theme ColorPalette) []byte {
	t.Helper()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})
	opt := makeFullLineChartOption()
	if len(opt.YAxis) == 0 {
		opt.YAxis = []YAxisOption{{}}
	}
	opt.YAxis[0].LabelCount = 5
	opt.YAxis[0].LabelSkipCount = 1
	if !fullChart {
		opt.Title.Show = Ptr(false)
		opt.Symbol = SymbolNone
	}
	opt.Theme = theme

	err := p.LineChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	return data
}

func TestThemeLight(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeLight)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"96\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"172\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"324\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"400\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"476\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"96\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"172\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"248\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"324\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"400\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"96\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"172\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"248\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"324\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"400\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"96\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"172\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"248\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"324\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"400\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"96\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"172\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"248\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"324\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"400\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"476\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/></svg>", svg)
}

func TestThemeDark(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeDark)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"96\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(40,40,40)\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"96\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(40,40,40)\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"96\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(40,40,40)\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"96\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(40,40,40)\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"96\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(40,40,40)\"/></svg>", svg)
}

func TestThemeVividLight(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeVividLight)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"96\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"172\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"324\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"400\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"476\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:white\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"96\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"172\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"248\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"324\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"400\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:white\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"96\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><circle cx=\"172\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><circle cx=\"248\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><circle cx=\"324\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><circle cx=\"400\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:white\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"96\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><circle cx=\"172\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><circle cx=\"248\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><circle cx=\"324\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><circle cx=\"400\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:white\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(154,100,180);fill:none\"/><circle cx=\"96\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/><circle cx=\"172\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/><circle cx=\"248\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/><circle cx=\"324\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/><circle cx=\"400\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/><circle cx=\"476\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:white\"/></svg>", svg)
}

func TestThemeVividDark(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeVividDark)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"96\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(40,40,40)\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"96\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(40,40,40)\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"96\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(40,40,40)\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"96\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(40,40,40)\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(154,100,180);fill:none\"/><circle cx=\"96\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/><circle cx=\"172\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/><circle cx=\"248\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/><circle cx=\"324\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/><circle cx=\"400\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/><circle cx=\"476\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(40,40,40)\"/></svg>", svg)
}

func TestThemeAnt(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeAnt)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(91,143,249);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(91,143,249);fill:rgb(91,143,249)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(90,216,166);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(90,216,166);fill:rgb(90,216,166)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(93,112,146);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(93,112,146);fill:rgb(93,112,146)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(246,189,22);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(246,189,22);fill:rgb(246,189,22)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(111,94,249);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(111,94,249);fill:rgb(111,94,249)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(91,143,249);fill:none\"/><circle cx=\"96\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><circle cx=\"172\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><circle cx=\"324\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><circle cx=\"400\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><circle cx=\"476\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(91,143,249);fill:white\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(90,216,166);fill:none\"/><circle cx=\"96\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><circle cx=\"172\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><circle cx=\"248\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><circle cx=\"324\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><circle cx=\"400\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(90,216,166);fill:white\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(93,112,146);fill:none\"/><circle cx=\"96\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><circle cx=\"172\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><circle cx=\"248\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><circle cx=\"324\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><circle cx=\"400\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(93,112,146);fill:white\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(246,189,22);fill:none\"/><circle cx=\"96\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><circle cx=\"172\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><circle cx=\"248\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><circle cx=\"324\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><circle cx=\"400\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(246,189,22);fill:white\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(111,94,249);fill:none\"/><circle cx=\"96\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/><circle cx=\"172\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/><circle cx=\"248\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/><circle cx=\"324\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/><circle cx=\"400\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/><circle cx=\"476\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(111,94,249);fill:white\"/></svg>", svg)
}

func TestThemeGrafana(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithThemeName(t, true, ThemeGrafana)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(31,29,29)\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(126,178,109);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(126,178,109);fill:rgb(126,178,109)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(234,184,57);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(234,184,57);fill:rgb(234,184,57)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(110,208,224);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(110,208,224);fill:rgb(110,208,224)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(239,132,60);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(239,132,60);fill:rgb(239,132,60)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(226,77,66);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(226,77,66);fill:rgb(226,77,66)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(68,67,67);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(68,67,67);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(68,67,67);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(68,67,67);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(216,217,218);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(126,178,109);fill:none\"/><circle cx=\"96\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><circle cx=\"172\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><circle cx=\"324\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><circle cx=\"400\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><circle cx=\"476\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(126,178,109);fill:rgb(31,29,29)\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(234,184,57);fill:none\"/><circle cx=\"96\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><circle cx=\"172\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><circle cx=\"248\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><circle cx=\"324\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><circle cx=\"400\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(234,184,57);fill:rgb(31,29,29)\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(110,208,224);fill:none\"/><circle cx=\"96\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><circle cx=\"172\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><circle cx=\"248\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><circle cx=\"324\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><circle cx=\"400\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(110,208,224);fill:rgb(31,29,29)\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(239,132,60);fill:none\"/><circle cx=\"96\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><circle cx=\"172\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><circle cx=\"248\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><circle cx=\"324\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><circle cx=\"400\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><circle cx=\"476\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(239,132,60);fill:rgb(31,29,29)\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(226,77,66);fill:none\"/><circle cx=\"96\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/><circle cx=\"172\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/><circle cx=\"248\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/><circle cx=\"324\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/><circle cx=\"400\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/><circle cx=\"476\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgb(226,77,66);fill:rgb(31,29,29)\"/></svg>", svg)
}

func TestLightThemeSeriesRepeat(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithTheme(t, false,
		&themeColorPalette{
			name:               t.Name(),
			isDarkMode:         false,
			xaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			yaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			axisSplitLineColor: Color{R: 224, G: 230, B: 242, A: 255},
			backgroundColor:    ColorWhite,
			legendTextColor:    Color{R: 70, G: 70, B: 70, A: 255},
			seriesColors: []Color{
				{R: 50, G: 50, B: 50, A: 255},
				{R: 200, G: 50, B: 50, A: 255},
			},
		})
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(50,50,50);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(50,50,50);fill:rgb(50,50,50)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(200,50,50);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(200,50,50);fill:rgb(200,50,50)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(95,95,95);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(95,95,95);fill:rgb(95,95,95)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(200,90,90);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(200,90,90);fill:rgb(200,90,90)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(141,141,141);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(141,141,141);fill:rgb(141,141,141)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(50,50,50);fill:none\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(200,50,50);fill:none\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(95,95,95);fill:none\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(200,90,90);fill:none\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(141,141,141);fill:none\"/></svg>", svg)
}

func TestDarkThemeSeriesRepeat(t *testing.T) {
	t.Parallel()

	svg := renderTestLineChartWithTheme(t, false,
		&themeColorPalette{
			name:               t.Name(),
			isDarkMode:         true,
			xaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			yaxisStrokeColor:   Color{R: 110, G: 112, B: 121, A: 255},
			axisSplitLineColor: Color{R: 40, G: 40, B: 60, A: 255},
			backgroundColor:    Color{R: 40, G: 40, B: 60, A: 255},
			legendTextColor:    Color{R: 70, G: 70, B: 70, A: 255},
			seriesColors: []Color{
				{R: 250, G: 250, B: 250, A: 255},
				{R: 200, G: 50, B: 50, A: 255},
			},
		})
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,60)\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(250,250,250);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,250,250);fill:rgb(250,250,250)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(200,50,50);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(200,50,50);fill:rgb(200,50,50)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(204,204,204);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(204,204,204);fill:rgb(204,204,204)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(156,52,52);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(156,52,52);fill:rgb(156,52,52)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(158,158,158);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(158,158,158);fill:rgb(158,158,158)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"19\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"209\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(40,40,60);fill:none\"/><path  d=\"M 59 123\nL 590 123\" style=\"stroke-width:1;stroke:rgb(40,40,60);fill:none\"/><path  d=\"M 59 202\nL 590 202\" style=\"stroke-width:1;stroke:rgb(40,40,60);fill:none\"/><path  d=\"M 59 281\nL 590 281\" style=\"stroke-width:1;stroke:rgb(40,40,60);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 333\nL 172 331\nL 248 338\nL 324 330\nL 400 340\nL 476 309\nL 552 313\" style=\"stroke-width:2;stroke:rgb(250,250,250);fill:none\"/><path  d=\"M 96 311\nL 172 320\nL 248 318\nL 324 308\nL 400 295\nL 476 286\nL 552 291\" style=\"stroke-width:2;stroke:rgb(200,50,50);fill:none\"/><path  d=\"M 96 327\nL 172 308\nL 248 315\nL 324 326\nL 400 318\nL 476 286\nL 552 268\" style=\"stroke-width:2;stroke:rgb(204,204,204);fill:none\"/><path  d=\"M 96 288\nL 172 286\nL 248 293\nL 324 285\nL 400 273\nL 476 286\nL 552 288\" style=\"stroke-width:2;stroke:rgb(156,52,52);fill:none\"/><path  d=\"M 96 176\nL 172 151\nL 248 158\nL 324 150\nL 400 70\nL 476 61\nL 552 63\" style=\"stroke-width:2;stroke:rgb(158,158,158);fill:none\"/></svg>", svg)
}

func TestWithAxisColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithXAxisColor(ColorBlack).WithYAxisColor(ColorBlack)

	assert.Equal(t, ColorWhite, whiteCP.xaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.yaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.backgroundColor)
	assert.Equal(t, ColorWhite, whiteCP.seriesColors[0])

	assert.NotEqual(t, ColorBlack, blackCP.GetAxisSplitLineColor())
	assert.Equal(t, ColorBlack, blackCP.GetXAxisStrokeColor())
	assert.Equal(t, ColorBlack, blackCP.GetYAxisStrokeColor())
	assert.Equal(t, ColorWhite, blackCP.GetBackgroundColor())
	assert.Equal(t, ColorWhite, blackCP.GetSeriesColor(0))
}

func TestWithTextColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		titleTextColor:     ColorWhite,
		markTextColor:      ColorWhite,
		labelTextColor:     ColorWhite,
		legendTextColor:    ColorWhite,
		xaxisTextColor:     ColorWhite,
		yaxisTextColor:     ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithTextColor(ColorBlack)

	assert.Equal(t, ColorWhite, whiteCP.axisSplitLineColor)
	assert.Equal(t, ColorWhite, whiteCP.xaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.yaxisStrokeColor)
	assert.Equal(t, ColorWhite, whiteCP.backgroundColor)
	assert.Equal(t, ColorWhite, whiteCP.GetTitleTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetLabelTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetLegendTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetMarkTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetXAxisTextColor())
	assert.Equal(t, ColorWhite, whiteCP.GetYAxisTextColor())
	assert.Equal(t, ColorWhite, whiteCP.seriesColors[0])

	assert.Equal(t, ColorBlack, blackCP.GetTitleTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetLabelTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetLegendTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetMarkTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetXAxisTextColor())
	assert.Equal(t, ColorBlack, blackCP.GetYAxisTextColor())
	assert.Equal(t, ColorWhite, blackCP.GetAxisSplitLineColor())
	assert.Equal(t, ColorWhite, blackCP.GetXAxisStrokeColor())
	assert.Equal(t, ColorWhite, blackCP.GetYAxisStrokeColor())
	assert.Equal(t, ColorWhite, blackCP.GetBackgroundColor())
	assert.Equal(t, ColorWhite, blackCP.GetSeriesColor(0))

	redCP := blackCP.WithTitleTextColor(ColorRed).WithLabelTextColor(ColorRed).
		WithMarkTextColor(ColorRed).WithLegendTextColor(ColorRed).
		WithXAxisTextColor(ColorRed).WithYAxisTextColor(ColorRed)

	assert.Equal(t, ColorRed, redCP.GetTitleTextColor())
	assert.Equal(t, ColorRed, redCP.GetLabelTextColor())
	assert.Equal(t, ColorRed, redCP.GetLegendTextColor())
	assert.Equal(t, ColorRed, redCP.GetMarkTextColor())
	assert.Equal(t, ColorRed, redCP.GetXAxisTextColor())
	assert.Equal(t, ColorRed, redCP.GetYAxisTextColor())
}

func TestWithSeriesColors(t *testing.T) {
	t.Parallel()

	yellowCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorYellow,
		xaxisStrokeColor:   ColorYellow,
		yaxisStrokeColor:   ColorYellow,
		backgroundColor:    ColorYellow,
		seriesColors:       []Color{ColorYellow},
	}

	t.Run("ignored", func(t *testing.T) {
		invalidCP := yellowCP.WithSeriesColors([]Color{})

		assert.Equal(t, yellowCP.GetSeriesColor(0), invalidCP.GetSeriesColor(0))
	})
	t.Run("updated", func(t *testing.T) {
		blackCP := yellowCP.WithSeriesColors([]Color{ColorBlack})

		assert.Equal(t, ColorYellow, blackCP.GetAxisSplitLineColor())
		assert.Equal(t, ColorYellow, blackCP.GetXAxisStrokeColor())
		assert.Equal(t, ColorYellow, blackCP.GetYAxisStrokeColor())
		assert.Equal(t, ColorYellow, blackCP.GetBackgroundColor())
		assert.Equal(t, ColorBlack, blackCP.GetSeriesColor(0))
	})
}

func TestWithBackgroundColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithBackgroundColor(ColorBlack)

	require.False(t, whiteCP.IsDark())
	assert.True(t, blackCP.IsDark())

	yellowCP := blackCP.WithBackgroundColor(ColorYellow)

	require.False(t, yellowCP.IsDark())
}

func TestWithLegendBorderColor(t *testing.T) {
	t.Parallel()

	whiteCP := &themeColorPalette{
		name:               t.Name(),
		isDarkMode:         false,
		axisSplitLineColor: ColorWhite,
		xaxisStrokeColor:   ColorWhite,
		yaxisStrokeColor:   ColorWhite,
		backgroundColor:    ColorWhite,
		legendBorderColor:  ColorWhite,
		seriesColors:       []Color{ColorWhite},
	}

	blackCP := whiteCP.WithLegendBorderColor(ColorBlack)
	assert.Equal(t, ColorBlack, blackCP.GetLegendBorderColor())
	assert.Equal(t, ColorWhite, whiteCP.GetLegendBorderColor())
}

func getThemeSeriesColors(name string) []Color {
	if value, ok := palettes.Load(name); ok {
		if cp, ok := value.(*themeColorPalette); ok {
			return cp.seriesColors
		}
	}
	return nil
}

const makeSeriesLoopColorSamples = false

func TestThemeLightSeriesLoopColors(t *testing.T) {
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, func(c Color) Color {
		return adjustSeriesColor(c, 1, false)
	}, getThemeSeriesColors(ThemeLight)...)
}

func TestThemeDarkSeriesLoopColors(t *testing.T) {
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, func(c Color) Color {
		return adjustSeriesColor(c, 1, true)
	}, getThemeSeriesColors(ThemeDark)...)
}

func TestThemeVividLightSeriesLoopColors(t *testing.T) {
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, func(c Color) Color {
		return adjustSeriesColor(c, 1, false)
	}, getThemeSeriesColors(ThemeVividLight)...)
}

func TestThemeVividDarkSeriesLoopColors(t *testing.T) {
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, func(c Color) Color {
		return adjustSeriesColor(c, 1, true)
	}, getThemeSeriesColors(ThemeVividDark)...)
}

func TestThemeAntSeriesLoopColors(t *testing.T) {
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, func(c Color) Color {
		return adjustSeriesColor(c, 1, false)
	}, getThemeSeriesColors(ThemeAnt)...)
}

func TestThemeGrafanaSeriesLoopColors(t *testing.T) {
	if !makeSeriesLoopColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, func(c Color) Color {
		return adjustSeriesColor(c, 1, true)
	}, getThemeSeriesColors(ThemeGrafana)...)
}

const makeTrendColorSamples = false

func TestThemeLightTrendColors(t *testing.T) {
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(ThemeLight)...)
}

func TestThemeDarkTrendColors(t *testing.T) {
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(ThemeDark)...)
}

func TestThemeVividLightTrendColors(t *testing.T) {
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(ThemeVividLight)...)
}

func TestThemeVividDarkTrendColors(t *testing.T) {
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(ThemeVividDark)...)
}

func TestThemeAntTrendColors(t *testing.T) {
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(ThemeAnt)...)
}

func TestThemeGrafanaTrendColors(t *testing.T) {
	if !makeTrendColorSamples {
		return // color samples are generated through a test failure
	}

	testColorShadeVariation(t, autoSeriesTrendColor, getThemeSeriesColors(ThemeGrafana)...)
}

func testColorShadeVariation(t *testing.T, mutateFunc func(Color) Color, colors ...Color) {
	t.Helper()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})

	sampleWidth := p.Width() / len(colors)
	for i, c := range colors {
		startX := i * sampleWidth
		endX := (i + 1) * sampleWidth
		if i == len(colors)-1 {
			endX = p.Width() // ensure edge is painted
		}
		p.FilledRect(startX, 0, endX, p.Height(),
			c, ColorTransparent, 0.0)
		margin := (endX - startX) / 4
		c2 := mutateFunc(c)
		p.FilledRect(startX+margin, 0, endX-margin, p.Height(),
			c2, ColorTransparent, 0.0)
	}

	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "", data)
}

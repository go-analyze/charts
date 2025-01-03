package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeBasicFunnelChartOption() FunnelChartOption {
	return FunnelChartOption{
		SeriesList: NewFunnelSeriesList([]float64{
			100, 80, 60, 40, 20,
		}),
		Legend: LegendOption{
			Data: []string{
				"Show", "Click", "Visit", "Inquiry", "Order",
			},
		},
		Title: TitleOption{
			Text: "Funnel",
		},
	}
}

func TestFunnelChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		defaultTheme bool
		makeOptions  func() FunnelChartOption
		result       string
	}{
		{
			name:         "default",
			defaultTheme: true,
			makeOptions:  makeBasicFunnelChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 87 9\nL 117 9\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"102\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"119\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path  d=\"M 177 9\nL 207 9\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"192\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"209\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path  d=\"M 263 9\nL 293 9\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"278\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"295\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path  d=\"M 346 9\nL 376 9\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"361\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"378\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path  d=\"M 445 9\nL 475 9\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"460\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"477\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><text x=\"0\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path  d=\"M 0 35\nL 600 35\nL 540 100\nL 60 100\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><text x=\"280\" y=\"67\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(100%)</text><path  d=\"M 60 102\nL 540 102\nL 480 167\nL 120 167\nL 60 102\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"284\" y=\"134\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(80%)</text><path  d=\"M 120 169\nL 480 169\nL 420 234\nL 180 234\nL 120 169\" style=\"stroke-width:0;stroke:none;fill:rgba(250,200,88,1.0)\"/><text x=\"284\" y=\"201\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(60%)</text><path  d=\"M 180 236\nL 420 236\nL 360 301\nL 240 301\nL 180 236\" style=\"stroke-width:0;stroke:none;fill:rgba(238,102,102,1.0)\"/><text x=\"284\" y=\"268\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(40%)</text><path  d=\"M 240 303\nL 360 303\nL 300 368\nL 300 368\nL 240 303\" style=\"stroke-width:0;stroke:none;fill:rgba(115,192,222,1.0)\"/><text x=\"284\" y=\"335\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(20%)</text></svg>",
		},
		{
			name:         "themed",
			defaultTheme: false,
			makeOptions:  makeBasicFunnelChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(40,40,40,1.0)\"/><path  d=\"M 87 9\nL 117 9\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"102\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><text x=\"119\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path  d=\"M 177 9\nL 207 9\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"192\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><text x=\"209\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path  d=\"M 263 9\nL 293 9\" style=\"stroke-width:3;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><circle cx=\"278\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><text x=\"295\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path  d=\"M 346 9\nL 376 9\" style=\"stroke-width:3;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><circle cx=\"361\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><text x=\"378\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path  d=\"M 445 9\nL 475 9\" style=\"stroke-width:3;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><circle cx=\"460\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><text x=\"477\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><text x=\"0\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path  d=\"M 0 35\nL 600 35\nL 540 100\nL 60 100\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><text x=\"280\" y=\"67\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(100%)</text><path  d=\"M 60 102\nL 540 102\nL 480 167\nL 120 167\nL 60 102\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/><text x=\"284\" y=\"134\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(80%)</text><path  d=\"M 120 169\nL 480 169\nL 420 234\nL 180 234\nL 120 169\" style=\"stroke-width:0;stroke:none;fill:rgba(100,180,210,1.0)\"/><text x=\"284\" y=\"201\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(60%)</text><path  d=\"M 180 236\nL 420 236\nL 360 301\nL 240 301\nL 180 236\" style=\"stroke-width:0;stroke:none;fill:rgba(64,160,110,1.0)\"/><text x=\"284\" y=\"268\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(40%)</text><path  d=\"M 240 303\nL 360 303\nL 300 368\nL 300 368\nL 240 303\" style=\"stroke-width:0;stroke:none;fill:rgba(154,100,180,1.0)\"/><text x=\"284\" y=\"335\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(20%)</text></svg>",
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		if tt.defaultTheme {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)

				validateFunnelChartRender(t, p, tt.makeOptions(), tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p, err := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				require.NoError(t, err)

				validateFunnelChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateFunnelChartRender(t, p, opt, tt.result)
			})
		}
	}
}

func validateFunnelChartRender(t *testing.T, p *Painter, opt FunnelChartOption, expectedResult string) {
	t.Helper()

	_, err := NewFunnelChart(p, opt).Render()
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
}

package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicFunnelChartOption() FunnelChartOption {
	return FunnelChartOption{
		SeriesList: NewSeriesListFunnel([]float64{
			100, 80, 60, 40, 20,
		}),
		Legend: LegendOption{
			SeriesNames: []string{"Show", "Click", "Visit", "Inquiry", "Order"},
		},
		Title: TitleOption{
			Text: "Funnel",
		},
	}
}

func TestNewFunnelChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewFunnelChartOptionWithData([]float64{12, 24, 48})

	assert.Len(t, opt.SeriesList, 3)
	assert.Equal(t, ChartTypeFunnel, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.FunnelChart(opt))
}

func TestFunnelChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() FunnelChartOption
		result      string
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicFunnelChartOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 87 3\nL 117 3\nL 117 16\nL 87 16\nL 87 3\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"119\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path  d=\"M 177 3\nL 207 3\nL 207 16\nL 177 16\nL 177 3\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"209\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path  d=\"M 263 3\nL 293 3\nL 293 16\nL 263 16\nL 263 3\" style=\"stroke:none;fill:rgb(100,180,210)\"/><text x=\"295\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path  d=\"M 346 3\nL 376 3\nL 376 16\nL 346 16\nL 346 3\" style=\"stroke:none;fill:rgb(64,160,110)\"/><text x=\"378\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path  d=\"M 445 3\nL 475 3\nL 475 16\nL 445 16\nL 445 3\" style=\"stroke:none;fill:rgb(154,100,180)\"/><text x=\"477\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><text x=\"0\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path  d=\"M 0 35\nL 600 35\nL 540 100\nL 60 100\nL 0 35\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"264\" y=\"67\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Show(100%)</text><path  d=\"M 60 102\nL 540 102\nL 480 167\nL 120 167\nL 60 102\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"269\" y=\"134\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Click(80%)</text><path  d=\"M 120 169\nL 480 169\nL 420 234\nL 180 234\nL 120 169\" style=\"stroke:none;fill:rgb(100,180,210)\"/><text x=\"271\" y=\"201\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Visit(60%)</text><path  d=\"M 180 236\nL 420 236\nL 360 301\nL 240 301\nL 180 236\" style=\"stroke:none;fill:rgb(64,160,110)\"/><text x=\"264\" y=\"268\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Inquiry(40%)</text><path  d=\"M 240 303\nL 360 303\nL 300 368\nL 300 368\nL 240 303\" style=\"stroke:none;fill:rgb(154,100,180)\"/><text x=\"268\" y=\"335\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Order(20%)</text></svg>",
		},
		{
			name: "custom_legend",
			makeOptions: func() FunnelChartOption {
				opt := makeBasicFunnelChartOption()
				opt.Legend.Symbol = SymbolDot
				opt.Legend.FontStyle = NewFontStyleWithSize(4.0)
				opt.Legend.Vertical = Ptr(true)
				opt.Legend.Offset = OffsetStr{
					Left: PositionRight,
					Top:  PositionBottom,
				}
				opt.Title.Show = Ptr(false)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 552 304\nL 582 304\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"567\" cy=\"304\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"584\" y=\"310\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Show</text><path  d=\"M 552 324\nL 582 324\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"567\" cy=\"324\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"584\" y=\"330\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Click</text><path  d=\"M 552 344\nL 582 344\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"567\" cy=\"344\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"584\" y=\"350\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Visit</text><path  d=\"M 552 364\nL 582 364\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"567\" cy=\"364\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"584\" y=\"370\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path  d=\"M 552 384\nL 582 384\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"567\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"584\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Order</text><path  d=\"M 0 0\nL 600 0\nL 540 72\nL 60 72\nL 0 0\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"264\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Show(100%)</text><path  d=\"M 60 74\nL 540 74\nL 480 146\nL 120 146\nL 60 74\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"269\" y=\"110\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Click(80%)</text><path  d=\"M 120 148\nL 480 148\nL 420 220\nL 180 220\nL 120 148\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"271\" y=\"184\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Visit(60%)</text><path  d=\"M 180 222\nL 420 222\nL 360 294\nL 240 294\nL 180 222\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"264\" y=\"258\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Inquiry(40%)</text><path  d=\"M 240 296\nL 360 296\nL 300 368\nL 300 368\nL 240 296\" style=\"stroke:none;fill:rgb(115,192,222)\"/><text x=\"268\" y=\"332\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Order(20%)</text></svg>",
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		if tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateFunnelChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateFunnelChartRender(t, p, opt, tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)

				validateFunnelChartRender(t, p, tt.makeOptions(), tt.result)
			})
		}
	}
}

func validateFunnelChartRender(t *testing.T, p *Painter, opt FunnelChartOption, expectedResult string) {
	t.Helper()

	err := p.FunnelChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
}

func TestFunnelChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() FunnelChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() FunnelChartOption {
				return NewFunnelChartOptionWithData([]float64{})
			},
			errorMsgContains: "empty series list",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.FunnelChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

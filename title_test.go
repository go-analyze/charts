package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTitleRenderer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			name: "no_content",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "",
					Subtext: "",
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"></svg>",
		},
		{
			name: "offset_number",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset: OffsetStr{
						Left: "20",
						Top:  "20",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"34\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"20\" y=\"50\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "offset_percent",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset: OffsetStr{
						Left: "20%",
						Top:  "20",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"134\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"120\" y=\"50\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "offset_right",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset:  OffsetRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"558\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"544\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "offset_center",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset:  OffsetCenter,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"286\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"272\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "offset_bottom",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset: OffsetStr{
						Top: PositionBottom,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"14\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"0\" y=\"400\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "offset_bottom_right",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset: OffsetStr{
						Top:  PositionBottom,
						Left: PositionRight,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"558\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"544\" y=\"400\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "offset_bottom_center",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:    "title",
					Subtext: "subTitle",
					Offset: OffsetStr{
						Top:  PositionBottom,
						Left: PositionCenter,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"286\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"272\" y=\"400\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "custom_font",
			render: func(p *Painter) ([]byte, error) {
				_, err := newTitlePainter(p, TitleOption{
					Text:             "title",
					Subtext:          "subTitle",
					FontStyle:        NewFontStyleWithSize(40.0).WithColor(ColorBlue),
					SubtextFontStyle: NewFontStyleWithSize(20.0).WithColor(ColorPurple),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"2\" y=\"51\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"0\" y=\"76\" style=\"stroke:none;fill:purple;font-size:25.6px;font-family:'Roboto Medium',sans-serif\">subTitle</text></svg>",
		},
		{
			name: "border",
			render: func(p *Painter) ([]byte, error) {
				theme := GetTheme(ThemeAnt).WithTitleBorderColor(ColorRed)
				_, err := newTitlePainter(p.Child(PainterThemeOption(theme), PainterPaddingOption(NewBoxEqual(100))),
					TitleOption{
						Text:        "title",
						Subtext:     "subTitle",
						BorderWidth: defaultStrokeWidth,
					}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"114\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">title</text><text x=\"100\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">subTitle</text><path  d=\"M 90 140\nL 90 90\nL 166 90\nL 166 140\nL 90 140\" style=\"stroke-width:2;stroke:red;fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)))
			data, err := tt.render(p)
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

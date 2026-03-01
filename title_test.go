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
			assertTestdataSVG(t, data)
		})
	}
}

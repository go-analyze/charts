package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLegend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		render func(*Painter) ([]byte, error)
	}{
		{
			name: "basic",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two", "Three"},
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
				_, err := newLegendPainter(p.Child(PainterPaddingOption(defaultPadding)), LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two", "Three"},
					BorderWidth: 2.0,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_border",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p.Child(PainterPaddingOption(defaultPadding)), LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two", "Three"},
					Vertical:    Ptr(true),
					BorderWidth: 2.0,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "position_left",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two", "Three"},
					Offset:      OffsetLeft,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "position_vertical_with_rect",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two", "Three"},
					Vertical:    Ptr(true),
					Symbol:      SymbolSquare,
					Offset: OffsetStr{
						Left: "10%",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "custom_padding_and_font",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two", "Three"},
					FontStyle:   NewFontStyleWithSize(20.0).WithColor(ColorBlue),
					Padding:     NewBox(20, 200, 20, 20),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "hidden",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"A", "B", "C"},
					Show:        Ptr(false),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "bottom_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
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
			name: "vertical_right_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
					Offset:      OffsetRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_bottom_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
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
			name: "vertical_right_bottom_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
					Offset: OffsetStr{
						Left: PositionRight,
						Top:  PositionBottom,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_right_position_custom_font_size",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
					Offset:      OffsetRight,
					FontStyle:   NewFontStyleWithSize(6.0),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_right_position_with_padding",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
					Offset:      OffsetRight,
					Padding:     NewBoxEqual(120),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "left_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme: GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetLeft,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "center_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme: GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetCenter,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "center_position_center_align_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme: GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetCenter,
					Align:  AlignCenter,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "50%_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme: GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetStr{
						Left: "50%",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_right_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme: GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Vertical: Ptr(true),
					Offset: OffsetStr{
						Left: "440",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "right_alignment",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Align:       AlignRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_right_alignment",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
					Align:       AlignRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
		},
		{
			name: "vertical_right_alignment_left_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Theme:       GetDefaultTheme(),
					SeriesNames: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical:    Ptr(true),
					Offset:      OffsetLeft,
					Align:       AlignRight,
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

func TestLegendCalculateBox(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opt         LegendOption
		expectedBox Box
	}{
		{
			name: "horizontal_center_default",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One", "Two", "Three"},
			},
			expectedBox: Box{
				Top:    -5,
				Bottom: 16,
				Left:   184,
				Right:  416,
				IsSet:  true,
			},
		},
		{
			name: "horizontal_left_position",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One", "Two", "Three"},
				Offset:      OffsetLeft,
			},
			expectedBox: Box{
				Top:    -5,
				Bottom: 16,
				Left:   0,
				Right:  232,
				IsSet:  true,
			},
		},
		{
			name: "horizontal_right_position",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One", "Two", "Three"},
				Offset:      OffsetRight,
			},
			expectedBox: Box{
				Top:    -5,
				Bottom: 16,
				Left:   368,
				Right:  600,
				IsSet:  true,
			},
		},
		{
			name: "horizontal_bottom_position",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One", "Two"},
				Offset: OffsetStr{
					Top: PositionBottom,
				},
			},
			expectedBox: Box{
				Top:    370,
				Bottom: 391,
				Left:   230,
				Right:  370,
				IsSet:  true,
			},
		},
		{
			name: "vertical_left_default",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One", "Two", "Three"},
				Vertical:    Ptr(true),
			},
			expectedBox: Box{
				Top:    -5,
				Bottom: 60,
				Left:   0,
				Right:  72,
				IsSet:  true,
			},
		},
		{
			name: "vertical_right_position",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One", "Two"},
				Vertical:    Ptr(true),
				Offset:      OffsetRight,
			},
			expectedBox: Box{
				Top:    -5,
				Bottom: 40,
				Left:   538,
				Right:  600,
				IsSet:  true,
			},
		},
		{
			name: "vertical_bottom_position",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"A", "B"},
				Vertical:    Ptr(true),
				Offset: OffsetStr{
					Top: PositionBottom,
				},
			},
			expectedBox: Box{
				Top:    350,
				Bottom: 395,
				Left:   0,
				Right:  43,
				IsSet:  true,
			},
		},
		{
			name: "with_custom_padding",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"X", "Y"},
				Padding:     NewBox(10, 20, 30, 40),
			},
			expectedBox: Box{
				Top:    -20,
				Bottom: 56,
				Left:   219,
				Right:  361,
				IsSet:  true,
			},
		},
		{
			name: "empty_hidden",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"One"},
				Show:        Ptr(false),
			},
			expectedBox: BoxZero,
		},
		{
			name: "empty_no_series",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{},
			},
			expectedBox: BoxZero,
		},
		{
			name: "numeric_offset_top",
			opt: LegendOption{
				Theme:       GetDefaultTheme(),
				SeriesNames: []string{"A", "B"},
				Offset: OffsetStr{
					Top: "50",
				},
			},
			expectedBox: Box{
				Top:    45,
				Bottom: 66,
				Left:   248,
				Right:  351,
				IsSet:  true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)))

			lp := newLegendPainter(p, tt.opt)
			box, err := lp.calculateBox()
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBox, box)
		})
	}
}

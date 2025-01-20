package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestNewLegend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			name: "basic",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two", "Three"},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 184 9\nL 214 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"199\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"216\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 264 9\nL 294 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"279\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"296\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two</text><path  d=\"M 346 9\nL 376 9\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"361\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"378\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three</text></svg>",
		},
		{
			name: "border",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p.Child(PainterPaddingOption(defaultPadding)), LegendOption{
					Data:        []string{"One", "Two", "Three"},
					BorderWidth: 2.0,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 184 29\nL 214 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"199\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"216\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 264 29\nL 294 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"279\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"296\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two</text><path  d=\"M 346 29\nL 376 29\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"361\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"378\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three</text><path  d=\"M 174 50\nL 174 10\nL 426 10\nL 426 50\nL 174 50\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "vertical_border",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p.Child(PainterPaddingOption(defaultPadding)), LegendOption{
					Data:        []string{"One", "Two", "Three"},
					Vertical:    True(),
					BorderWidth: 2.0,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"52\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"52\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two</text><path  d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"52\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three</text><path  d=\"M 10 95\nL 10 10\nL 102 10\nL 102 95\nL 10 95\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "position_left",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:   []string{"One", "Two", "Three"},
					Offset: OffsetLeft,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 9\nL 30 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"15\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"32\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 80 9\nL 110 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"95\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"112\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two</text><path  d=\"M 162 9\nL 192 9\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"177\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"194\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three</text></svg>",
		},
		{
			name: "position_vertical_with_rect",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two", "Three"},
					Vertical: True(),
					Icon:     IconRect,
					Offset: OffsetStr{
						Left: "10%",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 60 3\nL 90 3\nL 90 16\nL 60 16\nL 60 3\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"92\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 60 23\nL 90 23\nL 90 36\nL 60 36\nL 60 23\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"92\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two</text><path  d=\"M 60 43\nL 90 43\nL 90 56\nL 60 56\nL 60 43\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"92\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three</text></svg>",
		},
		{
			name: "custom_padding_and_font",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two", "Three"},
					FontStyle: FontStyle{
						FontSize:  20.0,
						FontColor: drawing.ColorBlue,
					},
					Padding: chartdraw.NewBox(200, 20, 20, 20),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 152 204\nL 182 204\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"167\" cy=\"204\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"184\" y=\"210\" style=\"stroke:none;fill:blue;font-size:25.6px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 250 204\nL 280 204\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"265\" cy=\"204\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"282\" y=\"210\" style=\"stroke:none;fill:blue;font-size:25.6px;font-family:'Roboto Medium',sans-serif\">Two</text><path  d=\"M 352 204\nL 382 204\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"367\" cy=\"204\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"384\" y=\"210\" style=\"stroke:none;fill:blue;font-size:25.6px;font-family:'Roboto Medium',sans-serif\">Three</text></svg>",
		},
		{
			name: "hidden",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"A", "B", "C"},
					Show: False(),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"></svg>",
		},
		{
			name: "bottom_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Offset: OffsetStr{
						Top: PositionBottom,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 27 384\nL 57 384\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"42\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"59\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 107 384\nL 137 384\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"122\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"139\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 229 384\nL 259 384\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"244\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"261\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 396 384\nL 426 384\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"411\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"428\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text></svg>",
		},
		{
			name: "vertical_right_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
					Offset:   OffsetRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 421 9\nL 451 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"436\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"453\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 421 29\nL 451 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"436\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"453\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 421 49\nL 451 49\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"436\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"453\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 421 69\nL 451 69\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"436\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"453\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text></svg>",
		},
		{
			name: "vertical_bottom_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
					Offset: OffsetStr{
						Top: PositionBottom,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 324\nL 30 324\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"15\" cy=\"324\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"32\" y=\"330\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 0 344\nL 30 344\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"15\" cy=\"344\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"32\" y=\"350\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 0 364\nL 30 364\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"15\" cy=\"364\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"32\" y=\"370\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 0 384\nL 30 384\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"15\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"32\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text></svg>",
		},
		{
			name: "vertical_right_bottom_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 421 324\nL 451 324\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"436\" cy=\"324\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"453\" y=\"330\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 421 344\nL 451 344\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"436\" cy=\"344\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"453\" y=\"350\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 421 364\nL 451 364\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"436\" cy=\"364\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"453\" y=\"370\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 421 384\nL 451 384\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"436\" cy=\"384\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"453\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text></svg>",
		},
		{
			name: "vertical_right_position_custom_font_size",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
					Offset:   OffsetRight,
					FontStyle: FontStyle{
						FontSize: 6.0,
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 494 9\nL 524 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"509\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"526\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:7.7px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 494 29\nL 524 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"509\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"526\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:7.7px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 494 49\nL 524 49\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"509\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"526\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:7.7px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 494 69\nL 524 69\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"509\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"526\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:7.7px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text></svg>",
		},
		{
			name: "vertical_right_position_with_padding",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
					Offset:   OffsetRight,
					Padding:  Box{Top: 120, Left: 120, Right: 120, Bottom: 120},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 301 124\nL 331 124\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"316\" cy=\"124\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"333\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 301 144\nL 331 144\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"316\" cy=\"144\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"333\" y=\"150\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 301 164\nL 331 164\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"316\" cy=\"164\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"333\" y=\"170\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 301 184\nL 331 184\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"316\" cy=\"184\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"333\" y=\"190\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text></svg>",
		},
		{
			name: "left_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetLeft,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 9\nL 30 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"15\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"32\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 80 9\nL 110 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"95\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"112\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 202 9\nL 232 9\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"217\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"234\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 369 9\nL 399 9\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"384\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"401\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 0 24\nL 30 24\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"15\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"32\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Five Words Is Even Longer</text><path  d=\"M 233 24\nL 263 24\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:none\"/><circle cx=\"248\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:rgb(59,162,114)\"/><text x=\"265\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Six Words Is The Longest Tested</text></svg>",
		},
		{
			name: "center_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetCenter,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 9\nL 30 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"15\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"32\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 80 9\nL 110 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"95\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"112\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 202 9\nL 232 9\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"217\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"234\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 369 9\nL 399 9\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"384\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"401\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 0 24\nL 30 24\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"15\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"32\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Five Words Is Even Longer</text><path  d=\"M 233 24\nL 263 24\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:none\"/><circle cx=\"248\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:rgb(59,162,114)\"/><text x=\"265\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Six Words Is The Longest Tested</text></svg>",
		},
		{
			name: "center_position_center_align_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Offset: OffsetCenter,
					Align:  AlignCenter,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 9\nL 30 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"15\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"32\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 80 9\nL 110 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"95\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"112\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 202 9\nL 232 9\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"217\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"234\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 369 9\nL 399 9\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"384\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"401\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 56 24\nL 86 24\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"71\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"88\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Five Words Is Even Longer</text><path  d=\"M 289 24\nL 319 24\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:none\"/><circle cx=\"304\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:rgb(59,162,114)\"/><text x=\"321\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Six Words Is The Longest Tested</text></svg>",
		},
		{
			name: "50%_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 300 9\nL 330 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"315\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"332\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 380 9\nL 410 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"395\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"412\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 300 24\nL 330 24\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"315\" cy=\"24\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"332\" y=\"30\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 300 39\nL 330 39\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"315\" cy=\"39\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"332\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 300 54\nL 330 54\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"315\" cy=\"54\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"332\" y=\"60\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Five Words Is Even Longer</text><path  d=\"M 300 69\nL 330 69\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:none\"/><circle cx=\"315\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:rgb(59,162,114)\"/><text x=\"332\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Six Words Is The Longest Tested</text></svg>",
		},
		{
			name: "vertical_right_position_overflow",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data: []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer",
						"Five Words Is Even Longer", "Six Words Is The Longest Tested"},
					Vertical: True(),
					Offset: OffsetStr{
						Left: "440",
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 440 9\nL 470 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"455\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"472\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 440 29\nL 470 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"455\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"472\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 440 49\nL 470 49\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"455\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"472\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 440 69\nL 470 69\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"455\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"472\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 440 89\nL 470 89\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"455\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"472\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Five Words Is Even Longer</text><path  d=\"M 440 109\nL 470 109\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:none\"/><circle cx=\"455\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgb(59,162,114);fill:rgb(59,162,114)\"/><text x=\"472\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Six Words Is The Longest Tested</text></svg>",
		},
		{
			name: "right_alignment",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:  []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Align: AlignRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"27\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 57 9\nL 87 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"72\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"107\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 179 9\nL 209 9\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"194\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"229\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 346 9\nL 376 9\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"361\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"396\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 545 9\nL 575 9\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"560\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/></svg>",
		},
		{
			name: "vertical_right_alignment",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
					Align:    AlignRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"540\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 570 9\nL 600 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"585\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"498\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 570 29\nL 600 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"585\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"453\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 570 49\nL 600 49\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"585\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"421\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 570 69\nL 600 69\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"585\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/></svg>",
		},
		{
			name: "vertical_right_alignment_left_position",
			render: func(p *Painter) ([]byte, error) {
				_, err := newLegendPainter(p, LegendOption{
					Data:     []string{"One", "Two Word", "Three Word Item", "Four Words Is Longer"},
					Vertical: True(),
					Offset:   OffsetLeft,
					Align:    AlignRight,
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><text x=\"119\" y=\"15\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">One</text><path  d=\"M 149 9\nL 179 9\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"164\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"77\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Two Word</text><path  d=\"M 149 29\nL 179 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"164\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"32\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Three Word Item</text><path  d=\"M 149 49\nL 179 49\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"164\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"0\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Four Words Is Longer</text><path  d=\"M 149 69\nL 179 69\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"164\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/></svg>",
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

package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAxisRender(t *testing.T) {
	t.Parallel()

	dayLabels := []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	letterLabels := []string{"A", "B", "C", "D", "E", "F", "G"}
	fs := NewFontStyleWithSize(18)
	fs.FontColor = ColorGreen

	tests := []struct {
		name          string
		optionFactory func(p *Painter) axisOption
		result        string
	}{
		{
			name: "x-axis",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					BoundaryGap: Ptr(true),
					FontStyle:   fs,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, 0, fs))
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 316\nL 550 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 321\nL 50 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 321\nL 121 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 321\nL 192 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 321\nL 264 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 321\nL 335 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 321\nL 407 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 321\nL 478 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 321\nL 550 316\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"62\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"136\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"205\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"279\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"358\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"425\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"494\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Sun</text></svg>",
		},
		{
			name: "x-axis_rotation45",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					BoundaryGap: Ptr(true),
					FontStyle:   fs,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, DegreesToRadians(45), fs))
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 285\nL 550 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 290\nL 50 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 290\nL 121 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 290\nL 192 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 290\nL 264 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 290\nL 335 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 290\nL 407 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 290\nL 478 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 290\nL 550 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"61\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,61,314)\">Mon</text><text x=\"135\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,135,314)\">Tue</text><text x=\"204\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,204,314)\">Wed</text><text x=\"278\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,278,314)\">Thu</text><text x=\"354\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,354,314)\">Fri</text><text x=\"423\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,423,314)\">Sat</text><text x=\"493\" y=\"314\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,493,314)\">Sun</text></svg>",
		},
		{
			name: "x-axis_rotation90",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					Labels:      dayLabels,
					BoundaryGap: Ptr(true),
					FontStyle:   fs,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, DegreesToRadians(90), fs))
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 282\nL 550 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 287\nL 50 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 287\nL 121 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 287\nL 192 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 287\nL 264 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 287\nL 335 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 287\nL 407 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 287\nL 478 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 287\nL 550 282\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"74\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,74,301)\">Mon</text><text x=\"145\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,145,301)\">Tue</text><text x=\"217\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,217,301)\">Wed</text><text x=\"288\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,288,301)\">Thu</text><text x=\"360\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,360,301)\">Fri</text><text x=\"431\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,431,301)\">Sat</text><text x=\"503\" y=\"301\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,503,301)\">Sun</text></svg>",
		},
		{
			name: "x-axis_splitline",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:        newTestRangeForLabels(letterLabels, 0, fs),
					splitLineShow: true,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 317\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 322\nL 50 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 322\nL 121 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 322\nL 192 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 322\nL 264 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 322\nL 335 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 322\nL 407 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 322\nL 478 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 322\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"77\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"149\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"220\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"291\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"365\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"436\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"506\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">G</text><path d=\"M 121 50\nL 121 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 192 50\nL 192 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 264 50\nL 264 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 335 50\nL 335 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 407 50\nL 407 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 478 50\nL 478 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 550 50\nL 550 317\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/></svg>",
		},
		{
			name: "y-axis_left",
			optionFactory: func(p *Painter) axisOption {
				opt := YAxisOption{
					Position:       PositionLeft,
					isCategoryAxis: true,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, 0, fs))
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 114 50\nL 114 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 50\nL 114 50\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 92\nL 114 92\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 135\nL 114 135\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 178\nL 114 178\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 221\nL 114 221\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 264\nL 114 264\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 307\nL 114 307\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 109 350\nL 114 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"64\" y=\"80\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Sun</text><text x=\"70\" y=\"122\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"77\" y=\"165\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"64\" y=\"207\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"58\" y=\"250\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"64\" y=\"292\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"58\" y=\"335\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Mon</text></svg>",
		},
		{
			name: "y-axis_right",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:        newTestRangeForLabels(dayLabels, 0, fs),
					position:      PositionRight,
					boundaryGap:   Ptr(false),
					splitLineShow: true,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 486 50\nL 486 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 50\nL 491 50\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 100\nL 491 100\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 150\nL 491 150\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 200\nL 491 200\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 250\nL 491 250\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 300\nL 491 300\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 486 350\nL 491 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"496\" y=\"59\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Sun</text><text x=\"496\" y=\"108\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"496\" y=\"158\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"496\" y=\"208\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"496\" y=\"257\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"496\" y=\"307\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"496\" y=\"357\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">Mon</text><path d=\"M 50 50\nL 486 50\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 50 100\nL 486 100\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 50 150\nL 486 150\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 50 200\nL 486 200\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 50 250\nL 486 250\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/><path d=\"M 50 300\nL 486 300\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/></svg>",
		},
		{
			name: "reduced_label_count",
			optionFactory: func(p *Painter) axisOption {
				aRange := newTestRangeForLabels(letterLabels, 0, fs)
				aRange.labelCount -= 2
				return axisOption{
					aRange:        aRange,
					splitLineShow: false,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 317\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 322\nL 50 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 322\nL 121 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 322\nL 192 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 322\nL 264 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 322\nL 335 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 322\nL 407 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 322\nL 478 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 322\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"149\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"291\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"365\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"534\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "label_start_offset",
			optionFactory: func(p *Painter) axisOption {
				aRange := newTestRangeForLabels(letterLabels, 0, fs)
				aRange.dataStartIndex = 2
				return axisOption{
					aRange: aRange,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 317\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 322\nL 192 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 322\nL 264 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 322\nL 335 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 322\nL 407 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 322\nL 478 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 322\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"220\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"291\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"365\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"436\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"506\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "custom_font",
			optionFactory: func(p *Painter) axisOption {
				fs := FontStyle{
					FontSize:  40.0,
					FontColor: ColorBlue,
				}
				return axisOption{
					aRange: newTestRangeForLabels(letterLabels, 0, fs),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 293\nL 550 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 298\nL 50 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 298\nL 121 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 298\nL 192 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 298\nL 264 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 298\nL 335 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 298\nL 407 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 298\nL 478 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 298\nL 550 293\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"68\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"140\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"211\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"282\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"357\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"428\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"497\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "boundary_gap_disable",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:      newTestRangeForLabels(letterLabels, 0, fs),
					boundaryGap: Ptr(false),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 317\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 322\nL 50 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 133 322\nL 133 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 216 322\nL 216 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 300 322\nL 300 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 383 322\nL 383 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 466 322\nL 466 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 322\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"132\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"215\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"299\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"382\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"465\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"534\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "boundary_gap_enable",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:      newTestRangeForLabels(letterLabels, 0, fs),
					boundaryGap: Ptr(true),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 317\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 322\nL 50 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 121 322\nL 121 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 192 322\nL 192 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 264 322\nL 264 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 335 322\nL 335 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 407 322\nL 407 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 478 322\nL 478 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 322\nL 550 317\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"77\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"149\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"220\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"291\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"365\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"436\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"506\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "dense_category_data",
			optionFactory: func(p *Painter) axisOption {
				const count = 1000
				labelLen := len(strconv.Itoa(count))
				labels := make([]string, count)
				tsl := testSeriesList{}
				for i := range labels {
					label := strconv.Itoa(i + 1)
					for len(label) < labelLen {
						label = "0" + label
					}
					labels[i] = label
					tsl = append(tsl, testSeries{values: []float64{float64(i)}})
				}
				return axisOption{
					aRange: calculateCategoryAxisRange(p, p.Width(), false, false, labels, 0,
						0, 0, 0, tsl, 0, fs),
					boundaryGap: Ptr(true),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 50 314\nL 550 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 50 319\nL 50 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 133 319\nL 133 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 216 319\nL 216 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 300 319\nL 300 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 383 319\nL 383 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 466 319\nL 466 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path d=\"M 550 319\nL 550 314\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">0001</text><text x=\"107\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">0168</text><text x=\"190\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">0334</text><text x=\"274\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">0501</text><text x=\"357\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">0668</text><text x=\"440\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">0834</text><text x=\"497\" y=\"347\" style=\"stroke:none;fill:green;font-size:23px;font-family:'Roboto Medium',sans-serif\">1000</text></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterPaddingOption(NewBoxEqual(50)))

			opt := tt.optionFactory(p)
			opt.axisColor = ColorBlue
			opt.axisSplitLineColor = ColorGray
			_, err := newAxisPainter(p, opt).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

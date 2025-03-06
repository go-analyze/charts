package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAxis(t *testing.T) {
	t.Parallel()

	dayLabels := []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	letterLabels := []string{"A", "B", "C", "D", "E", "F", "G"}

	tests := []struct {
		name          string
		optionFactory func() axisOption
		result        string
	}{
		{
			name: "x-axis",
			optionFactory: func() axisOption {
				opt := XAxisOption{
					Labels:      dayLabels,
					BoundaryGap: Ptr(true),
					FontStyle:   NewFontStyleWithSize(18),
				}
				return opt.toAxisOption(nil)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 314\nL 550 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 50 319\nL 50 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 121 319\nL 121 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 192 319\nL 192 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 264 319\nL 264 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 335 319\nL 335 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 407 319\nL 407 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 478 319\nL 478 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 550 319\nL 550 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"62\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"136\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"205\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"279\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"358\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"425\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"494\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\">Sun</text></svg>",
		},
		{
			name: "x-axis_rotation45",
			optionFactory: func() axisOption {
				opt := XAxisOption{
					Labels:        dayLabels,
					BoundaryGap:   Ptr(true),
					FontStyle:     NewFontStyleWithSize(18),
					LabelRotation: DegreesToRadians(45),
				}
				return opt.toAxisOption(nil)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 288\nL 550 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 50 293\nL 50 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 121 293\nL 121 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 192 293\nL 192 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 264 293\nL 264 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 335 293\nL 335 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 407 293\nL 407 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 478 293\nL 478 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 550 293\nL 550 288\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"61\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,61,347)\">Mon</text><text x=\"135\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,135,347)\">Tue</text><text x=\"204\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,204,347)\">Wed</text><text x=\"278\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,278,347)\">Thu</text><text x=\"354\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,354,347)\">Fri</text><text x=\"423\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,423,347)\">Sat</text><text x=\"493\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(45.00,493,347)\">Sun</text></svg>",
		},
		{
			name: "x-axis_rotation90",
			optionFactory: func() axisOption {
				opt := XAxisOption{
					Labels:        dayLabels,
					BoundaryGap:   Ptr(true),
					FontStyle:     NewFontStyleWithSize(18),
					LabelRotation: DegreesToRadians(90),
				}
				return opt.toAxisOption(nil)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 290\nL 550 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 50 295\nL 50 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 121 295\nL 121 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 192 295\nL 192 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 264 295\nL 264 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 335 295\nL 335 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 407 295\nL 407 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 478 295\nL 478 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 550 295\nL 550 290\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"74\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,74,347)\">Mon</text><text x=\"145\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,145,347)\">Tue</text><text x=\"217\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,217,347)\">Wed</text><text x=\"288\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,288,347)\">Thu</text><text x=\"360\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,360,347)\">Fri</text><text x=\"431\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,431,347)\">Sat</text><text x=\"503\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:23px;font-family:'Roboto Medium',sans-serif\" transform=\"rotate(90.00,503,347)\">Sun</text></svg>",
		},
		{
			name: "x-axis_splitline",
			optionFactory: func() axisOption {
				return axisOption{
					labels:        dayLabels,
					splitLineShow: true,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 321\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 326\nL 50 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 121 326\nL 121 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 192 326\nL 192 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 264 326\nL 264 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 335 326\nL 335 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 407 326\nL 407 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 478 326\nL 478 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 326\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"70\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"143\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"213\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"286\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"362\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"431\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"501\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 121 50\nL 121 321\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 192 50\nL 192 321\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 264 50\nL 264 321\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 335 50\nL 335 321\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 407 50\nL 407 321\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 478 50\nL 478 321\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 550 50\nL 550 321\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "x-axis_left",
			optionFactory: func() axisOption {
				return axisOption{
					labels:      dayLabels,
					boundaryGap: Ptr(false),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 321\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 326\nL 50 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 133 326\nL 133 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 216 326\nL 216 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 300 326\nL 300 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 383 326\nL 383 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 466 326\nL 466 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 326\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"132\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"215\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"299\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"382\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"523\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text></svg>",
		},
		{
			name: "y-axis_left",
			optionFactory: func() axisOption {
				opt := YAxisOption{
					Labels:         dayLabels,
					Position:       PositionLeft,
					isCategoryAxis: true,
				}
				return opt.toAxisOption(nil)
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 90 50\nL 90 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 50\nL 90 50\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 92\nL 90 92\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 135\nL 90 135\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 178\nL 90 178\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 221\nL 90 221\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 264\nL 90 264\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 307\nL 90 307\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 85 350\nL 90 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"49\" y=\"77\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"53\" y=\"119\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"49\" y=\"162\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"53\" y=\"204\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"62\" y=\"247\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"57\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"53\" y=\"332\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text></svg>",
		},
		{
			name: "y-axis_center",
			optionFactory: func() axisOption {
				return axisOption{
					labels:        dayLabels,
					position:      PositionLeft,
					boundaryGap:   Ptr(false),
					splitLineShow: true,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 90 50\nL 90 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 50\nL 90 50\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 100\nL 90 100\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 150\nL 90 150\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 200\nL 90 200\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 250\nL 90 250\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 300\nL 90 300\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 85 350\nL 90 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"56\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"53\" y=\"105\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"49\" y=\"155\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"53\" y=\"205\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"62\" y=\"254\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"57\" y=\"304\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"53\" y=\"354\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 86 50\nL 550 50\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 86 100\nL 550 100\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 86 150\nL 550 150\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 86 200\nL 550 200\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 86 250\nL 550 250\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 86 300\nL 550 300\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "y-axis_right",
			optionFactory: func() axisOption {
				return axisOption{
					labels:        dayLabels,
					position:      PositionRight,
					boundaryGap:   Ptr(false),
					splitLineShow: true,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 510 50\nL 510 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 50\nL 515 50\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 100\nL 515 100\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 150\nL 515 150\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 200\nL 515 200\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 250\nL 515 250\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 300\nL 515 300\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 510 350\nL 515 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"520\" y=\"56\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"520\" y=\"105\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"520\" y=\"155\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"520\" y=\"205\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"520\" y=\"254\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"520\" y=\"304\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"520\" y=\"354\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 50 50\nL 510 50\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 50 100\nL 510 100\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 50 150\nL 510 150\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 50 200\nL 510 200\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 50 250\nL 510 250\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 50 300\nL 510 300\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "top",
			optionFactory: func() axisOption {
				return axisOption{
					labels:    dayLabels,
					formatter: "{value} --",
					position:  PositionTop,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 350\nL 550 350\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 350\nL 50 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 121 350\nL 121 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 192 350\nL 192 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 264 350\nL 264 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 335 350\nL 335 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 407 350\nL 407 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 478 350\nL 478 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 350\nL 550 345\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"63\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon --</text><text x=\"136\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue --</text><text x=\"206\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed --</text><text x=\"279\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu --</text><text x=\"355\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri --</text><text x=\"424\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat --</text><text x=\"494\" y=\"325\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun --</text></svg>",
		},
		{
			name: "reduced_label_count",
			optionFactory: func() axisOption {
				return axisOption{
					labels:               letterLabels,
					splitLineShow:        false,
					labelCountAdjustment: -1,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 321\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 326\nL 50 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 121 326\nL 121 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 192 326\nL 192 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 335 326\nL 335 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 407 326\nL 407 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 478 326\nL 478 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 326\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"151\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"223\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"367\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"438\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"539\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "custom_unit",
			optionFactory: func() axisOption {
				return axisOption{
					labels:        letterLabels,
					splitLineShow: false,
					unit:          10,
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 321\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 326\nL 50 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 335 326\nL 335 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 326\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"367\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"539\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "custom_font",
			optionFactory: func() axisOption {
				return axisOption{
					labels: letterLabels,
					labelFontStyle: FontStyle{
						FontSize:  40.0,
						FontColor: ColorBlue,
					},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 285\nL 550 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 290\nL 50 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 121 290\nL 121 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 192 290\nL 192 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 264 290\nL 264 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 335 290\nL 335 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 407 290\nL 407 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 478 290\nL 478 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 290\nL 550 285\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"68\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"140\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"211\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"282\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"357\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"428\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"497\" y=\"347\" style=\"stroke:none;fill:blue;font-size:51.1px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "boundary_gap_disable",
			optionFactory: func() axisOption {
				return axisOption{
					labels:      letterLabels,
					boundaryGap: Ptr(false),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 321\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 326\nL 50 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 133 326\nL 133 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 216 326\nL 216 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 300 326\nL 300 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 383 326\nL 383 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 466 326\nL 466 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 326\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"49\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"132\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"215\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"299\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"382\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"465\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"539\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "boundary_gap_enable",
			optionFactory: func() axisOption {
				return axisOption{
					labels:      letterLabels,
					boundaryGap: Ptr(true),
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 50 321\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 50 326\nL 50 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 121 326\nL 121 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 192 326\nL 192 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 264 326\nL 264 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 335 326\nL 335 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 407 326\nL 407 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 478 326\nL 478 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 550 326\nL 550 321\" style=\"stroke-width:1;stroke:blue;fill:none\"/><text x=\"80\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"151\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"223\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"294\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"367\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"438\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"509\" y=\"347\" style=\"stroke:none;fill:green;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
	}

	axisTheme := MakeTheme(ThemeOption{
		IsDarkMode:         false,
		AxisStrokeColor:    ColorBlue,
		AxisSplitLineColor: ColorRed,
		BackgroundColor:    ColorWhite,
		TextColor:          ColorGreen,
	})
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(axisTheme), PainterPaddingOption(NewBoxEqual(50)))

			_, err := newAxisPainter(p, tt.optionFactory()).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

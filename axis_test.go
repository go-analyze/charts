package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestAxis(t *testing.T) {
	t.Parallel()

	dayLabels := []string{
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
		"Sun",
	}
	letterLabels := []string{"A", "B", "C", "D", "E", "F", "G"}

	tests := []struct {
		name   string
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			name: "x-axis_bottom",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:          dayLabels,
					SplitLineShow: true,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 375\nL 0 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 85 375\nL 85 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 171 375\nL 171 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 257 375\nL 257 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 342 375\nL 342 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 428 375\nL 428 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 514 375\nL 514 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 600 375\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 0 370\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"27\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"115\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"199\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"286\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"376\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"460\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"544\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 85 0\nL 85 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 171 0\nL 171 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 257 0\nL 257 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 342 0\nL 342 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 428 0\nL 428 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 514 0\nL 514 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 600 0\nL 600 370\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/></svg>",
		},
		{
			name: "x-axis_bottom_left",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:        dayLabels,
					BoundaryGap: False(),
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 375\nL 0 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 100 375\nL 100 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 200 375\nL 200 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 300 375\nL 300 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 400 375\nL 400 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 500 375\nL 500 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 600 375\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 0 370\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"-1\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"99\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"199\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"299\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"399\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"499\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"573\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text></svg>",
		},
		{
			name: "y-axis_left",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:     dayLabels,
					Position: PositionLeft,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 36 0\nL 41 0\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 57\nL 41 57\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 114\nL 41 114\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 171\nL 41 171\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 228\nL 41 228\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 285\nL 41 285\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 342\nL 41 342\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 400\nL 41 400\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 41 0\nL 41 400\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"0\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"4\" y=\"92\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"0\" y=\"149\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"4\" y=\"206\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"13\" y=\"263\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"8\" y=\"320\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"4\" y=\"378\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text></svg>",
		},
		{
			name: "y-axis_center",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:          dayLabels,
					Position:      PositionLeft,
					BoundaryGap:   False(),
					SplitLineShow: true,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 36 0\nL 41 0\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 66\nL 41 66\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 133\nL 41 133\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 200\nL 41 200\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 266\nL 41 266\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 333\nL 41 333\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 36 400\nL 41 400\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 41 0\nL 41 400\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"0\" y=\"7\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"4\" y=\"73\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"0\" y=\"140\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"4\" y=\"207\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"13\" y=\"273\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"8\" y=\"340\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"4\" y=\"407\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 41 0\nL 600 0\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 41 66\nL 600 66\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 41 133\nL 600 133\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 41 200\nL 600 200\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 41 266\nL 600 266\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 41 333\nL 600 333\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/></svg>",
		},
		{
			name: "y-axis_right",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:          dayLabels,
					Position:      PositionRight,
					BoundaryGap:   False(),
					SplitLineShow: true,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 559 0\nL 564 0\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 66\nL 564 66\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 133\nL 564 133\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 200\nL 564 200\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 266\nL 564 266\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 333\nL 564 333\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 400\nL 564 400\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 559 0\nL 559 400\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"569\" y=\"7\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"569\" y=\"73\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"569\" y=\"140\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"569\" y=\"207\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"569\" y=\"273\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"569\" y=\"340\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"569\" y=\"407\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 0 0\nL 559 0\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 0 66\nL 559 66\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 0 133\nL 559 133\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 0 200\nL 559 200\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 0 266\nL 559 266\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/><path  d=\"M 0 333\nL 559 333\" style=\"stroke-width:1;stroke:rgba(0,0,0,1.0);fill:none\"/></svg>",
		},
		{
			name: "top",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:      dayLabels,
					Formatter: "{value} --",
					Position:  PositionTop,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 380\nL 0 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 85 380\nL 85 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 171 380\nL 171 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 257 380\nL 257 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 342 380\nL 342 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 428 380\nL 428 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 514 380\nL 514 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 600 380\nL 600 375\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 0 380\nL 600 380\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"20\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon --</text><text x=\"108\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue --</text><text x=\"192\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed --</text><text x=\"279\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu --</text><text x=\"369\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri --</text><text x=\"453\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat --</text><text x=\"537\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun --</text></svg>",
		},
		{
			name: "reduced_label_count",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:                 letterLabels,
					SplitLineShow:        false,
					LabelCountAdjustment: -1,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 375\nL 0 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 85 375\nL 85 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 171 375\nL 171 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 342 375\nL 342 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 428 375\nL 428 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 514 375\nL 514 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 600 375\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 0 370\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"-1\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"123\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"209\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"381\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"467\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"589\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "custom_unit",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data:          letterLabels,
					SplitLineShow: false,
					Unit:          10,
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 375\nL 0 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 342 375\nL 342 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 600 375\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 0 370\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"-1\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"381\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"589\" y=\"395\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
		{
			name: "custom_font",
			render: func(p *Painter) ([]byte, error) {
				_, _ = NewAxisPainter(p, AxisOption{
					Data: letterLabels,
					FontStyle: FontStyle{
						FontSize:  40.0,
						FontColor: drawing.ColorBlue,
					},
				}).Render()
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 375\nL 0 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 85 375\nL 85 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 171 375\nL 171 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 257 375\nL 257 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 342 375\nL 342 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 428 375\nL 428 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 514 375\nL 514 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 600 375\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 0 370\nL 600 370\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"25\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"112\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"197\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"282\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"371\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"457\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"540\" y=\"431\" style=\"stroke-width:0;stroke:none;fill:rgba(0,0,255,1.0);font-size:51.1px;font-family:'Roboto Medium',sans-serif\">G</text></svg>",
		},
	}

	axisTheme := MakeTheme(ThemeOption{
		IsDarkMode:         false,
		AxisStrokeColor:    Color{R: 110, G: 112, B: 121, A: 255},
		AxisSplitLineColor: drawing.ColorBlack,
		BackgroundColor:    drawing.ColorWhite,
		TextColor:          Color{R: 70, G: 70, B: 70, A: 255},
	})
	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(axisTheme))
			require.NoError(t, err)
			data, err := tt.render(p)
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, string(data))
		})
	}
}

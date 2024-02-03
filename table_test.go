package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableChart(t *testing.T) {
	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				_, err := NewTableChart(p, TableChartOption{
					Header: []string{
						"Name",
						"Age",
						"Address",
						"Tag",
						"Action",
					},
					Spans: []int{
						1,
						1,
						2,
						1,
						// span and header do not match, and are automatically set to 1 at the end
						// 1,
					},
					Data: [][]string{
						{
							"John Brown",
							"32",
							"New York No. 1 Lake Park",
							"nice, developer",
							"Send Mail",
						},
						{
							"Jim Green	",
							"42",
							"London No. 1 Lake Park",
							"wow",
							"Send Mail",
						},
						{
							"Joe Black	",
							"32",
							"Sidney No. 1 Lake Park",
							"cool, teacher",
							"Send Mail",
						},
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 35\nL 0 35\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(240,240,240,1.0)\"/><path  d=\"M 0 35\nL 600 35\nL 600 90\nL 0 90\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 0 90\nL 600 90\nL 600 125\nL 0 125\nL 0 90\" style=\"stroke-width:0;stroke:none;fill:rgba(247,247,247,1.0)\"/><path  d=\"M 0 125\nL 600 125\nL 600 180\nL 0 180\nL 0 125\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"110\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"210\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"410\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"510\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John</text><text x=\"10\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brown</text><text x=\"110\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"210\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No. 1 Lake Park</text><text x=\"410\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"410\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"510\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"110\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"210\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1 Lake Park</text><text x=\"410\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"510\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"110\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"210\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1 Lake Park</text><text x=\"410\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool,</text><text x=\"410\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">teacher</text><text x=\"510\" y=\"147\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
		},
		{
			render: func(p *Painter) ([]byte, error) {
				_, err := NewTableChart(p, TableChartOption{
					Header: []string{
						"Name",
						"Age",
						"Address",
						"Tag",
						"Action",
					},
					Data: [][]string{
						{
							"John Brown",
							"32",
							"New York No. 1 Lake Park",
							"nice, developer",
							"Send Mail",
						},
						{
							"Jim Green	",
							"42",
							"London No. 1 Lake Park",
							"wow",
							"Send Mail",
						},
						{
							"Joe Black	",
							"32",
							"Sidney No. 1 Lake Park",
							"cool, teacher",
							"Send Mail",
						},
					},
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 35\nL 0 35\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(240,240,240,1.0)\"/><path  d=\"M 0 35\nL 600 35\nL 600 90\nL 0 90\nL 0 35\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 0 90\nL 600 90\nL 600 145\nL 0 145\nL 0 90\" style=\"stroke-width:0;stroke:none;fill:rgba(247,247,247,1.0)\"/><path  d=\"M 0 145\nL 600 145\nL 600 200\nL 0 200\nL 0 145\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Name</text><text x=\"130\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Age</text><text x=\"250\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Address</text><text x=\"370\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tag</text><text x=\"490\" y=\"22\" style=\"stroke-width:0;stroke:none;fill:rgba(98,105,118,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Action</text><text x=\"10\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">John Brown</text><text x=\"130\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">New York No.</text><text x=\"250\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1 Lake Park</text><text x=\"370\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">nice,</text><text x=\"370\" y=\"77\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">developer</text><text x=\"490\" y=\"57\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jim Green</text><text x=\"130\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"250\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">London No. 1</text><text x=\"250\" y=\"132\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">wow</text><text x=\"490\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text><text x=\"10\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Joe Black</text><text x=\"130\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"250\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sidney No. 1</text><text x=\"250\" y=\"187\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Lake Park</text><text x=\"370\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">cool, teacher</text><text x=\"490\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Send Mail</text></svg>",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				Type:   ChartOutputSVG,
				Width:  600,
				Height: 400,
			}, PainterThemeOption(defaultTheme))
			require.NoError(t, err)
			data, err := tt.render(p)
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, string(data))
		})
	}
}

package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBarChart(t *testing.T) {
	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				seriesList := NewSeriesListDataFromValues([][]float64{
					{
						2.0,
						4.9,
						7.0,
						23.2,
						25.6,
						76.7,
						135.6,
						162.2,
						32.6,
						20.0,
						6.4,
						3.3,
					},
					{
						2.6,
						5.9,
						9.0,
						26.4,
						28.7,
						70.7,
						175.6,
						182.2,
						48.7,
						18.8,
						6.0,
						2.3,
					},
				})
				for index := range seriesList {
					seriesList[index].Label.Show = true
				}
				_, err := NewBarChart(p, BarChartOption{
					Padding: Box{
						Left:   10,
						Top:    10,
						Right:  10,
						Bottom: 10,
					},
					SeriesList: seriesList,
					XAxis: NewXAxisOption([]string{
						"Jan",
						"Feb",
						"Mar",
						"Apr",
						"May",
						"Jun",
						"Jul",
						"Aug",
						"Sep",
						"Oct",
						"Nov",
						"Dec",
					}),
					YAxisOptions: NewYAxisOptions([]string{
						"Rainfall",
						"Evaporation",
					}),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"31\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"10\" y=\"75\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">166.39</text><text x=\"10\" y=\"133\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">132.79</text><text x=\"18\" y=\"192\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">99.19</text><text x=\"18\" y=\"250\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">65.59</text><text x=\"18\" y=\"308\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">31.99</text><text x=\"22\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">-1.60</text><path  d=\"M 68 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 68\nL 590 68\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 126\nL 590 126\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 243\nL 590 243\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 301\nL 590 301\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 365\nL 68 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 111 365\nL 111 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 155 365\nL 155 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 242 365\nL 242 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 285 365\nL 285 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 329 365\nL 329 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 372 365\nL 372 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 416 365\nL 416 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 503 365\nL 503 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 546 365\nL 546 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 68 360\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"68\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"120\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"162\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"248\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"294\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"340\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"380\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"424\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"510\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"563\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 73 354\nL 88 354\nL 88 359\nL 73 359\nL 73 354\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 116 349\nL 131 349\nL 131 359\nL 116 359\nL 116 349\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 160 346\nL 175 346\nL 175 359\nL 160 359\nL 160 346\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 203 317\nL 218 317\nL 218 359\nL 203 359\nL 203 317\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 247 313\nL 262 313\nL 262 359\nL 247 359\nL 247 313\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 290 225\nL 305 225\nL 305 359\nL 290 359\nL 290 225\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 334 122\nL 349 122\nL 349 359\nL 334 359\nL 334 122\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 377 76\nL 392 76\nL 392 359\nL 377 359\nL 377 76\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 421 301\nL 436 301\nL 436 359\nL 421 359\nL 421 301\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 464 323\nL 479 323\nL 479 359\nL 464 359\nL 464 323\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 508 347\nL 523 347\nL 523 359\nL 508 359\nL 508 347\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 551 352\nL 566 352\nL 566 359\nL 551 359\nL 551 352\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 91 353\nL 106 353\nL 106 359\nL 91 359\nL 91 353\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 134 347\nL 149 347\nL 149 359\nL 134 359\nL 134 347\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 178 342\nL 193 342\nL 193 359\nL 178 359\nL 178 342\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 221 312\nL 236 312\nL 236 359\nL 221 359\nL 221 312\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 265 308\nL 280 308\nL 280 359\nL 265 359\nL 265 308\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 308 235\nL 323 235\nL 323 359\nL 308 359\nL 308 235\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 352 53\nL 367 53\nL 367 359\nL 352 359\nL 352 53\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 395 41\nL 410 41\nL 410 359\nL 395 359\nL 395 41\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 439 273\nL 454 273\nL 454 359\nL 439 359\nL 439 273\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 482 325\nL 497 325\nL 497 359\nL 482 359\nL 482 325\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 526 347\nL 541 347\nL 541 359\nL 526 359\nL 526 347\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 569 354\nL 584 354\nL 584 359\nL 569 359\nL 569 354\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"77\" y=\"349\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"112\" y=\"344\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4.9</text><text x=\"164\" y=\"341\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"196\" y=\"312\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23.2</text><text x=\"240\" y=\"308\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25.6</text><text x=\"283\" y=\"220\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">76.7</text><text x=\"321\" y=\"117\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">135.6</text><text x=\"364\" y=\"71\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">162.2</text><text x=\"414\" y=\"296\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32.6</text><text x=\"462\" y=\"318\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"504\" y=\"342\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6.4</text><text x=\"547\" y=\"347\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3.3</text><text x=\"87\" y=\"348\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.6</text><text x=\"130\" y=\"342\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">5.9</text><text x=\"182\" y=\"337\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">9</text><text x=\"214\" y=\"307\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26.4</text><text x=\"258\" y=\"303\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28.7</text><text x=\"301\" y=\"230\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">70.7</text><text x=\"339\" y=\"48\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">175.6</text><text x=\"382\" y=\"36\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">182.2</text><text x=\"432\" y=\"268\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.7</text><text x=\"475\" y=\"320\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18.8</text><text x=\"530\" y=\"342\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"565\" y=\"349\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text></svg>",
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

package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/go-analyze/bulk"

	"github.com/go-analyze/charts"
)

/*
Used to render examples of our themes
*/

const (
	chartWidth  = 600
	chartHeight = 400
)

var allThemes = []string{charts.ThemeLight, charts.ThemeDark, charts.ThemeVividLight, charts.ThemeVividDark,
	charts.ThemeAnt, charts.ThemeGrafana, charts.ThemeNatureLight, charts.ThemeNatureDark,
	charts.ThemeRetro, charts.ThemeOcean, charts.ThemeSlate, charts.ThemeGray,
	charts.ThemeWinter, charts.ThemeSpring, charts.ThemeSummer, charts.ThemeFall}

func writeFile(buf []byte, filename string) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, filename)
	return os.WriteFile(file, buf, 0600)
}

func main() {
	for _, theme := range allThemes {
		renderMultiChart(theme)
	}

	renderChartGroup(allThemes, 4, "themes-all.png")
	darkThemes, lightTheme := bulk.SliceSplit(func(v string) bool {
		return charts.GetTheme(v).IsDark()
	}, allThemes)
	renderChartGroup(darkThemes, 3, "themes-dark.png")
	renderChartGroup(lightTheme, 3, "themes-light.png")
	renderChartGroup([]string{charts.ThemeVividLight, charts.ThemeAnt, charts.ThemeRetro,
		charts.ThemeGrafana, charts.ThemeNatureDark, charts.ThemeSlate},
		3, "themes-demo.png")
}

func renderChartGroup(themeNames []string, chartsPerRow int, filename string) {
	lineOpt := charts.NewLineChartOptionWithData([][]float64{
		{120, 132, 101, charts.GetNullValue(), 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
	})
	lineOpt.Title.FontStyle.FontSize = 12
	lineOpt.XAxis.Labels = []string{
		"A", "B", "C", "D", "E", "F", "G",
	}
	lineOpt.Legend.SeriesNames = []string{
		"Email", "Ads", "Direct", "Search",
	}
	lineOpt.Legend.Align = charts.AlignRight
	lineOpt.Legend.Offset = charts.OffsetRight
	lineOpt.LineStrokeWidth = 1.2

	rowCount := int(math.Ceil(float64(len(themeNames)) / float64(chartsPerRow)))
	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        chartWidth * chartsPerRow,
		Height:       chartHeight * rowCount,
	})
	var themeIndex int
	for r := 0; r < rowCount; r++ {
		startY := r * chartHeight
		for c := 0; c < chartsPerRow; c++ {
			startX := c * chartWidth
			if themeIndex >= len(themeNames) {
				p.FilledRect(startX, startY, startX+chartWidth, startY+chartHeight,
					charts.ColorWhite, charts.ColorWhite, 0.0)
				continue
			}
			themeName := themeNames[themeIndex]
			themeIndex++
			lineOpt.Title.Text = "Theme '" + themeName + "'"
			lineOpt.Theme = charts.GetTheme(themeName)
			p.Child(charts.PainterBoxOption(charts.NewBox(startX, startY, startX+chartWidth, startY+chartHeight))).
				LineChart(lineOpt)
		}
	}
	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf, filename); err != nil {
		panic(err)
	}
}

func renderMultiChart(themeName string) {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
	}
	theme := charts.GetTheme(themeName)
	p, err := charts.LineRender(
		values,
		charts.ThemeOptionFunc(theme),
		charts.TitleTextOptionFunc("Theme '"+themeName+"'"),
		charts.PaddingOptionFunc(charts.NewBox(24, 24, 24, 8)),
		charts.XAxisLabelsOptionFunc([]string{
			"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			SeriesNames: []string{
				"Email", "Video Ads", "Direct",
			},
			OverlayChart: charts.Ptr(false),
			Offset: charts.OffsetStr{
				Top:  charts.PositionBottom,
				Left: "20%",
			},
			Symbol: charts.SymbolDot,
		}),
		func(opt *charts.ChartOption) {
			opt.YAxis = []charts.YAxisOption{
				{
					Max:        charts.Ptr(1000.0),
					LabelCount: 6,
				},
			}
			opt.Symbol = charts.SymbolCircle
			opt.LineStrokeWidth = 1.2
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
		},
	)
	if err != nil {
		panic(err)
	}

	pieValues := []float64{
		84358845, 68070697, 58850717, 48059777, 36753736, 19051562, 17947406, 11754004,
		10827529, 10521556, 10467366, 10394055, 9597085, 9104772, 6447710, 5932654,
		5563970, 5428792, 5194336, 3850894, 2857279, 2116792, 1883008, 1373101,
		920701, 660809, 542051,
	}
	pieOpt := charts.PieChartOption{
		Theme: theme.WithBackgroundColor(charts.ColorTransparent),
		SeriesList: charts.NewSeriesListPie(pieValues, charts.PieSeriesOption{
			Label: charts.SeriesLabel{
				Show: charts.Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *charts.LabelStyle) {
					return name + "(" + charts.FormatValueHumanizeShort(val, 2, false) + ")", nil
				},
			},
		}),
		Radius:  "64",
		Padding: charts.NewBoxEqual(20),
		Legend: charts.LegendOption{
			SeriesNames: []string{
				"Germany",
				"France",
				"Italy",
				"Spain",
				"Poland",
				"Romania",
				"Netherlands",
				"Belgium",
				"Czech Republic",
				"Sweden",
				"Portugal",
				"Greece",
				"Hungary",
				"Austria",
				"Bulgaria",
				"Denmark",
				"Finland",
				"Slovakia",
				"Ireland",
				"Croatia",
				"Lithuania",
				"Slovenia",
				"Latvia",
				"Estonia",
				"Cyprus",
				"Luxembourg",
				"Malta",
			},
			Show: charts.Ptr(false),
		},
	}
	pieSum := pieOpt.SeriesList.SumSeries()
	for i := range pieOpt.SeriesList {
		pieOpt.SeriesList[i].Label.Show = charts.Ptr(pieOpt.SeriesList[i].Value/pieSum > 0.04)
		pieOpt.SeriesList[i].Label.FontStyle.FontColor = theme.GetSeriesColor(i)
	}
	p = p.Child(charts.PainterBoxOption(charts.NewBox(200, 0, 600, 200)))
	if err = p.PieChart(pieOpt); err != nil {
		panic(err)
	}
	filename := "theme-"
	if theme.IsDark() {
		filename += "dark-" + themeName + ".png"
	} else {
		filename += "light-" + themeName + ".png"
	}
	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf, filename); err != nil {
		panic(err)
	}
}

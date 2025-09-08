package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example showing how to use the built-in notosans font for Chinese characters.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "chinese-line-chart.png")
	return os.WriteFile(file, buf, 0600)
}

func main() {
	// Use the built-in notosans font which provides better CJK character support
	if err := charts.SetDefaultFont(charts.FontFamilyNotoSans); err != nil {
		panic(err)
	}
	// It's also possible to specify the font on the chart configuration (for example on the title, or legend specifically)

	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	p, err := charts.LineRender(
		values,
		charts.TitleTextOptionFunc("测试"),
		charts.XAxisLabelsOptionFunc([]string{
			"星期一",
			"星期二",
			"星期三",
			"星期四",
			"星期五",
			"星期六",
			"星期日",
		}),
		charts.LegendLabelsOptionFunc([]string{
			"邮件",
			"广告",
			"视频广告",
			"直接访问",
			"搜索引擎",
		}),
	)
	if err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

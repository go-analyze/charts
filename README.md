# go-analyze/charts

[![license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/go-analyze/charts/blob/master/LICENSE)
[![Build Status](https://github.com/go-analyze/charts/workflows/Test/badge.svg)](https://github.com/go-analyze/charts/actions)

Our library focuses on generating beautiful charts and graphs within Go. Graphs are used to show a lot of different types of data, needing to be represented in a unique in order to convey the meaning behind the data. This Go module attempts to use sophisticated defaults to try and render this data in a simple way, while still offering intuitive options to update the graph rendering as you see fit.

## Current Project Status

Forked from [vicanso/go-charts](https://github.com/vicanso/go-charts) and the archived [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart), our project introduces enhancements for rendering challenging datasets. We aim to build upon their solid foundation to offer a more versatile and user-friendly charting solution.

### API Stability

We're committed to refining the API, incorporating feedback and new ideas to enhance flexibility and ease of use.

Until the `v1.0.0` release, API changes should be anticipated. We detail needed API changes on our wiki [Version Migration Guide](https://github.com/go-analyze/charts/wiki/Version-Migration-Guide).

### Changes

Notable early improvements in our fork include:

* **Axis Improvements:** Significant enhancements to axis rendering, data range selection, and configuration simplification were made in PR [#3](https://github.com/go-analyze/charts/pull/3).
* **Theming:** In PR [#4](https://github.com/go-analyze/charts/pull/4) (and some subsequent changes) we introduced `vivid-light` and `vivid-dark` themes for more vibrant visualizations, alongside API changes for greater theme and font control. Long term we plan to make themes easier to mutate and define.
* **Configuration Simplification:** PR [#5](https://github.com/go-analyze/charts/pull/5) began our effort to streamline chart configuration, making names more descriptive and specific while focusing on a theme-centric approach. Documentation on configuration and use is also being improved. (See also [#15](https://github.com/go-analyze/charts/pull/15), [#20](https://github.com/go-analyze/charts/pull/20))
* **Expanded Testing:** Ongoing test coverage expansions have led to bug discoveries and fixes. This will continue to help ensure that our charts render perfectly for a wide range of configurations and use.

Our library is a work in progress, aiming to become a standout choice for Go developers seeking powerful, yet easy-to-use charting tools. We welcome contributions and feedback as we continue to enhance our library's functionality, configurability, and reliability.

#### `wcharczuk/go-chart` Changes

If you are a former user of `wcharczuk/go-chart`, you should be able to use this project with reasonable changes. The `wcharczuk/go-chart` project was forked under our `chartdraw` package. Any code changes necessary are documented on our [wcharczuk/go‐chart Migration Guide](https://github.com/go-analyze/charts/wiki/wcharczuk-go%E2%80%90chart-Migration-Guide).

## Functionality

### Chart Types

These chart types are supported: `line`, `bar`, `horizontal bar`, `pie`, `radar` or `funnel` and `table`.

Please see the [./examples/](./examples/) directory and the [README](./examples/README.md) within it to see a variety of implementations of our different chart types and configurations.

### Themes

Our built in themes are: `light`, `dark`, `vivid-light`, `vivid-dark`, `ant`, `grafana`

<p align="center">
    <img src="./assets/themes.png" alt="themes">
</p>

### Line Chart

<img src="./assets/chart-line.png" alt="Line Chart">

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	// values specified where the first index is for each data series or source, and the second index is for each sample.
	values := [][]float64{
		{
			120,
			132,
			101,
			134,
			90,
			230,
			210,
		},
		{
			// specify values for additional data series
		},
	}
	p, err := charts.LineRender(
		values,
		charts.TitleTextOptionFunc("Line Chart Demo"),
		charts.XAxisDataOptionFunc([]string{
			"Mon",	// notice the 7 labels here match to the 7 samples above
			"Tue",
			"Wed",
			"Thu",
			"Fri",
			"Sat",
			"Sun",
		}),
		charts.LegendLabelsOptionFunc([]string{
			"Email",
			"Search Engine",
		}),
		// other options as desired...
```

### Bar Chart

<img src="./assets/chart-bar.png" alt="Bar Chart">

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	// values specified where the first index is for each data series or source, and the second index is for each sample.
	values := [][]float64{
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
			// snip...	
		},
	}
	p, err := charts.BarRender(
		values,
		charts.XAxisDataOptionFunc([]string{
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
		charts.LegendLabelsOptionFunc([]string{
			"Rainfall",
			"Evaporation",
		}),
		charts.MarkLineOptionFunc(0, charts.SeriesMarkDataTypeAverage),
		charts.MarkPointOptionFunc(0, charts.SeriesMarkDataTypeMax,
			charts.SeriesMarkDataTypeMin),
		// custom option func
		func(opt *charts.ChartOption) {
			opt.SeriesList[1].MarkPoint = charts.NewMarkPoint(
				charts.SeriesMarkDataTypeMax,
				charts.SeriesMarkDataTypeMin,
			)
			opt.SeriesList[1].MarkLine = charts.NewMarkLine(
				charts.SeriesMarkDataTypeAverage,
			)
		},
		// other options as desired...
```

### Horizontal Bar Chart

<img src="./assets/chart-horizontal-bar.png" alt="Horizontal Bar Chart">

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	values := [][]float64{
		{
			18203,
			23489,
			29034,
			104970,
			131744,
			630230,
		},
		{
			// snip...	
		},
	}
	p, err := charts.HorizontalBarRender(
		values,
		charts.TitleTextOptionFunc("World Population"),
		charts.PaddingOptionFunc(charts.Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		charts.LegendLabelsOptionFunc([]string{
			"2011",
			"2012",
		}),
		charts.YAxisDataOptionFunc([]string{
			"Brazil",
			"Indonesia",
			"USA",
			"India",
			"China",
			"World",
		}),
		// other options as desired...
```

### Pie Chart

<img src="./assets/chart-pie.png" alt="Pie Chart">

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	values := []float64{
		1048,
		735,
		580,
		484,
		300,
	}
	p, err := charts.PieRender(
		values,
		charts.TitleOptionFunc(charts.TitleOption{
			Text:    "Rainfall vs Evaporation",
			Subtext: "Fake Data",
			Offset:  charts.OffsetCenter,
		}),
		charts.LegendOptionFunc(charts.LegendOption{
			Data: []string{
				"Search Engine",
				"Direct",
				"Email",
				"Union Ads",
				"Video Ads",
			},
		}),
		charts.SeriesShowLabel(true),
		// other options as desired...
```

### Radar Chart

<img src="./assets/chart-radar.png" alt="Radar Chart">

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	values := [][]float64{
		{
			4200,
			3000,
			20000,
			35000,
			50000,
			18000,
		},
		{
			// snip...
		},
	}
	p, err := charts.RadarRender(
		values,
		charts.TitleTextOptionFunc("Basic Radar Chart"),
		charts.LegendLabelsOptionFunc([]string{
			"Allocated Budget",
			"Actual Spending",
		}),
		charts.RadarIndicatorOptionFunc([]string{
			"Sales",
			"Administration",
			"Information Technology",
			"Customer Support",
			"Development",
			"Marketing",
		}, []float64{
			6500,
			16000,
			30000,
			38000,
			52000,
			25000,
		}),
		// other options as desired...
```

### Table

<img src="./assets/chart-table.png" alt="Table">

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	header := []string{
		"Name",
		"Age",
		"Address",
		"Tag",
		"Action",
	}
	data := [][]string{
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
	}
	spans := map[int]int{
		0: 2,
		1: 1,
		2: 3,
		3: 2,
		4: 2,
	}
	p, err := charts.TableRender(
		header,
		data,
		spans,
	)
	// snip...
```

### Funnel Chart

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	values := []float64{
		100,
		80,
		60,
		40,
		20,
	}
	p, err := charts.FunnelRender(
		values,
		charts.TitleTextOptionFunc("Funnel"),
		charts.LegendLabelsOptionFunc([]string{
			"Show",
			"Click",
			"Visit",
			"Inquiry",
			"Order",
		}),
	)
	// snip...
```

### ECharts Render

```go
import (
	"github.com/go-analyze/charts"
)

func main() {
	buf, err := charts.RenderEChartsToPNG(`{
		"title": {
			"text": "Line"
		},
		"xAxis": {
			"data": ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
		},
		"series": [
			{
				"data": [150, 230, 224, 218, 135, 147, 260]
			}
		]
	}`)
	// snip...
```

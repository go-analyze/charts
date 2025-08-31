# go-analyze/charts Examples

Examples are our primary method for demonstrating the starting point our use and configuration. Find an example close to your use case and use it as a starting place to find the relevant API's. Our API's include godocs, review the descriptions of each field to understand possible options.

Our library offers two primary ways to configure charts:
1. The "Painter API" allows you to initialize a "Painter" using `NewPainter`, then build chart configuration structs and apply them with function calls on the `Painter`. This API is designed to be easy to learn and navigate, with fields being named and formed to be natural for the chart type.
2. Also offered is `ChartOption`, which is built providing functions to modify the very generic chart struct. This API can be useful when representing the same data, and similar chart configuration, but changing the chart types.

For most use cases the `Painter API` is recommended. In our examples we demonstrate rendering chart types using both methods, see either examples in the [1-Painter](./1-Painter) directory or [2-OptionFun](./2-OptionFunc).

## `Painter` API Example List

* [bar_chart-1-basic](./1-Painter/bar_chart-1-basic) - Basic bar chart.
* [bar_chart-2-size_margin](./1-Painter/bar_chart-2-size_margin) - Showing the visual impact of different bar sizes and margins.
* [bar_chart-3-label_position-round_caps](./1-Painter/bar_chart-3-label_position-round_caps) - Showing the different label positions and rounded caps.
* [bar_chart-4-mark](./1-Painter/bar_chart-4-mark) - Bar chart with included mark points and mark lines.
* [bar_chart-5-stacked](./1-Painter/bar_chart-5-stacked) - A bar chart with "Stacked" series enabled, collapsing the bars into a single layered bar.
* [candlestick_chart-1-basic](./1-Painter/candlestick_chart-1-basic) - Basic candlestick chart.
* [candlestick_chart-2-multiple_series](./1-Painter/candlestick_chart-2-multiple_series) - Candlestick chart with multiple series and varied candle styles.
* [candlestick_chart-3-bollinger_bands](./1-Painter/candlestick_chart-3-bollinger_bands) - Candlestick chart with Bollinger Bands overlaid.
* [candlestick_chart-4-patterns](./1-Painter/candlestick_chart-4-patterns) - Candlestick chart highlighting core and custom candlestick patterns.
* [candlestick_chart-5-aggregation](./1-Painter/candlestick_chart-5-aggregation) - Candlestick data aggregation: 1-minute vs 5-minute with two stacked charts.
* [doughnut_chart-1-basic](./1-Painter/doughnut_chart-1-basic) - Basic doughnut chart, a variation on a pie chart with the center opened up, allowing labels or other values to be put in the middle to save space.
* [doughnut_chart-2-styles](./1-Painter/doughnut_chart-2-styles) - A variety of styles for doughnut charts shown.
* [funnel_chart-1-basic](./1-Painter/funnel_chart-1-basic) - Basic funnel chart.
* [heat_map-1-basic](./1-Painter/heat_map-1-basic) - Basic heat map chart.
* [horizontal_bar_chart-1-basic](./1-Painter/horizontal_bar_chart-1-basic) - Basic horizontal bar chart.
* [horizontal_bar_chart-2-size_margin](./1-Painter/horizontal_bar_chart-2-size_margin) - Showing the visual impact of different bar sizes and margins.
* [horizontal_bar_chart-3-mark](./1-Painter/horizontal_bar_chart-3-mark) - Horizontal bar chart with included mark lines.
* [horizontal_bar_chart-4-stacked](./1-Painter/horizontal_bar_chart-4-stacked) - A horizontal bar chart with "Stacked" series, collapsing the bars into a single layered bar.
* [line_chart-1-basic](./1-Painter/line_chart-1-basic) - Basic line chart with some simple styling changes and a demonstration of `null` values.
* [line_chart-2-symbols](./1-Painter/line_chart-2-symbols) - Basic line chart which sets a different symbol for each series item.
* [line_chart-3-smooth](./1-Painter/line_chart-3-smooth) - Basic line chart with thick smooth lines drawn.
* [line_chart-4-mark](./1-Painter/line_chart-4-mark) - Line chart with included mark points and mark lines.
* [line_chart-5-area](./1-Painter/line_chart-5-area) - Line chart with the area below the line shaded.
* [line_chart-6-stacked](./1-Painter/line_chart-6-stacked) - Line chart with "Stacked" series enabled, making each series a layer on the chart and the top line showing the sum.
* [line_chart-7-boundary_gap](./1-Painter/line_chart-7-boundary_gap) - Showing the visual difference on the line chart of enabling or disabling the x-axis boundary gap.
* [line_chart-8-dual_y_axis](./1-Painter/line_chart-8-dual_y_axis) - Basic line chart with two series, one rendered to the left axis and one to a second y axis on the right.
* [line_chart-9-custom](./1-Painter/line_chart-9-custom) - Line chart with dense data and most default rendering disabled, instead rendering labels manually on the Painter.
* [multiple_charts-1](./1-Painter/multiple_charts-1) - Example of manually building a painter so that you can render 4 charts on the same image.
* [multiple_charts-2](./1-Painter/multiple_charts-2) - Shows how to use a single set of data and demonstrate it with multiple chart types.
* [pie_chart-1-basic](./1-Painter/pie_chart-1-basic) - Pie chart with a variety of customization demonstrated including positioning the legend in the bottom right corner.
* [pie_chart-2-series_radius](./1-Painter/pie_chart-2-series_radius) - Pie chart which varies the series radius by the percentage of the series.
* [pie_chart-3-gap](./1-Painter/pie_chart-3-gap) - Pie chart with segment gaps between each slice.
* [radar_chart-1-basic](./1-Painter/radar_chart-1-basic) - Basic radar chart.
* [scatter_chart-1-basic](./1-Painter/scatter_chart-1-basic) - Basic scatter chart with some simple styling changes and a demonstration of `null` values.
* [scatter_chart-2-symbols](./1-Painter/scatter_chart-2-symbols) - Basic scatter chart showing per-series symbols.
* [scatter_chart-3-dense_data](./1-Painter/scatter_chart-3-dense_data) - Scatter chart with dense data, trend lines, and more custom styling configured.
* [table-1](./1-Painter/table-1) - Table with a variety of table specific configuration and styling demonstrated.

## `ChartOption` / `OptionFunc` Example List

* [bar_chart-1-basic](./2-OptionFunc/bar_chart-1-basic) - Bar chart with included mark points and mark lines.
* [candlestick_chart-1-basic](./2-OptionFunc/candlestick_chart-1-basic) - Basic candlestick chart using `ChartOption`.
* [doughnut_chart-1-basic](./2-OptionFunc/doughnut_chart-1-basic) - Doughnut chart, a variation on a pie chart with the center opened up, allowing labels or other values to be put in the middle to save space.
* [chinese](./2-OptionFunc/chinese) - Line chart with chinese characters that uses a custom font (must be downloaded by user, see comment in code).
* [funnel_chart-1-basic](./2-OptionFunc/funnel_chart-1-basic) - Basic funnel chart.
* [horizontal_bar_chart-1-basic](./2-OptionFunc/horizontal_bar_chart-1-basic) - Basic horizontal bar chart.
* [line_chart-1-basic](./2-OptionFunc/line_chart-1-basic) - Basic line chart with some simple styling changes and a demonstration of `null` values.
* [line_chart-2-dense_data](./2-OptionFunc/line_chart-2-dense_data) - Line chart with dense data and more custom styling configured.
* [line_chart-3-area](./2-OptionFunc/line_chart-3-area) - Line chart with the area below the line shaded.
* [multiple_charts-1](./2-OptionFunc/multiple_charts-1) - Combining two charts together by writting one chart over the other.
* [multiple_charts-2](./2-OptionFunc/multiple_charts-2) - An alternative API for overlaying two charts together in the same image.
* [pie_chart-1-basic](./2-OptionFunc/pie_chart-1-basic) - Pie chart with a variety of customization demonstrated including positioning the legend in the bottom right corner.
* [radar_chart-1-basic](./2-OptionFunc/radar_chart-1-basic) - Basic radar chart.
* [scatter_chart-1-basic](./2-OptionFunc/scatter_chart-1-basic) - Basic scatter chart with some simple styling changes and a demonstration of `null` values.
* [table-1](./2-OptionFunc/table-1) - Table with a variety of table specific configuration and styling demonstrated.
* [web-1](./2-OptionFunc/web-1) - Hosts an example http server which will render the charts to the web page.

## chartdraw/examples

The examples in the root [examples directory](https://github.com/go-analyze/charts/tree/main/examples) serves as our primary examples. These examples are the best representation of our library and what we are aiming to support and improve. If you're intereted in exploring the underline `chartdraw` implementation you can also check out [chartdraw/examples](https://github.com/go-analyze/charts/tree/main/chartdraw/examples). These are examples from the implementation based off [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart). These examples illustrate an alternative styling, which over time we aim to merge and unify with our `charts` package API.

If you find you prefer the `chartdraw` styling, configuration schema, or anything else, please open an [Issue](https://github.com/go-analyze/charts/issues) so that we can make sure we retain the best of both implementations as we seek unifying our API.

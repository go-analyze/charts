# go-analyze/charts Examples

Examples are our primary method for demonstrating the starting point our use and configuration. For more advanced configuration review the other fields and descriptions of the structs used in our examples.

## Example List

* [bar_chart-1](./bar_chart-1) - Bar chart with included mark points and mark lines.
* [chinese](./chinese) - Line chart with chinese characters that uses a custom font (must be downloaded by user, see comment in code).
* [funnel_chart-1](./funnel_chart-1) - Basic funnel chart.
* [horizontal_bar_chart-1](./horizontal_bar_chart-1) - Basic horizontal bar chart.
* [line_chart-1](./line_chart-1) - Basic line chart with some simple styling changes and a demonstration of `null` values.
* [line_chart-2](./line_chart-2) - Line chart with dense data and more custom styling configured.
* [line_chart-area](./line_chart-area) - Example line chart with the area below the line shaded.
* [multiple_charts-1](./multiple_charts-1) - Example of manually building a painter so that you can render 4 charts on the same image.
* [multiple_charts-2](./multiple_charts-2) - Combining two charts together by writting one chart over the other.
* [multiple_charts-3](./multiple_charts-3) - An alternative API for overlaying two charts together in the same image.
* [pie_chart-1](./pie_chart-1) - Pie chart with a variety of customization demonstrated including positioning the legend in the bottom right corner.
* [radar_chart-1](./radar_chart-1) - Basic radar chart.
* [table-1](./table-1) - Table with a variety of table specific configuration and styling demonstrated.
* [web-1](./web-1) - Hosts an example http server which will render the charts to the web page.

## chartdraw/examples

The examples in the root [examples directory](https://github.com/go-analyze/charts/tree/main/examples) serves as our primary examples. These examples are the best representation of our library and what we are aiming to support and improve. If you're intereted in exploring the underline `chartdraw` implementation you can also check out [chartdraw/examples](https://github.com/go-analyze/charts/tree/main/chartdraw/examples). These are examples from the implementation based off [wcharczuk/go-chart](https://github.com/wcharczuk/go-chart). These examples illustrate an alternative styling, which over time we aim to merge and unify with our `charts` package API.

If you find you prefer the `chartdraw` styling, configuration schema, or anything else, please open an [Issue](https://github.com/go-analyze/charts/issues) so that we can make sure we retain the best of both implementations as we seek unifying our API.

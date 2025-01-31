package chartdraw

// TODO - remove internal defaults from public API
const (
	// Deprecated: DefaultChartHeight is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultChartHeight is the default chart height.
	DefaultChartHeight = 400
	// Deprecated: DefaultChartWidth is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultChartWidth is the default chart width.
	DefaultChartWidth = 1024
	// Deprecated: DefaultStrokeWidth is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultStrokeWidth is the default chart stroke width.
	DefaultStrokeWidth = 0.0
	// Deprecated: DefaultDotWidth is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultDotWidth is the default chart dot width.
	DefaultDotWidth = 0.0
	// Deprecated: DefaultSeriesLineWidth is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultSeriesLineWidth is the default line width.
	DefaultSeriesLineWidth = 1.0
	// Deprecated: DefaultAxisLineWidth is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultAxisLineWidth is the line width of the axis lines.
	DefaultAxisLineWidth = 1.0
	// Deprecated: DefaultDPI is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultDPI is the default dots per inch for the chart.
	DefaultDPI = 92.0
	// Deprecated: DefaultFontSize is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultFontSize is the default font size.
	DefaultFontSize = 10.0
	// Deprecated: DefaultTitleFontSize is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultTitleFontSize is the default title font size.
	DefaultTitleFontSize = 18.0
	// Deprecated: DefaultAnnotationDeltaWidth is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultAnnotationDeltaWidth is the width of the left triangle out of annotations.
	DefaultAnnotationDeltaWidth = 10
	// Deprecated: DefaultAnnotationFontSize is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultAnnotationFontSize is the font size of annotations.
	DefaultAnnotationFontSize = 10.0
	// Deprecated: DefaultAxisFontSize is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultAxisFontSize is the font size of the axis labels.
	DefaultAxisFontSize = 10.0
	// Deprecated: DefaultTitleTop is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultTitleTop is the default distance from the top of the chart to put the title.
	DefaultTitleTop = 10

	// Deprecated: DefaultBackgroundStrokeWidth is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultBackgroundStrokeWidth is the default stroke on the chart background.
	DefaultBackgroundStrokeWidth = 0.0
	// Deprecated: DefaultCanvasStrokeWidth is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultCanvasStrokeWidth is the default stroke on the chart canvas.
	DefaultCanvasStrokeWidth = 0.0

	// Deprecated: DefaultLineSpacing is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultLineSpacing is the default vertical distance between lines of text.
	DefaultLineSpacing = 5

	// Deprecated: DefaultYAxisMargin is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultYAxisMargin is the default distance from the right of the canvas to the y-axis labels.
	DefaultYAxisMargin = 10
	// Deprecated: DefaultXAxisMargin is deprecated, it's not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	// DefaultXAxisMargin is the default distance from bottom of the canvas to the x-axis labels.
	DefaultXAxisMargin = 10

	// Deprecated: DefaultVerticalTickHeight is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultVerticalTickHeight is half the margin.
	DefaultVerticalTickHeight = DefaultXAxisMargin >> 1
	// Deprecated: DefaultHorizontalTickWidth is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultHorizontalTickWidth is half the margin.
	DefaultHorizontalTickWidth = DefaultYAxisMargin >> 1

	// Deprecated: DefaultTickCountSanityCheck is deprecated, it's not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	// DefaultTickCountSanityCheck is a hard limit on number of ticks to prevent infinite loops.
	DefaultTickCountSanityCheck = 1 << 10 //1024

	// Deprecated: DefaultMinimumTickHorizontalSpacing is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultMinimumTickHorizontalSpacing is the minimum distance between horizontal ticks.
	DefaultMinimumTickHorizontalSpacing = 20
	// Deprecated: DefaultMinimumTickVerticalSpacing is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultMinimumTickVerticalSpacing is the minimum distance between vertical ticks.
	DefaultMinimumTickVerticalSpacing = 20

	// DefaultDateFormat is the default date format.
	DefaultDateFormat = "2006-01-02"
	// DefaultDateHourFormat is the date format for hour timestamp formats.
	DefaultDateHourFormat = "01-02 3PM"
	// DefaultDateMinuteFormat is the date format for minute range timestamp formats.
	DefaultDateMinuteFormat = "01-02 3:04PM"
	// Deprecated: DefaultFloatFormat is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultFloatFormat is the default float format.
	DefaultFloatFormat = "%.2f"
	// Deprecated: DefaultPercentValueFormat is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultPercentValueFormat is the default percent format.
	DefaultPercentValueFormat = "%0.2f%%"

	// Deprecated: DefaultBarSpacing is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultBarSpacing is the default pixel spacing between bars.
	DefaultBarSpacing = 100
	// Deprecated: DefaultBarWidth is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultBarWidth is the default pixel width of bars in a bar chart.
	DefaultBarWidth = 50
)

var (
	// Deprecated: DefaultAnnotationPadding is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultAnnotationPadding is the padding around an annotation.
	DefaultAnnotationPadding = Box{Top: 5, Left: 5, Right: 5, Bottom: 5}

	// Deprecated: DefaultBackgroundPadding is deprecated, it's not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	// DefaultBackgroundPadding is the default canvas padding config.
	DefaultBackgroundPadding = Box{Top: 5, Left: 5, Right: 5, Bottom: 5}
)

const (
	// ContentTypePNG is the png mime type.
	ContentTypePNG = "image/png"

	// ContentTypeSVG is the svg mime type.
	ContentTypeSVG = "image/svg+xml"
)

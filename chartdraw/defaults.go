package chartdraw

const (
	// DefaultChartHeight is the default chart height.
	DefaultChartHeight = 400
	// DefaultChartWidth is the default chart width.
	DefaultChartWidth = 1024

	// Deprecated: DefaultStrokeWidth is not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	DefaultStrokeWidth = 0.0
	// Deprecated: DefaultDotWidth is not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	DefaultDotWidth = 0.0

	// DefaultSeriesLineWidth is the default line width.
	DefaultSeriesLineWidth = 1.0
	// DefaultAxisLineWidth is the line width of the axis lines.
	DefaultAxisLineWidth = 1.0

	// Deprecated: DefaultDPI is not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	DefaultDPI = 92.0

	// DefaultFontSize is the default font size.
	DefaultFontSize = 10.0
	// DefaultTitleFontSize is the default title font size.
	DefaultTitleFontSize = 18.0

	// Deprecated: DefaultAnnotationDeltaWidth is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultAnnotationDeltaWidth = 10

	// DefaultAnnotationFontSize is the font size of annotations.
	DefaultAnnotationFontSize = 10.0
	// DefaultAxisFontSize is the font size of the axis labels.
	DefaultAxisFontSize = 10.0

	// Deprecated: DefaultTitleTop is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultTitleTop = 10

	// Deprecated: DefaultBackgroundStrokeWidth is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultBackgroundStrokeWidth = 0.0
	// Deprecated: DefaultCanvasStrokeWidth is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultCanvasStrokeWidth = 0.0

	// DefaultLineSpacing is the default vertical distance between lines of text.
	DefaultLineSpacing = 5

	// Deprecated: DefaultYAxisMargin is not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	DefaultYAxisMargin = 10
	// Deprecated: DefaultXAxisMargin is not expected to be used externally. If you use this field,
	// open a new issue to prevent it from being made internal.
	DefaultXAxisMargin = 10

	// Deprecated: DefaultVerticalTickHeight is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultVerticalTickHeight = DefaultXAxisMargin >> 1
	// Deprecated: DefaultHorizontalTickWidth is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultHorizontalTickWidth = DefaultYAxisMargin >> 1

	// Deprecated: DefaultTickCountSanityCheck is not expected to be used externally. If you use this
	// field, open a new issue to prevent it from being made internal.
	DefaultTickCountSanityCheck = 1 << 10 //1024

	// Deprecated: DefaultMinimumTickHorizontalSpacing is not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	DefaultMinimumTickHorizontalSpacing = 20
	// Deprecated: DefaultMinimumTickVerticalSpacing is not expected to be used externally. If you
	// use this field, open a new issue to prevent it from being made internal.
	DefaultMinimumTickVerticalSpacing = 20

	// DefaultDateFormat is the default date format.
	DefaultDateFormat = "2006-01-02"
	// DefaultDateHourFormat is the date format for hour timestamp formats.
	DefaultDateHourFormat = "01-02 3PM"
	// DefaultDateMinuteFormat is the date format for minute range timestamp formats.
	DefaultDateMinuteFormat = "01-02 3:04PM"
	// DefaultFloatFormat is the default float format.
	DefaultFloatFormat = "%.2f"
	// DefaultPercentValueFormat is the default percent format.
	DefaultPercentValueFormat = "%0.2f%%"

	// DefaultBarSpacing is the default pixel spacing between bars.
	DefaultBarSpacing = 100
	// DefaultBarWidth is the default pixel width of bars in a bar chart.
	DefaultBarWidth = 50
)

var (
	// DefaultAnnotationPadding is the padding around an annotation.
	DefaultAnnotationPadding = Box{Top: 5, Left: 5, Right: 5, Bottom: 5}

	// DefaultBackgroundPadding is the default canvas padding config.
	DefaultBackgroundPadding = Box{Top: 5, Left: 5, Right: 5, Bottom: 5}
)

const (
	// ContentTypePNG is the png mime type.
	ContentTypePNG = "image/png"

	// ContentTypeSVG is the svg mime type.
	ContentTypeSVG = "image/svg+xml"
)

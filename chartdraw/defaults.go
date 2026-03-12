package chartdraw

const (
	// DefaultChartHeight is the default chart height.
	DefaultChartHeight = 400
	// DefaultChartWidth is the default chart width.
	DefaultChartWidth = 1024

	defaultStrokeWidth = 0.0
	defaultDotWidth    = 0.0

	// DefaultSeriesLineWidth is the default line width.
	DefaultSeriesLineWidth = 1.0
	// DefaultAxisLineWidth is the line width of the axis lines.
	DefaultAxisLineWidth = 1.0

	defaultDPI = 92.0

	// DefaultFontSize is the default font size.
	DefaultFontSize = 10.0
	// DefaultTitleFontSize is the default title font size.
	DefaultTitleFontSize = 18.0

	defaultAnnotationDeltaWidth = 10

	// DefaultAnnotationFontSize is the font size of annotations.
	DefaultAnnotationFontSize = 10.0
	// DefaultAxisFontSize is the font size of the axis labels.
	DefaultAxisFontSize = 10.0

	defaultTitleTop = 10

	// DefaultLineSpacing is the default vertical distance between lines of text.
	DefaultLineSpacing = 5

	defaultYAxisMargin = 10
	defaultXAxisMargin = 10

	defaultVerticalTickHeight  = defaultXAxisMargin >> 1
	defaultHorizontalTickWidth = defaultYAxisMargin >> 1

	defaultTickCountSanityCheck = 1 << 10 //1024

	defaultMinimumTickHorizontalSpacing = 20
	defaultMinimumTickVerticalSpacing   = 20

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

# ECharts Option Spec Reference

Reference for the **ECharts-compatible option format this library is built against**.

This is not the Apache ECharts option schema itself. The library accepts a restricted subset of ECharts-style JSON plus several library-specific fields and transformations. Fields not listed below are ignored or have no rendering effect.

The Apache ECharts reference is:

https://echarts.apache.org/en/option.html

## Compatibility model

The accepted JSON is best understood as:

**ECharts-style syntax + library-specific extensions + implementation-specific transformations.**

It is not safe to assume that an option documented by Apache ECharts will render identically, or even parse successfully, through this adapter.

In particular, current Apache ECharts supports substantially more axis types, components, data encodings, styling, and series options than this adapter.

## Top-level option object

| Field             | Type                   | Behavior                                                                                                                         |
| ----------------- | ---------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| `type`            | string                 | Parsed, but overwritten by `RenderEChartsToPNG/JPG/SVG`.                                                                         |
| `theme`           | string                 | Passed to `GetTheme`; unknown/empty names fall back to the default theme.                                                        |
| `fontFamily`      | string                 | Global chart font selection. If absent, the title text-style font family is used as a fallback when constructing the chart font. |
| `padding`         | number | array<number> | Library-level chart padding. See Padding. Values are parsed as integers.                                                         |
| `box`             | object                 | Library-level canvas/drawing box. See Box.                                                                                       |
| `width`           | int                    | Chart width.                                                                                                                     |
| `height`          | int                    | Chart height.                                                                                                                    |
| `backgroundColor` | string                 | Applied to the selected theme as the chart background color.                                                                     |
| `title`           | object                 | See Title.                                                                                                                       |
| `legend`          | object                 | See Legend.                                                                                                                      |
| `xAxis`           | object | array<object> | Parsed as a list; only the first x-axis entry is used.                                                                           |
| `yAxis`           | object | array<object> | Parsed as a list; all entries are converted to y-axis options.                                                                   |
| `radar`           | object                 | Radar indicator configuration.                                                                                                   |
| `series`          | object | array<object> | Parsed as a series list. Single-object input is normalized to a one-element list.                                                |
| `children`        | array<object>          | Recursively converted to child `ChartOption` objects.                                                                            |

The `backgroundColor`, `fontFamily`, `padding`, `box`, and `children` fields are library extensions rather than standard Apache ECharts top-level equivalents.

## Padding

`padding` accepts:

* one integer: all four sides;
* two integers: `[vertical, horizontal]`;
* three integers: `[top, horizontal, bottom]`;
* four integers: `[top, right, bottom, left]`.

The implementation converts the values directly into the library's `Box` representation. Floating-point padding values and arbitrary CSS-style values are not part of this adapter's accepted representation.

## Box

`box` is a library-specific object:

```json
{
  "top": 10,
  "bottom": 10,
  "left": 10,
  "right": 10,
  "isSet": true
}
```

The fields are:

| Field    | Type |
| -------- | ---- |
| `top`    | int  |
| `bottom` | int  |
| `left`   | int  |
| `right`  | int  |
| `isSet`  | bool |

`box` is not Apache ECharts' standard top-level layout syntax.

## Position (`EChartsPosition`)

Used by `title.left`, `title.top`, `legend.left`, and `legend.top`.

Accepted forms include string values such as:

```json
"left"
"center"
"right"
"20"
"20%"
```

and numeric JSON values such as:

```json
20
-10
1.5
```

The JSON decoder converts any numeric value (including negative and decimal literals) into its string representation. The native layer then interprets the resulting string as either a keyword, a pixel value, or a percentage.

Only `left` and `top` are represented by the adapter for title and legend positioning. Other ECharts positioning fields such as `right` and `bottom` are not part of this adapter.

## Title

| Field             | Type     | Rendering behavior                                                                    |
| ----------------- | -------- | ------------------------------------------------------------------------------------- |
| `show`            | bool     | Show/hide title. Defaults through the library when omitted.                           |
| `text`            | string   | Main title text.                                                                      |
| `subtext`         | string   | Subtitle text.                                                                        |
| `left`, `top`     | position | Title placement.                                                                      |
| `textStyle`       | object   | See Text Style.                                                                       |
| `subtextStyle`    | object   | See Text Style.                                                                       |
| `backgroundColor` | string   | Accepted but not rendered.                                                            |
| `borderColor`     | string   | Applied as the title border color; a default border stroke width is enabled when set. |

Apache ECharts itself supports many additional title properties including `right`, `bottom`, padding, border width/radius, shadows, rich text, alignment, and more; those are not represented by this adapter.

## Text Style

Used by title, subtitle, and legend text.

| Field        | Type   | Rendering behavior |
| ------------ | ------ | ------------------ |
| `color`      | string | Text color.        |
| `fontFamily` | string | Font family.       |
| `fontSize`   | number | Font size.         |

When a component's text style specifies no font family, the chart-level `fontFamily` (or title text-style family) is applied as a fallback. The same fallback covers x- and y-axis label fonts.

The adapter does not expose the broader ECharts text-style model such as `fontWeight`, `fontStyle`, text borders, shadows, rich text, overflow, etc.

## Legend

| Field             | Type           | Rendering behavior                                                              |
| ----------------- | -------------- | -------------------------------------------------------------------------------- |
| `show`            | bool           | Show/hide legend.                                                               |
| `data`            | array<string>  | Legend series names.                                                            |
| `align`           | string         | Horizontal legend alignment.                                                    |
| `orient`          | string         | Legend orientation, e.g. `"horizontal"` or `"vertical"`.                        |
| `padding`         | number | array | Space around the legend.                                                        |
| `left`, `top`     | position       | Legend placement.                                                               |
| `textStyle`       | object         | Color, font family, font size.                                                  |
| `backgroundColor` | string         | Accepted but not rendered.                                                      |
| `borderColor`     | string         | Applied as legend border color; a default border stroke width is enabled when set. |

The current Apache ECharts legend component has many additional options, including scrollable legends, selectors, icons, selection state, right/bottom placement, and more. They are not implemented here.

## X Axis

`xAxis` may be a single object or an array. The adapter parses all entries, but **only the first x-axis configuration is rendered**.

| Field         | Type          | Rendering behavior                                                                                                    |
| ------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `type`        | string        | `"value"` triggers the adapter's horizontal-chart transformation for bar and violin series. Other values have no axis-type semantics. |
| `boundaryGap` | bool          | Passed through; when omitted, the adapter explicitly forces `true`.                                                     |
| `splitNumber` | int           | Requested category-axis label count.                                                                                    |
| `data`        | array<string> | Category labels.                                                                                                        |
| `name`        | string        | Axis title; `nameTextStyle.color`/`.fontSize` map to the title font style.                                              |
| `axisLabel`   | object        | See Axis Label.                                                                                                         |
| `axisLine`    | object        | Color/opacity are applied. `show` and `width` are accepted but not rendered (see Axis Line).                            |

### Important `type` behavior

This adapter does **not** implement ECharts' general axis-type semantics.

When any x-axis entry has:

```json
"type": "value"
```

the adapter sets an internal horizontal-chart flag and converts:

* every `bar` series to the library's `horizontalBar` type;
* every `violin` series to the library's `horizontalViolin` type.

Thus `"value"` means more than "use a numeric ECharts axis" in this implementation; it is also an orientation switch for bars and violins.

Current Apache ECharts supports `value`, `time`, `category`, and `log` axis types, but those semantics are not implemented here.

### `boundaryGap`

If omitted, the adapter sets:

```go
BoundaryGap = true
```

unconditionally for the parsed x-axis. This is an **implementation default**, not a universal statement about Apache ECharts.

## Y Axis

`yAxis` may be a single object or an array. Unlike `xAxis`, every parsed y-axis entry becomes a native y-axis option.

| Field            | Type          | Rendering behavior                                                        |
| ---------------- | ------------- | --------------------------------------------------------------------------- |
| `min`            | number        | Explicit minimum.                                                           |
| `max`            | number        | Explicit maximum.                                                           |
| `data`           | array<string> | Category labels.                                                            |
| `name`           | string        | Axis title; `nameTextStyle.color`/`.fontSize` map to the title font style.  |
| `splitNumber`    | int           | Requested y-axis label count.                                               |
| `splitLine.show` | bool          | Split line visibility. True unless explicitly false, to match ECharts.      |
| `position`       | string        | `left`/`right` honored where supported. A secondary axis auto-places right. |
| `axisLabel`      | object        | See Axis Label.                                                             |
| `axisLine`       | object        | See Axis Line.                                                              |

### Y-axis formatter

When `axisLabel.formatter` is present, the adapter creates a formatter function and replaces every `{value}` token with the library's humanized numeric representation using two decimal places.

Example:

```json
{
  "axisLabel": {
    "formatter": "{value} kg"
  }
}
```

The library effectively produces strings such as:

```text
12.34 kg
```

The formatter is implemented specifically for the y-axis conversion path.

### Y-axis line color

When `axisLine.lineStyle.color` is supplied, the adapter applies the resulting color to both the y-axis stroke and y-axis text theme. `opacity` modifies the resulting alpha channel.

## Axis Line

Supported representation:

```json
{
  "show": true,
  "lineStyle": {
    "color": "#666",
    "opacity": 0.5,
    "width": 2
  }
}
```

| Field               | Type   | Rendering behavior                         |
| ------------------- | ------ | ------------------------------------------ |
| `show`              | bool   | Applied for the y-axis spine.              |
| `lineStyle.color`   | string | Applied to the relevant axis theme/stroke. |
| `lineStyle.opacity` | number | Multiplied into the color alpha.           |
| `lineStyle.width`   | int    | Accepted but not rendered.                 |

For the x-axis, `color` and `opacity` are used; `show` is not propagated by `ToOption`.

This is the **axis** line style. The **series** `lineStyle.width` (line stroke width) is a separate field, see Series fields.

## Axis Label

| Field       | Type   | Rendering behavior                                                   |
| ----------- | ------ | -------------------------------------------------------------------- |
| `show`      | bool   | When explicitly false, the adapter makes the label font transparent. |
| `color`     | string | Label text color.                                                    |
| `fontSize`  | int    | Label font size.                                                     |
| `formatter` | string | Rendered only for y-axis labels, using `{value}` substitution.       |
| `rotate`    | number | Label rotation in degrees.                                           |
| `interval`  | number | Label skip count. Value axes only; category axes ignore it.          |
| `margin`    | number | Approximated as the label offset, away from the axis.                |

The adapter's axis-label model is much smaller than Apache ECharts' full axis-label API.

## Radar

```json
{
  "radar": {
    "indicator": [
      {
        "name": "Speed",
        "max": 100,
        "min": 0
      }
    ]
  }
}
```

`radar.indicator` is an array of:

| Field  | Type   |
| ------ | ------ |
| `name` | string |
| `max`  | number |
| `min`  | number |

The indicators are converted into the library's native radar indicator objects.

## Series

`series` may be supplied as a single object or an array.

### Series fields

| Field                  | Type           | Behavior                                                                                                     |
| ---------------------- | -------------- | -------------------------------------------------------------------------------------------------------------- |
| `type`                 | string         | Library chart type. See Series type vocabulary.                                                               |
| `name`                 | string         | Series name.                                                                                                  |
| `radius`               | string         | Pie/doughnut radius. See `series.radius`.                                                                     |
| `yAxisIndex`           | int            | Selects the y-axis for cartesian series.                                                                      |
| `data`                 | array          | Series data. See Series data.                                                                                 |
| `label`                | object         | See Label.                                                                                                    |
| `itemStyle`            | object         | `color` and `opacity` for the series.                                                                         |
| `max`                  | number         | Accepted but not rendered.                                                                                    |
| `min`                  | number         | Accepted but not rendered.                                                                                    |
| `markPoint`            | object         | See Mark Point.                                                                                               |
| `markLine`             | object         | See Mark Line.                                                                                                |
| `stack`                | string         | Stack group label. Any stack label enables the library's global stacking; single stack group only, mixed stacked/unstacked configurations error rather than mis-render. Stacked lines force area fill and ignore smoothing. Stacked bars ignore bar margin, except a second y-axis places bars beside the stack. |
| `symbol`               | string         | Point symbol shape, chart-level (first series that defines it wins). `circle`, `rect`/`roundRect` (square), `diamond`, `none`/empty. `triangle`/`pin`/`arrow` and unknown names fall back to square. |
| `symbolSize`           | number         | Point symbol size, chart-level (first series that defines it wins). Plain numbers only.                        |
| `lineStyle.width`      | int            | Line stroke width, chart-level (first series wins). Series-line `lineStyle.type` (`dashed`/`dotted`) is not part of the spec. |
| `areaStyle`            | object         | Presence enables area fill. `opacity` (0-1) scales into the 0-255 alpha channel. `color` follows per-series color support. |
| `barWidth`             | string         | Bar thickness as a percent of the category slot; absolute pixel widths are ignored. Also sets candlestick candle width. |
| `barGap`               | string         | Spacing between grouped bars (ECharts default 30%), mapped to the library's bar margin ratio. Percent only; other forms are ignored. |
| `barCategoryGap`       | string         | No library counterpart; accepted but ignored. Bar width wins when both width and gap are present.              |
| `smooth`               | bool \| number | Line smoothing. Boolean true selects the default tension (0.5), numeric 0-1 maps directly and clamps at 1. Stacked lines ignore smoothing. |
| `trendLine`            | array<object>  | Dialect field (Apache ECharts has no native trend lines). `{type, period, ...style}` entries supporting `linear`, `cubic`, `sma`, `ema`, `bollinger_upper`, `bollinger_lower`, `rsi` for line, scatter, and candlestick. Applies to close price for candlestick. |

### Series type vocabulary

The library supports chart implementations including line, scatter, bar, horizontal bar, pie, doughnut, radar, heat map, candlestick, funnel, violin, and other chart types in its native API.

The ECharts adapter does not present that list as the Apache ECharts `series.type` enumeration.

In particular:

* `"bar"` is an ECharts-style series type, remapped to horizontal bar when `xAxis.type` is `"value"`. `"violin"` is remapped to horizontal violin the same way.
* `"heatMap"` and `"violin"` are library-specific type conventions rather than a claim of exact Apache ECharts syntax. Current Apache ECharts uses `"heatmap"` for the heatmap series.
* `"pie"` and `"doughnut"` expand one input series into one native series per data item. Radius: a single number/percent is the outer radius, the `[inner, outer]` array maps to ring and center radii. `center` placement arrays are ignored (circular charts always center).
* `"candlestick"` consumes `[open, close, lowest, highest]` tuples (standard ECharts order), converted to the library's internal O/H/L/C order at ingestion. Styling: wick width, per-series wick visibility, candle style (`filled`/`traditional`/`outline`), and candle margin. Optional `itemStyle.color`/`color0`/`borderColor` map to up/down body colors via a derived theme.
* Candlestick pattern detection is configured through a dialect object: enabled pattern list (`doji`, `hammer`, `engulfing_bull`, `morning_star`, ...), `dojiThreshold`, shadow tolerance/ratio, engulfing minimum size, and prefer-pattern-labels. Pattern formatter functions cannot come from JSON.
* `"scatter"` consumes `[x, y]` pairs when the x-axis is value-typed (or a point's tuple has both elements valid); points position by their own x against a numeric x range instead of category slots.
* `"heatMap"` uses a grid-based data model with `[x, y, value]` points.

## Series data

A data entry is either:

1. a JSON number (or a JSON `null`, or the string `"-"`); or
2. an object containing `value`, optionally with `name` and `itemStyle`.

Example:

```json
[
  10,
  20,
  null,
  {
    "value": 30,
    "name": "Example"
  }
]
```

### Value tuples

`value` may be a single number or an array of numbers. The tuple is consumed per chart type:

* single-value charts (`line`, `bar`, and the per-item pie/doughnut expansion) use the first element of each tuple;
* `scatter` consumes `[x, y]` pairs (see Series type vocabulary);
* `candlestick` consumes `[open, close, lowest, highest]` groups (see Series type vocabulary);
* `violin` consumes extent pairs, including flat interleaved form (`"data": [1.2, 3.4, ...]`, one number per item), decoded pairwise by the renderer;
* `radar` and `funnel` expand one input series per data item and retain every numeric element.

### `null` and `"-"`

A bare JSON `null`, the string `"-"`, and an empty data item all resolve to the library's null sentinel value (`GetNullValue()`, which is `math.MaxFloat64`):

* `First()` returns the sentinel when no values were parsed;
* numeric parsing also accepts scientific notation (e.g. `1e3`).

The null sentinel is excluded from axis extents, sums, averages, and other aggregations by the renderer, so such points read as gaps rather than zeros.

## Data item object

```json
{
  "value": 42,
  "name": "Example",
  "itemStyle": {
    "color": "#f00",
    "opacity": 0.5
  }
}
```

| Field       | Type                   | Rendering behavior                                                          |
| ----------- | ---------------------- | --------------------------------------------------------------------------- |
| `value`     | number | array<number> | Numeric value or tuple, consumed per chart type. See Value tuples.          |
| `name`      | string                 | Data point name, used for series/data-item naming, particularly pie/radar/funnel conversion. |
| `itemStyle` | object                 | `color` and `opacity` for the data item.                                    |

## Label

Series label configuration:

| Field       | Type   | Rendering behavior                        |
| ----------- | ------ | ----------------------------------------- |
| `show`     | bool   | Controls label display for normal series. |
| `color`    | string | Label text color.                         |
| `distance` | int    | Distance from the data element.           |
| `formatter` | string | Template string: `{a}` series name, `{b}` data name, `{c}` value, `{d}` percent of the series total. The library retains one name per series, so `{b}` resolves to that name. Functions and rich-text objects are not expressible in JSON and are ignored. |
| `position`  | string | `bottom`, `inside`, and `insideBottom` shift the label below the anchor through the label offset; other values keep the default above-placement. |

The broader ECharts label API such as rich text, rotation, background, borders, and emphasis is not part of the spec.

### Pie-label special case

For pie series, the adapter creates one native series per data item and explicitly forces its label to `show: true`. The supplied series-level label configuration is not preserved by this special conversion path.

## Pie, radar, and funnel data transformation

These series types are handled specially.

### Pie

A single input pie series such as:

```json
{
  "type": "pie",
  "data": [
    {"value": 1048, "name": "Search Engine"},
    {"value": 735, "name": "Direct"}
  ]
}
```

is converted into separate native series entries, one for each data item.

For each generated pie series:

* `name` comes from the data item's `name`;
* `value` comes from the data item's first numeric value;
* pie labels are forced on;
* `radius` is preserved;
* the original series-level `yAxisIndex`, `label`, `itemStyle`, and mark annotations are not preserved through this special branch.

### Radar and funnel

Radar and funnel also expand one input series into one native series per data item.

For these paths:

* data-item `name` becomes the generated series name;
* the full parsed numeric `value` array is retained by the intermediate generic series representation;
* label color, show flag, and distance are propagated;
* the general-series `yAxisIndex`, item style, and marks are not propagated through this special conversion path.

## `series.radius`

`radius` accepts a single number/percent string, used for pie rendering, or the `[inner, outer]` array form, mapping to the library's ring and center radii for doughnut. It is not a general faithful implementation of ECharts' radius syntax for every series type.

## Mark Point

Input form:

```json
{
  "symbolSize": 30,
  "data": [
    {"type": "max"},
    {"type": "min"},
    {"xAxis": 3},
    {"coord": [2, 45]}
  ]
}
```

| Field        | Type          | Rendering behavior                      |
| ------------ | ------------- | --------------------------------------- |
| `symbolSize` | int           | Point symbol size.                      |
| `data`       | array<object> | Converted into native mark definitions. |

Mark data entries use a named `type` (`max`, `min`, `average`, `median`; availability varies per chart) or positional forms.

### Positional mark fields

`xAxis`/`yAxis` values and `coord` pairs place marks at explicit coordinates. Positional values participate in axis range computation, so a mark outside the data range expands the axis instead of clipping. On category axes, `coord` x values are category indices.

## Mark Line

Input form:

```json
{
  "data": [
    {"type": "average"},
    {"type": "max"}
  ]
}
```

The adapter converts the mark-data entries into native mark-line entries. It does **not** reduce the list to only the first non-empty `type`; multiple mark definitions are preserved by the conversion.

Positional forms: vertical/horizontal lines at `xAxis`/`yAxis` values and lines through `coord` pairs, with the same range and coordinate semantics as positional mark points.

## Accepted but not rendered

The following are accepted by the adapter structures but are not part of the rendering contract:

* `title.backgroundColor`
* `legend.backgroundColor`
* `series.max`
* `series.min`
* `series.barCategoryGap`
* `axisLine.lineStyle.width`
* `xAxis.axisLabel.formatter`
* `xAxis.axisLine.show`
* `xAxis.axisLabel.interval` (category axes have no skip support)

The adapter also accepts many ECharts-style fields syntactically only because Go's JSON unmarshalling ignores unknown properties; that does **not** imply that those properties are supported.

## Explicitly excluded from the JSON surface

* `table` (no Apache ECharts analog).
* `axisLine.onZero`, `axisTick` styling, `minorTicks`, `minorSplitLine`.
* `xAxis.splitLine` (no category-axis split line control in the library).
* Series-line dash/dot `lineStyle.type`.
* Pie `center` placement arrays (circular charts always center).
* Formatter and pattern-label **functions** (JSON can only carry template strings).
* Library-only extras such as doughnut center value display modes, deferred as dialect fields.

## Semantics / defaults worth remembering

1. **Renderer output type wins.** `RenderEChartsToPNG`, `RenderEChartsToJPG`, and `RenderEChartsToSVG` overwrite the option's `type`.
2. **Only the first x-axis is used.** Multiple x-axis entries are parsed, but only index 0 is converted.
3. **All y-axes are converted.** Each parsed y-axis becomes a native y-axis option.
4. **Unset x-axis `boundaryGap` is forced to `true`.** This is an implementation behavior and should not be treated as a complete reproduction of ECharts' axis defaults.
5. **`xAxis.type: "value"` is an orientation switch.** Bar series become horizontal bar and violin series become horizontal violin. This is library behavior, not a general ECharts axis implementation.
6. **Y-axis `{value}` formatting is library-specific.** `{value}` is replaced with a humanized numeric value using two decimal places.
7. **Pie/doughnut/radar/funnel have special data conversion.** Their input data items become separate native series rather than remaining a single ordinary series.
8. **Value tuples are consumed per chart type.** Single-value charts use the first element of each tuple; scatter, candlestick, and violin consume the full tuple.
9. **Null data is a sentinel, not zero.** `null`, `"-"`, and empty items become `GetNullValue()` and are excluded from extents and aggregations.
10. **Apache ECharts compatibility is intentionally partial.** Current Apache ECharts (6.x) exposes a much broader option model than this adapter, including additional axis types and many more components and series types.


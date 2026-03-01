package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAxisRender(t *testing.T) {
	t.Parallel()

	dayLabels := []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	letterLabels := []string{"A", "B", "C", "D", "E", "F", "G"}
	fs := NewFontStyleWithSize(18)
	fs.FontColor = ColorGreen

	tests := []struct {
		name          string
		optionFactory func(p *Painter) axisOption
	}{
		{
			name: "x-axis",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					BoundaryGap: Ptr(true),
					FontStyle:   fs,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, 0, fs))
			},
		},
		{
			name: "x-axis_rotation45",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					BoundaryGap: Ptr(true),
					FontStyle:   fs,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, DegreesToRadians(45), fs))
			},
		},
		{
			name: "x-axis_rotation90",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					Labels:      dayLabels,
					BoundaryGap: Ptr(true),
					FontStyle:   fs,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, DegreesToRadians(90), fs))
			},
		},
		{
			name: "x-axis_splitline",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:        newTestRangeForLabels(letterLabels, 0, fs),
					splitLineShow: true,
				}
			},
		},
		{
			name: "y-axis_left",
			optionFactory: func(p *Painter) axisOption {
				opt := YAxisOption{
					Position:       PositionLeft,
					isCategoryAxis: true,
				}
				return opt.prep(nil).toAxisOption(newTestRangeForLabels(dayLabels, 0, fs))
			},
		},
		{
			name: "y-axis_right",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:        newTestRangeForLabels(dayLabels, 0, fs),
					position:      PositionRight,
					boundaryGap:   Ptr(false),
					splitLineShow: true,
				}
			},
		},
		{
			name: "reduced_label_count",
			optionFactory: func(p *Painter) axisOption {
				aRange := newTestRangeForLabels(letterLabels, 0, fs)
				aRange.labelCount -= 2
				return axisOption{
					aRange:        aRange,
					splitLineShow: false,
				}
			},
		},
		{
			name: "label_start_offset",
			optionFactory: func(p *Painter) axisOption {
				aRange := newTestRangeForLabels(letterLabels, 0, fs)
				aRange.dataStartIndex = 2
				return axisOption{
					aRange: aRange,
				}
			},
		},
		{
			name: "custom_font",
			optionFactory: func(p *Painter) axisOption {
				fs := FontStyle{
					FontSize:  40.0,
					FontColor: ColorBlue,
				}
				return axisOption{
					aRange: newTestRangeForLabels(letterLabels, 0, fs),
				}
			},
		},
		{
			name: "boundary_gap_disable",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:      newTestRangeForLabels(letterLabels, 0, fs),
					boundaryGap: Ptr(false),
				}
			},
		},
		{
			name: "boundary_gap_enable",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					aRange:      newTestRangeForLabels(letterLabels, 0, fs),
					boundaryGap: Ptr(true),
				}
			},
		},
		{
			name: "dense_category_data",
			optionFactory: func(p *Painter) axisOption {
				const count = 1000
				labelLen := len(strconv.Itoa(count))
				labels := make([]string, count)
				tsl := testSeriesList{}
				for i := range labels {
					label := strconv.Itoa(i + 1)
					for len(label) < labelLen {
						label = "0" + label
					}
					labels[i] = label
					tsl = append(tsl, testSeries{values: []float64{float64(i)}})
				}
				return axisOption{
					aRange: calculateCategoryAxisRange(p, p.Width(), false, false, labels, 0,
						0, 0, 0, tsl, 0, fs),
					boundaryGap: Ptr(true),
				}
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterPaddingOption(NewBoxEqual(50)))

			opt := tt.optionFactory(p)
			opt.axisColor = ColorBlue
			opt.axisSplitLineColor = ColorGray
			_, err := newAxisPainter(p, opt).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayoutByGrid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		cols            int
		rows            int
		setupCells      func(LayoutBuilderGrid) LayoutBuilderGrid
		expectedKeys    []string
		verifyLayout    func(*testing.T, map[string]*Painter)
		expectedDemoSVG string
	}{
		{
			name: "simple_2x2",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("topLeft", 0, 0).
					CellAt("topRight", 1, 0).
					CellAt("bottomLeft", 0, 1).
					CellAt("bottomRight", 1, 1)
			},
			expectedKeys: []string{"topLeft", "topRight", "bottomLeft", "bottomRight"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				// Each cell should be 300x200 (600/2 x 400/2)
				assert.Equal(t, 300, painters["topLeft"].Width())
				assert.Equal(t, 200, painters["topLeft"].Height())
				assert.Equal(t, 300, painters["topRight"].Width())
				assert.Equal(t, 200, painters["topRight"].Height())
				assert.Equal(t, 300, painters["bottomLeft"].Width())
				assert.Equal(t, 200, painters["bottomLeft"].Height())
				assert.Equal(t, 300, painters["bottomRight"].Width())
				assert.Equal(t, 200, painters["bottomRight"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 320 20\nL 320 220\nL 20 220\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 320 20\nL 620 20\nL 620 220\nL 320 220\nL 320 20\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 220\nL 320 220\nL 320 420\nL 20 420\nL 20 220\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 320 220\nL 620 220\nL 620 420\nL 320 420\nL 320 220\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "3x3_with_spanning",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("header", 0, 0).Span(3, 1).
					CellAt("sidebar", 0, 1).Span(1, 2).
					CellAt("main", 1, 1).Span(2, 2)
			},
			expectedKeys: []string{"header", "sidebar", "main"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				// header: full width (600), 1/3 height (133)
				assert.Equal(t, 600, painters["header"].Width())
				assert.Equal(t, 133, painters["header"].Height())

				// sidebar: 1/3 width (200), 2/3 height (267)
				assert.Equal(t, 200, painters["sidebar"].Width())
				assert.Equal(t, 267, painters["sidebar"].Height())

				// main: 2/3 width (400), 2/3 height (267)
				assert.Equal(t, 400, painters["main"].Width())
				assert.Equal(t, 267, painters["main"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 153\nL 20 153\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 153\nL 220 153\nL 220 420\nL 20 420\nL 20 153\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 220 153\nL 620 153\nL 620 420\nL 220 420\nL 220 153\" style=\"stroke-width:1;stroke:blue;fill:none\"/></svg>",
		},
		{
			name: "offset",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("normal", 0, 0).
					CellAt("offset", 1, 1).Offset("10", "-20")
			},
			expectedKeys: []string{"normal", "offset"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 300, painters["normal"].Width())
				assert.Equal(t, 200, painters["normal"].Height())

				// offset: same size but position is adjusted
				assert.Equal(t, 300, painters["offset"].Width())
				assert.Equal(t, 200, painters["offset"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 320 20\nL 320 220\nL 20 220\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 330 200\nL 630 200\nL 630 400\nL 330 400\nL 330 200\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "offset_percent",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("normal", 0, 0).
					CellAt("offset", 1, 1).Offset("10%", "-0.5%")
			},
			expectedKeys: []string{"normal", "offset"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 300, painters["normal"].Width())
				assert.Equal(t, 200, painters["normal"].Height())

				// offset: same size but position is adjusted
				assert.Equal(t, 300, painters["offset"].Width())
				assert.Equal(t, 200, painters["offset"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 320 20\nL 320 220\nL 20 220\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 350 219\nL 650 219\nL 650 419\nL 350 419\nL 350 219\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "12_column_responsive_grid",
			cols: 12,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("nav", 0, 0).Span(12, 1).
					CellAt("left", 0, 1).Span(3, 1).
					CellAt("center", 3, 1).Span(6, 1).
					CellAt("right", 9, 1).Span(3, 1).
					CellAt("footer", 0, 2).Span(12, 1)
			},
			expectedKeys: []string{"nav", "left", "center", "right", "footer"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				// nav: full width
				assert.Equal(t, 600, painters["nav"].Width())
				assert.Equal(t, 133, painters["nav"].Height())

				// left: 3/12 width = 150
				assert.Equal(t, 150, painters["left"].Width())
				assert.Equal(t, 133, painters["left"].Height())

				// center: 6/12 width = 300
				assert.Equal(t, 300, painters["center"].Width())
				assert.Equal(t, 133, painters["center"].Height())

				// right: 3/12 width = 150
				assert.Equal(t, 150, painters["right"].Width())
				assert.Equal(t, 133, painters["right"].Height())

				// footer: full width
				assert.Equal(t, 600, painters["footer"].Width())
				assert.Equal(t, 134, painters["footer"].Height()) // Rounding adjustment
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 153\nL 20 153\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 153\nL 170 153\nL 170 286\nL 20 286\nL 20 153\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 170 153\nL 470 153\nL 470 286\nL 170 286\nL 170 153\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 470 153\nL 620 153\nL 620 286\nL 470 286\nL 470 153\" style=\"stroke-width:1;stroke:black;fill:none\"/><path  d=\"M 20 286\nL 620 286\nL 620 420\nL 20 420\nL 20 286\" style=\"stroke-width:1;stroke:purple;fill:none\"/></svg>",
		},
		{
			name: "full_single_cell",
			cols: 4,
			rows: 4,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("full", 0, 0).Span(4, 4)
			},
			expectedKeys: []string{"full"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["full"].Width())
				assert.Equal(t, 400, painters["full"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 420\nL 20 420\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "overlapping_cells",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("base", 0, 0).Span(2, 2).
					CellAt("overlap", 1, 1).Span(2, 2) // Overlaps with base
			},
			expectedKeys: []string{"base", "overlap"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 400, painters["base"].Width())
				assert.Equal(t, 266, painters["base"].Height())
				assert.Equal(t, 400, painters["overlap"].Width())
				assert.Equal(t, 267, painters["overlap"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 420 20\nL 420 286\nL 20 286\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 220 153\nL 620 153\nL 620 420\nL 220 420\nL 220 153\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			const padding = 20
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600 + padding + padding,
				Height:       400 + padding + padding,
			})

			builder := p.Child(PainterPaddingOption(NewBoxEqual(padding))).LayoutByGrid(tc.cols, tc.rows)
			builder = tc.setupCells(builder)
			painters, err := builder.Build()
			require.NoError(t, err)

			assert.Len(t, painters, len(tc.expectedKeys))
			for _, key := range tc.expectedKeys {
				require.Contains(t, painters, key)
			}
			tc.verifyLayout(t, painters)

			// draw border around each painter and verify visually
			colors := []Color{
				ColorRed, ColorGreen, ColorBlue, ColorBlack, ColorPurple,
				ColorAqua, ColorChocolate, ColorSalmon, ColorGray,
			}
			for i, key := range tc.expectedKeys {
				painter := painters[key]
				painter.FilledRect(0, 0, painter.Width(), painter.Height(),
					ColorTransparent, colors[i%len(colors)], 1.0)
			}
			svg, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.expectedDemoSVG, svg)
		})
	}

	t.Run("span_and_offset_without_cellat", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			Width:  600,
			Height: 400,
		})

		// Test that Span and Offset without CellAt are handled gracefully
		painters, err := p.LayoutByGrid(2, 2).
			Span(2, 2).         // Should be no-op without preceding CellAt
			Offset("10", "20"). // Should be no-op without preceding CellAt
			CellAt("valid", 0, 0).
			Build()

		require.NoError(t, err)
		assert.Len(t, painters, 1)
		require.Contains(t, painters, "valid")

		// Verify the valid cell has default span (1,1) and no offset
		validPainter := painters["valid"]
		assert.Equal(t, 300, validPainter.Width())  // Half of 600
		assert.Equal(t, 200, validPainter.Height()) // Half of 400
	})

	t.Run("chained_cells_with_modifications", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			Width:  600,
			Height: 400,
		})

		painters, err := p.LayoutByGrid(3, 3).
			CellAt("first", 0, 0).Span(2, 1).Offset("5", "5").
			CellAt("second", 2, 0).
			CellAt("third", 0, 1).Span(1, 2).
			Build()

		require.NoError(t, err)
		require.Len(t, painters, 3)

		// first: 2 cols wide with offset
		assert.Equal(t, 400, painters["first"].Width())
		assert.Equal(t, 133, painters["first"].Height())

		// second: default 1x1
		assert.Equal(t, 200, painters["second"].Width())
		assert.Equal(t, 133, painters["second"].Height())

		// third: 2 rows tall
		assert.Equal(t, 200, painters["third"].Width())
		assert.Equal(t, 267, painters["third"].Height())
	})

	t.Run("multiple_offsets_on_same_cell", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			Width:  600,
			Height: 400,
		})

		// Last offset should win
		painters, err := p.LayoutByGrid(2, 2).
			CellAt("test", 0, 0).
			Offset("10", "10").
			Offset("20", "20"). // This should override the previous offset
			Build()

		require.NoError(t, err)
		require.Len(t, painters, 1)
		// We can't directly test the offset, but the cell should be created successfully
		assert.Equal(t, 300, painters["test"].Width())
		assert.Equal(t, 200, painters["test"].Height())
	})

	t.Run("multiple_spans_on_same_cell", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			Width:  600,
			Height: 400,
		})

		// Last span should win
		painters, err := p.LayoutByGrid(3, 3).
			CellAt("test", 0, 0).
			Span(1, 1).
			Span(2, 2). // This should override the previous span
			Build()

		require.NoError(t, err)
		require.Len(t, painters, 1)
		assert.Equal(t, 400, painters["test"].Width())  // 2 cols wide
		assert.Equal(t, 266, painters["test"].Height()) // 2 rows tall
	})
}

func TestLayoutByGridErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cols        int
		rows        int
		setupCells  func(LayoutBuilderGrid) LayoutBuilderGrid
		expectedErr string
	}{
		{
			name: "invalid_grid_dimensions_zero_cols",
			cols: 0,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("test", 0, 0)
			},
			expectedErr: "invalid grid dimensions: cols and rows must be positive",
		},
		{
			name: "invalid_grid_dimensions_zero_rows",
			cols: 2,
			rows: 0,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("test", 0, 0)
			},
			expectedErr: "invalid grid dimensions: cols and rows must be positive",
		},
		{
			name: "negative_cols",
			cols: -1,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("test", 0, 0)
			},
			expectedErr: "invalid grid dimensions: cols and rows must be positive",
		},
		{
			name: "negative_rows",
			cols: 2,
			rows: -1,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("test", 0, 0)
			},
			expectedErr: "invalid grid dimensions: cols and rows must be positive",
		},
		{
			name: "duplicate_cell_names",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.
					CellAt("duplicate", 0, 0).
					CellAt("duplicate", 1, 1).
					CellAt("normal", 0, 1)
			},
			expectedErr: "duplicate cell name: 'duplicate'",
		},
		{
			name: "cell_col_out_of_bounds",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("outOfBounds", 2, 0)
			},
			expectedErr: "cell 'outOfBounds' position (2, 0) exceeds grid dimensions (2, 2)",
		},
		{
			name: "cell_row_out_of_bounds",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("outOfBounds", 0, 2)
			},
			expectedErr: "cell 'outOfBounds' position (0, 2) exceeds grid dimensions (2, 2)",
		},
		{
			name: "negative_cell_col",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("negative", -1, 0)
			},
			expectedErr: "cell 'negative' position (-1, 0) exceeds grid dimensions (3, 3)",
		},
		{
			name: "negative_cell_row",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("negative", 0, -1)
			},
			expectedErr: "cell 'negative' position (0, -1) exceeds grid dimensions (3, 3)",
		},
		{
			name: "span_exceeds_cols",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("tooWide", 1, 1).Span(3, 1)
			},
			expectedErr: "cell 'tooWide' span extends beyond grid boundaries (4, 2) > (3, 3)",
		},
		{
			name: "span_exceeds_rows",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("tooTall", 0, 1).Span(1, 3)
			},
			expectedErr: "cell 'tooTall' span extends beyond grid boundaries (1, 4) > (3, 3)",
		},
		{
			name: "span_from_edge_exceeds",
			cols: 4,
			rows: 4,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("edge", 3, 3).Span(2, 2)
			},
			expectedErr: "cell 'edge' span extends beyond grid boundaries (5, 5) > (4, 4)",
		},
		{
			name: "zero_col_span",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("zeroWidth", 1, 1).Span(0, 1)
			},
			expectedErr: "cell 'zeroWidth' has invalid span (0, 1): spans must be positive",
		},
		{
			name: "zero_row_span",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("zeroHeight", 1, 1).Span(1, 0)
			},
			expectedErr: "cell 'zeroHeight' has invalid span (1, 0): spans must be positive",
		},
		{
			name: "negative_col_span",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("negWidth", 1, 1).Span(-1, 1)
			},
			expectedErr: "cell 'negWidth' has invalid span (-1, 1): spans must be positive",
		},
		{
			name: "negative_row_span",
			cols: 3,
			rows: 3,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("negHeight", 1, 1).Span(1, -2)
			},
			expectedErr: "cell 'negHeight' has invalid span (1, -2): spans must be positive",
		},
		{
			name: "both_spans_zero",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("bothZero", 0, 0).Span(0, 0)
			},
			expectedErr: "cell 'bothZero' has invalid span (0, 0): spans must be positive",
		},
		{
			name: "both_spans_negative",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b.CellAt("bothNeg", 0, 0).Span(-1, -1)
			},
			expectedErr: "cell 'bothNeg' has invalid span (-1, -1): spans must be positive",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				Width:  600,
				Height: 400,
			})

			builder := p.LayoutByGrid(tc.cols, tc.rows)
			builder = tc.setupCells(builder)
			_, err := builder.Build()

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

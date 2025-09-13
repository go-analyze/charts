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
			name: "no-op",
			cols: 2,
			rows: 2,
			setupCells: func(b LayoutBuilderGrid) LayoutBuilderGrid {
				return b
			},
			expectedKeys: []string{},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				// no-cells
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"></svg>",
		},
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
				ColorAqua, ColorChocolate, ColorGray, ColorSalmon,
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

func TestLayoutByRows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setupRows       func(LayoutBuilderRow) LayoutBuilderRow
		expectedKeys    []string
		verifyLayout    func(*testing.T, map[string]*Painter)
		expectedDemoSVG string
	}{
		{
			name: "no-op",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b
			},
			expectedKeys: []string{},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				// nothing
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"></svg>",
		},
		{
			name: "simple_rows_with_equal_columns",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("60").EqualCols("header").
					Row().Height("200").EqualCols("left", "right").
					Row().EqualCols("footer")
			},
			expectedKeys: []string{"header", "left", "right", "footer"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["header"].Width())
				assert.Equal(t, 60, painters["header"].Height())
				assert.Equal(t, 300, painters["left"].Width())
				assert.Equal(t, 200, painters["left"].Height())
				assert.Equal(t, 300, painters["right"].Width())
				assert.Equal(t, 200, painters["right"].Height())
				assert.Equal(t, 600, painters["footer"].Width())
				assert.Equal(t, 140, painters["footer"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 80\nL 20 80\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 80\nL 320 80\nL 320 280\nL 20 280\nL 20 80\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 320 80\nL 620 80\nL 620 280\nL 320 280\nL 320 80\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 280\nL 620 280\nL 620 420\nL 20 420\nL 20 280\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "mixed_column_widths",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("300").Col("sidebar", "150").Col("content", "").
					Row().EqualCols("footer")
			},
			expectedKeys: []string{"sidebar", "content", "footer"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 150, painters["sidebar"].Width())
				assert.Equal(t, 300, painters["sidebar"].Height())
				assert.Equal(t, 450, painters["content"].Width())
				assert.Equal(t, 300, painters["content"].Height())
				assert.Equal(t, 600, painters["footer"].Width())
				assert.Equal(t, 100, painters["footer"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 170 20\nL 170 320\nL 20 320\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 170 20\nL 620 20\nL 620 320\nL 170 320\nL 170 20\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 320\nL 620 320\nL 620 420\nL 20 420\nL 20 320\" style=\"stroke-width:1;stroke:blue;fill:none\"/></svg>",
		},
		{
			name: "percentage_widths_and_heights",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("30%").Col("left", "25%").Col("center", "50%").Col("right", "25%").
					Row().EqualCols("bottom")
			},
			expectedKeys: []string{"left", "center", "right", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 150, painters["left"].Width())
				assert.Equal(t, 120, painters["left"].Height())
				assert.Equal(t, 300, painters["center"].Width())
				assert.Equal(t, 120, painters["center"].Height())
				assert.Equal(t, 150, painters["right"].Width())
				assert.Equal(t, 120, painters["right"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 280, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 170 20\nL 170 140\nL 20 140\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 170 20\nL 470 20\nL 470 140\nL 170 140\nL 170 20\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 470 20\nL 620 20\nL 620 140\nL 470 140\nL 470 20\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 140\nL 620 140\nL 620 420\nL 20 420\nL 20 140\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "col_gap_at_start",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Height("100").ColGap("50").EqualCols("content")
			},
			expectedKeys: []string{"content"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 550, painters["content"].Width())
				assert.Equal(t, 100, painters["content"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 70 20\nL 620 20\nL 620 120\nL 70 120\nL 70 20\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "col_gap_at_end",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Height("100").EqualCols("content").ColGap("50")
			},
			expectedKeys: []string{"content"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 550, painters["content"].Width())
				assert.Equal(t, 100, painters["content"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 570 20\nL 570 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "col_gap_between_columns",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("left").ColGap("20").EqualCols("right")
			},
			expectedKeys: []string{"left", "right"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 290, painters["left"].Width())
				assert.Equal(t, 100, painters["left"].Height())
				assert.Equal(t, 290, painters["right"].Width())
				assert.Equal(t, 100, painters["right"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 310 20\nL 310 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 330 20\nL 620 20\nL 620 120\nL 330 120\nL 330 20\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "multiple_col_gaps_in_sequence",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Height("100").EqualCols("a").ColGap("10").ColGap("10").EqualCols("b")
			},
			expectedKeys: []string{"a", "b"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 290, painters["a"].Width())
				assert.Equal(t, 100, painters["a"].Height())
				assert.Equal(t, 290, painters["b"].Width())
				assert.Equal(t, 100, painters["b"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 310 20\nL 310 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 330 20\nL 620 20\nL 620 120\nL 330 120\nL 330 20\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "row_with_only_gaps",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					ColGap("100").ColGap("200").Height("50").
					Row().EqualCols("bottom")
			},
			expectedKeys: []string{"bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 350, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 70\nL 620 70\nL 620 420\nL 20 420\nL 20 70\" style=\"stroke-width:1;stroke:red;fill:none\"/></svg>",
		},
		{
			name: "row_gap_spacing",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("top").
					RowGap("50").
					Row().EqualCols("bottom")
			},
			expectedKeys: []string{"top", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["top"].Width())
				assert.Equal(t, 100, painters["top"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 250, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 170\nL 620 170\nL 620 420\nL 20 420\nL 20 170\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "empty_row_with_height",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("top").
					Row().Height("50"). // Empty row with height, same behavior expected as RowGap
					Row().EqualCols("bottom")
			},
			expectedKeys: []string{"top", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["top"].Width())
				assert.Equal(t, 100, painters["top"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 250, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 170\nL 620 170\nL 620 420\nL 20 420\nL 20 170\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "multiple_row_calls_idempotent",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("first").
					Row().Row().Row(). // Multiple Row() calls should have no impact
					EqualCols("second")
			},
			expectedKeys: []string{"first", "second"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["first"].Width())
				assert.Equal(t, 100, painters["first"].Height())
				assert.Equal(t, 600, painters["second"].Width())
				assert.Equal(t, 300, painters["second"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 120\nL 620 120\nL 620 420\nL 20 420\nL 20 120\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "overlapping_with_negative_offset",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("top").
					Row().Height("100").EqualCols("middle").Offset("0", "-20").
					Row().EqualCols("bottom")
			},
			expectedKeys: []string{"top", "middle", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["top"].Width())
				assert.Equal(t, 100, painters["top"].Height())
				assert.Equal(t, 600, painters["middle"].Width())
				assert.Equal(t, 100, painters["middle"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 200, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 100\nL 620 100\nL 620 200\nL 20 200\nL 20 100\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 220\nL 620 220\nL 620 420\nL 20 420\nL 20 220\" style=\"stroke-width:1;stroke:blue;fill:none\"/></svg>",
		},
		{
			name: "percentage_offset",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("normal").
					Row().Height("100").EqualCols("first", "offset1").Offset("10%", "-5%").
					Row().Height("100").Col("second", "250").ColGap("100").Col("offset2", "250").Offset("10%", "-5%")
			},
			expectedKeys: []string{"normal", "first", "offset1", "second", "offset2"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["normal"].Width())
				assert.Equal(t, 100, painters["normal"].Height())
				// Row with 10% x offset (60px) and -5% y offset (-5px)
				assert.Equal(t, 300, painters["first"].Width())
				assert.Equal(t, 100, painters["first"].Height())
				assert.Equal(t, 300, painters["offset1"].Width())
				assert.Equal(t, 100, painters["offset1"].Height())
				// Row with specific column widths and gap, with offset
				assert.Equal(t, 250, painters["second"].Width())
				assert.Equal(t, 100, painters["second"].Height())
				assert.Equal(t, 250, painters["offset2"].Width())
				assert.Equal(t, 100, painters["offset2"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 120\nL 320 120\nL 320 220\nL 20 220\nL 20 120\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 350 115\nL 650 115\nL 650 215\nL 350 215\nL 350 115\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 220\nL 270 220\nL 270 320\nL 20 320\nL 20 220\" style=\"stroke-width:1;stroke:black;fill:none\"/><path  d=\"M 395 215\nL 645 215\nL 645 315\nL 395 315\nL 395 215\" style=\"stroke-width:1;stroke:purple;fill:none\"/></svg>",
		},
		{
			name: "auto_height_distribution",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Row().EqualCols("auto1").
					Row().EqualCols("auto2").
					Row().EqualCols("auto3").
					Row().EqualCols("auto4")
			},
			expectedKeys: []string{"auto1", "auto2", "auto3", "auto4"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				const autoHeight = 100
				assert.Equal(t, 600, painters["auto1"].Width())
				assert.Equal(t, autoHeight, painters["auto1"].Height())
				assert.Equal(t, 600, painters["auto2"].Width())
				assert.Equal(t, autoHeight, painters["auto2"].Height())
				assert.Equal(t, 600, painters["auto3"].Width())
				assert.Equal(t, autoHeight, painters["auto3"].Height())
				assert.Equal(t, 600, painters["auto4"].Width())
				assert.Equal(t, autoHeight, painters["auto4"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 120\nL 620 120\nL 620 220\nL 20 220\nL 20 120\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 220\nL 620 220\nL 620 320\nL 20 320\nL 20 220\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 320\nL 620 320\nL 620 420\nL 20 420\nL 20 320\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "auto_height_mixed",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("20").EqualCols("fixed1").
					Row().EqualCols("auto1").
					RowGap("60").
					Row().EqualCols("auto2").
					Row().EqualCols("auto3").
					Row().Height("20").EqualCols("fixed2")
			},
			expectedKeys: []string{"fixed1", "auto1", "auto2", "auto3", "fixed2"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				const autoHeight = 100
				assert.Equal(t, 600, painters["fixed1"].Width())
				assert.Equal(t, 20, painters["fixed1"].Height())
				assert.Equal(t, 600, painters["auto1"].Width())
				assert.Equal(t, autoHeight, painters["auto1"].Height())
				assert.Equal(t, 600, painters["auto2"].Width())
				assert.Equal(t, autoHeight, painters["auto2"].Height())
				assert.Equal(t, 600, painters["auto3"].Width())
				assert.Equal(t, autoHeight, painters["auto3"].Height())
				assert.Equal(t, 600, painters["fixed2"].Width())
				assert.Equal(t, 20, painters["fixed2"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 40\nL 20 40\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 40\nL 620 40\nL 620 140\nL 20 140\nL 20 40\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 200\nL 620 200\nL 620 300\nL 20 300\nL 20 200\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 300\nL 620 300\nL 620 400\nL 20 400\nL 20 300\" style=\"stroke-width:1;stroke:black;fill:none\"/><path  d=\"M 20 400\nL 620 400\nL 620 420\nL 20 420\nL 20 400\" style=\"stroke-width:1;stroke:purple;fill:none\"/></svg>",
		},
		{
			name: "complex_dashboard_layout",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("40").EqualCols("title").
					RowGap("10").
					Row().Height("80").Col("kpi1", "").ColGap("15").Col("kpi2", "").ColGap("15").Col("kpi3", "").ColGap("15").Col("kpi4", "").ColGap("15").Col("kpi5", "").
					RowGap("10").
					Row().Height("40%").EqualCols("mainChart").
					Row().EqualCols("table")
			},
			expectedKeys: []string{"title", "kpi1", "kpi2", "kpi3", "kpi4", "kpi5", "mainChart", "table"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["title"].Width())
				assert.Equal(t, 40, painters["title"].Height())
				assert.Equal(t, 108, painters["kpi1"].Width())
				assert.Equal(t, 80, painters["kpi1"].Height())
				assert.Equal(t, 108, painters["kpi2"].Width())
				assert.Equal(t, 80, painters["kpi2"].Height())
				assert.Equal(t, 108, painters["kpi3"].Width())
				assert.Equal(t, 80, painters["kpi3"].Height())
				assert.Equal(t, 108, painters["kpi4"].Width())
				assert.Equal(t, 80, painters["kpi4"].Height())
				assert.Equal(t, 600, painters["mainChart"].Width())
				assert.Equal(t, 160, painters["mainChart"].Height()) // 40% of 400
				assert.Equal(t, 600, painters["table"].Width())
				assert.Equal(t, 100, painters["table"].Height()) // Remaining
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 60\nL 20 60\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 70\nL 128 70\nL 128 150\nL 20 150\nL 20 70\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 143 70\nL 251 70\nL 251 150\nL 143 150\nL 143 70\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 266 70\nL 374 70\nL 374 150\nL 266 150\nL 266 70\" style=\"stroke-width:1;stroke:black;fill:none\"/><path  d=\"M 389 70\nL 497 70\nL 497 150\nL 389 150\nL 389 70\" style=\"stroke-width:1;stroke:purple;fill:none\"/><path  d=\"M 512 70\nL 620 70\nL 620 150\nL 512 150\nL 512 70\" style=\"stroke-width:1;stroke:aqua;fill:none\"/><path  d=\"M 20 160\nL 620 160\nL 620 320\nL 20 320\nL 20 160\" style=\"stroke-width:1;stroke:rgb(210,105,30);fill:none\"/><path  d=\"M 20 320\nL 620 320\nL 620 420\nL 20 420\nL 20 320\" style=\"stroke-width:1;stroke:rgb(128,128,128);fill:none\"/></svg>",
		},
		{
			name: "multiple_empty_rows_different_heights",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					RowGap("50").                         // Empty row at start
					Row().Height("75").EqualCols("cell"). // Row with content
					Row().Height("25").                   // Empty row in middle
					Row().Height("100").EqualCols("bottom")
			},
			expectedKeys: []string{"cell", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["cell"].Width())
				assert.Equal(t, 75, painters["cell"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 100, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 70\nL 620 70\nL 620 145\nL 20 145\nL 20 70\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 170\nL 620 170\nL 620 270\nL 20 270\nL 20 170\" style=\"stroke-width:1;stroke:green;fill:none\"/></svg>",
		},
		{
			name: "column_width_pixels",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").Col("a", "100").Col("b", "200").Col("c", "150").Col("d", "150")
			},
			expectedKeys: []string{"a", "b", "c", "d"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 100, painters["a"].Width())
				assert.Equal(t, 200, painters["b"].Width())
				assert.Equal(t, 150, painters["c"].Width())
				assert.Equal(t, 150, painters["d"].Width())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 120 20\nL 120 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 120 20\nL 320 20\nL 320 120\nL 120 120\nL 120 20\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 320 20\nL 470 20\nL 470 120\nL 320 120\nL 320 20\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 470 20\nL 620 20\nL 620 120\nL 470 120\nL 470 20\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "column_width_mixed",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").Col("fixed", "200").Col("percent", "30%").Col("auto", "")
			},
			expectedKeys: []string{"fixed", "percent", "auto"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 200, painters["fixed"].Width())
				assert.Equal(t, 180, painters["percent"].Width()) // 30% of 600
				assert.Equal(t, 220, painters["auto"].Width())    // Remaining space
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 220 20\nL 220 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 220 20\nL 400 20\nL 400 120\nL 220 120\nL 220 20\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 400 20\nL 620 20\nL 620 120\nL 400 120\nL 400 20\" style=\"stroke-width:1;stroke:blue;fill:none\"/></svg>",
		},
		{
			name: "height_at_end",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					EqualCols("top").Height("20").
					Row().Col("left", "200").Col("right", "").Height("100").
					Row().EqualCols("bottom") // Auto height
			},
			expectedKeys: []string{"top", "left", "right", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["top"].Width())
				assert.Equal(t, 20, painters["top"].Height())
				assert.Equal(t, 200, painters["left"].Width())
				assert.Equal(t, 100, painters["left"].Height())
				assert.Equal(t, 400, painters["right"].Width())
				assert.Equal(t, 100, painters["right"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 280, painters["bottom"].Height()) // Remaining height
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 40\nL 20 40\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 40\nL 220 40\nL 220 140\nL 20 140\nL 20 40\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 220 40\nL 620 40\nL 620 140\nL 220 140\nL 220 40\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 140\nL 620 140\nL 620 420\nL 20 420\nL 20 140\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "column_gap_positioning_validation",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").Col("a", "100").ColGap("50").Col("b", "100").ColGap("50").Col("c", "100").
					Row().Height("100").ColGap("150").Col("d", "200").ColGap("150").Col("e", "100")
			},
			expectedKeys: []string{"a", "b", "c", "d", "e"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				// First row: 100 + 50 + 100 + 50 + 100 = 400 used, 200 remaining
				assert.Equal(t, 100, painters["a"].Width())
				assert.Equal(t, 100, painters["b"].Width())
				assert.Equal(t, 100, painters["c"].Width())
				// Second row: 150 + 200 + 150 + 100 = 600 total
				assert.Equal(t, 200, painters["d"].Width())
				assert.Equal(t, 100, painters["e"].Width())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 120 20\nL 120 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 170 20\nL 270 20\nL 270 120\nL 170 120\nL 170 20\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 320 20\nL 420 20\nL 420 120\nL 320 120\nL 320 20\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 170 120\nL 370 120\nL 370 220\nL 170 220\nL 170 120\" style=\"stroke-width:1;stroke:black;fill:none\"/><path  d=\"M 520 120\nL 620 120\nL 620 220\nL 520 220\nL 520 120\" style=\"stroke-width:1;stroke:purple;fill:none\"/></svg>",
		},
		{
			name: "row_offset_negative",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("first").
					Row().RowOffset("-20").Height("100").EqualCols("second").
					Row().RowOffset("-20").EqualCols("third").
					Row().Height("100").EqualCols("fourth")
			},
			expectedKeys: []string{"first", "second", "third", "fourth"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["first"].Width())
				assert.Equal(t, 100, painters["first"].Height())
				assert.Equal(t, 600, painters["second"].Width())
				assert.Equal(t, 100, painters["second"].Height())
				assert.Equal(t, 600, painters["third"].Width())
				assert.Equal(t, 100, painters["third"].Height())
				assert.Equal(t, 600, painters["fourth"].Width())
				assert.Equal(t, 100, painters["fourth"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 100\nL 620 100\nL 620 200\nL 20 200\nL 20 100\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 180\nL 620 180\nL 620 280\nL 20 280\nL 20 180\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 280\nL 620 280\nL 620 380\nL 20 380\nL 20 280\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "row_offset_positive",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("100").EqualCols("first").
					Row().RowOffset("40").Height("100").EqualCols("second").
					Row().Height("100").EqualCols("third")
			},
			expectedKeys: []string{"first", "second", "third"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["first"].Width())
				assert.Equal(t, 100, painters["first"].Height())
				assert.Equal(t, 600, painters["second"].Width())
				assert.Equal(t, 100, painters["second"].Height())
				assert.Equal(t, 600, painters["third"].Width())
				assert.Equal(t, 100, painters["third"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 160\nL 620 160\nL 620 260\nL 20 260\nL 20 160\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 260\nL 620 260\nL 620 360\nL 20 360\nL 20 260\" style=\"stroke-width:1;stroke:blue;fill:none\"/></svg>",
		},
		{
			name: "row_offset_percentage",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("25%").EqualCols("quarter1").
					Row().RowOffset("-5%").Height("25%").EqualCols("quarter2").
					Row().RowOffset("5%").Height("25%").EqualCols("quarter3").
					Row().EqualCols("quarter4")
			},
			expectedKeys: []string{"quarter1", "quarter2", "quarter3", "quarter4"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["quarter1"].Width())
				assert.Equal(t, 100, painters["quarter1"].Height())
				assert.Equal(t, 600, painters["quarter2"].Width())
				assert.Equal(t, 100, painters["quarter2"].Height())
				assert.Equal(t, 600, painters["quarter3"].Width())
				assert.Equal(t, 100, painters["quarter3"].Height())
				assert.Equal(t, 600, painters["quarter4"].Width())
				assert.Equal(t, 100, painters["quarter4"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 120\nL 20 120\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 100\nL 620 100\nL 620 200\nL 20 200\nL 20 100\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 220\nL 620 220\nL 620 320\nL 20 320\nL 20 220\" style=\"stroke-width:1;stroke:blue;fill:none\"/><path  d=\"M 20 320\nL 620 320\nL 620 420\nL 20 420\nL 20 320\" style=\"stroke-width:1;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "row_offset_with_gaps",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				const gap = "20"
				return b.
					Height("80").EqualCols("top").
					RowGap(gap).
					Row().Height("80").EqualCols("middle").RowOffset("-" + gap). // Offset into the gap
					Row().Height("80").EqualCols("bottom")
			},
			expectedKeys: []string{"top", "middle", "bottom"},
			verifyLayout: func(t *testing.T, painters map[string]*Painter) {
				assert.Equal(t, 600, painters["top"].Width())
				assert.Equal(t, 80, painters["top"].Height())
				assert.Equal(t, 600, painters["middle"].Width())
				assert.Equal(t, 80, painters["middle"].Height())
				assert.Equal(t, 600, painters["bottom"].Width())
				assert.Equal(t, 80, painters["bottom"].Height())
			},
			expectedDemoSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 640 440\"><path  d=\"M 20 20\nL 620 20\nL 620 100\nL 20 100\nL 20 20\" style=\"stroke-width:1;stroke:red;fill:none\"/><path  d=\"M 20 100\nL 620 100\nL 620 180\nL 20 180\nL 20 100\" style=\"stroke-width:1;stroke:green;fill:none\"/><path  d=\"M 20 180\nL 620 180\nL 620 260\nL 20 260\nL 20 180\" style=\"stroke-width:1;stroke:blue;fill:none\"/></svg>",
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

			builder := p.Child(PainterPaddingOption(NewBoxEqual(padding))).LayoutByRows()
			builder = tc.setupRows(builder)
			painters, err := builder.Build()
			require.NoError(t, err)

			assert.Len(t, painters, len(tc.expectedKeys))
			for _, key := range tc.expectedKeys {
				require.Contains(t, painters, key)
			}
			tc.verifyLayout(t, painters)

			// Draw border around each painter for visual verification
			colors := []Color{
				ColorRed, ColorGreen, ColorBlue, ColorBlack, ColorPurple,
				ColorAqua, ColorChocolate, ColorGray, ColorSalmon,
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

	t.Run("build_is_idempotent", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		})

		builder := p.LayoutByRows().
			Height("100").EqualCols("top").
			Row().Height("100").EqualCols("left", "right").
			Row().EqualCols("bottom")

		painters1, err := builder.Build()
		require.NoError(t, err)
		painters2, err := builder.Build()
		require.NoError(t, err)

		// Same number of painters and same keys
		assert.Len(t, painters2, len(painters1))
		for k, p1 := range painters1 {
			p2, ok := painters2[k]
			require.True(t, ok, "missing key on second build: "+k)
			assert.Equal(t, p1.Width(), p2.Width())
			assert.Equal(t, p1.Height(), p2.Height())
		}
	})
}

func TestLayoutByRowsErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupRows   func(LayoutBuilderRow) LayoutBuilderRow
		expectedErr string
	}{
		{
			name: "fixed_heights_exceed_total_no_auto",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					Height("300").EqualCols("a").
					Row().Height("200").EqualCols("b")
			},
			expectedErr: "exceed painter height",
		},
		{
			name: "duplicate_cell_names",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("dup", "dup")
			},
			expectedErr: "duplicate cell name: 'dup'",
		},
		{
			name: "negative_height",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("cell").Height("-10")
			},
			expectedErr: "negative height not allowed",
		},
		{
			name: "negative_width",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Col("cell", "-10")
			},
			expectedErr: "negative width not allowed",
		},
		{
			name: "column_percentages_exceed_100",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Col("a", "60%").Col("b", "50%")
			},
			expectedErr: "column percentages exceed 100%",
		},
		{
			name: "auto_rows_zero_height",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					EqualCols("cell1").Height("400").
					Row().EqualCols("cell2") // Auto height with no space left
			},
			expectedErr: "auto-distributed rows result in zero height",
		},
		{
			name: "invalid_height_format",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("cell").Height("abc")
			},
			expectedErr: "invalid height 'abc'",
		},
		{
			name: "invalid_width_format",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Col("cell", "xyz")
			},
			expectedErr: "invalid width 'xyz'",
		},
		{
			name: "invalid_col_gap_format",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("a").ColGap("bad").EqualCols("b")
			},
			expectedErr: "invalid width 'bad'",
		},
		{
			name: "invalid_row_gap_format",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("a").RowGap("bad")
			},
			expectedErr: "invalid height 'bad'",
		},
		{
			name: "invalid_offset_x_format",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("cell").Offset("bad", "0")
			},
			expectedErr: "invalid x offset 'bad'",
		},
		{
			name: "invalid_offset_y_format",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("cell").Offset("0", "bad")
			},
			expectedErr: "invalid y offset 'bad'",
		},
		{
			name: "all_rows_fixed_exceeding_total",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					EqualCols("a").Height("200").
					Row().EqualCols("b").Height("250").
					Row().EqualCols("c") // Need an auto row to trigger the error
			},
			expectedErr: "auto-distributed rows result in zero height",
		},
		{
			name: "percentages_totaling_over_100",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					EqualCols("a").Height("60%").
					Row().EqualCols("b").Height("50%").
					Row().EqualCols("c") // Need an auto row to trigger the error
			},
			expectedErr: "auto-distributed rows result in zero height",
		},
		{
			name: "negative_col_gap",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.EqualCols("a").ColGap("-10").EqualCols("b")
			},
			expectedErr: "negative width not allowed",
		},
		{
			name: "negative_row_gap",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.
					EqualCols("a").Height("100").
					RowGap("-10")
			},
			expectedErr: "negative height not allowed",
		},
		{
			name: "explicit_column_widths_exceed_row_width",
			setupRows: func(b LayoutBuilderRow) LayoutBuilderRow {
				return b.Height("100").Col("a", "400").Col("b", "400")
			},
			expectedErr: "exceed available row width",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				Width:  600,
				Height: 400,
			})

			builder := p.LayoutByRows()
			builder = tc.setupRows(builder)
			_, err := builder.Build()

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

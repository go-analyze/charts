package charts

import (
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed chartdraw/drawing/fonts/Roboto-Medium.ttf.gz
var testRobotoFont embed.FS

// getTestFontData loads the Roboto font for testing.
func getTestFontData(t *testing.T) []byte {
	t.Helper()

	compressed, err := testRobotoFont.ReadFile("chartdraw/drawing/fonts/Roboto-Medium.ttf.gz")
	require.NoError(t, err)

	r, err := gzip.NewReader(bytes.NewReader(compressed))
	require.NoError(t, err)
	defer func() { _ = r.Close() }()

	decompressed, err := io.ReadAll(r)
	require.NoError(t, err)
	return decompressed
}

func TestFontConstants(t *testing.T) {
	t.Parallel()

	fontList := []struct {
		font        string
		defaultFont bool
	}{
		{
			font:        FontFamilyRoboto,
			defaultFont: true,
		},
		{
			font: FontFamilyNotoSans,
		},
		{
			font: FontFamilyNotoSansBold,
		},
	}

	for _, tc := range fontList {
		font := GetFont(tc.font)
		assert.NotNil(t, font)
		if tc.defaultFont {
			assert.Equal(t, GetDefaultFont(), font)
		} else {
			assert.NotEqual(t, GetDefaultFont(), font)
		}
	}
}

func TestInstallGetFont(t *testing.T) {
	t.Parallel()

	fontFamily := "install-test"
	fontData := getTestFontData(t)
	err := InstallFont(fontFamily, fontData)
	require.NoError(t, err)

	font := GetFont(fontFamily)
	assert.NotNil(t, font)
}

func TestGetPreferredFont(t *testing.T) {
	t.Parallel()

	t.Run("nil_default", func(t *testing.T) {
		require.Equal(t, GetDefaultFont(), getPreferredFont(nil))
	})
}

func TestInstallFont(t *testing.T) {
	t.Parallel()

	t.Run("valid_font", func(t *testing.T) {
		fontData := getTestFontData(t)
		err := InstallFont("test-valid", fontData)
		require.NoError(t, err)

		font := GetFont("test-valid")
		assert.NotNil(t, font)
	})

	t.Run("invalid_font_data", func(t *testing.T) {
		err := InstallFont("test-invalid", []byte("invalid data"))
		require.Error(t, err)
	})

	t.Run("case_insensitive", func(t *testing.T) {
		fontData := getTestFontData(t)
		err := InstallFont("Test-Case", fontData)
		require.NoError(t, err)

		font1 := GetFont("test-case")
		font2 := GetFont("TEST-CASE")
		assert.Equal(t, font1, font2)
	})
}

func TestGetFont(t *testing.T) {
	t.Parallel()

	t.Run("nonexistent_font_returns_default", func(t *testing.T) {
		font := GetFont("nonexistent-font")
		defaultFont := GetDefaultFont()
		assert.Equal(t, defaultFont, font)
	})

	t.Run("embedded_font_lazy_loading", func(t *testing.T) {
		font := GetFont(FontFamilyNotoSans)
		assert.NotNil(t, font)
	})
}

func TestGetDefaultFont(t *testing.T) {
	t.Parallel()

	defaultFont := GetDefaultFont()
	assert.NotNil(t, defaultFont)
}

func TestSetDefaultFont(t *testing.T) {
	t.Run("set_existing_font", func(t *testing.T) {
		t.Cleanup(func() {
			_ = SetDefaultFont(FontFamilyRoboto)
		})
		err := SetDefaultFont(FontFamilyNotoSans) // loading with the same roboto will lead to assertNotEqual failing
		require.NoError(t, err)
		origDefault := GetDefaultFont()

		fontData := getTestFontData(t)
		err = InstallFont("test-set-default", fontData)
		require.NoError(t, err)

		err = SetDefaultFont("test-set-default")
		require.NoError(t, err)

		defaultFont := GetDefaultFont()

		assert.NotNil(t, defaultFont)
		assert.NotEqual(t, origDefault, defaultFont)
	})

	t.Run("set_embedded_font", func(t *testing.T) {
		t.Cleanup(func() {
			_ = SetDefaultFont(FontFamilyRoboto)
		})
		origDefault := GetDefaultFont()

		// Should work with embedded fonts that aren't loaded yet
		err := SetDefaultFont(FontFamilyNotoSans)
		require.NoError(t, err)

		defaultFont := GetDefaultFont()

		assert.NotNil(t, defaultFont)
		assert.NotEqual(t, origDefault, defaultFont)
	})

	t.Run("set_nonexistent_font", func(t *testing.T) {
		err := SetDefaultFont("nonexistent-font")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "font not found")
	})
}

func TestCustomFontSizeRender(t *testing.T) {
	t.Parallel()

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	}, PainterThemeOption(GetTheme(ThemeLight)))

	opt := makeBasicLineChartOption()
	opt.XAxis.FontStyle.FontSize = 4.0
	opt.YAxis = []YAxisOption{
		{
			FontStyle: NewFontStyleWithSize(4.0),
		},
	}
	opt.Title.FontStyle.FontSize = 4.0
	opt.Legend.FontStyle.FontSize = 4.0
	opt.Legend.Symbol = SymbolDot

	err := p.LineChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Line</text><path d=\"M 256 19\nL 286 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"271\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"288\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"9\" y=\"37\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"9\" y=\"73\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"9\" y=\"109\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"13\" y=\"146\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"13\" y=\"182\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"13\" y=\"218\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"13\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"13\" y=\"291\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"13\" y=\"327\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"19\" y=\"364\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 28 36\nL 590 36\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 72\nL 590 72\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 109\nL 590 109\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 145\nL 590 145\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 182\nL 590 182\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 218\nL 590 218\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 255\nL 590 255\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 291\nL 590 291\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 28 328\nL 590 328\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 32 365\nL 590 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 32 370\nL 32 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 111 370\nL 111 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 191 370\nL 191 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 271 370\nL 271 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 350 370\nL 350 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 430 370\nL 430 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 510 370\nL 510 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 590 370\nL 590 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"69\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"149\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"229\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"308\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"389\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"469\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"548\" y=\"381\" style=\"stroke:none;fill:rgb(70,70,70);font-size:5.1px;font-family:'Roboto Medium',sans-serif\">G</text><path d=\"M 71 338\nL 151 335\nL 231 342\nL 310 335\nL 390 345\nL 470 313\nL 550 318\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"71\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"151\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"231\" cy=\"342\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"310\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"390\" cy=\"345\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"470\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"550\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path d=\"M 71 178\nL 151 153\nL 231 160\nL 310 152\nL 390 71\nL 470 62\nL 550 64\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"71\" cy=\"178\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"151\" cy=\"153\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"231\" cy=\"160\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"310\" cy=\"152\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"390\" cy=\"71\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"470\" cy=\"62\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"550\" cy=\"64\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>", data)
}

func TestFontCapabilities(t *testing.T) {
	t.Parallel()

	// Test fonts and categories with expected minimum support levels per font
	availableFonts := []string{FontFamilyRoboto, FontFamilyNotoSans}

	type fontExpectation struct {
		minSupport    float64 // percentage of characters that should be supported
		shouldSupport []rune  // specific characters that must be supported
	}

	testCategories := []struct {
		name         string
		content      string
		expectations map[string]fontExpectation // fontFamily -> expected support
	}{
		{
			name:    "basic_latin",
			content: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 100, shouldSupport: []rune{'A', 'a', '1'}},
				FontFamilyNotoSans: {minSupport: 100, shouldSupport: []rune{'A', 'a', '1'}},
			},
		},
		{
			name:    "punctuation",
			content: "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 100, shouldSupport: []rune{'.', ',', '!', '?'}},
				FontFamilyNotoSans: {minSupport: 100, shouldSupport: []rune{'.', ',', '!', '?'}},
			},
		},
		{
			name:    "common_symbols",
			content: `©®™§¶†‡•…‰′″‹›«»""''–—`,
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 100},
				FontFamilyNotoSans: {minSupport: 100},
			},
		},
		{
			name:    "currency",
			content: "¢£¤¥€¦§¨©ª«¬­®¯°±²³´µ¶·¸¹º»¼½¾¿₠₡₢₣₤₥₦₧₨₩₪₫€₭₮₯₰₱₲₳₴₵₶₷₸₹₺₻₼₽₾₿＄￠￡￢￣￤￥￦",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 63, shouldSupport: []rune{'$', '€', '£', '¥'}},
				FontFamilyNotoSans: {minSupport: 88, shouldSupport: []rune{'$', '€', '£', '¥', '₹'}},
			},
		},
		{
			name:    "mathematical_operators",
			content: "±×÷√∞≈≠≤≥∑∏∂∫∆∇∈∉∀∃∄∅∆∇∈∉∊∋∌∍∎∏∐∑−∓∔∕∖∗∘∙∝∞∟∠∡∢∣∤∥∦∧∨∩∪",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 34},
				FontFamilyNotoSans: {minSupport: 5},
			},
		},
		{
			name:    "arrows",
			content: "←↑→↓↔↕↖↗↘↙↚↛↜↝↞↟↠↡↢↣↤↥↦↧↨↩↪↫↬↭↮↯↰↱↲↳↴↵↶↷↸↹↺↻↼↽↾↿⇀⇁⇂⇃⇄⇅⇆⇇⇈⇉⇊⇋⇌⇍⇎⇏⇐⇑⇒⇓⇔⇕⇖⇗⇘⇙⇚⇛⇜⇝⇞⇟⇠⇡⇢⇣⇤⇥⇦⇧⇨⇩⇪⇫⇬⇭⇮⇯⇰⇱⇲⇳⇴⇵⇶⇷⇸⇹⇺⇻⇼⇽⇾⇿",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 0},
				FontFamilyNotoSans: {minSupport: 0},
			},
		},
		{
			name:    "emoji_faces",
			content: "😀😃😄😁😆😅😂🤣😊😇🙂🙃😉😌😍🥰😘😗😙😚😋😛😝😜🤪🤨🧐🤓😎🤩🥳😏😒😞😔😟😕🙁☹😣😖😫😩🥺😢😭😤😠😡🤬🤯😳🥵🥶😱😨😰😥😓🤗🤔🤭🤫🤥😶😐😑😬🙄😯😦😧😮😲🥱😴🤤😪😵🤐🥴🤢🤮🤧😷🤒🤕🤑🤠😈👿👹👺🤡💩👻💀☠👽👾🤖🎃😺😸😹😻😼😽🙀😿😾",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 0},
				FontFamilyNotoSans: {minSupport: 0},
			},
		},
		{
			name:    "geometric_shapes",
			content: "■□▢▣▤▥▦▧▨▩▪▫▬▭▮▯▰▱▲△▴▵▶▷▸▹►▻▼▽▾▿◀◁◂◃◄◅◆◇◈◉◊○◌◍◎●◐◑◒◓◔◕◖◗◘◙◚◛◜◝◞◟◠◡◢◣◤◥",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:   {minSupport: 1},
				FontFamilyNotoSans: {minSupport: 1},
			},
		},
	}

	for _, fontFamily := range availableFonts {
		t.Run(fontFamily, func(t *testing.T) {
			font := GetFont(fontFamily)
			require.NotNil(t, font, "Font %s should be available", fontFamily)

			for _, category := range testCategories {
				t.Run(category.name, func(t *testing.T) {
					expectation, hasExpectation := category.expectations[fontFamily]
					require.True(t, hasExpectation, "No expectations defined for this font/category combination")

					var supportedCount int
					for _, r := range category.content {
						if font.Index(r) != 0 {
							supportedCount++
						}
					}
					totalCount := len([]rune(category.content))
					supportPercentage := float64(supportedCount) / float64(totalCount) * 100

					t.Logf("%s support: %d/%d characters (%.1f%%)",
						fontFamily, supportedCount, totalCount, supportPercentage)

					assert.GreaterOrEqual(t, supportPercentage, expectation.minSupport)
					assert.LessOrEqual(t, supportPercentage, expectation.minSupport+1.0)
					for _, requiredChar := range expectation.shouldSupport {
						assert.NotZero(t, font.Index(requiredChar))
					}
				})
			}
		})
	}
}

func TestFontDemo(t *testing.T) {
	t.Skip("Only used to demo font rendering")

	fontList := []string{FontFamilyRoboto, FontFamilyNotoSans, FontFamilyNotoSansBold}

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        1024,
		Height:       768,
	})
	p.FilledRect(0, 0, p.Width(), p.Height(), ColorWhite, ColorWhite, 0)
	const increment = 64
	font := FontStyle{
		FontSize:  16.0,
		FontColor: ColorBlack,
	}
	pos := 1
	for _, f := range fontList {
		font.Font = GetFont(f)
		p.Text("The quick brown fox jumped over the lazy dog.",
			increment, increment*pos, 0, font)
		pos++
		p.Text("🟢⭐❓💰▫●□▲▼◇★🂡🂢🂫🃄🃍🃘🃞🃟",
			increment, increment*pos, 0, font)
		pos++
	}

	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, 0x0, data)
}

func TestFontRenderPNG(t *testing.T) {
	t.Skip("Only used to debug font rendering")

	const text = "Charts'n'stuff"

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputPNG,
		Width:        400,
		Height:       400,
	})
	p.FilledRect(0, 0, p.Width(), p.Height(), ColorWhite, ColorWhite, 0)

	// write text in incrementing sizes from 4 to 32
	y := 10
	font := FontStyle{
		FontColor: ColorBlack,
		Font:      GetFont(FontFamilyRoboto),
	}
	for size := 4.0; size <= 32.0; size += 2.0 {
		font.FontSize = size
		fullText := fmt.Sprintf("%v: %s", size, text)
		p.Text(fullText, 10, y, 0, font)
		textSize := p.MeasureText(fullText, 0, font)
		y += textSize.Height() + 4
	}

	originalPNG, err := p.Bytes()
	require.NoError(t, err)

	// scale originalPNG 4x using nearest neighbor interpolation for easier detail inspection
	img, err := png.Decode(bytes.NewReader(originalPNG))
	require.NoError(t, err)
	bounds := img.Bounds()
	scaledWidth := bounds.Dx() * 4
	scaledHeight := bounds.Dy() * 4
	scaled := image.NewRGBA(image.Rect(0, 0, scaledWidth, scaledHeight))
	for y := 0; y < scaledHeight; y++ {
		for x := 0; x < scaledWidth; x++ {
			srcX := x / 4
			srcY := y / 4
			scaled.Set(x, y, img.At(srcX, srcY))
		}
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, scaled)
	require.NoError(t, err)
	data := buf.Bytes()

	assertEqualPNGCRC(t, 0x0, data)
}

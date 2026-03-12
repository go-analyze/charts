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
	opt.XAxis.LabelFontStyle.FontSize = 4.0
	opt.YAxis = []YAxisOption{
		{
			LabelFontStyle: NewFontStyleWithSize(4.0),
		},
	}
	opt.Title.FontStyle.FontSize = 4.0
	opt.Legend.FontStyle.FontSize = 4.0
	opt.Legend.Symbol = SymbolDot

	err := p.LineChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)
}

func TestFontCapabilities(t *testing.T) {
	t.Parallel()

	// Test fonts and categories with expected minimum support levels per font
	const fontFamilyNotoSansSymbols = "notosans-chartsymbols"
	availableFonts := []string{FontFamilyRoboto, FontFamilyNotoSans, fontFamilyNotoSansSymbols}

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
				FontFamilyRoboto:          {minSupport: 100, shouldSupport: []rune{'A', 'a', '1'}},
				FontFamilyNotoSans:        {minSupport: 100, shouldSupport: []rune{'A', 'a', '1'}},
				fontFamilyNotoSansSymbols: {minSupport: 0},
			},
		},
		{
			name:    "punctuation",
			content: "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 100, shouldSupport: []rune{'.', ',', '!', '?'}},
				FontFamilyNotoSans:        {minSupport: 100, shouldSupport: []rune{'.', ',', '!', '?'}},
				fontFamilyNotoSansSymbols: {minSupport: 0},
			},
		},
		{
			name:    "common_symbols",
			content: `©®™§¶†‡•…‰′″‹›«»""''–—`,
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 100},
				FontFamilyNotoSans:        {minSupport: 100},
				fontFamilyNotoSansSymbols: {minSupport: 0},
			},
		},
		{
			name:    "currency",
			content: "¢£¤¥€¦§¨©ª«¬­®¯°±²³´µ¶·¸¹º»¼½¾¿₠₡₢₣₤₥₦₧₨₩₪₫€₭₮₯₰₱₲₳₴₵₶₷₸₹₺₻₼₽₾₿＄￠￡￢￣￤￥￦",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 87, shouldSupport: []rune{'$', '€', '£', '¥'}},
				FontFamilyNotoSans:        {minSupport: 88, shouldSupport: []rune{'$', '€', '£', '¥', '₹'}},
				fontFamilyNotoSansSymbols: {minSupport: 0},
			},
		},
		{
			name:    "mathematical_operators",
			content: "±×÷√∞≈≠≤≥∑∏∂∫∆∇∈∉∀∃∄∅∆∇∈∉∊∋∌∍∎∏∐∑−∓∔∕∖∗∘∙∝∞∟∠∡∢∣∤∥∦∧∨∩∪",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 34},
				FontFamilyNotoSans:        {minSupport: 5},
				fontFamilyNotoSansSymbols: {minSupport: 3},
			},
		},
		{
			name:    "arrows",
			content: "←↑→↓↔↕↖↗↘↙↚↛↜↝↞↟↠↡↢↣↤↥↦↧↨↩↪↫↬↭↮↯↰↱↲↳↴↵↶↷↸↹↺↻↼↽↾↿⇀⇁⇂⇃⇄⇅⇆⇇⇈⇉⇊⇋⇌⇍⇎⇏⇐⇑⇒⇓⇔⇕⇖⇗⇘⇙⇚⇛⇜⇝⇞⇟⇠⇡⇢⇣⇤⇥⇦⇧⇨⇩⇪⇫⇬⇭⇮⇯⇰⇱⇲⇳⇴⇵⇶⇷⇸⇹⇺⇻⇼⇽⇾⇿",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 1},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 18, shouldSupport: []rune{'←', '→', '↔', '↕'}},
			},
		},
		{
			name:    "bold_arrows",
			content: "⬅⬆⬇⬈⬉⬊⬋⬌⬍⬒⬓⬔⬕⬖⬗⬘⬙⬚⬛⬜⬝⬞⬟⬠⬡⬢⬣⬤⬥⬦⬧⬨⬩⬪⬫⬬⬭⬮⬯⭐⭑⭒⭓⭔⭕",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 100, shouldSupport: []rune{'⬅', '⬆', '⬇'}},
			},
		},
		{
			name:    "emoji_faces",
			content: "😀😃😄😁😆😅😂🤣😊😇🙂🙃😉😌😍🥰😘😗😙😚😋😛😝😜🤪🤨🧐🤓😎🤩🥳😏😒😞😔😟😕🙁☹😣😖😫😩🥺😢😭😤😠😡🤬🤯😳🥵🥶😱😨😰😥😓🤗🤔🤭🤫🤥😶😐😑😬🙄😯😦😧😮😲🥱😴🤤😪😵🤐🥴🤢🤮🤧😷🤒🤕🤑🤠😈👿👹👺🤡💩👻💀☠👽👾🤖🎃😺😸😹😻😼😽🙀😿😾",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 2},
			},
		},
		{
			name:    "geometric_shapes",
			content: "■□▢▣▤▥▦▧▨▩▪▫▬▭▮▯▰▱▲△▴▵▶▷▸▹►▻▼▽▾▿◀◁◂◃◄◅◆◇◈◉◊○◌◍◎●◐◑◒◓◔◕◖◗◘◙◚◛◜◝◞◟◠◡◢◣◤◥",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 5},
				FontFamilyNotoSans:        {minSupport: 1},
				fontFamilyNotoSansSymbols: {minSupport: 92},
			},
		},
		{
			name:    "enclosed_numbers",
			content: "①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳⓪❶❷❸❹❺❻❼❽❾❿⓫⓬⓭⓮⓯⓰⓱⓲⓳⓴",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 100},
			},
		},
		{
			name:    "enclosed_letters",
			content: "ⒶⒷⒸⒹⒺⒻⒼⒽⒾⒿⓀⓁⓂⓃⓄⓅⓆⓇⓈⓉⓊⓋⓌⓍⓎⓏⓐⓑⓒⓓⓔⓕⓖⓗⓘⓙⓚⓛⓜⓝⓞⓟⓠⓡⓢⓣⓤⓥⓦⓧⓨⓩ",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 100},
			},
		},
		{
			name:    "technical",
			content: "⌀⌁⌂⌈⌉⌊⌋⌘⌚⌛⌨⎔⎖⎗⎘⏎⏏⏚⏛⏣",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 100},
			},
		},
		{
			name:    "alchemical",
			content: "🜀🜁🜂🜃🜄🜅🜆🜇🜈🜉🜊🜋🜌🜍🜎🜏🜐🜑🜒🜓🜔🜕🜖🜗🜘🜙🜚🜛🜜🜝",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 0},
			},
		},
		{
			name:    "religious_cultural",
			content: "☦☧☨☩☪☫☬☭☮☯☸✝✞✟✠✡",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 18},
			},
		},
		{
			name:    "misc",
			content: "☊☋☌☍☓☤☥☹☺☻☽☾☿♀♁♂♃♄♅♆♇♈♉♊♋♌♍♎♏♐♑♒♓",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 33},
			},
		},
		{
			name:    "dingbats",
			content: "✁✂✃✄✆✈✉✍✎✏✐✑✒✓✔✕✖✗✘✙✚✛✜✢✣✤✥✦✧✩✪✫✬✭✮✯✰✱✲✳✴✵✶✷✸✹✺✻✼✽✾✿❀❁❂❃❄❅❆❇❈❉❊❋❍❏❐❑❒❖❗❘❙❚❛❜❝❞❡❢❣❤❥❦❧❶❷❸❹❺❻❼❽❾❿➔➘➙➚➛➜➝➞➟➠➡➢➣➤➥➦➧➨➩➪➫➬➭➮➯➱➲➳➴➵➶➷➸➹➺➻➼➽➾",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 100, shouldSupport: []rune{'✓', '✔', '✗', '✘', '★'}},
			},
		},
		{
			name:    "pictographs",
			content: "🌍🌎🌏🌡💰💳📈📉📊📋📍🔍🔑🔒🔓🔔🔕🔗🔥🔧🔨🔩🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛🗃🗄🗑🗒🗓🗨🗳",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 78, shouldSupport: []rune{'📈', '📉', '📊'}},
			},
		},
		{
			name:    "geometric_ext",
			content: "🞀🞁🞂🞃🞄🞅🞆🞇🞈🞉🞊🞋🞌🞍🞎🞏🞐🞑🞒🞓🞔🞕🞖🞗🞘🞙🞚🞛🞜🞠🞡🞢🞣🞤🞥🞦🞧🞨🞩🞪🞫🞬🞭🞮🞯🞰🞱🞲🞳🞴🞵🞶🞷🞸🞹🞺🞻🞼🞽🞾🞿🟀🟁🟂🟃🟄🟅🟆🟇🟈🟉🟊🟋🟌🟍🟎🟏🟐🟑🟒🟓🟔🟕🟖🟗🟘🟰",
			expectations: map[string]fontExpectation{
				FontFamilyRoboto:          {minSupport: 0},
				FontFamilyNotoSans:        {minSupport: 0},
				fontFamilyNotoSansSymbols: {minSupport: 90},
			},
		},
	}

	for _, fontFamily := range availableFonts {
		t.Run(fontFamily, func(t *testing.T) {
			font := GetFont(fontFamily)
			require.NotNil(t, font)

			for _, category := range testCategories {
				t.Run(category.name, func(t *testing.T) {
					expectation, hasExpectation := category.expectations[fontFamily]
					require.True(t, hasExpectation, "Test case has no expectations defined")

					var supportedCount int
					for _, r := range category.content {
						if font.Index(r) != 0 {
							supportedCount++
						}
					}
					totalCount := len([]rune(category.content))
					supportPercentage := float64(supportedCount) / float64(totalCount) * 100

					t.Logf("%s support: %d/%d %s characters (%.1f%%)",
						fontFamily, supportedCount, totalCount, category.name, supportPercentage)

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

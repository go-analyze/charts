package chartdraw

import (
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageWriterWriteAndDecode(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.White)

	iw := &ImageWriter{}
	require.NoError(t, png.Encode(iw, img))

	decoded, err := iw.Image()
	require.NoError(t, err)
	assert.Equal(t, img.Bounds(), decoded.Bounds())
}

func TestImageWriterSetRGBA(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	iw := &ImageWriter{}
	iw.SetRGBA(img)

	got, err := iw.Image()
	require.NoError(t, err)
	assert.Equal(t, img, got)
}

func TestImageWriterEmpty(t *testing.T) {
	t.Parallel()

	iw := &ImageWriter{}
	_, err := iw.Image()
	assert.Error(t, err)
}

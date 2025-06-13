package chartdraw

import (
	"bytes"
	"image/png"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRasterRendererRotationAndSave(t *testing.T) {
	t.Parallel()

	rr := PNG(20, 20).(*rasterRenderer)
	x, y := rr.getCoords(5, 5)
	assert.Equal(t, 5, x)
	assert.Equal(t, 5, y)

	rr.SetTextRotation(math.Pi / 2)
	x, y = rr.getCoords(5, 5)
	assert.Zero(t, x)
	assert.Zero(t, y)

	iw := &ImageWriter{}
	require.NoError(t, rr.Save(iw))
	img, err := iw.Image()
	require.NoError(t, err)
	assert.Equal(t, 20, img.Bounds().Dx())
}

func TestRasterRendererSavePNG(t *testing.T) {
	t.Parallel()

	rr := PNG(10, 10).(*rasterRenderer)
	buf := bytes.Buffer{}
	require.NoError(t, rr.Save(&buf))
	img, err := png.Decode(bytes.NewReader(buf.Bytes()))
	require.NoError(t, err)
	assert.Equal(t, 10, img.Bounds().Dx())
}

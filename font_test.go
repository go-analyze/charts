package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wcharczuk/go-chart/v2/roboto"
)

func TestInstallFont(t *testing.T) {
	fontFamily := "test"
	err := InstallFont(fontFamily, roboto.Roboto)
	require.NoError(t, err)

	font, err := GetFont(fontFamily)
	require.NoError(t, err)
	assert.NotNil(t, font)
}

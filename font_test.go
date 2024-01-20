package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wcharczuk/go-chart/v2/roboto"
)

func TestInstallFont(t *testing.T) {
	assert := assert.New(t)

	fontFamily := "test"
	err := InstallFont(fontFamily, roboto.Roboto)
	assert.Nil(err)

	font, err := GetFont(fontFamily)
	assert.Nil(err)
	assert.NotNil(font)
}

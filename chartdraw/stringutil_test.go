package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitCSV(t *testing.T) {
	assert.Empty(t, SplitCSV(""))
	assert.Equal(t, []string{"foo"}, SplitCSV("foo"))
	assert.Equal(t, []string{"foo", "bar"}, SplitCSV("foo,bar"))
	assert.Equal(t, []string{"foo", "bar"}, SplitCSV("foo, bar"))
	assert.Equal(t, []string{"foo", "bar"}, SplitCSV(" foo , bar "))
	assert.Equal(t, []string{"foo", "bar", "baz"}, SplitCSV("foo,bar,baz"))
	assert.Equal(t, []string{"foo", "bar", "baz,buzz"}, SplitCSV("foo,bar,\"baz,buzz\""))
	assert.Equal(t, []string{"foo", "bar", "baz,'buzz'"}, SplitCSV("foo,bar,\"baz,'buzz'\""))
	assert.Equal(t, []string{"foo", "bar", "baz,'buzz"}, SplitCSV("foo,bar,\"baz,'buzz\""))
	assert.Equal(t, []string{"foo", "bar", "baz,\"buzz\""}, SplitCSV("foo,bar,'baz,\"buzz\"'"))
}

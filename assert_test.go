package charts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqualSVG(t *testing.T, expected, actual string) {
	t.Helper()

	if expected != actual {
		actualFile, err := writeTempFile(actual, t.Name()+"-actual", "svg")
		assert.NoError(t, err)

		if expected == "" {
			t.Fatalf("SVG written to %s", actualFile)
		} else {
			expectedFile, err := writeTempFile(expected, t.Name()+"-expected", "svg")
			assert.NoError(t, err)
			t.Fatalf("SVG content does not match. Expected file: %s, Actual file: %s",
				expectedFile, actualFile)
		}
	}
}

func writeTempFile(content, prefix, extension string) (string, error) {
	tmpFile, err := os.CreateTemp("", strings.ReplaceAll(prefix, string(os.PathSeparator), ".")+"-*."+extension)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return "", err
	}

	return filepath.Abs(tmpFile.Name())
}

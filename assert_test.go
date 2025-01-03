package charts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertEqualSVG(t *testing.T, expected string, actual []byte) {
	t.Helper()

	actualStr := string(actual)
	if expected != actualStr {
		actualFile, err := writeTempFile(actual, t.Name()+"-actual", "svg")
		assert.NoError(t, err)

		if expected == "" {
			t.Fatalf("SVG written to %s", actualFile)
		} else {
			expectedFile, err := writeTempFile([]byte(expected), t.Name()+"-expected", "svg")
			assert.NoError(t, err)
			t.Fatalf("SVG content does not match. Expected file: %s, Actual file: %s",
				expectedFile, actualFile)
		}
	}
}

func writeTempFile(content []byte, prefix, extension string) (string, error) {
	tmpFile, err := os.CreateTemp("", strings.ReplaceAll(prefix, string(os.PathSeparator), ".")+"-*."+extension)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		return "", err
	}

	return filepath.Abs(tmpFile.Name())
}

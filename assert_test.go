package charts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func assertEqualSVG(t *testing.T, expected string, actual []byte) {
	t.Helper()

	actualStr := string(actual)
	if expected != actualStr {
		actualFile, err1 := writeTempFile(actual, t.Name()+"-actual", "svg")

		if expected == "" {
			t.Errorf("SVG written to %s", actualFile)
		} else {
			expectedFile, err2 := writeTempFile([]byte(expected), t.Name()+"-expected", "svg")
			t.Errorf("SVG content does not match. Expected file: %s, Actual file: %s",
				expectedFile, actualFile)
			require.NoError(t, err2)
		}
		require.NoError(t, err1)
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

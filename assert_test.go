package charts

import (
	"hash/crc32"
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

func assertEqualPNGCRC(t *testing.T, expected uint32, actual []byte) {
	t.Helper()

	hash := crc32.ChecksumIEEE(actual)
	if expected != hash {
		actualFile, err := writeTempFile(actual, t.Name()+"-actual", "png")
		if expected == 0 {
			t.Errorf("PNG CRC32 0x%x written to %s", hash, actualFile)
		} else {
			t.Errorf("PNG CRC32 mismatch expected: 0x%x actual: 0x%x file: %s", expected, hash, actualFile)
		}
		require.NoError(t, err)
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

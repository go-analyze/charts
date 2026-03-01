package charts

import (
	"bytes"
	"hash/crc32"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func assertTestdataSVG(t *testing.T, actual []byte) {
	t.Helper()

	baseTestName := strings.Split(t.Name(), "/")[0]
	pcs := make([]uintptr, 64)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	var callerFile string
	for {
		frame, more := frames.Next()
		if strings.HasSuffix(frame.File, "_test.go") {
			if callerFile == "" {
				callerFile = frame.File
			}
			if strings.HasSuffix(frame.Function, "."+baseTestName) || strings.Contains(frame.Function, "."+baseTestName+".") {
				callerFile = frame.File
				break
			}
		}
		if !more {
			break
		}
	}
	require.NotEmpty(t, callerFile)

	fileBase := strings.TrimSuffix(filepath.Base(callerFile), filepath.Ext(callerFile))
	wd, err := os.Getwd()
	require.NoError(t, err)
	sanitizeFixturePathPart := func(s string) string {
		if s == "" {
			return "_"
		}
		buf := make([]rune, 0, len(s))
		for _, r := range s {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
				buf = append(buf, r)
				continue
			}
			buf = append(buf, '_')
		}
		if len(buf) == 0 {
			return "_"
		}
		return string(buf)
	}
	nameParts := strings.Split(t.Name(), "/")
	pathParts := append(make([]string, 0, len(nameParts)+4), wd, "testdata", "svg", fileBase)
	for _, part := range nameParts {
		pathParts = append(pathParts, sanitizeFixturePathPart(part))
	}
	expectedPath := filepath.Join(pathParts...) + ".svg"

	if os.Getenv("UPDATE_SVG_GOLDEN") == "1" {
		require.NoError(t, os.MkdirAll(filepath.Dir(expectedPath), 0o755))
		require.NoError(t, os.WriteFile(expectedPath, actual, 0o644))
		return
	}

	expected, _ := os.ReadFile(expectedPath)
	assertEqualSVG(t, expected, actual)
}

func assertEqualSVG(t *testing.T, expected, actual []byte) {
	t.Helper()

	if !bytes.Equal(expected, actual) {
		actualFile, err := writeTempFile(actual, t.Name()+"-actual", "svg")
		require.NoError(t, err)

		if len(expected) == 0 {
			t.Errorf("SVG written to %s", actualFile)
		} else {
			expectedFile, err := writeTempFile(expected, t.Name()+"-expected", "svg")
			require.NoError(t, err)
			t.Errorf("SVG content does not match. Expected file: %s, Actual file: %s", expectedFile, actualFile)
		}
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
	defer func() { _ = tmpFile.Close() }()

	if _, err := tmpFile.Write(content); err != nil {
		return "", err
	}

	return filepath.Abs(tmpFile.Name())
}

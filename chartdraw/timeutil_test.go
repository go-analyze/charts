package chartdraw

import (
	"testing"
	"time"
)

func TestTimeToFloat64(t *testing.T) {
	// zero time
	tf := TimeToFloat64(time.Time{})
	if tf != 0 {
		t.Errorf("Expected float64 representation of zero time to be 0, but got %f", tf)
	}

	// non-zero time
	tm := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedTF := float64(tm.UnixNano())
	tf = TimeToFloat64(tm)
	if tf != expectedTF {
		t.Errorf("Expected float64 representation of time %s to be %f, but got %f", tm, expectedTF, tf)
	}
}

func TestTimeFromFloat64(t *testing.T) {
	// zero float64
	expectedT := time.Time{}
	actualT := TimeFromFloat64(0)
	if actualT != expectedT {
		t.Errorf("Expected time from float64 representation of 0 to be zero time, but got %s", actualT)
	}

	// non-zero float64 represent nanoseconds
	expectedT = time.Date(2022, 1, 1, 0, 0, 0, 123456789, time.Local)
	nanosecondsFloat := float64(expectedT.UnixNano())
	actualT = TimeFromFloat64(nanosecondsFloat)
	if actualT.Equal(expectedT) {
		t.Errorf("Expected time from float64 representation %f to be %s, but got %s", nanosecondsFloat, expectedT, actualT)
	}
}

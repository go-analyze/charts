package chartdraw

import "time"

// TimeToFloat64 returns a float64 representation of a time.
func TimeToFloat64(t time.Time) float64 {
	if t.IsZero() {
		return 0
	}
	return float64(t.UnixNano())
}

// TimeFromFloat64 returns a time in nanosecond from a float64.
func TimeFromFloat64(tf float64) time.Time {
	if tf == 0 {
		return time.Time{}
	}
	return time.Unix(0, int64(tf))
}

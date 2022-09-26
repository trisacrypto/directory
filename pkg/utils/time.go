package utils

import "time"

// Helper method to return the latest string timestamp from the two RFC3339 timestamps
func Latest(a, b string) string {
	// Parse without checking errors - will use zero-valued ts for checks
	ats, _ := time.Parse(time.RFC3339Nano, a)
	bts, _ := time.Parse(time.RFC3339Nano, b)

	switch {
	case ats.IsZero() && bts.IsZero():
		return ""
	case !ats.IsZero() && ats.After(bts):
		return a
	case !bts.IsZero() && bts.After(ats):
		return b
	case !ats.IsZero() && !bts.IsZero() && ats.Equal(bts):
		return a
	}
	return ""
}

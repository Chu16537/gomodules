package util

import "time"

func TimeToUnix13(time time.Time) uint64 {
	if time.IsZero() {
		return 0
	}
	return uint64(time.UnixNano()) / 1000000
}

// Unix13ToTime unix 13碼
func Unix13ToTime(millisecond uint64) time.Time {
	return time.Unix(0, int64(millisecond)*1000000)
}

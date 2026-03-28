package util

import (
	"strconv"
	"time"
)

func TimePtr2String(t *time.Time) string {
	if t != nil {
		return t.String()
	}
	return ""
}

func Uint2String(n uint) string {
	return strconv.FormatUint(uint64(n), 10)
}

func String2Uint(s string) uint {
	u, _ := strconv.ParseUint(s, 10, 32)
	return uint(u)
}

func Uint2StringOrEmpty(n *uint) string {
	if n == nil {
		return ""
	}
	return strconv.FormatUint(uint64(*n), 10)
}

func ParseUint(s string) (uint, error) {
	uid, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(uid), nil
}

func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func FromTimestamp(ts string) (time.Time, error) {
	unix, err := strconv.ParseUint(ts, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(int64(unix)), nil
}

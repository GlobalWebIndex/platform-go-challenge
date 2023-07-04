package id

import (
	"strconv"
	"time"
)

func RevFromTime(t time.Time) string {
	return strconv.FormatInt(t.UnixNano(), 36)
}

func Rev() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

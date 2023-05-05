package db

import (
	"strings"
	"time"
)

const (
	DBMSTypeInstance = "instance"
	DBMSTypeCluster  = "cluster"
)

type DateTimeWithoutTZFormat struct {
	time.Time
}

func (t *DateTimeWithoutTZFormat) UnmarshalJSON(b []byte) (err error) {
	layout := "2006-01-02T15:04:05"
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	t.Time, err = time.Parse(layout, s)
	return
}

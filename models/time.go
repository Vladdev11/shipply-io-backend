package models

import (
	"strings"
	"time"
)

// Specify the layout for the time.Time for SingleDate
const SingleDateLayout = "2006-01-02"

// SingleDate is a custom type for time.Time that will be used for dates
type SingleDate struct {
	time.Time
}

func (sd *SingleDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		sd.Time = time.Time{}
		return
	}
	sd.Time, err = time.Parse(SingleDateLayout, s)
	return
}

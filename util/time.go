package util

import "time"

func ParseRFC3339Date(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}

func AddBusinessDays(date time.Time, days int) time.Time {
	for i := 0; i < days; {
		date = date.AddDate(0, 0, 1)
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			i++
		}
	}
	return date
}

package util

import "time"

func IfWeekendGetFirstWorkDay(t time.Time) time.Time {
	if t.Weekday() == 6 {
		t = t.AddDate(0, 0, 1)
	}
	if t.Weekday() == 0 {
		t = t.AddDate(0, 0, 1)
	}

	return t
}

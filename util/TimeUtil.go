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

func GetDaySerbian(time time.Time) string {
	dayInWeek := time.Weekday()
	switch dayInWeek {
	case 0:
		return "Недеља"
	case 1:
		return "Понедељал"
	case 2:
		return "Уторак"
	case 3:
		return "Среда"
	case 4:
		return "Четвртак"
	case 5:
		return "Петак"
	case 6:
		return "Субота"
	}

	return ""
}

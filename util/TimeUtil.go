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
		return "Nedelja"
	case 1:
		return "Ponedeljak"
	case 2:
		return "Utorak"
	case 3:
		return "Sreda"
	case 4:
		return "Četvrtak"
	case 5:
		return "Petak"
	case 6:
		return "Subota"
	}

	return ""
}

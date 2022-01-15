package tools

import "time"

func SameDate(date1 time.Time, date2 time.Time) bool {
	year1, month1, day1 := date1.Date()
	year2, month2, day2 := date2.Date()
	return year1 == year2 && month1 == month2 && day1 == day2
}
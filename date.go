package codekit

import "time"

func DateOnly(dt time.Time) time.Time {
	res := time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, dt.Location())
	return res
}

func TimeOnly(dt time.Time, useSecond, useNano bool) time.Time {
	var (
		second int
		nano   int
	)
	if useSecond {
		second = dt.Second()
	}
	if useNano {
		nano = dt.Nanosecond()
	}
	res := time.Date(1900, 1, 1, dt.Hour(), dt.Minute(), second, nano, dt.Location())
	return res
}

func FirstOfMonth(dt time.Time) time.Time {
	res := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
	return res
}

func EndOfMonth(dt time.Time) time.Time {
	fm := FirstOfMonth(dt)
	lm := fm.AddDate(0, 1, 0).Add(1 * time.Millisecond)
	return lm
}

package helper

import (
	"time"
)

const DateTimeFormat = "YYYY-MM-DD"

func BeginningOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func EndOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, -1, time.UTC)
}

func ConvertToNullPointDateTime(date *string) *time.Time {
	if date == nil {
		return nil
	}
	t, err := time.Parse(*date, DateTimeFormat)

	if err != nil {
		return nil
	}
	return &t
}

package utils

import (
	"errors"
	"time"
)

// ParseDate converts a string in DD/MM/YYYY format to a time.Time object
func ParseDate(dateStr string) (time.Time, error) {
	layout := "02/01/2006"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected DD/MM/YYYY")
	}
	return parsedDate, nil
}

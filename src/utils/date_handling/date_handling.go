package date_handling

import "time"

const (
	layoutDate = "2006-01-02"
	layoutTS   = "2006-01-02T15:04:05.000Z"
)

func ValidDate(s string) bool {
	_, err := time.Parse(layoutDate, s)
	if err != nil {
		return false
	}
	return true
}

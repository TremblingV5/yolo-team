package timeutils

import (
	"fmt"
	"time"
)

func Parse(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse time: %s", s)
}

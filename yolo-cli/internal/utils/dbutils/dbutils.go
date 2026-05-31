package dbutils

import "strings"

func IsDuplicate(err error) bool {
	s := err.Error()
	return strings.Contains(s, "UNIQUE") || strings.Contains(s, "unique") || strings.Contains(s, "Duplicate")
}

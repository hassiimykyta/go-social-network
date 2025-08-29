package helpers

import "strconv"

func ParseInt(s string, def, max int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return def
	}
	if n > max {
		return max
	}
	return n
}

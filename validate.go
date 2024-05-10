package main

import "regexp"

const (
	patternHM         = "^[0-9]{2}:[0-9]{2}$"
	patternClientName = "^[a-z0-9_-]+$"
)

func IsValidHour(h int) bool {
	return h >= 0 && h < 24
}

func IsValidMinute(m int) bool {
	return m >= 0 && m < 60
}

func MatchString(pattern, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}

var IsValidHMString = func(s string) bool { return MatchString(patternHM, s) }
var IsValidClientName = func(s string) bool { return MatchString(patternClientName, s) }

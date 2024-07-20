package service

import "regexp"

func IsValidText(text string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return regex.MatchString(text)
}

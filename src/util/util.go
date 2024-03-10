package util

import (
	"regexp"
)

func IsValidURL(url string) bool {
	pattern := `^(https?:\/\/www.)([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(\/[a-zA-Z0-9-_.~!*'();:@&=+$,%#]+)*\/?$`
	matched, _ := regexp.MatchString(pattern, url) //ignored the error since the regex is static.
	if !matched {
		return false
	}
	return true
}

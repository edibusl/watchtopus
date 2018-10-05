package infra

import "regexp"

var regexHostWithPort = regexp.MustCompile(`[:]\d+`)

func ParseHost(rawHost string) string {
	return regexHostWithPort.ReplaceAllString(rawHost, ``)
}

func FindInArray(haystack []string, needle string) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}

	return false
}

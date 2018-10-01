package infra

import "regexp"

var regexHostWithPort = regexp.MustCompile(`[:]\d+`)

func ParseHost(rawHost string) string {
	return regexHostWithPort.ReplaceAllString(rawHost, ``)
}

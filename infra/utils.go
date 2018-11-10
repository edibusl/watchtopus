package infra

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

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

func ParseResponseBody(resp *http.Response) map[string]*json.RawMessage {
	// Read body
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		// TODO
		return nil
	}

	// Unmarshal
	var msg map[string]*json.RawMessage
	err = json.Unmarshal(b, &msg)
	if err != nil {
		//TODO
		return nil
	}

	return msg
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

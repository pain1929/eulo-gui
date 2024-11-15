package fb_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// ...
func ParseAndPanic(message string) {
	error_regex := regexp.MustCompile("^(\\d{3} [a-zA-Z ]+)\n\n(.*?)($|\n)")
	err_matches := error_regex.FindAllStringSubmatch(message, 1)
	if len(err_matches) == 0 {
		panic("Unknown error")
	}
	panic(fmt.Errorf("%s: %s", err_matches[0][1], err_matches[0][2]))
}

// ...
func AssertAndParse[T any](resp *http.Response) T {
	if resp.StatusCode == 503 {
		panic("API server is down")
	}
	_body, _ := io.ReadAll(resp.Body)
	body := string(_body)
	if resp.StatusCode != 200 {
		ParseAndPanic(body)
	}
	var ret T
	err := json.Unmarshal([]byte(body), &ret)
	if err != nil {
		panic(fmt.Sprintf("Error parsing API response: %v", err))
	}
	return ret
}

package fb_client

import (
	"fmt"
	"net/http"
)

// ...
type SecretLoadingTransport struct {
	secret string
}

// ...
func (s *SecretLoadingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.secret))
	return http.DefaultTransport.RoundTrip(req)
}

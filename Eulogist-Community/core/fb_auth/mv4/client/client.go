package fb_client

import (
	"fmt"
	"io"
	"net/http"
)

// ...
type ClientInfo struct {
	FBUCUsername string
	GrowthLevel  int
	RespondTo    string
	Uid          string
}

// ...
type Client struct {
	ClientInfo ClientInfo
	HttpClient http.Client
	AuthServer string
}

// ...
func CreateClient(authServer string) *Client {
	secret_res, err := http.Get(fmt.Sprintf("%s/api/new", authServer))
	if err != nil {
		panic("Failed to contact with API")
	}

	_secret_body, _ := io.ReadAll(secret_res.Body)
	secret_body := string(_secret_body)
	if secret_res.StatusCode == 503 {
		panic("API server is down")
	} else if secret_res.StatusCode != 200 {
		ParseAndPanic(secret_body)
	}

	authclient := &Client{
		HttpClient: http.Client{Transport: &SecretLoadingTransport{
			secret: secret_body,
		}},
		AuthServer: authServer,
	}

	return authclient
}

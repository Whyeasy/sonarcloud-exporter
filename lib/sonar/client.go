package sonar

import (
	"net/http"
	"time"
)

//ListOptions provide general API request options.
type ListOptions struct {
	Page    int
	PerPage int
}

//Client holds the variables we want to use for the connection.
type Client struct {
	ListOptions
	sonarConnectionString string
	client                *http.Client
	organization          string
}

//NewClient creates a new SonarCloud client
func NewClient(token string, org string) *Client {

	uri := "https://" + token + "@sonarcloud.io/api"

	return &Client{
		sonarConnectionString: uri,
		client:                &http.Client{Timeout: time.Second * 10},
		organization:          org,
	}
}

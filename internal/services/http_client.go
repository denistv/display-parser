package services

import (
	"net/http"
	"time"
)

func NewDefaultHTTPClient(timeout time.Duration) *DefaultHTTPClient {
	c := http.Client{
		Timeout: timeout,
	}

	return &DefaultHTTPClient{
		client: &c,
	}
}

type DefaultHTTPClient struct {
	client *http.Client
}

func (d *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return d.client.Do(req)
}

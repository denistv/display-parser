package iface

import "net/http"

// HTTPClient с возможностью задавать таймаут соединения
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

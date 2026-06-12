package client

import "net/http"

func DefaultHTTPClient(string) (*http.Client, error) {
	return &http.Client{}, nil
}

package net

import (
	"net/http"
	"time"
)

// Client - Siego client
type Client struct {
	client *http.Client
}

// NewClient - Creates client
func NewClient(timeout int) *Client {
	c := Client{}

	c.client = &http.Client{
		Timeout:time.Duration(timeout) * time.Second,
	}

	return &c
}

// Do - Makes request
func (c *Client) Do(r *Request) *Response {
	start := time.Now()

	resp, err := c.client.Do(r.GetHttpRequest())

	response := Response{
		HttpResponse: resp,
		Duration:     time.Since(start),
		Error:        err,
	}

	return &response
}

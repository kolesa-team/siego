package net

import (
	"net/http"
	"time"
)

// Siego client
type Client struct {
	client *http.Client
}

func NewClient() *Client {
	c := Client{}

	c.client = http.DefaultClient

	return &c
}

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

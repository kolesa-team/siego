package net

import (
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	c := Client{}

	c.client = http.DefaultClient

	return &c
}

func (this *Client) Do(r *Request) *Response {
	start := time.Now()

	resp, err := this.client.Do(r.GetHttpRequest())

	response := Response{
		HttpResponse: resp,
		Duration:     time.Since(start),
		Error:        err,
	}

	return &response
}

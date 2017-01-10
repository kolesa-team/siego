package net

import (
	"net/http"
	"strings"
)

type Request struct {
	httpRequest                         *http.Request
	method, url, userAgent, contentType string
	headers                             []string
}

// Creates request object with method and url parts set
func NewRequest(method, url, params string) (*Request, error) {
	var err error

	r := Request{}
	r.httpRequest, err = http.NewRequest(method, url, strings.NewReader(params))

	return &r, err
}

// Sets User-Agent property
func (this *Request) UserAgent(userAgent string) {
	this.userAgent = userAgent
}

// Sets Content-Type property
func (this *Request) ContentType(contentType string) {
	this.contentType = contentType
}

// Sets request headers
func (this *Request) Headers(headers []string) {
	this.headers = headers
}

// Returns associated HTTP request with headers set
func (this *Request) GetHttpRequest() *http.Request {
	if this.userAgent != "" {
		this.httpRequest.Header.Set("User-Agent", this.userAgent)
	}

	if this.contentType != "" {
		this.httpRequest.Header.Set("Content-Type", this.contentType)
	}

	for _, header := range this.headers {
		split := strings.Split(header, ":")

		if len(split) == 2 {
			split[0] = strings.Trim(split[0], " ")
			split[1] = strings.Trim(split[1], " ")

			if split[0] != "" && split[1] != "" {
				this.httpRequest.Header.Set(split[0], split[1])
			}
		}
	}

	return this.httpRequest
}

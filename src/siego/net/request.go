package net

import (
	"net/http"
	"strings"
)

// HTTP request structure
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
func (rq *Request) UserAgent(userAgent string) {
	rq.userAgent = userAgent
}

// Sets Content-Type property
func (rq *Request) ContentType(contentType string) {
	rq.contentType = contentType
}

// Sets request headers
func (rq *Request) Headers(headers []string) {
	rq.headers = headers
}

// Returns associated HTTP request with headers set
func (rq *Request) GetHttpRequest() *http.Request {
	if rq.userAgent != "" {
		rq.httpRequest.Header.Set("User-Agent", rq.userAgent)
	}

	if rq.contentType != "" {
		rq.httpRequest.Header.Set("Content-Type", rq.contentType)
	}

	for _, header := range rq.headers {
		split := strings.Split(header, ":")

		if len(split) == 2 {
			split[0] = strings.Trim(split[0], " ")
			split[1] = strings.Trim(split[1], " ")

			if split[0] != "" && split[1] != "" {
				rq.httpRequest.Header.Set(split[0], split[1])
			}
		}
	}

	return rq.httpRequest
}

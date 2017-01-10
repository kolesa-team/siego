package net

import (
	"net/http"
	"time"
)

// HTTP response structure
type Response struct {
	HttpResponse *http.Response
	Duration     time.Duration
	Error        error
}

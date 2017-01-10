package net

import (
	"net/http"
	"time"
)

type Response struct {
	HttpResponse *http.Response
	Duration     time.Duration
	Error        error
}

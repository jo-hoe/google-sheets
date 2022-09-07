package client

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseSummery struct {
	ResponseCode int
	ResponseBody string
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewMockClient returns *http.Client with Transport replaced to avoid making real calls
func NewMockClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func CreateMockClient(summeries ...ResponseSummery) *http.Client {
	i := -1
	return NewMockClient(func(_ *http.Request) *http.Response {
		i = i + 1
		return &http.Response{
			StatusCode: summeries[i].ResponseCode,
			// Send response to be tested
			Body: io.NopCloser(bytes.NewBufferString(summeries[i].ResponseBody)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
}

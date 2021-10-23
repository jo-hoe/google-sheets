package reader

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestNewSheetReader(t *testing.T) {
	readerCloser, err := NewSheetReader(createMockClient(), "spreadSheatId", "sheetName")
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	} else {
		defer readerCloser.Close()
	}

	csv := csv.NewReader(readerCloser)

	actual, err := csv.ReadAll()
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	expected := [][]string{
		{"0", "1"},
		{"2", "3"},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected '%v' found '%v'", expected, actual)
	}
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

func createMockClient() *http.Client {
	return NewMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString("\"0\",\"1\"\n\"2\",\"3\"")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
}

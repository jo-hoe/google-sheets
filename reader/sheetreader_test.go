package reader

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func Test_NewSheetReader(t *testing.T) {
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
			Body: ioutil.NopCloser(bytes.NewBufferString("{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"values\":[[\"0\",\"1\"],[\"2\",\"3\"]]}")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
}

func Test_truncateExtraneousData(t *testing.T) {
	type args struct {
		reader io.ReadCloser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Positive Test",
			wantErr: false,
			args: args{
				reader: ioutil.NopCloser(strings.NewReader("{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"values\":[[\"a\",\"b\"],[\"1\",\"2\"]]}")),
			},
			want: "a,b\n1,2\n",
		}, {
			name:    "Not readable values",
			wantErr: true,
			args: args{
				reader: ioutil.NopCloser(strings.NewReader("{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"[not readable]}")),
			},
		}, {
			name:    "Read empty",
			wantErr: false,
			args: args{
				reader: ioutil.NopCloser(strings.NewReader("{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"values\":[]}")),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := truncateExtraneousData(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("truncateExtraneousData() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				// check if nil is returned in error case
				if got != nil {
					t.Errorf("expected nil got %v", got)
				}
			} else {
				// check if string in reader is expected
				buffer := new(bytes.Buffer)
				_, err = buffer.ReadFrom(got)
				stringOutput := buffer.String()
				if err != nil {
					t.Errorf("found error while reading to buffer %v", err)
				}
				if stringOutput != tt.want {
					t.Errorf("truncateExtraneousData() = %v, want %v", stringOutput, tt.want)
				}
			}
		})
	}
}

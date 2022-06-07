package apiwrapper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/jo-hoe/google-sheets/client"
)

func Test_GetSheetId(t *testing.T) {
	expected := 2047441944
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: fmt.Sprintf(`{
			"sheets": [{
					"properties": {
						"sheetId": 0,
						"title": "Sheet1"
					}
				}, {
					"properties": {
						"sheetId": %d,
						"title": "Sheet2"
					}
				}
			]
		}`, expected),
	}
	mockClient := client.CreateMockClient(mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	actual, err := wrappper.GetSheetId("spreadSheatId", "Sheet2")
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}
	if actual != expected {
		t.Errorf("expected %d but found %d", expected, actual)
	}
}

func Test_GetSheetData(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: "{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"values\":[[\"0\",\"1\"],[\"2\",\"3\"]]}",
	}
	wrappper := NewSheetsApiWrapper(client.CreateMockClient(mockResponse))

	actual, err := wrappper.GetSheetData("spreadSheatId", "sheetName")
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(actual)
	if err != nil {
		t.Errorf("found error %v", err)
	}
	stringActual := buf.String()
	expected := "0,1\n2,3\n"

	if !reflect.DeepEqual(expected, stringActual) {
		t.Errorf("expected '%v' found '%v'", expected, stringActual)
	}
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

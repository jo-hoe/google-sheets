package apiwrapper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jo-hoe/google-sheets/client"
)

func Test_GetSheetId(t *testing.T) {
	var expected int32 = 2047441944
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

func Test_CreateSheet(t *testing.T) {
	var expectedId int32 = 2047441944
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: fmt.Sprintf(`{
			"spreadsheetId": "spreadSheetId",
			"updatedSpreadsheet": {
			"sheets": [{
			"properties": {
				"sheetId": %d,
				"title": "Sheet1"
			}}]}
		}`, expectedId),
	}
	mockClient := client.CreateMockClient(mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	actual, err := wrappper.CreateSheet("spreadSheetId", int32(time.Now().UnixMilli()/1000), "Sheet1")
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}

	if err != nil {
		t.Error("expected no error but found", err)
	}

	if actual != expectedId {
		t.Errorf("expected %d but found %d", expectedId, actual)
	}
}

func Test_AutoResizeSheet(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
	}
	mockClient := client.CreateMockClient(mockResponse, mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	err := wrappper.AutoResizeSheet("spreadSheatId", 1)
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}

	if err != nil {
		t.Error("expected no error but found", err)
	}
}

func Test_WriteSheet(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
	}
	mockClient := client.CreateMockClient(mockResponse, mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	err := wrappper.WriteSheet("spreadSheatId", "spreadSheetName", [][]string{})
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}

	if err != nil {
		t.Error("expected no error but found", err)
	}
}


func Test_AppendToSheet(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
	}
	mockClient := client.CreateMockClient(mockResponse, mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	err := wrappper.AppendToSheet("spreadSheatId", "spreadSheetName", [][]string{})
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}

	if err != nil {
		t.Error("expected no error but found", err)
	}
}

func Test_UpdateSheetMetaData(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
	}
	mockClient := client.CreateMockClient(mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	err := wrappper.UpdateSheetMetaData("spreadSheatId", 1, "Sheet2")
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}

	if err != nil {
		t.Error("expected no error but found", err)
	}
}

func Test_Delete(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
	}
	mockClient := client.CreateMockClient(mockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	err := wrappper.DeleteSheet("spreadSheatId", 1)
	if err != nil {
		t.Errorf("found error while reading to buffer %v", err)
	}

	if err != nil {
		t.Error("expected no error but found", err)
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

func Test_ReplaceSheet(t *testing.T) {
	sheetsMockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: `{
			"sheets": [{
					"properties": {
						"sheetId": 1,
						"title": "sheetName"
					}
				}
			]
		}`,
	}
	sheetMockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: `{
			"properties": {
				"sheetId": 1,
				"title": "title"
			}
		}`,
	}
	mockClient := client.CreateMockClient(sheetsMockResponse, sheetMockResponse, sheetMockResponse, sheetMockResponse, sheetMockResponse, sheetMockResponse, sheetMockResponse)
	wrappper := NewSheetsApiWrapper(mockClient)
	err := wrappper.ReplaceSheetData("spreadSheatId", "sheetName", [][]string{})

	if err != nil {
		t.Error("expected no error but found", err)
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

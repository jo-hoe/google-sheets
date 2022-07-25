package reader

import (
	"encoding/csv"
	"reflect"
	"testing"

	"github.com/jo-hoe/google-sheets/internal/client"
)

func Test_NewSheetReader(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: "{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"values\":[[\"0\",\"1\"],[\"2\",\"3\"]]}",
	}
	mock := client.CreateMockClient(mockResponse)
	reader, err := NewSheetReader(mock, "spreadSheatId", "sheetName")
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	csv := csv.NewReader(reader)

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

func Test_NewSheetReader_Read_Partial(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: "{\"range\":\"Sheet2!A1:Z1000\",\"majorDimension\":\"ROWS\",\"values\":[[\"0\",\"1\"],[\"2\"]]}",
	}
	mock := client.CreateMockClient(mockResponse)
	reader, err := NewSheetReader(mock, "spreadSheatId", "sheetName")
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	csv := csv.NewReader(reader)
	csv.FieldsPerRecord = -1

	actual, err := csv.ReadAll()
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	expected := [][]string{
		{"0", "1"},
		{"2"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected '%v' found '%v'", expected, actual)
	}
}

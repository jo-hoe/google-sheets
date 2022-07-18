package writer

import (
	"encoding/csv"
	"testing"

	"github.com/jo-hoe/google-sheets/client"
)

func TestSheetWriter_Write(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: `{
			"sheets": [{
					"properties": {
						"sheetId": 0,
						"title": "sheetId"
					}
				}
			]
		}`,
	}
	mockClient := client.CreateMockClient(mockResponse, mockResponse)

	sheetWriter, err := NewSheetWriter(mockClient, "spreadSheetId", "sheetId", O_CREATE)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	writer := csv.NewWriter(sheetWriter)
	err = writer.WriteAll([][]string{
		{"0", "1"},
		{"2", "3"},
	})
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
}

func Test_hasFlag(t *testing.T) {
	flags := O_CREATE | O_EXCL | O_TRUNC
	if !hasFlag(flags, O_CREATE) {
		t.Errorf("%d did not have O_CREATE", flags)
	}
	if !hasFlag(flags, O_EXCL) {
		t.Errorf("%d did not have O_CREATE", flags)
	}
	if !hasFlag(flags, O_TRUNC) {
		t.Errorf("%d did not have O_CREATE", flags)
	}
}

func Test_hasFlag_without_flags(t *testing.T) {
	var flags int
	if hasFlag(flags, O_CREATE) {
		t.Errorf("%d did unexpectedly find O_CREATE", flags)
	}
	if hasFlag(flags, O_EXCL) {
		t.Errorf("%d did unexpectedly find O_CREATE", flags)
	}
	if hasFlag(flags, O_TRUNC) {
		t.Errorf("%d did unexpectedly find O_CREATE", flags)
	}
}

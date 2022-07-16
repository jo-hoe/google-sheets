package writer

import (
	"encoding/csv"
	"fmt"
	"testing"

	"github.com/jo-hoe/google-sheets/client"
)

func TestSheetWriter_Write(t *testing.T) {
	mockResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: fmt.Sprint(`{
			"sheets": [{
					"properties": {
						"sheetId": 0,
						"title": "sheetId"
					}
				}
			]
		}`),
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

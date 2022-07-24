package gs

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jo-hoe/google-sheets/api/client"
)

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

func Test_openSheetWithClient(t *testing.T) {
	// prepare
	expectedSpreadsheetId := "spreadsheetId"
	expectedSheetName := "sheetName"
	expectedSheetId := int32(1)

	// mock a scenario where the sheet exists already
	sheetResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: fmt.Sprintf(`{
			"sheets": [{
					"properties": {
						"sheetId": %d,
						"title": "sheetName"
					}
				}
			]
		}`, expectedSheetId),
	}
	// mock successful truncation of content
	truncatedSheetResponse := client.ResponseSummery{
		ResponseCode: 200,
	}
	mockClient := client.CreateMockClient(sheetResponse, truncatedSheetResponse)

	// test
	actual, err := openSheetWithClient(expectedSpreadsheetId, expectedSheetName, O_RDWR|O_TRUNC, mockClient)

	if err != nil {
		t.Errorf("found error %+v", err)
	}

	assertEqual(t, expectedSheetId, actual.Id())
	assertEqual(t, expectedSheetName, actual.Name())
	assertEqual(t, expectedSpreadsheetId, actual.SpreadSheetId())
}

func Test_openSheetWithClient_O_CREATE(t *testing.T) {
	// prepare
	expectedSpreadsheetId := "spreadsheetId"
	expectedSheetName := "sheetName"
	expectedSheetId := int32(1)

	// mock a scenario where the sheet does not exists
	sheetResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: `{"sheets": [{}]}`,
	}
	// mock the response after sheet creation
	creationResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: fmt.Sprintf(`{
			"updatedSpreadsheet": {
				"sheets": [{
						"properties": {
							"sheetId": %d,
							"title": "%s"
						}
					}
				]
			}
		}`, expectedSheetId, expectedSheetName),
	}
	mockClient := client.CreateMockClient(sheetResponse, creationResponse)

	// test
	actual, err := openSheetWithClient(expectedSpreadsheetId, expectedSheetName, O_CREATE, mockClient)

	if err != nil {
		t.Errorf("found error %+v", err)
	}

	assertEqual(t, expectedSheetId, actual.Id())
	assertEqual(t, expectedSheetName, actual.Name())
	assertEqual(t, expectedSpreadsheetId, actual.SpreadSheetId())
}

func Test_openSheetWithClient_O_EXCL(t *testing.T) {
	// prepare
	expectedSpreadsheetId := "spreadsheetId"
	expectedSheetName := "sheetName"
	expectedSheetId := int32(1)

	// mock a scenario where the sheet exists already
	sheetResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: fmt.Sprintf(`{
			"sheets": [{
					"properties": {
						"sheetId": %d,
						"title": "sheetName"
					}
				}
			]
		}`, expectedSheetId),
	}

	mockClient := client.CreateMockClient(sheetResponse)

	// test
	actual, err := openSheetWithClient(expectedSpreadsheetId, expectedSheetName, O_CREATE|O_EXCL, mockClient)

	if !errors.Is(err, ErrExist) {
		t.Errorf("expected '%v' but found '%v'", ErrExist, err)
	}
	if actual != nil {
		t.Errorf("expected that sheet is not created")
	}
}

func Test_openSheetWithClient_Non_Existing_File(t *testing.T) {
	// mock a scenario where the sheet exists already
	sheetResponse := client.ResponseSummery{
		ResponseCode: 200,
		ResponseBody: `{"sheets": [{}]}`,
	}

	mockClient := client.CreateMockClient(sheetResponse)

	// test
	actual, err := openSheetWithClient("spreadSheetId", "sheetName", O_RDONLY, mockClient)

	if !errors.Is(err, ErrNotExist) {
		t.Errorf("expected '%v' but found '%v'", ErrNotExist, err)
	}
	if actual != nil {
		t.Errorf("expected that sheet is not created")
	}
}

func Test_RemoveSheetWithClient(t *testing.T) {
	// mock a scenario where the sheet exists already
	sheetResponse := client.ResponseSummery{
		ResponseCode: 200,
	}

	mockClient := client.CreateMockClient(sheetResponse)

	err := removeSheetWithClient("spreadSheetId", 1, mockClient)

	if err != nil {
		t.Errorf("found error '%+v'", err)
	}
}

func Test_createClient(t *testing.T) {
	client, err := createClient(context.Background(), O_CREATE, []byte(`{"type":"service_account"}`))

	if err != nil {
		t.Errorf("found error '%+v'", err)
	}
	if client == nil {
		t.Error("expected client not to be nil")
	}
}

func Test_CSV_Writer_Interface_Support(t *testing.T) {
	testSheet := &Sheet{}
	if csv.NewWriter(testSheet) == nil {
		t.Errorf("expected writer not to be nil")
	}
}

func Test_CSV_Reader_Interface_Support(t *testing.T) {
	testSheet := &Sheet{}
	if csv.NewReader(testSheet) == nil {
		t.Errorf("expected writer not to be nil")
	}
}

func assertEqual(t *testing.T, expected any, actual any) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected '%+v' but found '%+v'", expected, actual)
	}
}

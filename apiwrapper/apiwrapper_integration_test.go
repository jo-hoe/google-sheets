package apiwrapper

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/jo-hoe/google-sheets/client"
)

func Test_Integration_Replace(t *testing.T) {
	wrapper, spreadSheetId := createWrapper(t)
	id := createTestSheet(t, wrapper, spreadSheetId)

	err := wrapper.ReplaceSheetData(spreadSheetId, fmt.Sprint(id), [][]string{
		{"0", "1"},
		{"2", "3"},
	})
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
	deleteTestSheet(t, wrapper, spreadSheetId, id)
}

func Test_Integration_AppendToSheet(t *testing.T) {
	wrapper, spreadSheetId := createWrapper(t)
	id := createTestSheet(t, wrapper, spreadSheetId)

	err := wrapper.AppendToSheet(spreadSheetId, fmt.Sprint(id), [][]string{
		{"0", "1"},
		{"2", "3"},
	})
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
	err = wrapper.AppendToSheet(spreadSheetId, fmt.Sprint(id), [][]string{
		{"4", "5"},
		{"6", "7"},
	})
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}

	deleteTestSheet(t, wrapper, spreadSheetId, id)
}

func Test_Integration_Create_Delete(t *testing.T) {
	wrapper, spreadSheetId := createWrapper(t)
	expectedId := int32(time.Now().UnixMilli() / 1000)

	actualId := createTestSheetWithId(t, wrapper, spreadSheetId, expectedId)
	if expectedId != actualId {
		t.Errorf("Expected Id %d but found %d", expectedId, actualId)
	}

	deleteTestSheet(t, wrapper, spreadSheetId, actualId)
}

func createTestSheet(t *testing.T, wrapper *SheetsApiWrapper, spreadSheetId string) int32 {
	expectedId := int32(time.Now().UnixMilli() / 1000)
	return createTestSheetWithId(t, wrapper, spreadSheetId, expectedId)
}

func createTestSheetWithId(t *testing.T, wrapper *SheetsApiWrapper, spreadSheetId string, id int32) int32 {
	result, err := wrapper.CreateSheet(spreadSheetId, id, fmt.Sprint(id))
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
	return result
}

func deleteTestSheet(t *testing.T, wrapper *SheetsApiWrapper, spreadSheetId string, id int32) {
	err := wrapper.DeleteSheet(spreadSheetId, id)
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
}

func createWrapper(t *testing.T) (wrapper *SheetsApiWrapper, spreadSheetId string) {
	filePath := os.Getenv("CREDENTIALS_FILE_PATH")
	if filePath == "" {
		t.Skip("No credentials found for intergration test, skipping test")
	}
	sheetId := os.Getenv("SPREADSHEET_ID")
	if sheetId == "" {
		t.Skip("No spread sheet Id found for intergration test, skipping test")
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Skipf("Could not read file %+v", err)
	}
	httpClient, err := client.NewReadWriteScopesServiceAccountClient(context.Background(), string(content))
	if err != nil {
		t.Skipf("Could not create client %+v", err)
	}
	return NewSheetsApiWrapper(httpClient), sheetId
}

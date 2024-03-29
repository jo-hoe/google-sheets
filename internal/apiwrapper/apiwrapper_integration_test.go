package apiwrapper

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jo-hoe/google-sheets/internal/client"
)

func Test_Integration_AppendToSheet(t *testing.T) {
	wrapper, spreadSheetId := createWrapper(t)
	sheetName, id := createTestSheet(t, wrapper, spreadSheetId)

	err := wrapper.AppendToSheet(spreadSheetId, sheetName, [][]string{
		{"0", "1"},
		{"2", "3"},
	})
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
	err = wrapper.AppendToSheet(spreadSheetId, sheetName, [][]string{
		{"4", "5"},
		{"6", "7"},
	})
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}

	deleteTestSheet(t, wrapper, spreadSheetId, id)
}

func Test_Integration_Create_Get_Delete(t *testing.T) {
	wrapper, spreadSheetId := createWrapper(t)
	sheetName := fmt.Sprint(time.Now().UnixMilli() / 1000)

	createdId := createTestSheetWithId(t, wrapper, spreadSheetId, sheetName)
	foundId, err := wrapper.GetSheetId(spreadSheetId, sheetName)
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
	if createdId != foundId {
		t.Errorf("Expected Id %d but found %d", createdId, foundId)
	}

	deleteTestSheet(t, wrapper, spreadSheetId, createdId)
}

func createTestSheet(t *testing.T, wrapper *SheetsApiWrapper, spreadSheetId string) (string, int32) {
	sheetName := fmt.Sprint(time.Now().UnixMilli() / 1000)
	return sheetName, createTestSheetWithId(t, wrapper, spreadSheetId, sheetName)
}

func createTestSheetWithId(t *testing.T, wrapper *SheetsApiWrapper, spreadSheetId string, sheetName string) int32 {
	result, err := wrapper.CreateSheet(spreadSheetId, sheetName)
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
		t.Skip("No credentials found for integration test, skipping test")
	}
	sheetId := os.Getenv("SPREADSHEET_ID")
	if sheetId == "" {
		t.Skip("No spread sheet Id found for integration test, skipping test")
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Skipf("Could not read file %+v", err)
	}
	httpClient, err := client.NewReadWriteClient(context.Background(), string(content))
	if err != nil {
		t.Skipf("Could not create client %+v", err)
	}
	return NewSheetsApiWrapper(httpClient), sheetId
}

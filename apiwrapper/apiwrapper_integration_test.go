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

	err := wrapper.ReplaceSheetData(spreadSheetId, "Sheet4", [][]string{
		{"0", "1"},
		{"2", "3"},
	})
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
}

func Test_Integration_Create_Delete(t *testing.T) {
	wrapper, spreadSheetId := createWrapper(t)

	expectedId := time.Now().UnixMilli() / 1000
	actualId, err := wrapper.CreateSheet(spreadSheetId, int32(expectedId), fmt.Sprint(expectedId))
	if err != nil {
		t.Errorf("Found error during sheet creation %+v", err)
	}
	if int32(expectedId) != actualId {
		t.Errorf("Expected Id %d but found %d", expectedId, actualId)
	}

	err = wrapper.DeleteSheet(spreadSheetId, actualId)
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

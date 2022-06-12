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
	sheetTestName := fmt.Sprintf("%s-%d", "test", time.Now().UnixMilli())
	_, err := wrapper.CreateSheet(spreadSheetId, sheetTestName)
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
	httpClient, err := client.NewServiceAccountClient(context.Background(), string(content))
	if err != nil {
		t.Skipf("Could not create client %+v", err)
	}
	return NewSheetsApiWrapper(httpClient), sheetId
}

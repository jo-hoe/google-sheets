package writer

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jo-hoe/google-sheets/apiwrapper"
	"github.com/jo-hoe/google-sheets/client"
)

func TestSheetWriter_Integration_Write(t *testing.T) {
	// setup
	client, spreadSheetId := createClient(t)
	wrapper := apiwrapper.NewSheetsApiWrapper(client)
	sheetName := time.Now().UnixMilli() / 1000

	// test
	writer, err := NewSheetWriter(client, spreadSheetId, fmt.Sprint(sheetName), O_CREATE)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	csvWriter := csv.NewWriter(writer)
	err = csvWriter.WriteAll([][]string{
		{"0", "1"},
		{"2", "3"},
	})
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	err = csvWriter.Write([]string{"4", "5"})
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	csvWriter.Flush()

	// clean-up
	sheetId, err := wrapper.GetSheetId(spreadSheetId, fmt.Sprint(sheetName))
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	err = wrapper.DeleteSheet(spreadSheetId, sheetId)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
}

func createClient(t *testing.T) (httpClient *http.Client, spreadSheetId string) {
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
	httpClient, err = client.NewReadWriteScopesServiceAccountClient(context.Background(), string(content))
	if err != nil {
		t.Skipf("Could not create client %+v", err)
	}
	return httpClient, sheetId
}

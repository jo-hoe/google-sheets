package gs

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSheet_Integration_Write(t *testing.T) {
	// setup
	fileContent, spreadSheetId := getPrerequisites(t)
	sheetName := fmt.Sprint(time.Now().UnixMilli() / 1000)

	// test
	sheet, err := OpenSheet(context.Background(), spreadSheetId, sheetName, O_CREATE, fileContent)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	csvWriter := csv.NewWriter(sheet)
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

	csvReader := csv.NewReader(sheet)
	actual, err := csvReader.ReadAll()
	if err != nil {
		t.Errorf("Found error %+v", err)
	}

	expected := [][]string{{"0", "1"}, {"2", "3"}, {"4", "5"}}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v but got %+v", expected, actual)
	}

	// clean-up
	err = Remove(context.Background(), sheet.SpreadSheetId(), sheet.Id(), fileContent)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
}

func TestSheet_Integration_Read_Incomplete(t *testing.T) {
	// setup
	fileContent, spreadSheetId := getPrerequisites(t)
	sheetName := fmt.Sprint(time.Now().UnixMilli() / 1000)

	// test
	sheet, err := OpenSheet(context.Background(), spreadSheetId, sheetName, O_CREATE, fileContent)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	csvWriter := csv.NewWriter(sheet)
	err = csvWriter.WriteAll([][]string{
		{"0", "1"},
	})
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	err = csvWriter.Write([]string{"2"})
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
	csvWriter.Flush()

	csvReader := csv.NewReader(sheet)
	csvReader.FieldsPerRecord = -1
	actual, err := csvReader.ReadAll()
	if err != nil {
		t.Errorf("Found error %+v", err)
	}

	expected := [][]string{{"0", "1"}, {"2"}}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %+v but got %+v", expected, actual)
	}

	// clean-up
	err = Remove(context.Background(), sheet.SpreadSheetId(), sheet.Id(), fileContent)
	if err != nil {
		t.Errorf("Found error %+v", err)
	}
}

func getPrerequisites(t *testing.T) (fileContent []byte, spreadSheetId string) {
	filePath := os.Getenv("CREDENTIALS_FILE_PATH")
	if filePath == "" {
		t.Skip("No credentials found for integration test, skipping test")
	}
	spreadSheetId = os.Getenv("SPREADSHEET_ID")
	if spreadSheetId == "" {
		t.Skip("No spread sheet Id found for integration test, skipping test")
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Skipf("Could not read file %+v", err)
	}

	return fileContent, spreadSheetId
}

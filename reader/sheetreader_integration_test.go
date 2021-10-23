package reader

import (
	"encoding/csv"
	"net/http"
	"testing"
)

const baseUrl = "https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/"
const spreadSheatId = "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
const sheetName = "Class Data"

func TestIntegrationNewSheetReader(t *testing.T) {
	readerCloser, err := NewSheetReader(&http.Client{}, spreadSheatId, sheetName)
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	} else {
		defer readerCloser.Close()
	}

	csv := csv.NewReader(readerCloser)

	items, err := csv.ReadAll()
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	if len(items) == 0 {
		t.Errorf("no items where found in '%s'", baseUrl+spreadSheatId)
	}
}

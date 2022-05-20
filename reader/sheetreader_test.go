package reader

import (
	"encoding/csv"
	"reflect"
	"testing"

	"github.com/jo-hoe/google-sheets/apiwrapper"
)

func Test_NewSheetReader(t *testing.T) {
	readerCloser, err := NewSheetReader(apiwrapper.CreateMockClient(), "spreadSheatId", "sheetName")
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	} else {
		defer readerCloser.Close()
	}

	csv := csv.NewReader(readerCloser)

	actual, err := csv.ReadAll()
	if err != nil {
		t.Errorf("error found during http reqest %v", err)
	}

	expected := [][]string{
		{"0", "1"},
		{"2", "3"},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected '%v' found '%v'", expected, actual)
	}
}

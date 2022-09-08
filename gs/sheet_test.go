package gs

import (
	"testing"
)

func TestSheet_Id(t *testing.T) {
	expected := 1
	sheet := Sheet{
		id: int32(expected),
	}

	if sheet.Id() != int32(expected){
		t.Errorf("expected %d actual %d", expected, sheet.Id())
	}
}

func TestSheet_SpreadSheetId(t *testing.T) {
	expected := "testName"
	sheet := Sheet{
		spreadSheetId: expected,
	}

	if sheet.SpreadSheetId() != expected{
		t.Errorf("expected %s actual %s", expected, sheet.SpreadSheetId())
	}
}

func TestSheet_Name(t *testing.T) {
	expected := "testName"
	sheet := Sheet{
		sheetName: expected,
	}

	if sheet.Name() != expected{
		t.Errorf("expected %s actual %s", expected, sheet.Name())
	}
}

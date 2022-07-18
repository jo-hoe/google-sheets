package writer

import (
	"bytes"
	"encoding/csv"
	"io"
	"net/http"

	"github.com/jo-hoe/google-sheets/api/apiwrapper"
)

type SheetWriter struct {
	io.Writer
	wrapper       *apiwrapper.SheetsApiWrapper
	spreadSheetId string
	sheetName     string
}

func NewSheetWriter(client *http.Client, spreadSheetId string, sheetName string) (*SheetWriter, error) {
	wrapper := apiwrapper.NewSheetsApiWrapper(client)

	return &SheetWriter{
		wrapper:       wrapper,
		spreadSheetId: spreadSheetId,
		sheetName:     sheetName,
	}, nil
}

func (service *SheetWriter) Write(byteData []byte) (n int, err error) {
	csvReader := csv.NewReader(bytes.NewReader(byteData))
	data, err := csvReader.ReadAll()
	if err != nil {
		return 0, err
	}

	err = service.wrapper.AppendToSheet(service.spreadSheetId, service.sheetName, data)
	if err != nil {
		return 0, err
	}

	return len(byteData), nil
}

func hasFlag(flags int, flag int) bool {
	return flags&flag != 0
}

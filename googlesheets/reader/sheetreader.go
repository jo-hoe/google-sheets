package reader

import (
	"io"
	"net/http"

	"github.com/jo-hoe/google-sheets/api/apiwrapper"
)

type SheetReader struct {
	io.Reader
	reader        io.Reader
	spreadSheetId string
	sheetName     string
	wrapper       *apiwrapper.SheetsApiWrapper
}

func NewSheetReader(client *http.Client, spreadSheetId string, sheetName string) (*SheetReader, error) {
	return &SheetReader{
		wrapper:       apiwrapper.NewSheetsApiWrapper(client),
		spreadSheetId: spreadSheetId,
		sheetName:     sheetName,
	}, nil
}

func (service *SheetReader) Read(p []byte) (n int, err error) {
	if service.reader == nil {
		service.reader, err = service.wrapper.GetSheetData(service.spreadSheetId, service.sheetName)
		if err != nil {
			return -1, err
		}
	}

	return service.reader.Read(p)
}

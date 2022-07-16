package reader

import (
	"io"
	"net/http"

	"github.com/jo-hoe/google-sheets/apiwrapper"
)

// wrapper for io.ReaderCloser
type SheetReader struct {
	io.Reader
	spreadSheetId string
	sheetName     string
	wrapper       *apiwrapper.SheetsApiWrapper
}

func NewSheetReader(client *http.Client, spreadSheetId string, sheetName string) (*SheetReader, error) {
	return &SheetReader{
		wrapper: apiwrapper.NewSheetsApiWrapper(client),
	}, nil
}

func (service *SheetReader) Read(p []byte) (n int, err error) {
	reader, err := service.wrapper.GetSheetData(service.spreadSheetId, service.sheetName)
	if err != nil {
		return -1, err
	}
	return reader.Read(p)
}

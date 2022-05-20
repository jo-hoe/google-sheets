package reader

import (
	"io"
	"net/http"

	"github.com/jo-hoe/google-sheets/apiwrapper"
)

// wrapper for io.ReaderCloser
type SheetReadCloser struct {
	io.ReadCloser
	readerCloser io.ReadCloser
}

func NewSheetReader(client *http.Client, spreadSheetId string, sheetName string) (*SheetReadCloser, error) {
	wrapper := apiwrapper.NewSheetsApiWrapper(client)
	readerCloser, err := wrapper.GetSheetData(spreadSheetId, sheetName)
	if err != nil {
		return nil, err
	}
	return &SheetReadCloser{
		readerCloser: readerCloser,
	}, nil
}

func (sheetReadCloser *SheetReadCloser) Read(p []byte) (n int, err error) {
	return sheetReadCloser.readerCloser.Read(p)
}

func (sheetReadCloser *SheetReadCloser) Close() error {
	return sheetReadCloser.readerCloser.Close()
}

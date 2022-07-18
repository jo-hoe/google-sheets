package sheet

import (
	"io"
	"net/http"

	"github.com/jo-hoe/google-sheets/sheet/reader"
	"github.com/jo-hoe/google-sheets/sheet/writer"
)

type Sheet struct {
	io.Reader
	io.Writer
	id            int32
	sheetName     string
	spreadSheetId string
	client        *http.Client
	writer        *writer.SheetWriter
	reader        *reader.SheetReader
}

func (service *Sheet) Write(byteData []byte) (n int, err error) {
	return service.writer.Write(byteData)
}

func (service *Sheet) Read(p []byte) (n int, err error) {
	return service.reader.Read(p)
}

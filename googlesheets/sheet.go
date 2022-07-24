package google

import (
	"io"

	"github.com/jo-hoe/google-sheets/googlesheets/reader"
	"github.com/jo-hoe/google-sheets/googlesheets/writer"
)

type Sheet struct {
	io.ReadWriter
	id            int32
	sheetName     string
	spreadSheetId string
	writer        *writer.SheetWriter
	reader        *reader.SheetReader
}

func (service *Sheet) Write(byteData []byte) (n int, err error) {
	return service.writer.Write(byteData)
}

func (service *Sheet) Read(p []byte) (n int, err error) {
	return service.reader.Read(p)
}

func (service *Sheet) Id() int32 {
	return service.id
}

func (service *Sheet) SpreadSheetId() string {
	return service.spreadSheetId
}

func (service *Sheet) Name() string {
	return service.sheetName
}

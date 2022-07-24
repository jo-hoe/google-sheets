package gs

import (
	"io"

	"github.com/jo-hoe/google-sheets/gs/reader"
	"github.com/jo-hoe/google-sheets/gs/writer"
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

// Returns the ID of the sheet
func (service *Sheet) Id() int32 {
	return service.id
}

// Returns the spreadsheet ID
func (service *Sheet) SpreadSheetId() string {
	return service.spreadSheetId
}

// Returns the name of the Sheet
func (service *Sheet) Name() string {
	return service.sheetName
}

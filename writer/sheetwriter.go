package writer

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/jo-hoe/google-sheets/apiwrapper"
)

type SheetWriter struct {
	io.Writer
	wrapper       *apiwrapper.SheetsApiWrapper
	spreadSheetId string
	sheetName     string
	flag          int
}

const (
	// The remaining values may be or'ed in to control behavior.
	O_CREATE int = syscall.O_CREAT // create a new sheet if none exists.
	O_EXCL   int = syscall.O_EXCL  // used with O_CREATE, sheet must not exist.
	O_TRUNC  int = syscall.O_TRUNC // truncate regular writable sheet when opened.
)

func NewSheetWriter(client *http.Client, spreadSheetId string, sheetName string, flag int) (*SheetWriter, error) {
	wrapper := apiwrapper.NewSheetsApiWrapper(client)

	// check if file exists
	_, err := wrapper.GetSheetId(spreadSheetId, sheetName)

	if err == nil && hasFlag(flag, O_EXCL) && hasFlag(flag, O_CREATE) {
		return nil, fmt.Errorf("sheet %s already exists in spreadsheet %s", sheetName, spreadSheetId)
	}
	if err == nil && hasFlag(flag, O_TRUNC) {
		err = wrapper.ClearSheet(spreadSheetId, sheetName)
		if err != nil {
			return nil, err
		}
	}
	if err != nil && hasFlag(flag, O_CREATE) {
		// create new with an id = current timestamp
		_, err = wrapper.CreateSheet(spreadSheetId, int32(time.Now().UnixMilli()/1000), sheetName)
		if err != nil {
			return nil, err
		}
	}

	return &SheetWriter{
		wrapper:       wrapper,
		spreadSheetId: spreadSheetId,
		sheetName:     sheetName,
		flag:          flag,
	}, nil
}

func (service *SheetWriter) Write(p []byte) (n int, err error) {
	stringRepresentation := string(p)

	csvReader := csv.NewReader(strings.NewReader(stringRepresentation))
	data, err := csvReader.ReadAll()
	if err != nil {
		return 0, err
	}

	err = service.wrapper.AppendToSheet(service.spreadSheetId, service.sheetName, data)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func hasFlag(flags int, flag int) bool {
	return flags&flag != 0
}

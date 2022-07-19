package sheet

import (
	"context"
	"fmt"
	"net/http"
	"syscall"

	"github.com/jo-hoe/google-sheets/api/apiwrapper"
	"github.com/jo-hoe/google-sheets/api/client"
	"github.com/jo-hoe/google-sheets/sheet/reader"
	"github.com/jo-hoe/google-sheets/sheet/writer"
)

const (
	// Exactly one of O_RDONLY or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the sheet read-only.
	O_RDWR   int = syscall.O_RDWR   // open the sheet read-write.
	// The remaining values may be or'ed in to control behavior.
	O_CREATE int = syscall.O_CREAT // create a new sheet if none exists.
	O_EXCL   int = syscall.O_EXCL  // used with O_CREATE, sheet must not exist.
	O_TRUNC  int = syscall.O_TRUNC // truncate regular writable sheet when opened.
)

func Remove(ctx context.Context, spreadSheetId string, sheetId int32, clientCredentialsJson []byte) error {
	client, err := createClient(ctx, O_RDWR, clientCredentialsJson)
	if err != nil {
		return err
	}
	return RemoveSheetWithClient(spreadSheetId, sheetId, client)
}

func RemoveSheetWithClient(spreadSheetId string, sheetId int32, client *http.Client) error {
	wrapper := apiwrapper.NewSheetsApiWrapper(client)

	return wrapper.DeleteSheet(spreadSheetId, sheetId)
}

func OpenSheet(ctx context.Context, spreadSheetId string, sheetName string, flag int, clientCredentialsJson []byte) (*Sheet, error) {
	client, err := createClient(ctx, flag, clientCredentialsJson)
	if err != nil {
		return nil, err
	}
	return OpenSheetWithClient(spreadSheetId, sheetName, flag, client)
}

func OpenSheetWithClient(spreadSheetId string, sheetName string, flag int, client *http.Client) (*Sheet, error) {
	wrapper := apiwrapper.NewSheetsApiWrapper(client)

	// check if file exists
	id, err := wrapper.GetSheetId(spreadSheetId, sheetName)
	sheetExists := err == nil

	if sheetExists {
		if hasFlag(flag, O_EXCL) && hasFlag(flag, O_CREATE) {
			// if file exist and should not -> return error
			return nil, fmt.Errorf("sheet %s already exists in spreadsheet %s", sheetName, spreadSheetId)
		}
		if hasFlag(flag, O_TRUNC) {
			// if file exists and content should be truncated -> clear sheet
			err = wrapper.ClearSheet(spreadSheetId, sheetName)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if hasFlag(flag, O_CREATE) {
			// create new with an id = current timestamp
			id, err = wrapper.CreateSheet(spreadSheetId, sheetName)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("sheet with name '%s' not found in spreadsheet '%s'", sheetName, spreadSheetId)
		}
	}

	reader, err := reader.NewSheetReader(client, spreadSheetId, sheetName)
	if err != nil {
		return nil, err
	}

	writer, err := writer.NewSheetWriter(client, spreadSheetId, sheetName)
	if err != nil {
		return nil, err
	}

	return &Sheet{
		id:            id,
		sheetName:     sheetName,
		spreadSheetId: spreadSheetId,
		reader:        reader,
		writer:        writer,
	}, nil
}

func hasFlag(flags int, flag int) bool {
	return flags&flag != 0
}

func createClient(ctx context.Context, flag int, clientCredentialsJson []byte) (*http.Client, error) {
	var scope string
	if hasFlag(flag, O_RDONLY) {
		scope = client.ReadOnlyScopes
	} else {
		scope = client.ReadWriteScopes
	}

	client, err := client.NewServiceAccountClient(ctx, clientCredentialsJson, scope)
	if err != nil {
		return nil, err
	}
	return client, err
}

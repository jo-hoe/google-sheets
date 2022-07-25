package gs

import (
	"context"
	"errors"
	"net/http"
	"syscall"

	"github.com/jo-hoe/google-sheets/gs/reader"
	"github.com/jo-hoe/google-sheets/gs/writer"
	"github.com/jo-hoe/google-sheets/internal/apiwrapper"
	"github.com/jo-hoe/google-sheets/internal/client"
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

var (
	ErrInvalid  = errors.New("invalid argument")     // "invalid argument"
	ErrExist    = errors.New("sheet already exists") // "file already exists"
	ErrNotExist = errors.New("sheet does not exist") // "file does not exist"
)

// Remove removes the sheet in a given spreadspeed.
func Remove(ctx context.Context, spreadSheetId string, sheetId int32, clientCredentialsJson []byte) error {
	client, err := createClient(ctx, O_RDWR, clientCredentialsJson)
	if err != nil {
		return err
	}
	return removeSheetWithClient(spreadSheetId, sheetId, client)
}

// OpenSheet is the generalized open call. It opens the sheet with specified flag (O_RDONLY etc.).
// If the sheet does not exist, and the O_CREATE flag is passed, it is created.
// If successful, methods on the returned Sheet can be used for csv I/O.
//
// Can also be used to check if a given file exists.
// To do so analysis the returned error like so errors.Is(err, gs.ErrExist).
func OpenSheet(ctx context.Context, spreadSheetId string, sheetName string, flag int, clientCredentialsJson []byte) (*Sheet, error) {
	client, err := createClient(ctx, flag, clientCredentialsJson)
	if err != nil {
		return nil, err
	}
	return openSheetWithClient(spreadSheetId, sheetName, flag, client)
}

func removeSheetWithClient(spreadSheetId string, sheetId int32, client *http.Client) error {
	wrapper := apiwrapper.NewSheetsApiWrapper(client)

	return wrapper.DeleteSheet(spreadSheetId, sheetId)
}

func openSheetWithClient(spreadSheetId string, sheetName string, flag int, client *http.Client) (*Sheet, error) {
	if client == nil {
		return nil, ErrInvalid
	}

	wrapper := apiwrapper.NewSheetsApiWrapper(client)

	// check if file exists
	id, err := wrapper.GetSheetId(spreadSheetId, sheetName)
	if err != nil {
		return nil, err
	}
	sheetExists := id > -1

	if sheetExists {
		if hasFlag(flag, O_EXCL) && hasFlag(flag, O_CREATE) {
			// if file exist and should not -> return error
			return nil, ErrExist
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
			return nil, ErrNotExist
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

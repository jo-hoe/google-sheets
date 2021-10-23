package reader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SheetReadCloser struct {
	io.ReadCloser
	readerCloser io.ReadCloser
}

// good description of how the URL is build can be found in this post https://stackoverflow.com/a/33727897/11951869
// offical information sources:
// - https://developers.google.com/chart/interactive/docs/spreadsheets
// - https://developers.google.com/chart/interactive/docs/querylanguage
// - https://developers.google.com/chart/interactive/docs/dev/implementing_data_source
const csvUrlTemplate = "https://docs.google.com/spreadsheets/d/%s/gviz/tq?tqx=out:csv&sheet=%s"

func NewSheetReader(client *http.Client, spreadSheatId string, sheetName string) (*SheetReadCloser, error) {
	readerCloser, err := getFile(client, spreadSheatId, sheetName)
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

func getFile(httpClient *http.Client, spreadSheatId string, sheetName string) (io.ReadCloser, error) {
	// escape sheet name, since it may contain spaces and other URL incompatible characters
	encodedSheetName := url.QueryEscape(sheetName)
	url := fmt.Sprintf(csvUrlTemplate, spreadSheatId, encodedSheetName)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not get sheet: %d - %s", resp.StatusCode, resp.Status)
	}

	return resp.Body, nil
}

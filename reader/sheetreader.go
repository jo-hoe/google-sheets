package reader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// wrapper for io.ReaderCloser
type SheetReadCloser struct {
	io.ReadCloser
	readerCloser io.ReadCloser
}

type partialSheetResult struct {
	Values [][]string `json:"values"`
}

// url is reverse engineered from:
// https://github.com/googleapis/google-api-go-client/blob/bc181c33247b7fe3d06d2d7139da0fa06fabbd71/sheets/v4/sheets-gen.go#L14283
// but it is also described here:
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
const csvUrlTemplate = "https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?alt=json&prettyPrint=false"

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
		return nil, fmt.Errorf("could not get sheet: %d: %s", resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	return truncateExtraneousData(resp.Body)
}

// the api returns something like:
// {"range":"Sheet2!A1:Z1000","majorDimension":"ROWS","values":[["a","b"],["1","2"]]}
// only the values field will be contained in the reader
func truncateExtraneousData(reader io.ReadCloser) (io.ReadCloser, error) {
	// not the fastest way to do thing but easy to read and maintain

	// read complete string
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	// unmarshal into struct
	stringOutput := buffer.String()
	result := partialSheetResult{}
	json.Unmarshal([]byte(stringOutput), &result)

	// convert string back to reader closer
	jsonBytes, err := json.Marshal(result.Values)
	if err != nil {
		return nil, err
	}
	if result.Values == nil {
		return nil, fmt.Errorf("could not read 'values' api answer from %v", stringOutput)
	}
	return ioutil.NopCloser(bytes.NewReader(jsonBytes)), nil
}

package reader

import (
	"bytes"
	"encoding/csv"
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

func NewSheetReader(client *http.Client, spreadSheetId string, sheetName string) (*SheetReadCloser, error) {
	readerCloser, err := getFile(client, spreadSheeId, sheetName)
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

func getFile(httpClient *http.Client, spreadSheetId string, sheetName string) (io.ReadCloser, error) {
	// escape sheet name, since it may contain spaces and other URL incompatible characters
	encodedSheetName := url.QueryEscape(sheetName)
	url := fmt.Sprintf(csvUrlTemplate, spreadSheetId, encodedSheetName)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not get sheet from url '%s'\nerror %d: %s", url, resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	return truncateExtraneousData(resp.Body)
}

// The api returns something like:
// {"range":"Sheet2!A1:Z1000","majorDimension":"ROWS","values":[["a","b"],["1","2"]]}
// only the 'values' field is relevant. 
// The returned reader will contain it in an 'encoding/csv' readable format
func truncateExtraneousData(reader io.ReadCloser) (io.ReadCloser, error) {
	// not the fastest way to do things but easy to read and maintain

	// read complete string
	buffer := new(bytes.Buffer)
	_, err := buffer.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	// unmarshal to struct
	result := partialSheetResult{}
	err = json.Unmarshal(buffer.Bytes(), &result)
	if err != nil {
		return nil, err
	}
	if result.Values == nil {
		return nil, fmt.Errorf("could not read 'values' api answer from %v", buffer.String())
	}

	// write slices to csv data
	output := &bytes.Buffer{}
	writer := csv.NewWriter(output)
	for _, value := range result.Values {
		err := writer.Write(value)
		if err != nil {
			return nil, err
		}
	}
	writer.Flush()
	err = writer.Error()
	if err != nil {
		return nil, err
	}

	// return data as reader
	return ioutil.NopCloser(bytes.NewReader(output.Bytes())), nil
}

package apiwrapper

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

const baseUrl = "https://sheets.googleapis.com/v4/spreadsheets/%s"

// url is reverse engineered from:
// https://github.com/googleapis/google-api-go-client/blob/bc181c33247b7fe3d06d2d7139da0fa06fabbd71/sheets/v4/sheets-gen.go#L14283
// but it is also described here:
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
const csvUrlTemplate = baseUrl + "/values/%s?alt=json&prettyPrint=false"

type partialSheetResult struct {
	Values [][]string `json:"values"`
}

type spreadSheets struct {
	Sheets []struct {
		Properties struct {
			SheetID int    `json:"sheetId"`
			Title   string `json:"title"`
		} `json:"properties"`
	} `json:"sheets"`
}

type SheetsApiWrapper struct {
	httpClient *http.Client
}

func NewSheetsApiWrapper(httpClient *http.Client) *SheetsApiWrapper {
	return &SheetsApiWrapper{
		httpClient: httpClient,
	}
}

func (wrapper SheetsApiWrapper) GetSheetData(spreadSheetId string, sheetName string) (io.ReadCloser, error) {
	// escape sheet name, since it may contain spaces and other URL incompatible characters
	encodedSheetName := url.QueryEscape(sheetName)
	url := fmt.Sprintf(csvUrlTemplate, spreadSheetId, encodedSheetName)
	resp, err := wrapper.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("could not get sheet from url '%s'\nerror %d: %s", url, resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	return truncateExtraneousData(resp.Body)
}

func (wrapper SheetsApiWrapper) GetSheetId(spreadSheetId string, sheetName string) (int, error) {
	url := fmt.Sprintf(baseUrl, spreadSheetId)
	resp, err := wrapper.httpClient.Get(url)
	if err != nil {
		return -1, err
	}

	result, err := deserialize[spreadSheets](resp.Body, spreadSheets{})
	if err != nil {
		return -1, err
	}
	
	for _, sheet := range result.(spreadSheets).Sheets {
		if sheet.Properties.Title == sheetName{
			return sheet.Properties.SheetID, nil
		}
	}
	return -1, nil
}

func deserialize[T any](reader io.ReadCloser, in any) (out any, err error) {
	defer reader.Close()
	
	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	// unmarshal to struct
	result := spreadSheets{}
	err = json.Unmarshal(buffer.Bytes(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
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

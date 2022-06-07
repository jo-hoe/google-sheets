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

const baseUrl = "https://sheets.googleapis.com/v4/SpreadSheets/%s"

// url is reverse engineered from:
// https://github.com/googleapis/google-api-go-client/blob/bc181c33247b7fe3d06d2d7139da0fa06fabbd71/sheets/v4/sheets-gen.go#L14283
// but it is also described here:
// https://developers.google.com/sheets/api/reference/rest/v4/SpreadSheets.values/get
const csvUrlTemplate = baseUrl + "/values/%s?alt=json&prettyPrint=false"
const updateUrl = "https://sheets.googleapis.com/v4/SpreadSheets/%s:batchUpdate"

type partialSheetResult struct {
	Values [][]string `json:"values"`
}

type SpreadSheet struct {
	Sheets []Sheet `json:"sheets"`
}

type Sheet struct {
	Properties SpreadSheetProperties `json:"properties"`
}

type SpreadSheetProperties struct {
	SheetID int    `json:"sheetId"`
	Title   string `json:"title"`
}

type SheetsApiWrapper struct {
	httpClient *http.Client
}

func NewSheetsApiWrapper(httpClient *http.Client) *SheetsApiWrapper {
	return &SheetsApiWrapper{
		httpClient: httpClient,
	}
}

func (wrapper SheetsApiWrapper) CreateSheet(SpreadSheetId string, sheetName string) (out *Sheet, err error) {
	body := SpreadSheet{}
	body.Sheets = append(body.Sheets, Sheet{
		SpreadSheetProperties{
			Title: sheetName,
		},
	})

	payloadBuf := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuf).Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf(updateUrl, SpreadSheetId), payloadBuf)
	if err != nil {
		return nil, err
	}

	response, err := wrapper.httpClient.Do(req)
	result := Sheet{}
	err = deserialize[SpreadSheet](response.Body, &result)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Response was '%d': %s", response.StatusCode, response.Status)
	}

	return &result, nil
}

func (wrapper SheetsApiWrapper) GetSheetData(SpreadSheetId string, sheetName string) (io.ReadCloser, error) {
	// escape sheet name, since it may contain spaces and other URL incompatible characters
	encodedSheetName := url.QueryEscape(sheetName)
	url := fmt.Sprintf(csvUrlTemplate, SpreadSheetId, encodedSheetName)
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

func (wrapper SheetsApiWrapper) GetSheetId(SpreadSheetId string, sheetName string) (int, error) {
	url := fmt.Sprintf(baseUrl, SpreadSheetId)
	resp, err := wrapper.httpClient.Get(url)
	if err != nil {
		return -1, err
	}

	result := SpreadSheet{}
	err = deserialize[SpreadSheet](resp.Body, &result)
	if err != nil {
		return -1, err
	}

	for _, sheet := range result.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetID, nil
		}
	}
	return -1, fmt.Errorf("sheet was not found")
}

func deserialize[T any](reader io.ReadCloser, in any) (err error) {
	defer reader.Close()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(reader)
	if err != nil {
		return err
	}
	// unmarshal to struct
	err = json.Unmarshal(buffer.Bytes(), in)
	if err != nil {
		return err
	}
	return nil
}

// The api returns something like:
// {"range":"Sheet2!A1:Z1000","majorDimension":"ROWS","values":[["a","b"],["1","2"]]}
// only the 'values' field is relevant.
// The returned reader will contain it in an 'encoding/csv' readable format
func truncateExtraneousData(reader io.ReadCloser) (io.ReadCloser, error) {
	// not the fastest way to do things but easy to read and maintain
	result := partialSheetResult{}
	err := deserialize[partialSheetResult](reader, &result)
	if err != nil {
		return nil, err
	}
	if result.Values == nil {
		return nil, fmt.Errorf("could not read 'values' api")
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

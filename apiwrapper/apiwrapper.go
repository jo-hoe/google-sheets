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
const updateUrl = "https://sheets.googleapis.com/v4/spreadsheets/%s:batchUpdate"

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

type duplicateSheet struct {
	DuplicateSheet struct {
		SourceSheetId    int    `json:"sourceSheetId"`
		InsertSheetIndex int    `json:"insertSheetIndex"`
		NewSheetName     string `json:"newSheetName"`
	} `json:"duplicateSheet"`
}

type SheetsApiWrapper struct {
	httpClient *http.Client
}

func NewSheetsApiWrapper(httpClient *http.Client) *SheetsApiWrapper {
	return &SheetsApiWrapper{
		httpClient: httpClient,
	}
}

func (wrapper SheetsApiWrapper) DuplicateSheet(spreadSheetId string, sheetId int, newSheetName string) error {
	body := duplicateSheet{}
	body.DuplicateSheet.InsertSheetIndex = 0
	body.DuplicateSheet.SourceSheetId = sheetId
	body.DuplicateSheet.NewSheetName = newSheetName

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf(updateUrl, spreadSheetId), payloadBuf)

	response, err := wrapper.httpClient.Do(req)

	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("Response was '%d': %s", response.StatusCode, response.Status)
	}

	return nil
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

	result := spreadSheets{}
	err = deserialize[spreadSheets](resp.Body, &result)
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

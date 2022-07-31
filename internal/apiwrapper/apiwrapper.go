package apiwrapper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const baseUrl = "https://sheets.googleapis.com/v4/spreadsheets/%s"

// url is reverse engineered from:
// https://github.com/googleapis/google-api-go-client/blob/bc181c33247b7fe3d06d2d7139da0fa06fabbd71/sheets/v4/sheets-gen.go#L14283
// but it is also described here:
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
const csvUrlTemplate = baseUrl + "/values/%s?alt=json&prettyPrint=false"
const updateSheetUrl = baseUrl + ":batchUpdate"
const clearSheetUrl = baseUrl + "/values/%s:clear"
const appendSheetUrl = baseUrl + "/values/%s:append"

const majorDimension = "ROWS"

// https://developers.google.com/sheets/api/reference/rest/v4/ValueInputOption
const valueInputOption = "RAW"

type values struct {
	Values [][]string `json:"values"`
}

type spreadSheet struct {
	Sheets []sheet `json:"sheets"`
}

type sheet struct {
	Properties spreadSheetProperties `json:"properties"`
}

type spreadSheetProperties struct {
	SheetID int32  `json:"sheetId,omitempty"`
	Title   string `json:"title,omitempty"`
}

type updateRequest struct {
	Request                      []batchRequest `json:"requests,omitempty"`
	IncludeSpreadsheetInResponse bool           `json:"includeSpreadsheetInResponse"`
	ResponseRanges               []string       `json:"responseRanges,omitempty"`
	ResponseIncludeGridData      bool           `json:"responseIncludeGridData"`
}

type batchRequest struct {
	DeleteSheet *deleteSheet `json:"deleteSheet,omitempty"`
	AddSheet    *addSheet    `json:"addSheet,omitempty"`
}

type batchResponse struct {
	UpdatedSpreadsheet updatedSpreadsheet `json:"updatedSpreadsheet,omitempty"`
}

type updatedSpreadsheet struct {
	Sheets []sheet `json:"sheets,omitempty"`
}

type addSheet struct {
	Properties spreadSheetProperties `json:"properties,omitempty"`
}

type deleteSheet struct {
	SheetId int32 `json:"sheetId"`
}

type valueRange struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

type SheetsApiWrapper struct {
	httpClient *http.Client
}

func NewSheetsApiWrapper(httpClient *http.Client) *SheetsApiWrapper {
	return &SheetsApiWrapper{
		httpClient: httpClient,
	}
}

func (wrapper SheetsApiWrapper) DeleteSheet(spreadSheetId string, sheetId int32) (err error) {
	body := updateRequest{}
	body.Request = []batchRequest{{
		DeleteSheet: &deleteSheet{
			SheetId: sheetId,
		}}}

	response, err := wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheetId), body)
	if response != nil {
		response.Close()
	}

	if err != nil {
		return err
	}

	return nil
}

func (wrapper SheetsApiWrapper) GetSheetData(spreadSheetId string, sheetName string) (io.Reader, error) {
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

func (wrapper SheetsApiWrapper) CreateSheet(spreadSheetId string, sheetName string) (id int32, err error) {
	body := updateRequest{}
	body.IncludeSpreadsheetInResponse = true
	body.Request = []batchRequest{{
		AddSheet: &addSheet{
			Properties: spreadSheetProperties{
				Title: sheetName,
			},
		}}}

	response, err := wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheetId), body)
	if err != nil {
		return -1, err
	}

	result := batchResponse{}
	err = deserialize[spreadSheet](response, &result)
	if err != nil {
		return -1, err
	}

	return wrapper.findSheetIdInResponse(result.UpdatedSpreadsheet.Sheets, sheetName)
}

// Returns the id of a sheet with a given spreadSheetId and sheetName
// If the sheet does not exist, the sheetId will be -1 and err will be nil.
// In case an issue with the API or deserizalization occurs, the error is returned.
func (wrapper SheetsApiWrapper) GetSheetId(spreadSheetId string, sheetName string) (sheetId int32, err error) {
	url := fmt.Sprintf(baseUrl, spreadSheetId)
	resp, err := wrapper.httpClient.Get(url)

	if err != nil {
		return -1, err
	}
	if resp.StatusCode != 200 {
		return -1, fmt.Errorf("could not get sheet from url '%s'\nerror %d: %s", url, resp.StatusCode, resp.Status)
	}

	result := spreadSheet{}
	err = deserialize[spreadSheet](resp.Body, &result)
	if err != nil {
		return -1, err
	}

	return wrapper.findSheetIdInResponse(result.Sheets, sheetName)
}

func (wrapper SheetsApiWrapper) AppendToSheet(spreadSheetId string, sheetName string, data [][]string) (err error) {
	body := valueRange{}
	body.Range = sheetName
	body.MajorDimension = majorDimension
	body.Values = data

	queryParameters := make(map[string]string)
	queryParameters["valueInputOption"] = valueInputOption

	response, err := wrapper.postSheetRequestQueryParameter(fmt.Sprintf(appendSheetUrl, spreadSheetId, sheetName), body, queryParameters)
	if response != nil {
		response.Close()
	}
	if err != nil {
		return err
	}

	return nil
}

// delete all data from a sheet
func (wrapper SheetsApiWrapper) ClearSheet(spreadSheetId string, sheetName string) (err error) {
	_, err = wrapper.postSheetRequest(fmt.Sprintf(clearSheetUrl, spreadSheetId, sheetName), nil)
	return err
}

func (wrapper SheetsApiWrapper) postSheetRequest(url string, body any) (out io.ReadCloser, err error) {
	return wrapper.postSheetRequestQueryParameter(url, body, map[string]string{})
}

func (wrapper SheetsApiWrapper) postSheetRequestQueryParameter(url string, body any, queryParams map[string]string) (out io.ReadCloser, err error) {
	var request *http.Request
	if body != nil {
		request, err = wrapper.createJSONPostRequest(url, body)
	} else {
		request, err = http.NewRequest("POST", url, nil)
	}
	if err != nil {
		return nil, err
	}

	// add query parameters
	if len(queryParams) > 0 {
		query := request.URL.Query()
		for key, value := range queryParams {
			query.Add(key, value)
		}
		request.URL.RawQuery = query.Encode()
	}

	response, err := wrapper.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Response was '%s': %s", response.Status, string(responseBody))
	}

	return response.Body, nil
}

func (wrapper SheetsApiWrapper) findSheetIdInResponse(allSheets []sheet, sheetName string) (id int32, err error) {
	for _, sheet := range allSheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetID, nil
		}
	}
	return -1, nil
}

func (wrapper SheetsApiWrapper) createJSONPostRequest(url string, body any) (request *http.Request, err error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if body != nil {
		request.Header.Add("Content-Type", "application/json")
	}
	if err != nil {
		return nil, err
	}

	return request, err
}

func deserialize[T any](reader io.ReadCloser, in any) (err error) {
	defer reader.Close()
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	// unmarshal to struct
	err = json.Unmarshal(bytes, in)
	if err != nil {
		return err
	}
	return nil
}

// The api returns something like:
// {"range":"Sheet2!A1:Z1000","majorDimension":"ROWS","values":[["a","b"],["1","2"]]}
// only the 'values' field is relevant.
// The returned reader will contain it in an 'encoding/csv' readable format
func truncateExtraneousData(reader io.ReadCloser) (io.Reader, error) {
	// not the fastest way to do things but easy to read and maintain
	result := values{}
	err := deserialize[values](reader, &result)
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
	return bytes.NewReader(output.Bytes()), nil
}

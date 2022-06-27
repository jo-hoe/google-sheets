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
const updateSheetUrl = baseUrl + ":batchUpdate"
const updateValuesUrl = baseUrl + "/values:batchUpdate"
const copySheetUrl = baseUrl + "/sheets/%d:copyTo"
const clearSheetUrl = baseUrl + "/values/%s:clear"

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
	UpdateSheetProperties *updateSheetProperties `json:"updateSheetProperties,omitempty"`
	DeleteSheet           *deleteSheet           `json:"deleteSheet,omitempty"`
	AutoResizeDimensions  *autoResizeDimensions  `json:"autoResizeDimensions,omitempty"`
	AddSheet              *addSheet              `json:"addSheet,omitempty"`
}

type copyRequest struct {
	DestinationSpreadsheetId string `json:"destinationSpreadsheetId,omitempty"`
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

type updateSheetProperties struct {
	Properties     spreadSheetProperties `json:"properties"`
	FieldsToUpdate string                `json:"fields"`
}

type deleteSheet struct {
	SheetId int32 `json:"sheetId"`
}

type autoResizeDimensions struct {
	Dimensions dimensions `json:"dimensions"`
}

type dimensions struct {
	SheetId   int32  `json:"sheetId"`
	Dimension string `json:"dimension"`
}

type vauleInput struct {
	ValueInputOption          string     `json:"valueInputOption"`
	ValueRange                valueRange `json:"data"`
	IncludeValuesInResponse   bool       `json:"includeValuesInResponse"`
	ResponseValueRenderOption string     `json:"responseValueRenderOption"`
	DateTimeRenderOption      string     `json:"responseDateTimeRenderOption"`
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

func (wrapper SheetsApiWrapper) UpdateSheetMetaData(spreadSheetId string, sheetId int32, newSheetName string) (err error) {
	body := updateRequest{}
	body.Request = []batchRequest{{
		UpdateSheetProperties: &updateSheetProperties{
			FieldsToUpdate: "title",
			Properties: spreadSheetProperties{
				Title:   newSheetName,
				SheetID: sheetId,
			}}}}

	response, err := wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheetId), body)
	if response != nil {
		response.Close()
	}
	if err != nil {
		return err
	}

	return nil
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

func (wrapper SheetsApiWrapper) AutoResizeSheet(spreadSheetId string, sheetId int32) (err error) {
	body := updateRequest{}
	body.Request = []batchRequest{{
		AutoResizeDimensions: &autoResizeDimensions{
			Dimensions: dimensions{
				SheetId:   sheetId,
				Dimension: "ROWS",
			},
		}}}

	response, err := wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheetId), body)
	if response != nil {
		response.Close()
	}
	if err != nil {
		return err
	}
	body.Request[0].AutoResizeDimensions.Dimensions.Dimension = "COLUMNS"
	response, err = wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheetId), body)
	if response != nil {
		response.Close()
	}
	if err != nil {
		return err
	}

	return nil
}

func (wrapper SheetsApiWrapper) postSheetRequest(url string, body any) (out io.ReadCloser, err error) {
	var request *http.Request
	if body != nil {
		request, err = wrapper.createJSONPostRequest(url, body)
	} else {
		request, err = http.NewRequest("POST", url, nil)
	}
	if err != nil {
		return nil, err
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

func (wrapper SheetsApiWrapper) CreateSheet(spreadSheetId string, sheetId int32, sheetName string) (id int32, err error) {
	body := updateRequest{}
	body.IncludeSpreadsheetInResponse = true
	body.Request = []batchRequest{{
		AddSheet: &addSheet{
			Properties: spreadSheetProperties{
				Title:   sheetName,
				SheetID: sheetId,
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

func (wrapper SheetsApiWrapper) GetSheetId(spreadSheetId string, sheetName string) (int32, error) {
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

func (wrapper SheetsApiWrapper) findSheetIdInResponse(allSheets []sheet, sheetName string) (id int32, err error) {
	for _, sheet := range allSheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetID, nil
		}
	}
	return -1, fmt.Errorf("sheet was not found")
}

func (wrapper SheetsApiWrapper) WriteSheet(spreadSheetId string, sheetName string, data [][]string) (err error) {
	body := vauleInput{}
	body.ValueRange.Range = sheetName
	body.ValueInputOption = "USER_ENTERED"
	body.ValueRange.MajorDimension = "COLUMNS"
	body.DateTimeRenderOption = "FORMATTED_STRING"
	body.IncludeValuesInResponse = false
	body.ResponseValueRenderOption = "UNFORMATTED_VALUE"
	body.ValueRange.Values = data

	response, err := wrapper.postSheetRequest(fmt.Sprintf(updateValuesUrl, spreadSheetId), body)
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

// copies a sheet
func (wrapper SheetsApiWrapper) CopySheet(spreadSheetId string, sheetId int32, destinationSpreadsheetId string) (id int32, err error) {
	body := copyRequest{
		DestinationSpreadsheetId: destinationSpreadsheetId,
	}
	response, err := wrapper.postSheetRequest(fmt.Sprintf(copySheetUrl, spreadSheetId, sheetId), body)
	if err != nil {
		return -1, err
	}

	result := spreadSheetProperties{}
	err = deserialize[spreadSheet](response, &result)
	if err != nil {
		return -1, err
	}

	return result.SheetID, nil
}

// Deletes data in the sheet and replaces it with new data.
// Takes a backup of the current data during the processing of the request. This 
// backup it deleted in the last set of this function.
func (wrapper SheetsApiWrapper) ReplaceSheetData(spreadSheetId string, initialSheetName string, data [][]string) (err error) {
	initialSheetId, err := wrapper.GetSheetId(spreadSheetId, initialSheetName)
	if err != nil {
		return err
	}

	// create a backup
	backupId, err := wrapper.CopySheet(spreadSheetId, initialSheetId, spreadSheetId)
	if err != nil {
		return err
	}

	err = wrapper.ClearSheet(spreadSheetId, initialSheetName)
	if err != nil {
		return err
	}
	err = wrapper.WriteSheet(spreadSheetId, initialSheetName, data)
	if err != nil {
		return err
	}
	err = wrapper.AutoResizeSheet(spreadSheetId, initialSheetId)
	if err != nil {
		return err
	}
	err = wrapper.DeleteSheet(spreadSheetId, backupId)
	if err != nil {
		return err
	}
	return err
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
func truncateExtraneousData(reader io.ReadCloser) (io.ReadCloser, error) {
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
	return ioutil.NopCloser(bytes.NewReader(output.Bytes())), nil
}

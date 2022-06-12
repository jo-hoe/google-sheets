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
	"time"
)

const baseUrl = "https://sheets.googleapis.com/v4/spreadsheets/%s"

// url is reverse engineered from:
// https://github.com/googleapis/google-api-go-client/blob/bc181c33247b7fe3d06d2d7139da0fa06fabbd71/sheets/v4/sheets-gen.go#L14283
// but it is also described here:
// https://developers.google.com/sheets/api/reference/rest/v4/spreadsheets.values/get
const csvUrlTemplate = baseUrl + "/values/%s?alt=json&prettyPrint=false"
const updateSheetUrl = baseUrl + ":batchUpdate"
const updateValuesUrl = baseUrl + "/values:batchUpdate"

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
	SheetID int    `json:"sheetId"`
	Title   string `json:"title"`
}

type batchRequest struct {
	UpdateSheetProperties updateSheetProperties `json:"updateSheetProperties"`
	DeleteSheet           deleteSheet           `json:"deleteSheet"`
	AutoResizeDimensions  autoResizeDimensions  `json:"autoResizeDimensions"`
}

type updateSheetProperties struct {
	Properties     spreadSheetProperties `json:"properties"`
	FieldsToUpdate string                `json:"fields"`
}

type deleteSheet struct {
	SheetId int `json:"sheetId"`
}

type autoResizeDimensions struct {
	Dimensions dimensions `json:"dimensions"`
}

type dimensions struct {
	SheetId   int    `json:"sheetId"`
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
	Range          string `json:"range"`
	MajorDimension string `json:"majorDimension"`
	Values         values `json:"values"`
}

type SheetsApiWrapper struct {
	httpClient *http.Client
}

func NewSheetsApiWrapper(httpClient *http.Client) *SheetsApiWrapper {
	return &SheetsApiWrapper{
		httpClient: httpClient,
	}
}

func (wrapper SheetsApiWrapper) CreateSheet(SpreadSheetId string, sheetName string) (out *sheet, err error) {
	body := spreadSheet{}
	body.Sheets = append(body.Sheets, sheet{
		spreadSheetProperties{
			Title:   sheetName,
			SheetID: int(time.Now().UnixMilli()),
		},
	})

	response, err := wrapper.postSheetRequest(fmt.Sprintf(baseUrl, ""), body)
	if err != nil {
		return nil, err
	}
	result := sheet{}
	err = deserialize[spreadSheet](response, &result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (wrapper SheetsApiWrapper) UpdateSheetMetaData(spreadSheatId string, oldSheetId int, newSheetId int, newSheetName string) (err error) {
	body := batchRequest{}
	body.UpdateSheetProperties = updateSheetProperties{
		FieldsToUpdate: "sheetId,title",
		Properties: spreadSheetProperties{
			Title:   newSheetName,
			SheetID: newSheetId,
		},
	}

	_, err = wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheatId), body)

	if err != nil {
		return err
	}

	return nil
}

func (wrapper SheetsApiWrapper) DeleteSheet(spreadSheatId string, sheetId int) (err error) {
	body := batchRequest{}
	body.DeleteSheet = deleteSheet{
		SheetId: sheetId,
	}

	_, err = wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheatId), body)

	if err != nil {
		return err
	}

	return nil
}

func (wrapper SheetsApiWrapper) AutoResizeSheet(spreadSheatId string, sheetId int) (err error) {
	body := batchRequest{}
	body.AutoResizeDimensions = autoResizeDimensions{
		Dimensions: dimensions{
			SheetId:   sheetId,
			Dimension: "ROWS",
		},
	}

	_, err = wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheatId), body)
	if err != nil {
		return err
	}
	body.AutoResizeDimensions.Dimensions.Dimension = "COLUMNS"
	_, err = wrapper.postSheetRequest(fmt.Sprintf(updateSheetUrl, spreadSheatId), body)
	if err != nil {
		return err
	}

	return nil
}

func (wrapper SheetsApiWrapper) postSheetRequest(url string, body any) (out io.ReadCloser, err error) {
	payloadBuf := new(bytes.Buffer)
	err = json.NewEncoder(payloadBuf).Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		return nil, err
	}

	response, err := wrapper.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Response was '%d': %s", response.StatusCode, response.Status)
	}

	return response.Body, nil
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

func (wrapper SheetsApiWrapper) GetSheetId(spreadSheetId string, sheetName string) (int, error) {
	url := fmt.Sprintf(baseUrl, spreadSheetId)
	resp, err := wrapper.httpClient.Get(url)
	if err != nil {
		return -1, err
	}

	result := spreadSheet{}
	err = deserialize[spreadSheet](resp.Body, &result)
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

func (wrapper SheetsApiWrapper) WriteSheet(spreadSheatId string, sheetName string, data [][]string) (err error) {
	body := vauleInput{}
	body.ValueRange.Range = sheetName
	body.ValueInputOption = "USER_ENTERED"
	body.ValueRange.MajorDimension = "COLUMS"
	body.DateTimeRenderOption = "FORMATTED_STRING"
	body.IncludeValuesInResponse = false
	body.ValueRange.Values = values{
		Values: data,
	}

	_, err = wrapper.postSheetRequest(fmt.Sprintf(updateValuesUrl, spreadSheatId), body)
	if err != nil {
		return err
	}

	return nil
}

func (wrapper SheetsApiWrapper) ReplaceSheet(spreadSheatId string, initialSheetName string, data [][]string) (err error) {
	initialSheetId, err := wrapper.GetSheetId(spreadSheatId, initialSheetName)
	if err != nil {
		return err
	}
	newSheet, err := wrapper.CreateSheet(spreadSheatId, fmt.Sprintf("%s-%d", initialSheetName, time.Now().UnixMilli()))
	if err != nil {
		return err
	}
	err = wrapper.WriteSheet(spreadSheatId, initialSheetName, data)
	if err != nil {
		return err
	}
	err = wrapper.AutoResizeSheet(spreadSheatId, newSheet.Properties.SheetID)
	if err != nil {
		return err
	}
	err = wrapper.DeleteSheet(spreadSheatId, initialSheetId)
	if err != nil {
		return err
	}
	err = wrapper.UpdateSheetMetaData(spreadSheatId, newSheet.Properties.SheetID, initialSheetId, initialSheetName)
	if err != nil {
		return err
	}
	return err
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

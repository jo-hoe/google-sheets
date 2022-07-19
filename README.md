# Google Sheets

[![Test Status](https://github.com/jo-hoe/google-sheets/workflows/test/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=test)
[![Lint Status](https://github.com/jo-hoe/google-sheets/workflows/lint/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=lint)
[![CodeQL Status](https://github.com/jo-hoe/google-sheets/workflows/CodeQL/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=CodeQL)

Provides an idiomatic way to read and write data from google sheets.

## Example Useage

```golang
// Creating a http client with credentials of a gcp service account.
jsonServiceAccount, err := ioutil.ReadFile("path\\to\\serivce_account_file")
if err != nil {
  log.Print(err.Error())
  return
}

// spread sheet id can be taken from the URL
// example URL: https://docs.google.com/spreadsheets/d/c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T/edit#gid=0
// the spreadsheet ID would be "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T"
sheet, err := sheet.OpenSheet(context.Background(), "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T", "Sheet1", sheet.O_CREATE|sheet.O_RDWR, jsonServiceAccount)
if err != nil {
  log.Print(err.Error())
  return
}
csvWriter := csv.NewWriter(sheet)
err = csvWriter.WriteAll([][]string{
  {"0", "1"},
  {"2", "3"},
})

csvReader := csv.NewReader(sheet)
csvResult, err := csvReader.ReadAll()
if err != nil {
  log.Print(err.Error())
  return
}
fmt.Printf("results: %v", csvResult)
```

## Google Sheets Authorization

The offical documentation can be found here: <https://developers.google.com/sheets/api/guides/authorizing>.
Note, that there is no possiblity to reduce the API access to only a specific sheet.
To mitigate that, consider to use a dedicated service account for your google sheets.

After creating the json key for your service account do not forget to enable the google project in which the key residesit for the sheet api. You may do so using this url scheme;
<https://console.cloud.google.com/apis/library/sheets.googleapis.com?project=>[project id]
  
## Linting

Project used `golangci-lint` for linting.

### Installation

<https://golangci-lint.run/usage/install/>

### Execution

Run the linting locally by executing

```cli
golangci-lint run ./...
```

in the working directory

## Testing

The project contains both unit and integrations tests.

### Unit Test Execution

The unit test can be excuted using the default golang commands. To run all test execute the following in the parent folder of the repository.

```powershell
go test ./...
```

### Integration Test Execution

A credentials file and a google spreadsheet needed as prerequisite for the integration tests. You may use the following launch.json file in VSCode to run the tests.

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "API Wrapper Integration Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/api/apiwrapper/apiwrapper_integration_test.go",
            "env": {
                "CREDENTIALS_FILE_PATH": "C:\\Folder\\file-name-352919-3f8fa23b9bba.json",
                "SPREADSHEET_ID": "1yxmv2lTtOtvpkBi-5hSMq86CHFMfYq6kdjfasudfasih"
            },
        },{
            "name": "Sheets Integration Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/sheet/sheets_integraton_test.go",
            "env": {
                "CREDENTIALS_FILE_PATH": "C:\\Folder\\file-name-352919-3f8fa23b9bba.json",
                "SPREADSHEET_ID": "1yxmv2lTtOtvpkBi-5hSMq86CHFMfYq6kdjfasudfasih"
            },
        },
    ]
}
```

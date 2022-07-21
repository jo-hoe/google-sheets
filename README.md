# Google Sheets

[![Test Status](https://github.com/jo-hoe/google-sheets/workflows/test/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=test)
[![Lint Status](https://github.com/jo-hoe/google-sheets/workflows/lint/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=lint)
[![CodeQL Status](https://github.com/jo-hoe/google-sheets/workflows/CodeQL/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=CodeQL)

Provides an idiomatic way to read and write data from google sheets.

## Example Useage

```golang
// open the key file for your GCP service account (see below how to create that key file)
jsonServiceAccount, err := ioutil.ReadFile("path\\to\\service_account_key.json")
if err != nil {
  log.Print(err.Error())
  return
}

// spreadsheet id can be taken from the URL
// example URL: https://docs.google.com/spreadsheets/d/c8ACvfAd4X09Hi9mCl4qcBidP635S8z5luk-vGG54N5T/edit#gid=0
// the spreadsheet ID would be "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T"
sheet, err := sheet.OpenSheet(context.Background(), "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5luk-vGG54N5T", "Sheet1", sheet.O_CREATE|sheet.O_RDWR, jsonServiceAccount)
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

sheet.Remove(context.Background(), sheet.SpreadSheetId(), sheet.Id(), jsonServiceAccount)
```

## Google Sheets AuthN/AuthZ

### General

The offical documentation can be found here: <https://developers.google.com/sheets/api/guides/authorizing>.

### Creating the key file

1. [create a gcp service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts#creating)
2. after creating the service account, ensure that google project in which the service account resides, is enabled to use the sheet api. You verifiy or enable the API using this url scheme <https://console.cloud.google.com/apis/library/sheets.googleapis.com?project=>[my gcp project id]
3. after the service account is created, take the mail address of that account and [share your spreadsheet with that mail address](https://support.google.com/a/users/answer/9305987?hl=en#)
4. [create a json key for your GCP service account](https://cloud.google.com/iam/docs/creating-managing-service-account-keys#creating)
  
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

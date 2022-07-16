# Google Sheets

[![Test Status](https://github.com/jo-hoe/google-sheets/workflows/test/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=test)
[![Lint Status](https://github.com/jo-hoe/google-sheets/workflows/lint/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=lint)

Provides an idiomatic way to read and write data from google sheets.

## Example Useage

```golang
// Creating a http client with credentials of service account.
clientCredentialsJson := os.GetEnv("myClientCredentialJsonString")
myClient := client.NewServiceAccountClient(ctx, clientCredentialsJson)

// spread sheet id can be taken from the URL
// example URL: https://docs.google.com/spreadsheets/d/c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T/edit#gid=0
// the spreadsheet ID would be "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T"
writer, err := NewSheetWriter(myClient, "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T", "Sheet1", O_CREATE)
csvWriter := csv.NewWriter(writer)
err = csvWriter.WriteAll([][]string{
  {"0", "1"},
  {"2", "3"},
})

reader, err := reader.NewSheetReader(myClient, "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T", "Sheet1")
if err != nil {
  return nil, err
}
csvReader := csv.NewReader(reader)
csvResult, err := csvReader.ReadAll()
if err != nil {
  return nil, err
}
fmt.Printf("results: %v", csvResult)
```

## Google Sheets Authorization

The offical documentation can be found here: <https://developers.google.com/sheets/api/guides/authorizing>.
Note, that there is no possiblity to reduce the API access to only a specific file.
To mitigate that, consider to use a dedicated service account.

After creating the key do not forget to enable it for the sheet api
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

### Integration Test Exception

A credentials file and a google spreadsheet needed as prerequisite for the integration tests. You may use the following launch.json file to run the tests.

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug API Wrapper Tests",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/apiwrapper/apiwrapper_integration_test.go",
            "env": {
                "CREDENTIALS_FILE_PATH": "C:\\Folder\\file-name-352919-3f8fa23b9bba.json",
                "SPREADSHEET_ID": "1yxmv2lTtOtvpkBi-5hSMq86CHFMfYq6kdjfasudfasih"
            },
        },{
            "name": "Debug Writer Tests",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/writer/sheetwriter_integration_test.go",
            "env": {
                "CREDENTIALS_FILE_PATH": "C:\\Folder\\file-name-352919-3f8fa23b9bba.json",
                "SPREADSHEET_ID": "1yxmv2lTtOtvpkBi-5hSMq86CHFMfYq6kdjfasudfasih"
            },
        },
    ]
}
```

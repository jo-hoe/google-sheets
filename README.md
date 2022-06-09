# Google Sheets

[![Test Status](https://github.com/jo-hoe/google-sheets/workflows/test/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=test)
[![Lint Status](https://github.com/jo-hoe/google-sheets/workflows/lint/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=lint)

Provides an idiomatic way to read data from google sheets.

## Example Useage

```golang
// Creating a http client with credentials of service account.
clientCredentialsJson := os.GetEnv("myClientCredentialJsonString")
// If sheet is public, a regular client can also be used e.g.:
// myClient := http.Client{}
myClient := client.NewServiceAccountClient(ctx, clientCredentialsJson)

// spread sheet id can be taken from the URL
// example URL: https://docs.google.com/spreadsheets/d/c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T/edit#gid=0
// the spreadsheet ID would be "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T"
readerCloser, err := reader.NewSheetReader(myClient, "c8ACvfAd4X09Hi9mCl4qcBidP635S8z5lukxvGG54N5T", "Sheet1")
defer closeReader(readerCloser)
if err != nil {
  return nil, err
}

csv := csv.NewReader(readerCloser)
csvResult, err := csv.ReadAll()
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

## Todo

- read sheet ✔
- create sheet ? (not integration tested)
- rename sheet ? (not integration tested)
- delete sheet ? (not integration tested)
- write in a sheet
- fit column length ? (not integration tested)
  
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

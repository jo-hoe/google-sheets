# Google Sheets

[![Test Status](https://github.com/jo-hoe/google-sheets/workflows/test/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=test)
[![Lint Status](https://github.com/jo-hoe/google-sheets/workflows/lint/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=lint)

Provides an idiomatic way to read data from google sheets.

## Example Useage

```golang
clientCredentialsJson := os.GetEnv("myClientCredentialJsonString")
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

## Linting

Project used `golangci-lint` for linting. You can download it by executing

```cli
go get github.com/golangci/golangci-lint/cmd/golangci-lint
```

and run the linting locally by executing

```cli
golangci-lint run ./...
```

in the working directory

# Google Sheets

[![Test Status](https://github.com/jo-hoe/google-sheets/workflows/test/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=test)
[![Lint Status](https://github.com/jo-hoe/google-sheets/workflows/lint/badge.svg)](https://github.com/jo-hoe/google-sheets/actions?workflow=lint)

Idiomatic way to read data from google sheets.

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

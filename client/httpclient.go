package client

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google client creation reference to
// https://developers.google.com/sheets/api/quickstart/go
// for more details

const Scopes = "https://www.googleapis.com/auth/spreadsheets.readonly"

// NewSpreadsheetClient creates a http client to access non-public spreedsheets.
func NewSpreadsheetClient(clientCredentialsJson string, token *oauth2.Token, saveToken func(token *oauth2.Token)) (*http.Client, error) {
	oauthConfig, err := google.ConfigFromJSON([]byte(clientCredentialsJson), Scopes)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	tokenSource := oauthConfig.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if saveToken != nil {
		saveToken(newToken)
	}
	return oauth2.NewClient(ctx, tokenSource), nil
}

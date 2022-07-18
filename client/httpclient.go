package client

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"
)

// Google client creation reference to
// https://developers.google.com/sheets/api/quickstart/go
// for more details

const ReadOnlyScopes = "https://www.googleapis.com/auth/spreadsheets.readonly"
const ReadWriteScopes = "https://www.googleapis.com/auth/spreadsheets"

// NewReadClient creates a http client to access non-public spreedsheets.
// Account will only have read access
func NewReadClient(ctx context.Context, clientCredentialsJson string) (*http.Client, error) {
	return newServiceAccountClient(ctx, clientCredentialsJson, ReadOnlyScopes)
}

func NewReadWriteClient(ctx context.Context, clientCredentialsJson string) (*http.Client, error) {
	return newServiceAccountClient(ctx, clientCredentialsJson, ReadWriteScopes)
}

func newServiceAccountClient(ctx context.Context, clientCredentialsJson string, scopes string) (*http.Client, error) {
	clientCredentials := []byte(clientCredentialsJson)
	config, err := google.JWTConfigFromJSON(clientCredentials, scopes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	return config.Client(ctx), nil
}

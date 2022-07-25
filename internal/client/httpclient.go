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
	return NewServiceAccountClient(ctx, []byte(clientCredentialsJson), ReadOnlyScopes)
}

func NewReadWriteClient(ctx context.Context, clientCredentialsJson string) (*http.Client, error) {
	return NewServiceAccountClient(ctx, []byte(clientCredentialsJson), ReadWriteScopes)
}

func NewServiceAccountClient(ctx context.Context, clientCredentials []byte, scopes string) (*http.Client, error) {
	config, err := google.JWTConfigFromJSON(clientCredentials, scopes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	return config.Client(ctx), nil
}

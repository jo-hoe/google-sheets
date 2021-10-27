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

const Scopes = "https://www.googleapis.com/auth/spreadsheets.readonly"

// NewServiceAccountClient creates a http client to access non-public spreedsheets.
func NewServiceAccountClient(ctx context.Context, clientCredentialsJson string) (*http.Client, error) {
	clientCredentials := []byte(clientCredentialsJson)
	config, err := google.JWTConfigFromJSON(clientCredentials, Scopes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	return config.Client(ctx), nil
}

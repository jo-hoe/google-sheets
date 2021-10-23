package example

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jo-hoe/google-sheets/client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const configEnvKey = "GoogleOAuthClient"

// creates a token based on an OAuth configuration in the enviroment variables
func main() {
	oauthConfig, exist := os.LookupEnv(configEnvKey)
	if !exist {
		log.Fatalf("Did not find env %s", configEnvKey)
	}

	config, err := google.ConfigFromJSON([]byte(oauthConfig), client.Scopes)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	jsonToken, err := json.Marshal(token)
	if err != nil {
		log.Fatalf("Unable to parse token: %v", err)
	}
	fmt.Print(jsonToken)
}

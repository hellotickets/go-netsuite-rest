package netsuite_test

import (
	"context"
	"log"
	"os"
	"testing"

	netsuite "github.com/hellotickets/go-netsuite-rest"
	"golang.org/x/oauth2"
)

var (
	client *netsuite.Client
)

func TestMain(m *testing.M) {
	authType := os.Getenv("AUTH_TYPE")
	baseURL := os.Getenv("BASE_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	tokenID := os.Getenv("TOKEN_ID")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	tokenURL := os.Getenv("TOKEN_URL")
	companyID := os.Getenv("COMPANY_ID")
	contentLanguage := os.Getenv("CONTENT_LANGUAGE")
	debug := os.Getenv("DEBUG")

	client = netsuite.NewClient(nil)
	if authType == "oauth" {
		oauthConfig := netsuite.NewOauth2Config(companyID)
		oauthConfig.ClientID = clientID
		oauthConfig.ClientSecret = clientSecret

		// set alternative token url
		if tokenURL != "" {
			oauthConfig.Endpoint.TokenURL = tokenURL
		}

		// b, _ := json.MarshalIndent(oauthConfig, "", "  ")
		// log.Fatal(string(b))

		token := &oauth2.Token{
			RefreshToken: refreshToken,
		}

		// get http client with automatic oauth logic
		httpClient := oauthConfig.Client(context.Background(), token)

		client = netsuite.NewClient(httpClient)
		client.SetCompanyID(companyID)
	} else if authType == "token" {
		client.SetUseTokenAuth(true)
		client.SetClientID(clientID)
		client.SetClientSecret(clientSecret)
		client.SetTokenID(tokenID)
		client.SetTokenSecret(tokenSecret)
		client.SetCompanyID(companyID)
	} else {
		log.Fatalf("Unknown auth type: %s", authType)
	}
	if debug != "" {
		client.SetDebug(true)
	}

	if baseURL != "" {
		client.SetBaseURL(baseURL)
	}

	client.SetContentLanguage(contentLanguage)
	client.SetDisallowUnknownFields(true)
	m.Run()
}

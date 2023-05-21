package v1

import (
	"os"
	"strings"
)

// NewAuthenticatedBasicClient creates a new BasicClient with credentials from environment variables
func NewClient() *Client {
	username := os.Getenv("BITBUCKET_USERNAME")
	password := os.Getenv("BITBUCKET_PASSWORD")
	oauthClientId := os.Getenv("BITBUCKET_OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("BITBUCKET_OAUTH_CLIENT_SECRET")
	authMethod := os.Getenv("BITBUCKET_AUTH_METHOD")

	// no detailed check necessary, it was already performed by provider_test.go
	if strings.EqualFold(authMethod, "oauth") {
		return NewOAuthClient(oauthClientId, oauthClientSecret)
	} else {
		return NewBasicAuthClient(username, password)
	}
}

package bitbucket

import (
	"os"
	"strings"

	gobb "github.com/ktrysmt/go-bitbucket"
)

// NewAuthenticatedBasicClient creates a new BasicClient with credentials from environment variables
func NewClient() *gobb.Client {
	username := os.Getenv("BITBUCKET_USERNAME")
	password := os.Getenv("BITBUCKET_PASSWORD")
	oauthClientId := os.Getenv("BITBUCKET_OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("BITBUCKET_OAUTH_CLIENT_SECRET")
	authMethod := os.Getenv("BITBUCKET_AUTH_METHOD")

	// no detailed check necessary, it was already performed by provider_test.go
	if strings.EqualFold(authMethod, "oauth") {
		return gobb.NewOAuthClientCredentials(oauthClientId, oauthClientSecret)
	} else {
		return gobb.NewBasicAuth(username, password)
	}
}

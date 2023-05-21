package v1

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicAuthClient(t *testing.T) {
	client := NewBasicAuthClient("test", "password")

	assert.Equal(t, "https://api.bitbucket.org/1.0", client.ApiBaseUrl.String())
	assert.IsType(t, &BasicAuth{}, client.Auth)
	assert.Equal(t, client.Auth.(*BasicAuth).Username, "test")
	assert.Equal(t, client.Auth.(*BasicAuth).Password, "password")
	assert.IsType(t, &Groups{}, client.Groups)
	assert.IsType(t, &http.Client{}, client.HttpClient)
}

func TestNewOAuthClient(t *testing.T) {
	client := NewOAuthClient(os.Getenv("BITBUCKET_OAUTH_CLIENT_ID"), os.Getenv("BITBUCKET_OAUTH_CLIENT_SECRET"))

	assert.Equal(t, "https://api.bitbucket.org/1.0", client.ApiBaseUrl.String())
	assert.IsType(t, &BearerAuth{}, client.Auth)
	assert.NotEmpty(t, client.Auth.(*BearerAuth).Token)
	assert.IsType(t, &Groups{}, client.Groups)
	assert.IsType(t, &http.Client{}, client.HttpClient)
}

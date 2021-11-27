package v1

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	auth := &Auth{
		Username: "test",
		Password: "test",
	}
	client := NewClient(auth)

	assert.Equal(t, "https://api.bitbucket.org/1.0", client.ApiBaseUrl.String())
	assert.Equal(t, auth, client.Auth)
	assert.IsType(t, &Groups{}, client.Groups)
	assert.IsType(t, &GroupPrivileges{}, client.GroupPrivileges)
	assert.IsType(t, &http.Client{}, client.HttpClient)
}

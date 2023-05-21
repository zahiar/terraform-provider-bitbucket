package v1

import (
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/bitbucket"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	Auth Auth

	ApiBaseUrl *url.URL
	HttpClient *http.Client

	Groups       *Groups
	GroupMembers *GroupMembers
}

type Auth interface {
	SetRequestAuth(request *http.Request)
}

type BasicAuth struct {
	Username string
	Password string
}

type BearerAuth struct {
	Token string
}

func (auth *BasicAuth) SetRequestAuth(request *http.Request) {
	request.SetBasicAuth(auth.Username, auth.Password)
}

func (auth *BearerAuth) SetRequestAuth(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+auth.Token)
}

func NewOAuthClient(clientId string, clientSecret string) *Client {
	ctx := context.Background()
	conf := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     bitbucket.Endpoint.TokenURL,
	}

	tok, err := conf.Token(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return newClient(&BearerAuth{Token: tok.AccessToken})
}

func NewBasicAuthClient(username string, password string) *Client {
	return newClient(&BasicAuth{Username: username, Password: password})
}

func newClient(auth Auth) *Client {
	apiBaseUrl, err := url.Parse("https://api.bitbucket.org/1.0")
	if err != nil {
		log.Fatal(err)
	}

	client := &Client{
		Auth:       auth,
		ApiBaseUrl: apiBaseUrl,
	}
	client.Groups = &Groups{client: client}
	client.GroupMembers = &GroupMembers{client: client}
	client.HttpClient = new(http.Client)

	return client
}

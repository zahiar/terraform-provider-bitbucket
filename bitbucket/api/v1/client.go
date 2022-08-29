package v1

import (
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	Auth *Auth

	ApiBaseUrl *url.URL
	HttpClient *http.Client

	Groups       *Groups
	GroupMembers *GroupMembers
}

type Auth struct {
	Username string
	Password string
}

func NewClient(auth *Auth) *Client {
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

package v1

// Implements: https://support.atlassian.com/bitbucket-cloud/docs/groups-endpoint/#Overview

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Groups struct {
	client *Client
}

type Group struct {
	Owner struct {
		Uuid string `json:"uuid"`
	}
	Name       string `json:"name"`
	AutoAdd    bool   `json:"auto_add"`
	Slug       string `json:"slug"`
	Permission string `json:"permission"`
}

type GroupOptions struct {
	OwnerUuid  string
	Name       string
	AutoAdd    bool
	Slug       string
	Permission string
}

func (g *Groups) Get(gro *GroupOptions) (*Group, error) {
	url := fmt.Sprintf("%s/groups?group=%s/%s", g.client.ApiBaseUrl, gro.OwnerUuid, gro.Slug)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(g.client.Auth.Username, g.client.Auth.Password)

	response, err := g.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := make([]Group, 1)
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no group found")
	}

	return &result[0], nil
}

func (g *Groups) Create(gro *GroupOptions) (*Group, error) {
	url := fmt.Sprintf("%s/groups/%s", g.client.ApiBaseUrl, gro.OwnerUuid)
	body := strings.NewReader(fmt.Sprintf("name=%s", gro.Name))
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(g.client.Auth.Username, g.client.Auth.Password)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := g.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := &Group{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}

	return result, nil
}

func (g *Groups) Update(gro *GroupOptions) (*Group, error) {
	url := fmt.Sprintf("%s/groups/%s/%s", g.client.ApiBaseUrl, gro.OwnerUuid, gro.Slug)

	requestBody := struct {
		Name       string `json:"name,omitempty"`
		AutoAdd    bool   `json:"auto_add,omitempty"`
		Permission string `json:"permission,omitempty"`
	}{
		Name:       gro.Name,
		AutoAdd:    gro.AutoAdd,
		Permission: gro.Permission,
	}
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	requestBodyJsonString := strings.NewReader(string(requestBodyJson))
	request, err := http.NewRequest("PUT", url, requestBodyJsonString)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(g.client.Auth.Username, g.client.Auth.Password)
	request.Header.Set("Content-Type", "application/json")

	response, err := g.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := &Group{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}

	return result, nil
}

func (g *Groups) Delete(gro *GroupOptions) error {
	url := fmt.Sprintf("%s/groups/%s/%s", g.client.ApiBaseUrl, gro.OwnerUuid, gro.Slug)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.SetBasicAuth(g.client.Auth.Username, g.client.Auth.Password)

	response, err := g.client.HttpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code, reponse was %s", response.Body)
	}

	return nil
}

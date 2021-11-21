package v1

// Implements: https://support.atlassian.com/bitbucket-cloud/docs/group-privileges-endpoint/#Overview

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type GroupPrivileges struct {
	client *Client
}

type GroupPrivilege struct {
	Privilege string `json:"privilege"`
	Group     struct {
		Owner struct {
			Uuid string `json:"uuid"`
		}
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	Repository struct {
		Owner struct {
			Uuid string `json:"uuid"`
		}
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
}

type GroupPrivilegeOptions struct {
	WorkspaceId string
	RepoSlug    string
	GroupOwner  string
	GroupSlug   string
	Privilege   string
}

func (gp *GroupPrivileges) Get(gpo *GroupPrivilegeOptions) (*GroupPrivilege, error) {
	url := fmt.Sprintf("%s/group-privileges/%s/%s/%s/%s", gp.client.ApiBaseUrl, gpo.WorkspaceId, gpo.RepoSlug, gpo.GroupOwner, gpo.GroupSlug)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(gp.client.Auth.Username, gp.client.Auth.Password)

	response, err := gp.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := make([]GroupPrivilege, 1)
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}

	return &result[0], nil
}

func (gp *GroupPrivileges) Create(gpo *GroupPrivilegeOptions) (*GroupPrivilege, error) {
	url := fmt.Sprintf("%s/group-privileges/%s/%s/%s/%s", gp.client.ApiBaseUrl, gpo.WorkspaceId, gpo.RepoSlug, gpo.GroupOwner, gpo.GroupSlug)
	body := strings.NewReader(gpo.Privilege)
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(gp.client.Auth.Username, gp.client.Auth.Password)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := gp.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := make([]GroupPrivilege, 1)
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}

	return &result[0], nil
}

func (gp *GroupPrivileges) Delete(gpo *GroupPrivilegeOptions) error {
	url := fmt.Sprintf("%s/group-privileges/%s/%s/%s/%s", gp.client.ApiBaseUrl, gpo.WorkspaceId, gpo.RepoSlug, gpo.GroupOwner, gpo.GroupSlug)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	request.SetBasicAuth(gp.client.Auth.Username, gp.client.Auth.Password)

	response, err := gp.client.HttpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("invalid status code, reponse was %s", response.Body)
	}

	return nil
}

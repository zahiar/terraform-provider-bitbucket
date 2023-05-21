package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Implements: https://support.atlassian.com/bitbucket-cloud/docs/groups-endpoint/#GET-the-group-members

type GroupMembers struct {
	client *Client
}

type GroupMember struct {
	DisplayName string `json:"display_name"`
	AccountID   string `json:"account_id"`
	IsActive    bool   `json:"is_active"`
	IsTeam      bool   `json:"is_team"`
	IsStaff     bool   `json:"is_staff"`
	Avatar      string `json:"avatar"`
	ResourceURI string `json:"resource_uri"`
	Nickname    string `json:"nickname"`
	UUID        string `json:"uuid"`
}

type GroupMemberOptions struct {
	OwnerUuid string
	Slug      string
	UserUuid  string
}

func (gm *GroupMembers) Get(gmo *GroupMemberOptions) ([]GroupMember, error) {
	url := fmt.Sprintf("%s/groups/%s/%s/members", gm.client.ApiBaseUrl, gmo.OwnerUuid, gmo.Slug)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	gm.client.Auth.SetRequestAuth(request)

	response, err := gm.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := make([]GroupMember, 1)
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}

	return result, nil
}

func (gm *GroupMembers) Create(gmo *GroupMemberOptions) (*GroupMember, error) {
	url := fmt.Sprintf("%s/groups/%s/%s/members/%s", gm.client.ApiBaseUrl, gmo.OwnerUuid, gmo.Slug, gmo.UserUuid)
	request, err := http.NewRequest("PUT", url, strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	gm.client.Auth.SetRequestAuth(request)
	request.Header.Set("Content-Type", "application/json")

	response, err := gm.client.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code was not 200")
	}

	result := new(GroupMember)
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Println("Could not unmarshal JSON payload")
		return nil, err
	}

	return result, nil
}

func (gm *GroupMembers) Delete(gmo *GroupMemberOptions) error {
	url := fmt.Sprintf("%s/groups/%s/%s/members/%s", gm.client.ApiBaseUrl, gmo.OwnerUuid, gmo.Slug, gmo.UserUuid)
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	gm.client.Auth.SetRequestAuth(request)

	response, err := gm.client.HttpClient.Do(request)
	if err != nil {
		return nil
	}
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("response status code was not 204")
	}

	return nil
}

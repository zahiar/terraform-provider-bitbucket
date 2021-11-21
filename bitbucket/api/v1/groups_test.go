package v1

import (
	"os"
	"testing"
)

func TestGroups(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("ENV TF_ACC=1 not set")
	}

	c := NewClient(&Auth{
		Username: os.Getenv("BITBUCKET_USERNAME"),
		Password: os.Getenv("BITBUCKET_PASSWORD"),
	})

	var groupResourceSlug string

	name := "tf-bb-group-test"

	t.Run("create", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: c.Auth.Username,
			Name:      name,
		}

		group, err := c.Groups.Create(opt)
		if err != nil {
			t.Error(err)
		}

		if group.Name != name {
			t.Error("The Group `name` attribute does not match the expected value.")
		}
		if group.AutoAdd != false {
			t.Error("The Group `auto_add` attribute does not match the expected value.")
		}
		if group.Permission != "" {
			t.Error("The Group `permission` attribute does not match the expected value.")
		}

		groupResourceSlug = group.Slug
	})

	t.Run("get", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: c.Auth.Username,
			Slug:      groupResourceSlug,
		}
		group, err := c.Groups.Get(opt)
		if err != nil {
			t.Error(err)
		}

		if group.Name != name {
			t.Error("The Group `name` attribute does not match the expected value.")
		}
		if group.AutoAdd != false {
			t.Error("The Group `auto_add` attribute does not match the expected value.")
		}
		if group.Permission != "" {
			t.Error("The Group `permission` attribute does not match the expected value.")
		}
		if group.Slug != groupResourceSlug {
			t.Error("The Group `slug` attribute does not match the expected value.")
		}
	})

	t.Run("update", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid:  c.Auth.Username,
			Slug:       groupResourceSlug,
			AutoAdd:    true,
			Permission: "write",
		}
		group, err := c.Groups.Update(opt)
		if err != nil {
			t.Error(err)
		}

		if group.Name != name {
			t.Error("The Group `name` attribute does not match the expected value.")
		}
		if group.AutoAdd != true {
			t.Error("The Group `auto_add` attribute does not match the expected value.")
		}
		if group.Permission != "write" {
			t.Error("The Group `permission` attribute does not match the expected value.")
		}
		if group.Slug != groupResourceSlug {
			t.Error("The Group `slug` attribute does not match the expected value.")
		}
	})

	t.Run("delete", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: c.Auth.Username,
			Slug:      groupResourceSlug,
		}
		err := c.Groups.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})
}

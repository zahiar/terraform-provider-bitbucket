package v1

import (
	"os"
	"testing"
)

func TestGroupMembers(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("ENV TF_ACC=1 not set")
	}

	c := NewClient(&Auth{
		Username: os.Getenv("BITBUCKET_USERNAME"),
		Password: os.Getenv("BITBUCKET_PASSWORD"),
	})

	var group *Group

	t.Run("setup", func(t *testing.T) {
		group, _ = c.Groups.Create(
			&GroupOptions{
				OwnerUuid: c.Auth.Username,
				Name:      "tf-bb-group-members-test",
			},
		)
		if group == nil {
			t.Error("The Group could not be created.")
		}
	})

	t.Run("Create", func(t *testing.T) {
		result, err := c.GroupMembers.Create(
			&GroupMemberOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
				UserUuid:  group.Owner.Uuid,
			},
		)

		if err != nil {
			t.Error(err)
		}

		if result == nil {
			t.Error("Expected result")
		}
		if (*result).UUID != group.Owner.Uuid {
			t.Error("The GroupMember list contains an unexpected member.")
		}
	})

	t.Run("Get", func(t *testing.T) {
		members, err := c.GroupMembers.Get(
			&GroupMemberOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
			},
		)
		if err != nil {
			t.Error(err)
		}

		if len(members) != 1 {
			t.Error("The GroupMember list contains unexpected members.")
		}
		if members[0].UUID != group.Owner.Uuid {
			t.Error("The GroupMember list contains an unexpected member.")
		}
	})

	t.Run("teardown", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: c.Auth.Username,
			Slug:      group.Slug,
		}
		err := c.Groups.Delete(opt)
		if err != nil {
			t.Error(err)
		}
	})
}

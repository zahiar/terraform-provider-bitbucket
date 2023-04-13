package v1

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/stretchr/testify/assert"
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
				Name:      "tf-bb-group-members-test" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum),
			},
		)
		assert.NotNil(t, group, "The Group could not be created")
	})

	t.Run("create", func(t *testing.T) {
		result, err := c.GroupMembers.Create(
			&GroupMemberOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
				UserUuid:  group.Owner.Uuid,
			},
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, group.Owner.Uuid, result.UUID, "The GroupMember list contains an unexpected member")
	})

	t.Run("get", func(t *testing.T) {
		members, err := c.GroupMembers.Get(
			&GroupMemberOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(members), "The GroupMember list contains unexpected members")
		assert.Equal(t, group.Owner.Uuid, members[0].UUID, "The GroupMember list contains an unexpected member")
	})

	t.Run("delete", func(t *testing.T) {
		err := c.GroupMembers.Delete(
			&GroupMemberOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
				UserUuid:  group.Owner.Uuid,
			},
		)
		assert.NoError(t, err)

		members, err := c.GroupMembers.Get(
			&GroupMemberOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(members), "The GroupMember list contains unexpected members after deleting the member")
	})

	t.Run("teardown", func(t *testing.T) {
		opt := &GroupOptions{
			OwnerUuid: c.Auth.Username,
			Slug:      group.Slug,
		}
		err := c.Groups.Delete(opt)
		assert.NoError(t, err)
	})
}

package v1

import (
	"os"
	"testing"

	gobb "github.com/ktrysmt/go-bitbucket"
	"github.com/stretchr/testify/assert"
)

func TestGroupPrivileges(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("ENV TF_ACC=1 not set")
	}

	c := NewClient(&Auth{
		Username: os.Getenv("BITBUCKET_USERNAME"),
		Password: os.Getenv("BITBUCKET_PASSWORD"),
	})

	gobbClient := gobb.NewBasicAuth(
		os.Getenv("BITBUCKET_USERNAME"),
		os.Getenv("BITBUCKET_PASSWORD"),
	)

	var group *Group
	var repo *gobb.Repository
	var project *gobb.Project

	t.Run("setup", func(t *testing.T) {
		group, _ = c.Groups.Create(
			&GroupOptions{
				OwnerUuid: c.Auth.Username,
				Name:      "tf-bb-group-test",
			},
		)
		assert.NotNil(t, group)

		project, _ = gobbClient.Workspaces.CreateProject(
			&gobb.ProjectOptions{
				Owner:     c.Auth.Username,
				Name:      "tf-bb-proj-test",
				Key:       "TF_BB_PROJ_TEST",
				IsPrivate: true,
			},
		)
		assert.NotNil(t, project)

		repo, _ = gobbClient.Repositories.Repository.Create(
			&gobb.RepositoryOptions{
				Owner:      c.Auth.Username,
				RepoSlug:   "tf-bb-repo-test",
				ForkPolicy: "no_forks",
				Project:    project.Key,
				IsPrivate:  "true",
			},
		)
		assert.NotNil(t, repo)
	})

	t.Run("create", func(t *testing.T) {
		opt := &GroupPrivilegeOptions{
			WorkspaceId: c.Auth.Username,
			RepoSlug:    repo.Slug,
			GroupOwner:  group.Owner.Uuid,
			GroupSlug:   group.Slug,
			Privilege:   "write",
		}
		groupPrivilege, err := c.GroupPrivileges.Create(opt)

		assert.NoError(t, err)
		assert.Equal(t, "write", groupPrivilege.Privilege)
		assert.Equal(t, group.Slug, groupPrivilege.Group.Slug)
		assert.Equal(t, repo.Slug, groupPrivilege.Repository.Slug)
	})

	t.Run("get", func(t *testing.T) {
		opt := &GroupPrivilegeOptions{
			WorkspaceId: c.Auth.Username,
			RepoSlug:    repo.Slug,
			GroupOwner:  group.Owner.Uuid,
			GroupSlug:   group.Slug,
		}
		groupPrivilege, err := c.GroupPrivileges.Get(opt)

		assert.NoError(t, err)
		assert.Equal(t, "write", groupPrivilege.Privilege)
		assert.Equal(t, group.Slug, groupPrivilege.Group.Slug)
		assert.Equal(t, repo.Slug, groupPrivilege.Repository.Slug)
	})

	t.Run("delete", func(t *testing.T) {
		err := c.GroupPrivileges.Delete(
			&GroupPrivilegeOptions{
				WorkspaceId: c.Auth.Username,
				RepoSlug:    repo.Slug,
				GroupOwner:  group.Owner.Uuid,
				GroupSlug:   group.Slug,
			},
		)
		assert.NoError(t, err)
	})

	t.Run("teardown", func(t *testing.T) {
		err := c.Groups.Delete(
			&GroupOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
			},
		)
		assert.NoError(t, err)

		_, err = gobbClient.Repositories.Repository.Delete(
			&gobb.RepositoryOptions{
				Owner:    c.Auth.Username,
				RepoSlug: repo.Slug,
				Project:  project.Key,
			},
		)
		assert.NoError(t, err)

		_, err = gobbClient.Workspaces.DeleteProject(
			&gobb.ProjectOptions{
				Owner: c.Auth.Username,
				Name:  project.Name,
				Key:   project.Key,
			},
		)
		assert.NoError(t, err)
	})
}

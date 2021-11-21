package v1

import (
	"os"
	"testing"

	gobb "github.com/ktrysmt/go-bitbucket"
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
		if group == nil {
			t.Error("The Group could not be created.")
		}

		project, _ = gobbClient.Workspaces.CreateProject(
			&gobb.ProjectOptions{
				Owner:     c.Auth.Username,
				Name:      "tf-bb-proj-test",
				Key:       "TF_BB_PROJ_TEST",
				IsPrivate: true,
			},
		)
		if project == nil {
			t.Error("The Project could not be created.")
		}

		repo, _ = gobbClient.Repositories.Repository.Create(
			&gobb.RepositoryOptions{
				Owner:      c.Auth.Username,
				RepoSlug:   "tf-bb-repo-test",
				ForkPolicy: "no_forks",
				Project:    project.Key,
				IsPrivate:  "true",
			},
		)
		if repo == nil {
			t.Error("The Repository could not be created.")
		}
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
		if err != nil {
			t.Error(err)
		}

		if groupPrivilege.Privilege != "write" {
			t.Error("The Group Privilege `privilege` attribute does not match the expected value.")
		}
		if groupPrivilege.Group.Slug != group.Slug {
			t.Error("The Group Privilege `group.slug` attribute does not match the expected value.")
		}
		if groupPrivilege.Repository.Slug != repo.Slug {
			t.Error("The Group Privilege `repo.slug` attribute does not match the expected value.")
		}
	})

	t.Run("get", func(t *testing.T) {
		opt := &GroupPrivilegeOptions{
			WorkspaceId: c.Auth.Username,
			RepoSlug:    repo.Slug,
			GroupOwner:  group.Owner.Uuid,
			GroupSlug:   group.Slug,
		}
		groupPrivilege, err := c.GroupPrivileges.Get(opt)
		if err != nil {
			t.Error(err)
		}

		if groupPrivilege.Privilege != "write" {
			t.Error("The Group Privilege `privilege` attribute does not match the expected value.")
		}
		if groupPrivilege.Group.Slug != group.Slug {
			t.Error("The Group Privilege `group.slug` attribute does not match the expected value.")
		}
		if groupPrivilege.Repository.Slug != repo.Slug {
			t.Error("The Group Privilege `repo.slug` attribute does not match the expected value.")
		}
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
		if err != nil {
			t.Error(err)
		}

		err = c.Groups.Delete(
			&GroupOptions{
				OwnerUuid: c.Auth.Username,
				Slug:      group.Slug,
			},
		)
		if err != nil {
			t.Error(err)
		}

		_, err = gobbClient.Repositories.Repository.Delete(
			&gobb.RepositoryOptions{
				Owner:    c.Auth.Username,
				RepoSlug: repo.Slug,
				Project:  project.Key,
			},
		)
		if err != nil {
			t.Error(err)
		}

		_, err = gobbClient.Workspaces.DeleteProject(
			&gobb.ProjectOptions{
				Owner: c.Auth.Username,
				Name:  project.Name,
				Key:   project.Key,
			},
		)
		if err != nil {
			t.Error(err)
		}
	})
}

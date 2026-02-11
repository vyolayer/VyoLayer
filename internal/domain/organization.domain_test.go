package domain_test

import (
	"testing"
	"worklayer/internal/domain"
)

var domainUser, _ = domain.NewUser(
	"worklayer@test.com",
	"worklayer_password",
	"worklayer fullName",
)

func TestNewOrganization(t *testing.T) {
	org := domain.NewOrganization(
		domainUser,
		"worklayer organization",
		"worklayer organization description",
		nil,
		nil,
	)

	t.Log("Organization ID: ", org.ID)
	t.Log("Organization Name: ", org.Name)
	t.Log("Organization Description: ", org.Description)
	t.Log("Organization Slug: ", org.Slug)
	t.Log("Organization Owner ID: ", org.OwnerID)
	t.Log("Organization IsActive: ", org.IsActive)
	t.Log("Organization Max Projects: ", org.MaxProjects)
	t.Log("Organization Max Members: ", org.MemberInfo.MaxNoOfMembers)
	t.Log("Organization NoOfMembers: ", org.MemberInfo.NoOfMembers)
}

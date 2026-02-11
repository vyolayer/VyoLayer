package domain

import (
	"log"
	"time"

	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"
	"worklayer/pkg/utils"
)

var (
	DefaultMaxProjects = 1
	DefaultMaxMembers  = 5
)

type memberInfo struct {
	isInfoLoaded   bool
	MaxNoOfMembers int
	NoOfMembers    int
	Members        []OrganizationMember
}

type Organization struct {
	ID            types.OrganizationID
	Name          string
	Slug          string
	Description   string
	OwnerID       types.UserID
	IsActive      bool
	DeactivatedBy *types.UserID
	DeactivatedAt *time.Time
	MaxProjects   int

	MemberInfo memberInfo
}

func NewOrganization(
	owner *User,
	name string,
	description string,
	maxProjects *int,
	maxMembers *int,
) *Organization {
	id := types.NewOrganizationID()

	member := NewOrganizationMember(id, nil, nil, owner)
	slugify := utils.ToSlug(name).Slugify().String()

	if maxProjects == nil {
		maxProjects = &DefaultMaxProjects
	}

	if maxMembers == nil {
		maxMembers = &DefaultMaxMembers
	}

	return &Organization{
		ID:          id,
		Name:        name,
		Slug:        slugify,
		Description: description,
		OwnerID:     owner.ID,
		IsActive:    true,
		MaxProjects: *maxProjects,
		MemberInfo: memberInfo{
			MaxNoOfMembers: *maxMembers,
			NoOfMembers:    1,
			Members:        []OrganizationMember{*member},
		},
	}
}

// ReconstructOrganization reconstructs an organization from database data
func ReconstructOrganization(
	id types.OrganizationID,
	name, slug, description string,
	ownerID types.UserID,
	isActive bool,
	deactivatedBy *types.UserID,
	deactivatedAt *time.Time,
	maxProjects, maxMembers int,
) *Organization {
	return &Organization{
		ID:            id,
		Name:          name,
		Slug:          slug,
		Description:   description,
		OwnerID:       ownerID,
		IsActive:      isActive,
		DeactivatedBy: deactivatedBy,
		DeactivatedAt: deactivatedAt,
		MaxProjects:   maxProjects,
		MemberInfo: memberInfo{
			isInfoLoaded:   false,
			MaxNoOfMembers: maxMembers,
			NoOfMembers:    0,
			Members:        []OrganizationMember{},
		},
	}
}

// Deactivate deactivates the organization
func (o *Organization) Deactivate(deactivatedBy types.UserID) *errors.AppError {
	if !o.IsActive {
		return OrganizationNotActiveError()
	}

	now := time.Now()
	o.IsActive = false
	o.DeactivatedBy = &deactivatedBy
	o.DeactivatedAt = &now

	return nil
}

// Reactivate reactivates the organization
func (o *Organization) Reactivate() *errors.AppError {
	if o.IsActive {
		return nil // Already active
	}

	o.IsActive = true
	o.DeactivatedBy = nil
	o.DeactivatedAt = nil

	return nil
}

// UpdateName updates the organization name and regenerates slug
func (o *Organization) UpdateName(name string) {
	o.Name = name
	o.Slug = utils.ToSlug(name).Slugify().String()
}

// UpdateDescription updates the organization description
func (o *Organization) UpdateDescription(description string) {
	o.Description = description
}

// UpdateMaxProjects updates the maximum number of projects
func (o *Organization) UpdateMaxProjects(maxProjects int) {
	o.MaxProjects = maxProjects
}

// UpdateMaxMembers updates the maximum number of members
func (o *Organization) UpdateMaxMembers(maxMembers int) {
	o.MemberInfo.MaxNoOfMembers = maxMembers
}

// AddMember adds a member to the organization
func (o *Organization) AddMember(member *OrganizationMember) *errors.AppError {
	if !o.MemberInfo.isInfoLoaded {
		return OrganizationMembersNotLoadedError()
	}

	if !o.CanAddMember() {
		return OrganizationFullError()
	}

	// Check if member already exists
	for _, m := range o.MemberInfo.Members {
		if m.UserID == member.UserID {
			return OrganizationMemberAlreadyExistsError(member.UserID.String())
		}
	}

	o.MemberInfo.Members = append(o.MemberInfo.Members, *member)
	o.MemberInfo.NoOfMembers++

	return nil
}

// RemoveMember removes a member from the organization
func (o *Organization) RemoveMember(memberID types.OrganizationMemberID) *errors.AppError {
	if !o.MemberInfo.isInfoLoaded {
		return OrganizationMembersNotLoadedError()
	}

	for i, member := range o.MemberInfo.Members {
		if member.ID.String() == memberID.String() {
			// Cannot remove the owner
			if member.UserID == o.OwnerID {
				return OrganizationCannotRemoveOwnerError()
			}

			// Remove member
			o.MemberInfo.Members = append(o.MemberInfo.Members[:i], o.MemberInfo.Members[i+1:]...)
			o.MemberInfo.NoOfMembers--
			return nil
		}
	}

	return OrganizationMemberNotFoundError(memberID.String())
}

// CanAddMember checks if a new member can be added
func (o *Organization) CanAddMember() bool {
	return o.MemberInfo.NoOfMembers < o.MemberInfo.MaxNoOfMembers
}

// GetMembers returns all members if loaded
func (o *Organization) GetMembers() ([]OrganizationMember, *errors.AppError) {
	if !o.MemberInfo.isInfoLoaded {
		return nil, OrganizationMembersNotLoadedError()
	}
	return o.MemberInfo.Members, nil
}

// LoadMembers loads members into the organization
func (o *Organization) LoadMembers(members []OrganizationMember) {
	o.MemberInfo.Members = members
	o.MemberInfo.NoOfMembers = len(members)
	o.MemberInfo.isInfoLoaded = true
}

// IsMember checks if a user is a member of the organization
func (o *Organization) IsMember(userID types.UserID) bool {
	if !o.MemberInfo.isInfoLoaded {
		return false
	}

	for _, member := range o.MemberInfo.Members {
		if member.UserID.String() == userID.String() {
			log.Println("Member found")
			return true
		}
	}
	log.Println("Member not found")
	return false
}

// IsOwner checks if a user is the owner of the organization
func (o *Organization) IsOwner(userID types.UserID) bool {
	return o.OwnerID == userID
}

// GetMemberByUserID retrieves a member by user ID
func (o *Organization) GetMemberByUserID(userID types.UserID) (*OrganizationMember, *errors.AppError) {
	if !o.MemberInfo.isInfoLoaded {
		return nil, OrganizationMembersNotLoadedError()
	}

	for _, member := range o.MemberInfo.Members {
		if member.UserID.String() == userID.String() {
			return &member, nil
		}
	}

	return nil, OrganizationMemberNotFoundError(userID.String())
}

// GetActiveMemberCount returns the count of active members
func (o *Organization) GetActiveMemberCount() int {
	if !o.MemberInfo.isInfoLoaded {
		return 0
	}

	count := 0
	for _, member := range o.MemberInfo.Members {
		if member.IsActive {
			count++
		}
	}
	return count
}

// Validate validates the organization
func (o *Organization) Validate() *errors.AppError {
	if o.Name == "" {
		return ValidationError("Organization name is required")
	}

	if o.Slug == "" {
		return ValidationError("Organization slug is required")
	}

	if o.MaxProjects < 0 {
		return ValidationError("Max projects cannot be negative")
	}

	if o.MemberInfo.MaxNoOfMembers < 1 {
		return ValidationError("Max members must be at least 1")
	}

	return nil
}

package mapper

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
)

// ToDomainProject converts a Project model to a Project domain object.
func ToDomainProject(m *models.Project) *domain.Project {
	if m == nil {
		return nil
	}

	return domain.ReconstructProject(
		m.PublicID(),
		m.OrganizationPublicID(),
		m.Name,
		m.Slug,
		m.Description,
		m.IsActive,
		m.Creator.PublicID(),
		m.MaxApiKeys,
		m.MaxMembers,
	)
}

// ToDomainProjectWithMembers converts with members loaded.
func ToDomainProjectWithMembers(m *models.Project) *domain.Project {
	p := ToDomainProject(m)
	if p == nil {
		return nil
	}

	members := make([]domain.ProjectMember, 0, len(m.Members))
	for _, mm := range m.Members {
		if pm := ToDomainProjectMember(&mm); pm != nil {
			members = append(members, *pm)
		}
	}
	p.LoadMembers(members)
	return p
}

// ToDomainProjectMember converts a ProjectMember model to domain.
func ToDomainProjectMember(m *models.ProjectMember) *domain.ProjectMember {
	if m == nil {
		return nil
	}

	pmID, _ := types.ReconstructProjectMemberID(m.ID.String())
	projectID, _ := types.ReconstructProjectID(m.ProjectID.String())
	userID, _ := types.ReconstructUserID(m.UserID.String())
	addedByID, _ := types.ReconstructUserID(m.AddedBy.String())

	var removedBy *types.UserID
	if m.RemovedBy != nil {
		id, _ := types.ReconstructUserID(m.RemovedBy.String())
		removedBy = &id
	}

	return domain.ReconstructProjectMember(
		pmID,
		projectID,
		userID,
		m.Role,
		m.User.Email,
		m.User.FullName,
		m.IsActive(),
		addedByID,
		m.JoinedAt,
		m.RemovedAt,
		removedBy,
	)
}

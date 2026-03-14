package mapper

import (
	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
)

// ToDomainOrganizationMemberInvitation converts an OrganizationMemberInvitation model to domain object
func ToDomainOrganizationMemberInvitation(invitationModel *models.OrganizationMemberInvitation) *domain.OrganizationMemberInvitation {
	if invitationModel == nil {
		return nil
	}

	// Parse role IDs from JSON
	roleIDs, _ := domain.UnmarshalRoleIDs(invitationModel.RoleIDs)

	var deletedBy *types.OrganizationMemberID
	if invitationModel.DeletedBy != nil {
		id, _ := types.ReconstructOrganizationMemberID(invitationModel.DeletedBy.String())
		deletedBy = &id
	}

	return domain.ReconstructOrganizationMemberInvitation(
		invitationModel.PublicID(),
		invitationModel.OrganizationPublicID(),
		invitationModel.InvitedByPublicID(),
		invitationModel.Email,
		invitationModel.Token,
		roleIDs,
		invitationModel.InvitedAt,
		invitationModel.IsAccepted,
		invitationModel.AcceptedAt,
		invitationModel.ExpiredAt,
		deletedBy,
	)
}

// ToModelOrganizationMemberInvitation converts a domain invitation to model
func ToModelOrganizationMemberInvitation(inv *domain.OrganizationMemberInvitation) *models.OrganizationMemberInvitation {
	if inv == nil {
		return nil
	}

	// Marshal role IDs to JSON
	roleIDsJSON, _ := domain.MarshalRoleIDs(inv.ToRoleIDsString())

	var deletedBy *uuid.UUID
	if inv.DeletedBy != nil {
		id := (*inv.DeletedBy).InternalID().ID()
		deletedBy = &id
	}

	return &models.OrganizationMemberInvitation{
		BaseModel: models.BaseModel{
			ID: inv.ID.InternalID().ID(),
		},
		OrganizationID: inv.OrganizationID.InternalID().ID(),
		InvitedBy:      inv.InvitedBy.InternalID().ID(),
		Email:          inv.Email,
		Token:          inv.Token,
		RoleIDs:        roleIDsJSON,
		InvitedAt:      inv.InvitedAt,
		IsAccepted:     inv.IsAccepted,
		AcceptedAt:     inv.AcceptedAt,
		ExpiredAt:      inv.ExpiredAt,
		DeletedBy:      deletedBy,
	}
}

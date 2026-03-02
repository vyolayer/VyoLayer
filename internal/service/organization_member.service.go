package service

import (
	"fmt"
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
	"worklayer/pkg/cache"
	"worklayer/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

// currentMemberCacheKey returns a unique cache key scoped to an org+user pair.
func currentMemberCacheKey(orgID types.OrganizationID, userID types.UserID) string {
	return fmt.Sprintf("current_member:%s:%s", orgID.String(), userID.String())
}

type OrganizationMemberService interface {
	ListByOrgAndUserId(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		userID types.UserID,
	) ([]dto.OrganizationMemberDTO, error)

	GetOrgMemberByMemberID(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		memberID types.OrganizationMemberID,
	) (dto.OrganizationMemberDTO, error)

	GetCurrentMember(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		userID types.UserID,
	) (domain.OrganizationMemberWithRBAC, error)

	InvalidateCurrentMember(
		orgID types.OrganizationID,
		userID types.UserID,
	)

	RemoveMember(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		actorUserID types.UserID,
		targetMemberID types.OrganizationMemberID,
	) *errors.AppError

	ChangeMemberRole(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		actorUserID types.UserID,
		targetMemberID types.OrganizationMemberID,
		newRoleID types.OrganizationRoleID,
	) *errors.AppError

	LeaveOrganization(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		userID types.UserID,
	) *errors.AppError

	TransferOwnership(
		ctx *fiber.Ctx,
		orgID types.OrganizationID,
		currentOwnerUserID types.UserID,
		newOwnerMemberID types.OrganizationMemberID,
	) *errors.AppError
}

type organizationMemberService struct {
	orgMemberRepo repository.OrganizationMemberRepository
	auditLogRepo  repository.AuditLogRepository
	rbacRepo      repository.OrganizationRBACRepository
	memberCache   cache.Cache[domain.OrganizationMemberWithRBAC]
}

func NewOrganizationMemberService(
	orgMemberRepo repository.OrganizationMemberRepository,
	auditLogRepo repository.AuditLogRepository,
	rbacRepo repository.OrganizationRBACRepository,
) OrganizationMemberService {
	return &organizationMemberService{
		orgMemberRepo: orgMemberRepo,
		auditLogRepo:  auditLogRepo,
		rbacRepo:      rbacRepo,
		memberCache:   cache.NewMemoryCache[domain.OrganizationMemberWithRBAC](),
	}
}

const currentMemberTTL = 5 * time.Minute

func (service *organizationMemberService) ListByOrgAndUserId(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	userID types.UserID,
) ([]dto.OrganizationMemberDTO, error) {
	members, err := service.orgMemberRepo.GetByOrgID(ctx.Context(), userID, orgID)
	if err != nil {
		return nil, WrapRepositoryError(err, "listing organization members")
	}

	membersDto := make([]dto.OrganizationMemberDTO, 0, len(members))
	for _, member := range members {
		membersDto = append(membersDto, dto.FromDomainOrganizationMember(&member))
	}

	return membersDto, nil
}

func (service *organizationMemberService) GetCurrentMember(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	userID types.UserID,
) (domain.OrganizationMemberWithRBAC, error) {
	cacheKey := currentMemberCacheKey(orgID, userID)

	if cached, ok := ctx.Locals(cacheKey).(domain.OrganizationMemberWithRBAC); ok {
		return cached, nil
	}

	if cached, ok := service.memberCache.Get(cacheKey); ok {
		ctx.Locals(cacheKey, cached)
		return cached, nil
	}

	member, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if err != nil {
		return domain.OrganizationMemberWithRBAC{}, WrapRepositoryError(err, "getting organization member")
	}

	service.memberCache.Set(cacheKey, *member, currentMemberTTL)
	ctx.Locals(cacheKey, *member)

	return *member, nil
}

func (service *organizationMemberService) InvalidateCurrentMember(
	orgID types.OrganizationID,
	userID types.UserID,
) {
	service.memberCache.Delete(currentMemberCacheKey(orgID, userID))
}

func (service *organizationMemberService) GetOrgMemberByMemberID(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	memberID types.OrganizationMemberID,
) (dto.OrganizationMemberDTO, error) {
	member, err := service.orgMemberRepo.GetByOrgIDAndMemberID(ctx.Context(), orgID, memberID)
	if err != nil {
		return dto.OrganizationMemberDTO{}, WrapRepositoryError(err, "getting organization member")
	}

	return dto.FromDomainOrganizationMember(member), nil
}

// RemoveMember removes a member from the organization.
func (service *organizationMemberService) RemoveMember(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	actorUserID types.UserID,
	targetMemberID types.OrganizationMemberID,
) *errors.AppError {
	actor, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, actorUserID)
	if err != nil {
		return WrapRepositoryError(err, "get actor member")
	}
	if !actor.IsAdmin() {
		return errors.Forbidden("You don't have permission to remove members")
	}

	target, targetErr := service.orgMemberRepo.GetByOrgIDAndMemberID(ctx.Context(), orgID, targetMemberID)
	if targetErr != nil {
		return WrapRepositoryError(targetErr, "get target member")
	}

	if target.UserID.String() == actorUserID.String() {
		return errors.BadRequest("Use the leave endpoint to remove yourself from the organization")
	}

	targetRBAC, targetRBACErr := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, target.UserID)
	if targetRBACErr != nil {
		return WrapRepositoryError(targetRBACErr, "get target member RBAC")
	}

	if targetRBAC.IsOwner() && !actor.IsOwner() {
		return errors.Forbidden("Admin cannot remove an Owner")
	}

	if targetRBAC.IsOwner() {
		ownerCount, countErr := service.orgMemberRepo.CountOwners(ctx.Context(), orgID)
		if countErr != nil {
			return WrapRepositoryError(countErr, "count owners")
		}
		if ownerCount <= 1 {
			return domain.OrganizationLastOwnerError()
		}
	}

	if deleteErr := service.orgMemberRepo.Delete(ctx.Context(), targetMemberID); deleteErr != nil {
		return WrapRepositoryError(deleteErr, "remove member")
	}

	service.InvalidateCurrentMember(orgID, target.UserID)

	_ = service.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        actor.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "member.removed",
		ResourceType:   "member",
		ResourceID:     targetMemberID.InternalID().ID(),
		Metadata: map[string]interface{}{
			"email": target.Email,
		},
	})

	return nil
}

// ChangeMemberRole changes a member's role.
func (service *organizationMemberService) ChangeMemberRole(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	actorUserID types.UserID,
	targetMemberID types.OrganizationMemberID,
	newRoleID types.OrganizationRoleID,
) *errors.AppError {
	actor, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, actorUserID)
	if err != nil {
		return WrapRepositoryError(err, "get actor member")
	}
	if !actor.IsAdmin() {
		return errors.Forbidden("You don't have permission to change roles")
	}

	target, targetErr := service.orgMemberRepo.GetByOrgIDAndMemberID(ctx.Context(), orgID, targetMemberID)
	if targetErr != nil {
		return WrapRepositoryError(targetErr, "get target member")
	}

	targetRBAC, targetRBACErr := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, target.UserID)
	if targetRBACErr != nil {
		return WrapRepositoryError(targetRBACErr, "get target member RBAC")
	}

	if targetRBAC.IsOwner() && !actor.IsOwner() {
		return errors.Forbidden("Admin cannot change the role of an Owner")
	}

	roles, rolesErr := service.rbacRepo.GetAllRoles(orgID)
	if rolesErr != nil {
		return WrapRepositoryError(rolesErr, "get all roles")
	}

	var newRoleName string
	for _, r := range roles {
		if r.PublicID().String() == newRoleID.String() {
			newRoleName = r.Name
			break
		}
	}

	if newRoleName == "" {
		return errors.BadRequest("Invalid role ID")
	}

	if newRoleName == "Owner" && !actor.IsOwner() {
		return errors.Forbidden("Only owners can assign the Owner role")
	}

	if targetRBAC.IsOwner() && newRoleName != "Owner" {
		ownerCount, countErr := service.orgMemberRepo.CountOwners(ctx.Context(), orgID)
		if countErr != nil {
			return WrapRepositoryError(countErr, "count owners")
		}
		if ownerCount <= 1 {
			return domain.OrganizationLastOwnerError()
		}
	}

	if revokeErr := service.orgMemberRepo.RevokeAllRoles(ctx.Context(), targetMemberID, orgID); revokeErr != nil {
		return WrapRepositoryError(revokeErr, "revoke existing roles")
	}

	if assignErr := service.orgMemberRepo.AssignRole(ctx.Context(), targetMemberID, orgID, newRoleID, actor.ID); assignErr != nil {
		return WrapRepositoryError(assignErr, "assign new role")
	}

	service.InvalidateCurrentMember(orgID, target.UserID)

	_ = service.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        actor.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "member.role_changed",
		ResourceType:   "member",
		ResourceID:     targetMemberID.InternalID().ID(),
		Metadata: map[string]interface{}{
			"new_role": newRoleName,
		},
	})

	return nil
}

// LeaveOrganization allows a member to leave. Owner cannot leave if last owner.
func (service *organizationMemberService) LeaveOrganization(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	userID types.UserID,
) *errors.AppError {
	member, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if err != nil {
		return WrapRepositoryError(err, "get current member")
	}

	if member.IsOwner() {
		ownerCount, countErr := service.orgMemberRepo.CountOwners(ctx.Context(), orgID)
		if countErr != nil {
			return WrapRepositoryError(countErr, "count owners")
		}
		if ownerCount <= 1 {
			return errors.BadRequest("You are the last owner. Transfer ownership before leaving.")
		}
	}

	if deleteErr := service.orgMemberRepo.Delete(ctx.Context(), member.ID); deleteErr != nil {
		return WrapRepositoryError(deleteErr, "leave organization")
	}

	service.InvalidateCurrentMember(orgID, userID)

	_ = service.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        member.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "member.left",
		ResourceType:   "member",
		ResourceID:     member.ID.InternalID().ID(),
	})

	return nil
}

// TransferOwnership transfers ownership to another member.
func (service *organizationMemberService) TransferOwnership(
	ctx *fiber.Ctx,
	orgID types.OrganizationID,
	currentOwnerUserID types.UserID,
	newOwnerMemberID types.OrganizationMemberID,
) *errors.AppError {
	actor, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, currentOwnerUserID)
	if err != nil {
		return WrapRepositoryError(err, "get current owner member")
	}
	if !actor.IsOwner() {
		return errors.Forbidden("Only the owner can transfer ownership")
	}

	target, targetErr := service.orgMemberRepo.GetByOrgIDAndMemberID(ctx.Context(), orgID, newOwnerMemberID)
	if targetErr != nil {
		return WrapRepositoryError(targetErr, "get target member for ownership transfer")
	}

	roles, rolesErr := service.rbacRepo.GetAllRoles(orgID)
	if rolesErr != nil {
		return WrapRepositoryError(rolesErr, "get all roles")
	}

	var ownerRoleID, memberRoleID types.OrganizationRoleID
	for _, r := range roles {
		if r.Name == "Owner" {
			ownerRoleID = r.PublicID()
		} else if r.Name == "Member" {
			memberRoleID = r.PublicID()
		}
	}

	if ownerRoleID == nil || memberRoleID == nil {
		return errors.Internal("Owner or Member role not found in system")
	}

	// Downgrade current owner to Member
	if revokeErr := service.orgMemberRepo.RevokeAllRoles(ctx.Context(), actor.ID, orgID); revokeErr != nil {
		return WrapRepositoryError(revokeErr, "revoke current owner roles")
	}
	if assignErr := service.orgMemberRepo.AssignRole(ctx.Context(), actor.ID, orgID, memberRoleID, actor.ID); assignErr != nil {
		return WrapRepositoryError(assignErr, "assign member role to current owner")
	}

	// Upgrade new owner to Owner
	if revokeErr := service.orgMemberRepo.RevokeAllRoles(ctx.Context(), newOwnerMemberID, orgID); revokeErr != nil {
		return WrapRepositoryError(revokeErr, "revoke new owner roles")
	}
	if assignErr := service.orgMemberRepo.AssignRole(ctx.Context(), newOwnerMemberID, orgID, ownerRoleID, actor.ID); assignErr != nil {
		return WrapRepositoryError(assignErr, "assign owner role to new owner")
	}

	service.InvalidateCurrentMember(orgID, currentOwnerUserID)
	service.InvalidateCurrentMember(orgID, target.UserID)

	_ = service.auditLogRepo.Create(ctx.Context(), &repository.AuditLogEntry{
		OrganizationID: orgID.InternalID().ID(),
		ActorID:        actor.ID.InternalID().ID(),
		ActorType:      "member",
		Action:         "org.ownership_transferred",
		ResourceType:   "member",
		ResourceID:     newOwnerMemberID.InternalID().ID(),
		Severity:       "warning",
		Metadata: map[string]interface{}{
			"previous_owner": actor.ID.String(),
			"new_owner":      newOwnerMemberID.String(),
		},
	})

	return nil
}

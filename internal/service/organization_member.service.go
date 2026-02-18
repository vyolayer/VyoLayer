package service

import (
	"fmt"
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
	"worklayer/pkg/cache"

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

	// InvalidateCurrentMember removes the cached member entry.
	// Call after role changes, membership updates, or removal.
	InvalidateCurrentMember(
		orgID types.OrganizationID,
		userID types.UserID,
	)
}

type organizationMemberService struct {
	orgMemberRepo repository.OrganizationMemberRepository
	// memberCache caches current-member lookups.
	// Defaults to an in-memory cache; can be swapped for any Cache[V] implementation.
	memberCache cache.Cache[domain.OrganizationMemberWithRBAC]
}

// NewOrganizationMemberService creates the service with a default in-memory cache (5 min TTL).
func NewOrganizationMemberService(orgMemberRepo repository.OrganizationMemberRepository) OrganizationMemberService {
	return &organizationMemberService{
		orgMemberRepo: orgMemberRepo,
		memberCache:   cache.NewMemoryCache[domain.OrganizationMemberWithRBAC](),
	}
}

// NewOrganizationMemberServiceWithCache creates the service with a custom cache backend.
func NewOrganizationMemberServiceWithCache(
	orgMemberRepo repository.OrganizationMemberRepository,
	memberCache cache.Cache[domain.OrganizationMemberWithRBAC],
) OrganizationMemberService {
	return &organizationMemberService{
		orgMemberRepo: orgMemberRepo,
		memberCache:   memberCache,
	}
}

// currentMemberTTL is the TTL for the current member cache.
// Keep short since roles/permissions can change.
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

	// 1. Check per-request Locals first (zero-cost within same request)
	if cached, ok := ctx.Locals(cacheKey).(domain.OrganizationMemberWithRBAC); ok {
		return cached, nil
	}

	// 2. Check the shared cache (in-memory with TTL, or custom backend)
	if cached, ok := service.memberCache.Get(cacheKey); ok {
		// Warm the per-request local so subsequent calls in this request are free
		ctx.Locals(cacheKey, cached)
		return cached, nil
	}

	// 3. Cache miss — fetch from DB
	member, err := service.orgMemberRepo.GetCurrentMember(ctx.Context(), orgID, userID)
	if err != nil {
		return domain.OrganizationMemberWithRBAC{}, WrapRepositoryError(err, "getting organization member")
	}

	// Store in both layers
	service.memberCache.Set(cacheKey, *member, currentMemberTTL)
	ctx.Locals(cacheKey, *member)

	return *member, nil
}

// InvalidateCurrentMember removes a member from the cache.
// Call this after role changes, membership updates, or removal.
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

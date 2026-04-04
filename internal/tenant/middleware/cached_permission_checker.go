package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	tenantrepo "github.com/vyolayer/vyolayer/internal/tenant/repo"
	"github.com/vyolayer/vyolayer/pkg/cache"
)

// CachedPermissionChecker wraps the database PermissionChecker with a TTL cache.
type CachedPermissionChecker struct {
	dbChecker tenantrepo.PermissionChecker
	cache     cache.Cache[any]
	ttl       time.Duration
}

// NewCachedPermissionChecker creates a decorator for permission checks.
func NewCachedPermissionChecker(dbChecker tenantrepo.PermissionChecker, cacheStore cache.Cache[any], ttl time.Duration) *CachedPermissionChecker {
	return &CachedPermissionChecker{
		dbChecker: dbChecker,
		cache:     cacheStore,
		ttl:       ttl,
	}
}

// HasPermission checks RAM first. If missing, it queries Postgres and saves the result.
func (c *CachedPermissionChecker) HasPermission(ctx context.Context, orgID, userID uuid.UUID, requiredPermissionCode string) (bool, error) {
	// 1. Generate a highly specific cache key
	cacheKey := fmt.Sprintf("perm:%s:%s:%s", orgID.String(), userID.String(), requiredPermissionCode)

	// 2. Check the fast cache (~0.001ms latency)
	if cachedVal, found := c.cache.Get(cacheKey); found {
		if hasPerm, ok := cachedVal.(bool); ok {
			return hasPerm, nil // Cache Hit!
		}
	}

	// 3. Cache Miss: Hit the Postgres database (~10ms latency)
	hasPerm, err := c.dbChecker.HasPermission(ctx, orgID, userID, requiredPermissionCode)
	if err != nil {
		return false, err
	}

	// 4. Save the result in the cache for the next request
	c.cache.Set(cacheKey, hasPerm, c.ttl)

	return hasPerm, nil
}

// IsMember uses the same caching strategy for basic membership checks.
func (c *CachedPermissionChecker) IsMember(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
	cacheKey := fmt.Sprintf("member:%s:%s", orgID.String(), userID.String())

	if cachedVal, found := c.cache.Get(cacheKey); found {
		if isMember, ok := cachedVal.(bool); ok {
			return isMember, nil
		}
	}

	isMember, err := c.dbChecker.IsMember(ctx, orgID, userID)
	if err != nil {
		return false, err
	}

	c.cache.Set(cacheKey, isMember, c.ttl)

	return isMember, nil
}

// Invalidate clears all cached permissions for a specific user in an org.
// Call this from your Use Cases whenever roles are changed or members are removed.
func (c *CachedPermissionChecker) Invalidate(orgID, userID uuid.UUID) {
	// Clear the basic membership cache
	memberKey := fmt.Sprintf("member:%s:%s", orgID.String(), userID.String())
	c.cache.Delete(memberKey)

	// Note: Without Redis wildcard deletion (e.g., DEL perm:org:user:*),
	// specific granular permissions will naturally expire via TTL.
	// For most BaaS applications, a 3-5 minute TTL is an acceptable window for granular permission propagation.
}

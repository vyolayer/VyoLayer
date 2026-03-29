// Package grpcctx provides a minimal adapter to produce a *fiber.Ctx that
// is safe to pass into service-layer functions which only use ctx.Context()
// (and optionally ctx.Locals() for request-scoped caching).
//
// Background: the current internal/service layer was written against Fiber's
// HTTP handler context. During the gRPC migration we reuse those services
// without refactoring them.  Once services accept context.Context directly
// this package can be removed.
package grpcctx

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Wrap returns a *fiber.Ctx backed by the supplied context.Context.
// Only ctx.Context() is guaranteed to work correctly; ctx.Locals() stores
// values in request-local memory (the fasthttp RequestCtx is nil in this
// shim, so Locals will be a no-op – the service's in-memory cache still
// handles short-term caching).
func Wrap(ctx context.Context) *fiber.Ctx {
	// fiber.Ctx.SetUserContext stores the context and Context() returns it.
	// A zero-value *fiber.Ctx is sufficient for services that only call
	// .Context() and .Locals().
	fctx := new(fiber.Ctx)
	fctx.SetUserContext(ctx)
	return fctx
}

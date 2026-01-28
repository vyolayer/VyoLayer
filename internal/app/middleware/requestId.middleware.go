package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	XRequestID = "X-Request-ID"
	RequestID  = "requestid"
)

type RequestIDMiddleware struct {
	requestID string
}

func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{requestID: ""}
}

func (r *RequestIDMiddleware) RequestIDMiddleware(c *fiber.Ctx) error {
	id := r.getIdFromHeader(c)
	if id == "" && r.requestID != "" {
		id = r.requestID
	}

	if id == "" {
		id = uuid.New().String()
	}

	r.setRequestId(c, id)
	return c.Next()
}

func (r *RequestIDMiddleware) getIdFromHeader(c *fiber.Ctx) string {
	return c.Get(XRequestID)
}

func (r *RequestIDMiddleware) setRequestId(ctx *fiber.Ctx, id string) {
	ctx.Locals(RequestID, id)
	ctx.Set(XRequestID, id)
}

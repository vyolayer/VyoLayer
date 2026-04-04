package grpc

import (
	"time"

	"github.com/vyolayer/vyolayer/pkg/pagination"
)

func getPageSize(count int32) int {
	limit := int(count)
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	return limit
}

func getOffset(pageToken string) int {
	return pagination.DecodePageToken(pageToken)
}

func strPtr(s string) *string {
	return &s
}

func timeToStringPtr(t time.Time) *string {
	if t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

func timePtrToStringPtr(t *time.Time) *string {
	if t == nil || t.IsZero() {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}

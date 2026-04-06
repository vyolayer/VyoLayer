package tenant

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	prefix = "tenant_"
)

type DatabaseSchema interface {
	GenerateSchema() (string, error)
}

type DatabaseSchemaImpl struct {
	tenantID   string
	tenantSlug string
}

func NewDatabaseSchema(tenantID, tenantSlug string) DatabaseSchema {
	return &DatabaseSchemaImpl{
		tenantID:   tenantID,
		tenantSlug: tenantSlug,
	}
}

func (d *DatabaseSchemaImpl) GenerateSchema() (string, error) {
	if d.tenantID == "" || d.tenantSlug == "" {
		return "", errors.New("tenantID and tenantSlug are required")
	}

	var (
		t, s string
	)

	t = d.tenantID
	if len(t) > 6 {
		t = t[:6]
	}

	s = d.tenantSlug
	if len(s) > 6 {
		s = s[:6]
	}

	hashInput := d.tenantID + ":" + d.tenantSlug
	hash := sha1.Sum([]byte(hashInput))
	h := hex.EncodeToString(hash[:])[:6]

	schema := fmt.Sprintf("%s%s_%s_%s", prefix, t, s, h)
	return schema, nil
}

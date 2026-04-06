package tenant

import (
	"errors"
	"regexp"
)

var (
	validSchemaRegex = regexp.MustCompile(`^[a-z][a-z0-9_]{0,62}$`)
)

func ValidateSchemaName(schema string) error {
	if schema == "" {
		return errors.New("schema is empty")
	}
	if len(schema) > 63 {
		return errors.New("schema exceeds 63 characters")
	}
	if !validSchemaRegex.MatchString(schema) {
		return errors.New("invalid schema format (must match ^[a-z][a-z0-9_]{0,62}$)")
	}
	return nil
}

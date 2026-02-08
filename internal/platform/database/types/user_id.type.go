package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UserID struct {
	PublicID
}

// NewUserID creates a new UserID.
func NewUserID() UserID {
	return UserID{PublicID: NewPublicID(UserPrefix)}
}

// ReconstructUserID reconstructs a UserID from a string.
func ReconstructUserID(s string) (*UserID, error) {
	if s == "" {
		return nil, errors.New("id cannot be empty")
	}

	// 1. Clean the string
	// This handles both "user_123..." AND "123..." automatically.
	// We check specifically for your prefix to ensure we don't accidentally
	// parse an 'org_' id as a 'user_' id.
	prefix := UserPrefix.String() + "_"

	// Safety Check: If it has a prefix, but it's the WRONG prefix (e.g. "org_...")
	if strings.Contains(s, "_") && !strings.HasPrefix(s, prefix) {
		return nil, errors.New("id type mismatch")
	}

	// Remove the prefix if it exists
	cleanUUID := strings.TrimPrefix(s, prefix)

	// 2. Parse safely (No MustParse)
	id, err := uuid.Parse(cleanUUID)
	if err != nil {
		return nil, errors.New("invalid uuid format")
	}

	// 3. Construct and return
	return &UserID{
		PublicID: PublicID{
			prefix: UserPrefix,
			uuid:   ReconstructInternalID(id),
		},
	}, nil
}

func ParseUserID(s string) (UserID, error) {
	public, err := ParsePublicID(s)
	if err != nil {
		return UserID{}, err
	}
	if public.prefix != UserPrefix {
		return UserID{}, fmt.Errorf("invalid prefix: %s", public.prefix)
	}
	return UserID{PublicID: public}, nil
}

func (id UserID) MarshalJSON() ([]byte, error) {
	public := PublicID{prefix: UserPrefix, uuid: id.PublicID.uuid}
	return public.MarshalJSON()
}

func (id *UserID) UnmarshalJSON(b []byte) error {
	var public PublicID
	if err := public.UnmarshalJSON(b); err != nil {
		return err
	}

	// Validation: Ensure the ID sent matches the expected type
	if public.prefix != UserPrefix {
		return fmt.Errorf("id type mismatch: expected %s but got %s", UserPrefix, public.prefix)
	}

	// Validation: Ensure the ID sent matches the expected type
	id.PublicID.prefix = public.prefix
	id.PublicID.uuid = public.uuid
	return nil
}

func (id *UserID) String() string {
	public := PublicID{prefix: UserPrefix, uuid: id.PublicID.uuid}
	return public.String()
}

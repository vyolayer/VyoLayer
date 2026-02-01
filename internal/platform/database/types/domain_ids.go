package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UserID is the type for user IDs.
type UserID struct {
	PublicID
}

func NewUserID() UserID {
	return UserID{PublicID: NewPublicID(UserPrefix)}
}
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
			Prefix: UserPrefix,
			UUID:   ReconstructInternalID(id),
		},
	}, nil
}

func (id UserID) MarshalJSON() ([]byte, error) {
	public := PublicID{Prefix: UserPrefix, UUID: id.PublicID.UUID}
	return public.MarshalJSON()
}

func (id *UserID) UnmarshalJSON(b []byte) error {
	var public PublicID
	if err := public.UnmarshalJSON(b); err != nil {
		return err
	}

	// Validation: Ensure the ID sent matches the expected type
	if public.Prefix != UserPrefix {
		return fmt.Errorf("id type mismatch: expected %s but got %s", UserPrefix, public.Prefix)
	}

	id.PublicID.UUID = public.UUID
	return nil
}

func (id *UserID) String() string {
	public := PublicID{Prefix: UserPrefix, UUID: id.PublicID.UUID}
	return public.String()
}

package types

import (
	"fmt"
)

// User id
type userID struct {
	PublicID[IDPrefix]
}

type UserID interface {
	VyoLayerPublicID
}

// NewUserID creates a new UserID.
func NewUserID() UserID {
	return &userID{PublicID: NewPublicID(UserPrefix)}
}

// ReconstructUserID reconstructs a UserID from a string.
func ReconstructUserID(s string) (UserID, error) {
	public, err := ReconstructPublicID(UserPrefix, s)
	if err != nil {
		return nil, err
	}
	return &userID{PublicID: public}, nil
}

func ParseUserID(s string) (UserID, error) {
	public, err := ParsePublicID(s)
	if err != nil {
		return nil, err
	}
	if public.prefix != UserPrefix {
		return nil, fmt.Errorf("invalid prefix: %s", public.prefix)
	}
	return &userID{PublicID: public}, nil
}

func (id *userID) MarshalJSON() ([]byte, error) {
	public := PublicID[IDPrefix]{prefix: UserPrefix, uuid: id.PublicID.uuid}
	return public.MarshalJSON()
}

func (id *userID) UnmarshalJSON(b []byte) error {
	var public PublicID[IDPrefix]
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

func (id *userID) String() string {
	public := PublicID[IDPrefix]{prefix: UserPrefix, uuid: id.PublicID.uuid}
	return public.String()
}

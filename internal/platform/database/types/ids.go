package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type IDPrefix string

const (
	PrefixSeparator = "_"

	UserPrefix              IDPrefix = "user"
	OrgPrefix               IDPrefix = "org"
	OrgRolePrefix           IDPrefix = "org_role"
	OrgPermissionPrefix     IDPrefix = "org_permission"
	OrgMemberPrefix         IDPrefix = "org_member"
	ProjectPrefix           IDPrefix = "project"
	ProjectMemberPrefix     IDPrefix = "project_member"
	ProjectRolePrefix       IDPrefix = "project_role"
	ProjectPermissionPrefix IDPrefix = "project_permission"
)

var validPrefixes = map[IDPrefix]bool{
	UserPrefix:              true,
	OrgPrefix:               true,
	OrgRolePrefix:           true,
	OrgPermissionPrefix:     true,
	OrgMemberPrefix:         true,
	ProjectPrefix:           true,
	ProjectMemberPrefix:     true,
	ProjectRolePrefix:       true,
	ProjectPermissionPrefix: true,
}

func (p IDPrefix) IsValid() bool {
	return validPrefixes[p]
}

func (p IDPrefix) String() string {
	return string(p)
}

// Internal id for database
type InternalID uuid.UUID

func NewInternalID() InternalID {
	return InternalID(uuid.New())
}

func ReconstructInternalID(uuid uuid.UUID) InternalID {
	return InternalID(uuid)
}

func ParseInternalID(s string) (InternalID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return InternalID{}, err
	}
	return InternalID(id), nil
}

func (id InternalID) ID() uuid.UUID {
	return uuid.UUID(id)
}

func (id InternalID) String() string {
	return uuid.UUID(id).String()
}

func (id InternalID) IsNil() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id InternalID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *InternalID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}

	*id = InternalID(parsed)
	return nil
}

func (id InternalID) Value() (driver.Value, error) {
	return id.String(), nil
}

func (id *InternalID) Scan(value any) error {
	if value != nil {
		*id = InternalID(uuid.Nil)
		return nil
	}
	switch v := value.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = InternalID(parsed)

	case []byte:
		parsed, err := uuid.ParseBytes(v)
		if err != nil {
			return err
		}
		*id = InternalID(parsed)

	default:
		return fmt.Errorf("cannot scan %T into InternalID", value)
	}
	return nil
}

// Public id for API/external use
type PublicID struct {
	Prefix IDPrefix
	UUID   InternalID
}

func NewPublicID(prefix IDPrefix) PublicID {
	return PublicID{
		Prefix: prefix,
		UUID:   NewInternalID(),
	}
}

func NewPublicIDFromInternalID(prefix IDPrefix, id InternalID) PublicID {
	return PublicID{
		Prefix: prefix,
		UUID:   id,
	}
}

func ParsePublicID(s string) (PublicID, error) {
	parts := strings.SplitN(s, PrefixSeparator, 2)
	if len(parts) != 2 {
		return PublicID{}, fmt.Errorf("invalid public id format: %s", s)
	}

	prefix := IDPrefix(parts[0])
	if !prefix.IsValid() {
		return PublicID{}, fmt.Errorf("invalid prefix: %s", prefix)
	}

	uuid, err := ParseInternalID(parts[1])
	if err != nil {
		return PublicID{}, err
	}

	return PublicID{
		Prefix: prefix,
		UUID:   uuid,
	}, nil
}

func (id PublicID) String() string {
	return string(id.Prefix) + PrefixSeparator + id.UUID.String()
}

func (id PublicID) IsNil() bool {
	return id.Prefix == "" && id.UUID.IsNil()
}

func (id PublicID) InternalID() InternalID {
	return id.UUID
}

func (id PublicID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PublicID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed, err := ParsePublicID(s)
	if err != nil {
		return err
	}

	*id = parsed
	return nil
}

// Creates a PublicID from a raw uuid.UUID with the given prefix
func NewPublicIDFromUUID(prefix IDPrefix, id uuid.UUID) (PublicID, error) {
	if id == uuid.Nil {
		return PublicID{}, fmt.Errorf("invalid uuid: %s", id)
	}
	return PublicID{
		Prefix: prefix,
		UUID:   InternalID(id),
	}, nil
}

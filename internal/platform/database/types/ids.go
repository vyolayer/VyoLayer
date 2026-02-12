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
	OrgMemberRolePrefix     IDPrefix = "org_member_role"
	OrgInvitationPrefix     IDPrefix = "org_invitation"
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
	OrgInvitationPrefix:     true,
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

// Interface for all InternalID
type WorkLayerInternalID interface {
	ID() uuid.UUID
	String() string
	IsNil() bool
}

// Interface for all the ID types
type WorkLayerPublicID interface {
	String() string
	IsNil() bool
	InternalID() WorkLayerInternalID
}

// Internal id for database
type InternalID uuid.UUID

func NewInternalID() WorkLayerInternalID {
	return InternalID(uuid.New())
}

func ReconstructInternalID(uuid uuid.UUID) WorkLayerInternalID {
	return InternalID(uuid)
}

func ParseInternalID(s string) (WorkLayerInternalID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
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
type PublicID[T IDPrefix] struct {
	prefix T
	uuid   WorkLayerInternalID
}

func NewPublicID[T IDPrefix](prefix T) PublicID[T] {
	return PublicID[T]{
		prefix: prefix,
		uuid:   NewInternalID(),
	}
}

func NewPublicIDFromInternalID[T IDPrefix](prefix T, id WorkLayerInternalID) PublicID[T] {
	return PublicID[T]{
		prefix: prefix,
		uuid:   id,
	}
}

func ParsePublicID[T IDPrefix](s string) (PublicID[T], error) {
	parts := strings.SplitN(s, PrefixSeparator, 2)
	if len(parts) != 2 {
		return PublicID[T]{}, fmt.Errorf("invalid public id format: %s", s)
	}

	prefix := IDPrefix(parts[0])
	if !prefix.IsValid() {
		return PublicID[T]{}, fmt.Errorf("invalid prefix: %s", prefix)
	}

	uuid, err := ParseInternalID(parts[1])
	if err != nil {
		return PublicID[T]{}, err
	}

	return PublicID[T]{
		prefix: T(prefix),
		uuid:   uuid,
	}, nil
}

// ReconstructPublicID reconstructs a PublicID from a string.
// This handles both "prefix_uuid" and "uuid" formats.
// It validates that if a prefix exists in the string, it matches the expected prefix.
func ReconstructPublicID[T IDPrefix](prefix T, s string) (PublicID[T], error) {
	if s == "" {
		return PublicID[T]{}, fmt.Errorf("id cannot be empty")
	}

	// Construct the expected prefix string
	expectedPrefix := string(prefix) + PrefixSeparator

	// Safety Check: If it has a prefix, but it's the WRONG prefix (e.g. "org_..." when expecting "user_")
	if strings.Contains(s, PrefixSeparator) && !strings.HasPrefix(s, expectedPrefix) {
		return PublicID[T]{}, fmt.Errorf("id type mismatch: expected prefix %s", prefix)
	}

	// Remove the prefix if it exists
	cleanUUID := strings.TrimPrefix(s, expectedPrefix)

	// Parse the UUID
	parsedUUID, err := uuid.Parse(cleanUUID)
	if err != nil {
		return PublicID[T]{}, fmt.Errorf("invalid uuid format: %w", err)
	}

	// Construct and return
	return PublicID[T]{
		prefix: prefix,
		uuid:   ReconstructInternalID(parsedUUID),
	}, nil
}

func (id PublicID[T]) String() string {
	return string(id.prefix) + PrefixSeparator + id.uuid.String()
}

func (id PublicID[T]) IsNil() bool {
	return id.prefix == "" && id.uuid.IsNil()
}

func (id PublicID[T]) InternalID() WorkLayerInternalID {
	return id.uuid
}

func (id PublicID[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *PublicID[T]) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsed, err := ParsePublicID[T](s)
	if err != nil {
		return err
	}

	*id = parsed
	return nil
}

// Creates a PublicID from a raw uuid.UUID with the given prefix
func NewPublicIDFromUUID[T IDPrefix](prefix T, id uuid.UUID) (PublicID[T], error) {
	if id == uuid.Nil {
		return PublicID[T]{}, fmt.Errorf("invalid uuid: %s", id)
	}
	return PublicID[T]{
		prefix: prefix,
		uuid:   ReconstructInternalID(id),
	}, nil
}

// -- Organization ID --
// OrgPrefix is the prefix for organization IDs
type organizationID struct {
	PublicID[IDPrefix]
}

type OrganizationID interface {
	WorkLayerPublicID
}

// NewOrganizationID creates a new organization ID
func NewOrganizationID() OrganizationID {
	return organizationID{PublicID: NewPublicID(OrgPrefix)}
}

// ReconstructOrganizationID reconstructs an organization ID from a string
func ReconstructOrganizationID(s string) (OrganizationID, error) {
	public, err := ReconstructPublicID(OrgPrefix, s)
	if err != nil {
		return nil, err
	}
	return organizationID{PublicID: public}, nil
}

// --- Organization Member ID ---
type organizationMemberID struct {
	PublicID[IDPrefix]
}

type OrganizationMemberID interface {
	WorkLayerPublicID
}

func NewOrganizationMemberID() OrganizationMemberID {
	return organizationMemberID{PublicID: NewPublicID(OrgMemberPrefix)}
}

func ReconstructOrganizationMemberID(s string) (OrganizationMemberID, error) {
	public, err := ReconstructPublicID(OrgMemberPrefix, s)
	if err != nil {
		return nil, err
	}
	return organizationMemberID{PublicID: public}, nil
}

// --- Organization Role ID ---
type organizationRoleID struct {
	PublicID[IDPrefix]
}

type OrganizationRoleID interface {
	WorkLayerPublicID
}

func NewOrganizationRoleID() OrganizationRoleID {
	return organizationRoleID{PublicID: NewPublicID(OrgRolePrefix)}
}

func ReconstructOrganizationRoleID(s string) (OrganizationRoleID, error) {
	public, err := ReconstructPublicID(OrgRolePrefix, s)
	if err != nil {
		return nil, err
	}
	return organizationRoleID{PublicID: public}, nil
}

// --- Organization Permission ID ---
type organizationPermissionID struct {
	PublicID[IDPrefix]
}

type OrganizationPermissionID interface {
	WorkLayerPublicID
}

func NewOrganizationPermissionID() OrganizationPermissionID {
	return organizationPermissionID{PublicID: NewPublicID(OrgPermissionPrefix)}
}

func ReconstructOrganizationPermissionID(s string) (OrganizationPermissionID, error) {
	public, err := ReconstructPublicID(OrgPermissionPrefix, s)
	if err != nil {
		return nil, err
	}
	return organizationPermissionID{PublicID: public}, nil
}

// --- Member Organization Role ID ---
type memberOrganizationRoleID struct {
	PublicID[IDPrefix]
}

type MemberOrganizationRoleID interface {
	WorkLayerPublicID
}

func NewMemberOrganizationRoleID() MemberOrganizationRoleID {
	return memberOrganizationRoleID{PublicID: NewPublicID(OrgMemberRolePrefix)}
}

func ReconstructMemberOrganizationRoleID(s string) (MemberOrganizationRoleID, error) {
	public, err := ReconstructPublicID(OrgMemberRolePrefix, s)
	if err != nil {
		return nil, err
	}
	return memberOrganizationRoleID{PublicID: public}, nil
}

// --- Organization Member Invitation ID ---
type organizationMemberInvitationID struct {
	PublicID[IDPrefix]
}

type OrganizationMemberInvitationID interface {
	WorkLayerPublicID
}

func NewOrganizationMemberInvitationID() OrganizationMemberInvitationID {
	return organizationMemberInvitationID{PublicID: NewPublicID(OrgInvitationPrefix)}
}

func ReconstructOrganizationMemberInvitationID(s string) (OrganizationMemberInvitationID, error) {
	public, err := ReconstructPublicID(OrgInvitationPrefix, s)
	if err != nil {
		return nil, err
	}
	return organizationMemberInvitationID{PublicID: public}, nil
}

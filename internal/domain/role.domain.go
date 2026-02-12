package domain

import "encoding/json"

// MarshalRoleIDs converts role IDs to JSON string for database storage
func MarshalRoleIDs(roleIDs []string) (string, error) {
	if len(roleIDs) == 0 {
		return "[]", nil
	}
	bytes, err := json.Marshal(roleIDs)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// UnmarshalRoleIDs converts JSON string from database to role IDs slice
func UnmarshalRoleIDs(roleIDsJSON string) ([]string, error) {
	if roleIDsJSON == "" {
		return []string{}, nil
	}
	var roleIDs []string
	err := json.Unmarshal([]byte(roleIDsJSON), &roleIDs)
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

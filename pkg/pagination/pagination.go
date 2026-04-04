package pagination

import (
	"encoding/base64"
	"strconv"
)

// decodePageToken converts a base64 opaque token back into an integer offset
func DecodePageToken(token string) int {
	if token == "" {
		return 0
	}
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0 // Default to 0 if token is invalid
	}
	offset, err := strconv.Atoi(string(decoded))
	if err != nil {
		return 0
	}
	return offset
}

// encodePageToken converts an integer offset into a base64 opaque token
func EncodePageToken(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(offset)))
}

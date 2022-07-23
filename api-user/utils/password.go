package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// EncryptPassword encrypts the client password on server side
func EncryptPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

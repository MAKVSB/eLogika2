package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func HashJSON(raw json.RawMessage) (string, error) {
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}

	// stringify like JSON.stringify in JS
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

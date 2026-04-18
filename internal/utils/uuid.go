package utils

import (
	"encoding/hex"
	"errors"
	"strings"
)

func UUIDStringToBytes(s string) ([16]byte, error) {
	var b [16]byte

	s = strings.ReplaceAll(s, "-", "")

	if len(s) != 32 {
		return b, errors.New("invalid length")
	}

	_, err := hex.Decode(b[:], []byte(s))
	return b, err
}

func BytesToUUIDString(b [16]byte) string {
	return hex.EncodeToString(b[0:4]) + "-" +
		hex.EncodeToString(b[4:6]) + "-" +
		hex.EncodeToString(b[6:8]) + "-" +
		hex.EncodeToString(b[8:10]) + "-" +
		hex.EncodeToString(b[10:16])
}

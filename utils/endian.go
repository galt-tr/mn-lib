package utils

import (
	"encoding/hex"
)

func BigToLittleEndian(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	ret := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		ret = append(ret, bytes[i])
	}
	return hex.EncodeToString(ret), nil
}

package core

import "encoding/base64"

// Base64Encode encodes string to base64
func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Base64Decode decodes base64 string to original form
func Base64Decode(data string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// Base64EncodeByte encodes bytes to base64
func Base64EncodeByte(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64DecodeByte decodes base64 string to bytes
func Base64DecodeByte(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
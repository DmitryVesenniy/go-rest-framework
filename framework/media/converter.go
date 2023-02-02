package media

import (
	"encoding/base64"
	"strings"
)

func ConvertBase64StringToByte(b64 string) ([]byte, string, error) {
	_base64, prefix := clearBase64Str(b64)
	b, err := base64.StdEncoding.DecodeString(_base64)
	if err != nil {
		return nil, prefix, err
	}
	return b, prefix, nil
}

func ConvertByteToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func clearBase64Str(txt string) (string, string) {
	b64Str := strings.Split(txt, ",")
	if len(b64Str) > 1 {
		return strings.TrimSpace(b64Str[1]), strings.TrimSpace(b64Str[0])
	}
	return strings.TrimSpace(b64Str[0]), ""
}

package crypto

import (
	b "encoding/base64"
)

type base64 struct{}

func init() {
	Register("base64", newBase64)
}

func newBase64() Crypto {
	return &base64{}
}

func (im *base64) Encrypt(text string) (string, error) {
	encoded := b.StdEncoding.EncodeToString([]byte(text))
	return encoded, nil
}

func (im *base64) Decrypt(encryptText string) (string, error) {
	decoded, err := b.StdEncoding.DecodeString(encryptText)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

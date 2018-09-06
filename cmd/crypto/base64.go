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

func (im *base64) Encrypt(text []byte) ([]byte, error) {
	encoded := b.StdEncoding.EncodeToString(text)
	return []byte(encoded), nil
}

func (im *base64) Decrypt(encryptText []byte) ([]byte, error) {
	decoded, err := b.StdEncoding.DecodeString(string(encryptText))
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

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
	encoded := make([]byte, b.StdEncoding.EncodedLen(len(text)))
	b.StdEncoding.Encode(encoded, text)

	return encoded, nil
}

func (im *base64) Decrypt(encryptText []byte) ([]byte, error) {
	decoded := make([]byte, b.StdEncoding.DecodedLen(len(encryptText)))
	_, err := b.StdEncoding.Decode(decoded, encryptText)
	if err != nil {
		return nil, err
	}
	// To remove base64 padding 0
	for decoded[len(decoded)-1] == 0 {
		decoded = decoded[:len(decoded)-1]
	}
	return decoded, nil
}

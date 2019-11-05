package crypto

import (
	b "encoding/base64"
	"log"

	a "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awskms "github.com/aws/aws-sdk-go/service/kms"
	"github.com/spf13/viper"
)

type aws struct {
	client *session.Session
}

func init() {
	Register("aws", newAWS)
}

func newAWS() Crypto {
	return &aws{client: newSession()}
}

func newSession() *session.Session {
	region := viper.GetString("AWSlocationID")
	sess, err := session.NewSession(&a.Config{
		Region: a.String(region),
	})

	if err != nil {
		log.Fatal(err)
	}

	return sess
}

func (im *aws) Encrypt(input []byte) ([]byte, error) {
	keyID := viper.GetString("AWScryptoKeyID")
	client := awskms.New(im.client)
	params := &awskms.EncryptInput{
		KeyId:     a.String(keyID),
		Plaintext: input,
	}

	req, resp := client.EncryptRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	encodedPlaintext := make([]byte, b.StdEncoding.EncodedLen(len(resp.CiphertextBlob)))
	b.StdEncoding.Encode(encodedPlaintext, resp.CiphertextBlob)

	return encodedPlaintext, nil
}

func (im *aws) Decrypt(input []byte) ([]byte, error) {
	text, err := b.StdEncoding.DecodeString(string(input))
	if err != nil {
		return nil, err
	}

	client := awskms.New(im.client)
	params := &awskms.DecryptInput{
		CiphertextBlob: []byte(text),
	}

	req, resp := client.DecryptRequest(params)

	err = req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Plaintext, nil
}

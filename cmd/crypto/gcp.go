package crypto

import (
	"context"
	b "encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type gcp struct {
	client *http.Client
}

func init() {
	Register("gcp", newgcp)
}

func newgcp() Crypto {
	return &gcp{client: newAuthenticatedClient()}
}

func (im *gcp) Encrypt(plaintext string) (string, error) {
	cloudkmsService, err := cloudkms.New(im.client)
	if err != nil {
		return "", err
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		viper.GetString("GCPprojectID"),
		viper.GetString("GCPlocationID"),
		viper.GetString("GCPkeyRingID"),
		viper.GetString("GCPcryptoKeyID"),
	)

	req := &cloudkms.EncryptRequest{
		Plaintext: b.StdEncoding.EncodeToString([]byte(plaintext)),
	}

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(parentName, req).Do()
	if err != nil {
		return "", err
	}

	return resp.Ciphertext, err
}

func (im *gcp) Decrypt(ciphertext string) (string, error) {
	cloudkmsService, err := cloudkms.New(im.client)
	if err != nil {
		return "", err
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		viper.GetString("GCPprojectID"),
		viper.GetString("GCPlocationID"),
		viper.GetString("GCPkeyRingID"),
		viper.GetString("GCPcryptoKeyID"),
	)

	if err != nil {
		log.Fatal(err)
	}
	req := &cloudkms.DecryptRequest{
		Ciphertext: ciphertext,
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parentName, req).Do()
	if err != nil {
		return "", err
	}
	text, err := b.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func newAuthenticatedClient() *http.Client {
	var client *http.Client
	oAuthLocation := viper.GetString("oAuthLocation")
	if len(oAuthLocation) > 0 {
		data, err := ioutil.ReadFile(oAuthLocation)

		if err != nil {
			log.Fatal("unable to read JSON key file", err)
		}
		conf, err := google.JWTConfigFromJSON(data, cloudkms.CloudPlatformScope)
		if err != nil {
			log.Fatal("unable to parse JSON key file", err)
		}
		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client = conf.Client(oauth2.NoContext)
	} else {
		ctx := context.Background()
		defaultClient, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
		if err != nil {
			log.Fatal(err)
		}
		client = defaultClient
	}
	return client
}

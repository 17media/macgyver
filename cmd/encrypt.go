package cmd

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/17media/macgyver/cmd/gcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt entire flags",
	Run:   encrypt,
	Args:  cobra.NoArgs,
}

var reEncryptFlag = regexp.MustCompile(encryptFlagRegexp)

const (
	encryptFlagRegexp = `^\-(\w*)=(.*)$`
)

func init() {
	encryptCmd.MarkFlagRequired("flags")
	encryptCmd.MarkFlagRequired("GCPprojectID")
	encryptCmd.MarkFlagRequired("GCPlocationID")
	encryptCmd.MarkFlagRequired("GCPkeyRingID")
	encryptCmd.MarkFlagRequired("GCPcryptoKeyID")
	RootCmd.AddCommand(encryptCmd)
}

func encrypt(cmd *cobra.Command, args []string) {
	var originalFlags []*env
	var client *http.Client
	client = gcp.NewAuthenticatedClient()
	splitFlags := strings.Split(viper.GetString("flags"), " ")

	for _, value := range splitFlags {
		encryptText := ""
		match := reEncryptFlag.FindStringSubmatch(value)

		flag := &env{
			key:   match[1],
			value: match[2],
		}
		encryptText, err := gcp.Encrypt(flag.value, client)
		if err != nil {
			log.Fatal(err)
		}
		flag.value = "<kms>" + encryptText
		originalFlags = append(originalFlags, flag)
	}

	// Convert encrypted flags back to string
	encryptedFlags := covertFlags(originalFlags)

	fmt.Println(encryptedFlags)
}

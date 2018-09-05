package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/17media/macgyver/cmd/gcp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt entire flags",
	Run:   decrypt,
	Args:  cobra.NoArgs,
}

var reDecryptFlag = regexp.MustCompile(decryptFlagRegexp)

const (
	decryptFlagRegexp = `^\-(\w*)=((?:<kms>))?(.*)$`
)

func init() {
	decryptCmd.MarkFlagRequired("flags")
	decryptCmd.MarkFlagRequired("GCPprojectID")
	decryptCmd.MarkFlagRequired("GCPlocationID")
	decryptCmd.MarkFlagRequired("GCPkeyRingID")
	decryptCmd.MarkFlagRequired("GCPcryptoKeyID")
	RootCmd.AddCommand(decryptCmd)
}

func decrypt(cmd *cobra.Command, args []string) {
	var originalFlags []*env
	client := gcp.NewAuthenticatedClient()
	splitFlags := strings.Split(viper.GetString("flags"), " ")
	for _, value := range splitFlags {
		match := reDecryptFlag.FindStringSubmatch(value)
		flag := &env{
			key:   match[1],
			value: match[3],
		}

		// if it needs to be decrypted
		if match[2] == "<kms>" {
			decryptText, err := gcp.Decrypt(flag.value, client)
			if err != nil {
				log.Fatal(err)
			}
			flag.value = string(decryptText)
		}
		originalFlags = append(originalFlags, flag)
	}

	// Convert decrypted flags back to string
	decryptedFlags := covertFlags(originalFlags)
	fmt.Println(decryptedFlags)
}

func covertFlags(decrypt []*env) string {
	var result string
	for _, flag := range decrypt {
		result += fmt.Sprintf(" -%s=%s", flag.key, flag.value)
	}
	return result
}

package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/17media/macgyver/cmd/crypto"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt entire flags",
	Run:   decrypt,
	Args:  cobra.NoArgs,
}

func init() {
	decryptCmd.MarkFlagRequired("flags")
	decryptCmd.MarkFlagRequired("GCPprojectID")
	decryptCmd.MarkFlagRequired("GCPlocationID")
	decryptCmd.MarkFlagRequired("GCPkeyRingID")
	decryptCmd.MarkFlagRequired("GCPcryptoKeyID")
	RootCmd.AddCommand(decryptCmd)
}

func decrypt(cmd *cobra.Command, args []string) {
	crypto.Init(cryptoProvide)
	var originalFlags []*env
	splitFlags := strings.Split(flags, " ")
	p := crypto.Providers[cryptoProvide]

	decryptFlagRegexp := `^\-(\w*)=((?:` + Perfix + `))?(.*)$`
	var reDecryptFlag = regexp.MustCompile(decryptFlagRegexp)

	for _, value := range splitFlags {
		match := reDecryptFlag.FindStringSubmatch(value)
		flag := &env{
			key:   match[1],
			value: match[3],
		}

		// if it needs to be decrypted
		if match[2] == Perfix {
			decryptText, err := p.Decrypt(flag.value)
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
	return strings.TrimLeft(result, " ")
}

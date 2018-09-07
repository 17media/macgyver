package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/17media/macgyver/cmd/crypto"
	"github.com/spf13/cobra"
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
	encryptCmd.MarkFlagRequired("cryptoProvider")
	encryptCmd.MarkFlagRequired("GCPprojectID")
	encryptCmd.MarkFlagRequired("GCPlocationID")
	encryptCmd.MarkFlagRequired("GCPkeyRingID")
	encryptCmd.MarkFlagRequired("GCPcryptoKeyID")

	RootCmd.AddCommand(encryptCmd)
}

func encrypt(cmd *cobra.Command, args []string) {
	crypto.Init(cryptoProvider)
	var originalFlags []*env
	splitFlags := strings.Split(flags, " ")
	p := crypto.Providers[cryptoProvider]

	for _, value := range splitFlags {
		var encryptText []byte
		match := reEncryptFlag.FindStringSubmatch(value)

		flag := &env{
			key:   match[1],
			value: match[2],
		}
		encryptText, err := p.Encrypt([]byte(flag.value))
		if err != nil {
			log.Fatal(err)
		}
		flag.value = Perfix + string(encryptText)
		originalFlags = append(originalFlags, flag)
	}

	// Convert encrypted flags back to string
	encryptedFlags := covertFlags(originalFlags)

	fmt.Println(strings.TrimLeft(encryptedFlags, " "))
}

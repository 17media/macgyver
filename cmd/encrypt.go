package cmd

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/17media/macgyver/cmd/crypto"
	"github.com/17media/macgyver/cmd/keys"
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
	p := crypto.Providers[cryptoProvider]

	k, ok := keys.Types[cryptoType]
	if !ok {
		panic("Without support " + cryptoType + " encrypt")
	}

	keyFlags, err := k.Import(flags, Perfix)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range keyFlags {
		encryptText, err := p.Encrypt([]byte(v.Value))
		if err != nil {
			log.Fatal(err)
		}
		keyFlags[i].Value = string(encryptText)
	}

	// Convert encrypted flags back to string
	encryptedFlags := covertFlags(keyFlags)

	fmt.Println(strings.TrimLeft(encryptedFlags, " "))
}

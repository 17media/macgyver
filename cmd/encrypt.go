package cmd

import (
	"log"
	"regexp"

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
	var keyFlags []keys.Key
	crypto.Init(cryptoProvider)
	p := crypto.Providers[cryptoProvider]

	if cryptoType == CryptoTypeName[0] {
		keyFlags = keys.FlagsImporter(flags, Perfix)
	} else {
		panic("Without support " + cryptoType + " cryptoType")
	}

	for i, v := range keyFlags {
		encryptText, err := p.Encrypt([]byte(v.Value))
		if err != nil {
			log.Fatal(err)
		}
		keyFlags[i].Value = string(encryptText)
	}

	// Convert encrypted flags back to string
	keys.FlagsExporter(keyFlags)
}

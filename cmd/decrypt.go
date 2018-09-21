package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/17media/macgyver/cmd/crypto"
	"github.com/17media/macgyver/cmd/keys"
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
	encryptCmd.MarkFlagRequired("cryptoProvider")
	decryptCmd.MarkFlagRequired("GCPprojectID")
	decryptCmd.MarkFlagRequired("GCPlocationID")
	decryptCmd.MarkFlagRequired("GCPkeyRingID")
	decryptCmd.MarkFlagRequired("GCPcryptoKeyID")
	decryptCmd.MarkFlagRequired("cryptoType")
	RootCmd.AddCommand(decryptCmd)
}

func decrypt(cmd *cobra.Command, args []string) {
	var keyFlags []keys.Key
	crypto.Init(cryptoProvider)
	p := crypto.Providers[cryptoProvider]

	if cryptoType == CryptoTypeName[0] {
		keyFlags = keys.FlagsImporter(flags, Perfix)
	} else if cryptoType == CryptoTypeName[1] {
		keyFlags = keys.EnvsImporter(Perfix)
	} else {
		panic("Without support " + cryptoType + " cryptoType")
	}

	for i, v := range keyFlags {
		if v.IsEncrypted {
			decryptText, err := p.Decrypt([]byte(v.Value))
			if err != nil {
				log.Fatal(err)
			}
			keyFlags[i].Value = string(decryptText)
		}
	}

	if cryptoType == CryptoTypeName[0] {
		// Convert decrypted flags back to string
		decryptedFlags := covertFlags(keyFlags)
		fmt.Println(decryptedFlags)
	} else if cryptoType == CryptoTypeName[1] {
		keys.EnvsOutputer(keyFlags)
	}
}

func covertFlags(decrypt []keys.Key) string {
	var result string
	for _, flag := range decrypt {
		result += fmt.Sprintf(" -%s=%s", flag.Key, flag.Value)
	}
	return strings.TrimLeft(result, " ")
}

package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/17media/macgyver/cmd/crypto"
	"github.com/17media/macgyver/cmd/keys"
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
	decryptCmd.MarkFlagRequired("keysType")
	RootCmd.AddCommand(decryptCmd)
}

func decrypt(cmd *cobra.Command, args []string) {
	crypto.Init(cryptoProvider)
	inputs := map[keys.Type][]string{
		keys.TypeText: strings.Split(flags, " "),
		keys.TypeEnv:  os.Environ(),
	}
	k, err := keys.Get(keysType)
	if err != nil {
		log.Panic(err)
	}
	keyFlags := k.Import(inputs[keysType], SecretTag)

	// Decrype all secrets that are encrypted of each key
	p := crypto.Providers[cryptoProvider]
	for _, k := range keyFlags {
		for _, s := range k.Secrets {
			if !s.IsEncrypted {
				continue
			}
			decryptText, err := p.Decrypt([]byte(s.Text))
			if err != nil {
				log.Panic(err)
			}
			s.Text = string(decryptText)
			s.IsEncrypted = false
		}
	}

	// Convert decrypted keys back to string
	k.Export(keyFlags, SecretTag, os.Stdout)
}

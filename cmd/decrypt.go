package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	decryptCmd.MarkFlagRequired("cryptoProvider")
	decryptCmd.MarkFlagRequired("keysType")

	switch viper.GetString("cryptoProvider") {
	case "gcp":
		decryptCmd.MarkFlagRequired("GCPprojectID")
		decryptCmd.MarkFlagRequired("GCPlocationID")
		decryptCmd.MarkFlagRequired("GCPkeyRingID")
		decryptCmd.MarkFlagRequired("GCPcryptoKeyID")
	case "aws":
		decryptCmd.MarkFlagRequired("AWSlocationID")
		decryptCmd.MarkFlagRequired("AWScryptoKeyID")
	default:

	}

	RootCmd.AddCommand(decryptCmd)
}

func decrypt(cmd *cobra.Command, args []string) {
	crypto.Init(cryptoProvider)
	k, err := keys.Get(keysType)
	if err != nil {
		log.Panic(err)
	}

	// use file, flags or nev
	switch keysType {
	case keys.TypeText:
		keyFlags := k.Import(strings.Split(flags, " "), SecretTag)
		// Decrypt all secrets that are encrypted of each key
		p := crypto.Providers[cryptoProvider]
		for _, keyFlag := range keyFlags {
			for _, s := range keyFlag.Secrets {
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
	case keys.TypeFile:
		log.Panic("file type is not ready")
	case keys.TypeEnv:
		log.Panic("env type is not ready")
	default:
		log.Panicf("keysType does not support %s", keysType)
	}
}

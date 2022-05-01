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

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt entire flags",
	Run:   encrypt,
	Args:  cobra.NoArgs,
}

const (
	encryptFlagRegexp = `^\-(\w+)=(.+)$`
)

func init() {
	encryptCmd.MarkFlagRequired("flags")
	encryptCmd.MarkFlagRequired("cryptoProvider")

	switch viper.GetString("cryptoProvider") {
	case "gcp":
		encryptCmd.MarkFlagRequired("GCPprojectID")
		encryptCmd.MarkFlagRequired("GCPlocationID")
		encryptCmd.MarkFlagRequired("GCPkeyRingID")
		encryptCmd.MarkFlagRequired("GCPcryptoKeyID")
	case "aws":
		encryptCmd.MarkFlagRequired("AWSlocationID")
		encryptCmd.MarkFlagRequired("AWScryptoKeyID")
	default:

	}

	RootCmd.AddCommand(encryptCmd)
}

func encrypt(cmd *cobra.Command, args []string) {
	crypto.Init(cryptoProvider)
	inputs := map[keys.Type][]string{
		keys.TypeText: strings.Split(flags, " "),
		keys.TypeEnv:  os.Environ(),
	}
	k, err := keys.Get(keysType)
	if err != nil {
		log.Panic(err)
	}

	// use file, flags oe env
	switch keysType {
	case keys.TypeText, keys.TypeEnv:
		keyFlags := k.Import(inputs[keysType], SecretTag)

		// Encrypt all secrets that are not encrypted of each key
		p := crypto.Providers[cryptoProvider]
		for _, keyFlag := range keyFlags {
			for _, s := range keyFlag.Secrets {
				if s.IsEncrypted {
					continue
				}
				encryptText, err := p.Encrypt([]byte(s.Text))
				if err != nil {
					log.Panic(err)
				}
				s.Text = string(encryptText)
				s.IsEncrypted = true
			}
		}
		// Convert decrypted keys back to string
		k.Export(keyFlags, SecretTag, os.Stdout)
	case keys.TypeFile:
		p := crypto.Providers[cryptoProvider]
		values := k.Encrypt(file, SecretTag, p)
		if err := k.ReplaceOriginFile(file, values); err != nil {
			log.Panic(err)
		}
	default:
		log.Panicf("keysType does not support %s", keysType)
	}
}

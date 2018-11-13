package cmd

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/17media/macgyver/cmd/crypto"
	"github.com/17media/macgyver/cmd/keys"
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
	inputs := map[keys.Type][]string{
		keys.TypeText: strings.Split(flags, " "),
		keys.TypeEnv:  os.Environ(),
	}

	k, err := keys.Get(keysType)
	if err != nil {
		log.Panic(err)
	}
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

	// Convert encrypted keys back to string
	k.Export(keyFlags, SecretTag, os.Stdout)
}

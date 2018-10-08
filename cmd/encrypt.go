package cmd

import (
	"log"
	"os"
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
	inputs := map[keys.Type][]string{
		keys.TypeText: strings.Split(flags, " "),
		keys.TypeEnv:  os.Environ(),
	}

	k, err := keys.Get(keysType)
	if err != nil {
		panic(err)
	}
	keyFlags := k.Import(inputs[keysType], Prefix)

	p := crypto.Providers[cryptoProvider]
	for i, v := range keyFlags {
		encryptText := []byte(v.Value)
		if !v.IsEncrypted {
			encryptText, err = p.Encrypt([]byte(v.Value))
			if err != nil {
				log.Fatal(err)
			}
		}
		keyFlags[i].Value = Prefix + string(encryptText)
	}

	// Convert encrypted flags back to string
	k.Export(keyFlags)
}

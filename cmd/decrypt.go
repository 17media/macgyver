package cmd

import (
	"log"
	"os"
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
		panic(err)
	}
	keyFlags := k.Import(inputs[keysType], Prefix)

	p := crypto.Providers[cryptoProvider]
	for i, v := range keyFlags {
		if v.IsEncrypted {
			decryptText, err := p.Decrypt([]byte(v.Value))
			if err != nil {
				log.Fatal(err)
			}
			keyFlags[i].Value = string(decryptText)
		}
	}

	k.Export(keyFlags)
}

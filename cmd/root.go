package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/17media/macgyver/cmd/keys"
)

var (
	cfgFile        string
	cryptoProvider string
	oAuthLocation  string
	flags          string
	keysType       keys.Type
	SecretTag      string
	GCPprojectID   string
	GCPlocationID  string
	GCPkeyRingID   string
	GCPcryptoKeyID string
	AWSlocationID  string
	AWScryptoKeyID string
	AWSprofileName string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "macgyver",
	Short: "A tool of decrypt and encrypt in GCP KMS",
	Long: `A tool for decrypting and encrypting strings in GCP / AWS by using KMS,
The tool is golang's flags friendly.
For example:
$ go run main.go decrypt                 \
                --cryptoProvider=gcp     \
                --GCPprojectID="demo"    \
                --GCPlocationID="global" \
                --GCPkeyRingID="foo"     \
                --GCPcryptoKeyID="bar"   \
                --flags="-a=kms_asda`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.macgyver.yaml)")

	// var cryptoProvider string
	RootCmd.PersistentFlags().StringVar(&cryptoProvider, "cryptoProvider", "", "Which type you using encrypto and encryto")
	viper.BindPFlag("cryptoProvider", RootCmd.PersistentFlags().Lookup("cryptoProvider"))

	// var keysType string
	RootCmd.PersistentFlags().StringVar((*string)(&keysType), "keysType", "text", "Which input type you using for encrypto and encryto")
	viper.BindPFlag("keysType", RootCmd.PersistentFlags().Lookup("keysType"))

	// var SecretTag string
	RootCmd.PersistentFlags().StringVar(&SecretTag, "secretTag", "secret_tag", "the prefix of secret")
	viper.BindPFlag("secretTag", RootCmd.PersistentFlags().Lookup("secretTag"))

	// var oAuthLocation string
	RootCmd.PersistentFlags().StringVar(&oAuthLocation, "oAuthLocation", "", "location of the JSON key credentials file. If empty then use the Google Application Defaults.")
	viper.BindPFlag("oAuthLocation", RootCmd.PersistentFlags().Lookup("oAuthLocation"))

	// var flags string
	RootCmd.PersistentFlags().StringVar(&flags, "flags", "", "the sort code of the contact account")
	viper.BindPFlag("flags", RootCmd.PersistentFlags().Lookup("flags"))

	// var GCPprojectID string
	RootCmd.PersistentFlags().StringVar(&GCPprojectID, "GCPprojectID", "", "the projectID of GCP")
	viper.BindPFlag("GCPprojectID", RootCmd.PersistentFlags().Lookup("GCPprojectID"))

	// var GCPlocationID string
	RootCmd.PersistentFlags().StringVar(&GCPlocationID, "GCPlocationID", "", "the locationID of GCP")
	viper.BindPFlag("GCPlocationID", RootCmd.PersistentFlags().Lookup("GCPlocationID"))

	// var GCPkeyRingID string
	RootCmd.PersistentFlags().StringVar(&GCPkeyRingID, "GCPkeyRingID", "", "the keyRingID of GCP")
	viper.BindPFlag("GCPkeyRingID", RootCmd.PersistentFlags().Lookup("GCPkeyRingID"))

	// var GCPcryptoKeyID string
	RootCmd.PersistentFlags().StringVar(&GCPcryptoKeyID, "GCPcryptoKeyID", "", "the cryptoKeyID of GCP")
	viper.BindPFlag("GCPcryptoKeyID", RootCmd.PersistentFlags().Lookup("GCPcryptoKeyID"))

	// var AWSlocationID string
	RootCmd.PersistentFlags().StringVar(&AWSlocationID, "AWSlocationID", "", "the cryptoKeyID of AWS")
	viper.BindPFlag("AWSlocationID", RootCmd.PersistentFlags().Lookup("AWSlocationID"))

	// var AWScryptoKeyID string
	RootCmd.PersistentFlags().StringVar(&AWScryptoKeyID, "AWScryptoKeyID", "", "the cryptoKeyID of AWS")
	viper.BindPFlag("AWScryptoKeyID", RootCmd.PersistentFlags().Lookup("AWScryptoKeyID"))

	// var AWSprofile string
	RootCmd.PersistentFlags().StringVar(&AWSprofileName, "AWSprofileName", "", "the profile name used for AWS authentication")
	viper.BindPFlag("AWSprofileName", RootCmd.PersistentFlags().Lookup("AWSprofileName"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".macgyver") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

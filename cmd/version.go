package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Macgyver",
	Long:  `All software has versions. This is Macgyver's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString("version"))
	},
}

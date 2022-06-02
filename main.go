package main

import (
	"fmt"
	"os"

	"github.com/17media/macgyver/cmd"
	"github.com/spf13/viper"
)

var version = "dev"

func main() {
	viper.SetDefault("version", version)

	if err := cmd.RootCmd.Execute(); err != nil {
	    fmt.Println(err)
		os.Exit(1)
	}
}

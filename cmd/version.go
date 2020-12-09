package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display client version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("v%s", viper.GetString("app.version")))
		},
	}
}

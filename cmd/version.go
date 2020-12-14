package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"runtime"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display client version",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Name: \t\t%s\n", viper.GetString("app.name"))
			fmt.Printf("Version: \t%s\n", viper.GetString("app.version"))
			fmt.Printf("Git Commit: \t%s\n", viper.GetString("git.commit"))
			fmt.Printf("OS/ARCH: \t%s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}

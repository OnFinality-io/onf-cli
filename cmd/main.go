package main

import (
	"errors"
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/networkspec"
	"log"
	"path"

	"github.com/OnFinality-io/onf-cli/cmd/info"
	"github.com/OnFinality-io/onf-cli/cmd/node"
	"github.com/OnFinality-io/onf-cli/cmd/workspace"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var profile string
var version string
var gitCommit string

func init() {
	loadConfig()
}

func loadConfig() {
	home, _ := homedir.Dir()
	viper.SetConfigType("ini")
	viper.SetConfigName("credentials")
	viper.AddConfigPath(path.Join(home, ".onf"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("app.name", "onf-cli")
	viper.SetDefault("app.version", version)
	viper.SetDefault("git.commit", gitCommit)
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "onf",
		Version: viper.GetString("app.version"),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			accessKey := viper.GetString(fmt.Sprintf("%s.onf_access_key", profile))
			secretKey := viper.GetString(fmt.Sprintf("%s.onf_secret_key", profile))
			if accessKey == "" || secretKey == "" {
				return errors.New("invalid accessKey or secretKey")
			}
			service.Init(accessKey, secretKey)
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "profile scope in the credentials file")

	rootCmd.AddCommand(
		VersionCmd(),

		workspace.ListCmd(),
		workspace.MemberCmd(),
		workspace.InviteCmd(),

		node.New(),
		networkspec.New(),
		info.NewCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

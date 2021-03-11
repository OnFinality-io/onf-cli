package main

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/OnFinality-io/onf-cli/cmd/network"
	"github.com/OnFinality-io/onf-cli/cmd/networkspec"

	"github.com/OnFinality-io/onf-cli/cmd/info"
	"github.com/OnFinality-io/onf-cli/cmd/node"
	"github.com/OnFinality-io/onf-cli/cmd/setup"
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
	viper.SetDefault("app.name", "onf-cli")
	viper.SetDefault("app.version", version)
	viper.SetDefault("git.commit", gitCommit)
	viper.SetDefault("base_url", "https://api.onfinality.io/api/v1")
	viper.SetEnvPrefix("onf")
	viper.BindEnv("base_url")
}

func checkSetup() bool {
	credentialFile := &setup.CredentialFile{}
	return credentialFile.IsExistAtOnfAtHomeDir()
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
}

func main() {
	rootCmd := &cobra.Command{
		Use:           "onf",
		Version:       viper.GetString("app.version"),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if !checkSetup() {
				return errors.New("please run `onf setup` to initial the configurations")
			}
			loadConfig()
			accessKey := viper.GetString(fmt.Sprintf("%s.onf_access_key", profile))
			secretKey := viper.GetString(fmt.Sprintf("%s.onf_secret_key", profile))
			baseURL := viper.GetString("base_url")
			if accessKey == "" || secretKey == "" {
				return errors.New("invalid accessKey or secretKey")
			}
			service.Init(accessKey, secretKey, baseURL)
			return nil
		},
	}
	rootCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "profile scope in the credentials file")

	rootCmd.AddCommand(
		VersionCmd(),

		workspace.ListCmd(),
		workspace.MemberCmd(),
		workspace.InviteCmd(),

		node.NewCmd(),
		networkspec.NewCmd(),
		info.NewCmd(),
		setup.NewCmd(),
		network.NewCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

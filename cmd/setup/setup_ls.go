package setup

import (
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "Show all profile",
		Run: func(cmd *cobra.Command, args []string) {
			CredentialFile := &CredentialFile{}
			configArray := CredentialFile.GetAllSetup()
			for _, config := range configArray {
				if config.Section != "DEFAULT" {
					printer.New().PrintWithTitle(config.Section, config.Credential)
				}
			}
		},
	}
}

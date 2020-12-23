package setup

import (
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "ls",
		Short: "Show all profiles",
		Run: func(cmd *cobra.Command, args []string) {
			CredentialFile := &CredentialFile{}
			configArray := CredentialFile.GetAllSetup()
			for _, config := range configArray {
				if config.Section != "DEFAULT" {
					printer.NewWithPrintFlag(printFlags).PrintWithTitle(config.Section, config.Credential)
				}
			}
		},
	}
	printFlags.AddFlags(c)
	return c
}

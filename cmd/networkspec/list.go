package networkspec

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "list",
		Short: "list all the network specs in workspace",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			specs, err := service.GetNetworkSpecs(wsID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).Print(specs)
		},
	}
	printFlags.AddFlags(c)
	return c
}

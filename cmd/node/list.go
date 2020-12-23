package node

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the dedicated nodes in the workspace",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			nodes, err := service.GetNodeList(wsID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).Print(nodes)
		},
	}
	printFlags.AddFlags(cmd)
	return cmd
}

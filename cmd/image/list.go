package image

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	var networkID string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all the image in the current network specification",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ret, err := service.GetImage(wsID, networkID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).Print(ret)
		},
	}
	printFlags.AddFlags(cmd)
	cmd.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	_ = cmd.MarkFlagRequired("network")
	return cmd
}

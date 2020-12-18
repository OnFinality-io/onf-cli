package networkspec

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ShowCmd() *cobra.Command {
	var networkID string
	printFlags := printer.NewPrintFlags()
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show network specs in the OnFinality platform",
		Run: func(cmd *cobra.Command, args []string) {
			if networkID != "" {
				wsID, err := helpers.GetWorkspaceID(cmd)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				specs, err := service.GetNetworkSpec(wsID, networkID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				printer.NewWithPrintFlag(printFlags).Print(specs)
			}
		},
	}
	printFlags.AddFlags(cmd)
	cmd.PersistentFlags().Int64VarP(&wsID, "workspace", "w", 0, "Workspace ID")
	cmd.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	_ = cmd.MarkFlagRequired("network")
	return cmd
}

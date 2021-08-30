package networkspec

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func UpdateCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "update",
		Short: "update your network specs",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := &service.UpdateNetworkSpecPayload{}
			if filePath != "" {
				err = helpers.ApplyDefinitionFile(filePath, payload)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
			if payload.Config == nil {
				err = fmt.Errorf("metadata is required but not found")
				fmt.Println(err)
				return
			}
			err = service.UpdateNetworkSpec(wsID, networkID, payload)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Successfully update Network, #ID:", networkID)
		},
	}
	c.Flags().StringVarP(&filePath, "file", "f", "", "definition file for update node, yaml or json")
	c.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	_ = c.MarkFlagRequired("network")
	printFlags.AddFlags(c)
	return c
}

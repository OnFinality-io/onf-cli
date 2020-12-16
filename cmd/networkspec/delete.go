package networkspec

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	var networkID string
	c := &cobra.Command{
		Use:   "delete",
		Short: "delete a network",
		Run: func(cmd *cobra.Command, args []string) {
			if networkID != "" {
				wsID, err := helpers.GetWorkspaceID(cmd)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				err = service.DeleteNetworkSpecs(wsID, networkID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println("network", networkID, "is deleted")
			}
		},
	}
	c.Flags().StringVarP(&networkID, "network", "n", "", "network id")
	_ = c.MarkFlagRequired("network")
	return c
}

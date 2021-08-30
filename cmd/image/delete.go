package image

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func deleteCmd() *cobra.Command {
	var networkID string
	var id uint64
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a image version.",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = service.DeleteImage(wsID, networkID, id)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		},
	}
	cmd.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	_ = cmd.MarkFlagRequired("network")
	cmd.Flags().Uint64VarP(&id, "id", "i", 0, "Image ID")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

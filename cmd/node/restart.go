package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func restartCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "restart",
		Short: "restart a running node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = service.RestartNode(wsID, nodeID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("node", nodeID, "is restarted")
		},
	}
	c.Flags().Uint64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	return c
}

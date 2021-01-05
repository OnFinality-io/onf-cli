package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

var storageSize uint64

func expandStorageCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "expand-storage",
		Short: "Expand the storage for a given node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = service.ExpandNodeStorage(wsID, nodeID, storageSize)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		},
	}
	c.Flags().Uint64VarP(&nodeID, "node", "n", 0, "node id")
	c.Flags().Uint64VarP(&storageSize, "size", "s", 0, "storage size (Gi)")
	_ = c.MarkFlagRequired("node")
	_ = c.MarkFlagRequired("size")
	return c
}

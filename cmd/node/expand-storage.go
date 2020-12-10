package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

var storageSize int64

func expandStorageCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "expand-storage",
		Short: "expand storage for a given node",
		Run: func(cmd *cobra.Command, args []string) {
			err := service.ExpandNodeStorage(wsID, nodeID, storageSize)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		},
	}
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	c.Flags().Int64VarP(&storageSize, "size", "s", 0, "storage size (Gi)")
	_ = c.MarkFlagRequired("node")
	_ = c.MarkFlagRequired("size")
	return c
}

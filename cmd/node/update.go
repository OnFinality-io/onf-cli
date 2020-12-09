package node

import (
	"fmt"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "update",
		Short: "update a given node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("todo")
		},
	}
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	return c
}

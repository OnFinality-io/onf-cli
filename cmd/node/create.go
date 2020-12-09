package node

import (
	"fmt"
	"github.com/spf13/cobra"
)

func createCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "create a new dedicate node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("todo")
		},
	}
	return c
}

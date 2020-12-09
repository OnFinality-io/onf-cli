package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list all the dedicated nodes in workspace",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := service.GetWorkspaceList()
			fmt.Println(list, err)
		},
	}
}

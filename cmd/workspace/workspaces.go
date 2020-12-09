package workspace

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "workspaces",
		Short: "list all workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := service.GetWorkspaceList()
			if err != nil {
				return err[0]
			}
			fmt.Println(list)
			return nil
		},
	}
}

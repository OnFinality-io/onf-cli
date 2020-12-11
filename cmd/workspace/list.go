package workspace

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "workspaces",
		Short: "list all workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := service.GetWorkspaceList()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.New().Print(list)
		},
	}
}

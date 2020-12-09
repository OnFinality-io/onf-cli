package workspace

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/rodaine/table"
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
			t := table.New("ID", "Name", "Plan")
			for _, row := range list {
				t.AddRow(row.ID, row.Name, row.Plan)
			}
			t.Print()
		},
	}
}

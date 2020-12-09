package workspace

import (
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/rodaine/table"
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
			t := table.New("ID", "Name", "Plan")
			for _, row := range list {
				t.AddRow(row.ID, row.Name, row.Plan)
			}
			t.Print()
			return nil
		},
	}
}

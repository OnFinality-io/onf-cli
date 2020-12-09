package workspace

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func MemberCmd() *cobra.Command {
	var wsID int64
	c := &cobra.Command{
		Use:   "members",
		Short: "list all members in a given workspace",
		Run: func(cmd *cobra.Command, args []string) {
			members, err := service.GetMembers(wsID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			t := table.New("ID", "Name", "Email", "Role")
			for _, m := range members {
				t.AddRow(m.ID, m.Name, m.Email, m.Role)
			}
			t.Print()
		},
	}
	c.Flags().Int64VarP(&wsID, "workspace", "w", 0, "workspace id")
	_ = c.MarkFlagRequired("workspace")
	return c
}

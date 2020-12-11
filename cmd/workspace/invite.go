package workspace

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

var email string
var userRole string

func InviteCmd() *cobra.Command {
	var wsID int64
	c := &cobra.Command{
		Use:   "invite",
		Short: "invite a new member to join the workspace",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if len(args) == 0 {
				fmt.Println("Error: Please give the user email address")
				return
			}
			err = service.InviteMember(wsID, &service.InviteMemberPayload{
				Email: email,
				Role:  userRole,
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println()
		},
	}
	c.Flags().Int64VarP(&wsID, "workspace", "w", 0, "workspace id")
	c.Flags().StringVarP(&email, "email", "e", "", "email address")
	_ = c.MarkFlagRequired("email")
	c.Flags().StringVarP(&userRole, "role", "r", "member", "user role")
	return c
}

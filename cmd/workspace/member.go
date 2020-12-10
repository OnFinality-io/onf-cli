package workspace

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/utils/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
	"strconv"
	"sync"
)

type memberView struct {
	ID     string `header:"ID"`
	Name   string `header:"Name"`
	Email  string `header:"Email"`
	Role   string `header:"Role"`
	Status string `header:"Status"`
}

func MemberCmd() *cobra.Command {
	var wsID int64
	c := &cobra.Command{
		Use:   "members",
		Short: "list all members in a given workspace",
		Run: func(cmd *cobra.Command, args []string) {
			var members []service.Member
			var invites []service.InviteLog
			wg := sync.WaitGroup{}

			wg.Add(1)
			go func() {
				result, err := service.GetMembers(wsID)
				members = result
				if err != nil {
					fmt.Println(err.Error())
				}
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				result, err := service.GetInvitations(wsID)
				invites = result
				if err != nil {
					fmt.Println(err.Error())
				}
				wg.Done()
			}()

			wg.Wait()

			if members == nil || invites == nil {
				return
			}

			var view []memberView
			for _, m := range members {
				view = append(view, memberView{
					ID:     strconv.FormatInt(m.ID, 10),
					Name:   m.Name,
					Email:  m.Email,
					Role:   m.Role,
					Status: "active",
				})
			}
			for _, m := range invites {
				view = append(view, memberView{
					ID:     "",
					Name:   "",
					Email:  m.Email,
					Role:   m.Role,
					Status: "pending",
				})
			}
			printer.New().Print(view)
		},
	}
	c.Flags().Int64VarP(&wsID, "workspace", "w", 0, "workspace id")
	_ = c.MarkFlagRequired("workspace")
	return c
}

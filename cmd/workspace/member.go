package workspace

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

type memberView struct {
	ID     string `header:"ID"`
	Name   string `header:"Name"`
	Email  string `header:"Email"`
	Role   string `header:"Role"`
	Status string `header:"Status"`
}

func MemberCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	var wsID uint64
	c := &cobra.Command{
		Use:   "members",
		Short: "List all members in a given workspace",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

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
					ID:     strconv.FormatUint(m.ID, 10),
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
			printer.NewWithPrintFlag(printFlags).Print(view)
		},
	}
	printFlags.AddFlags(c)
	c.Flags().Uint64VarP(&wsID, "workspace", "w", 0, "workspace id")
	return c
}

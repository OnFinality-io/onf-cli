package workspace

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "workspaces",
		Short: "list all workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			list, err := service.GetWorkspaceList()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).Print(list)
		},
	}
	printFlags.AddFlags(c)
	return c
}

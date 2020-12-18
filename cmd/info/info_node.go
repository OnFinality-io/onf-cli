package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"

	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func NodeSpecCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "node-spec",
		Short: "Show node specs",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).Print(result.NodeSpecs)

		},
	}
	printFlags.AddFlags(c)

	c.AddCommand(
		NodeRecommendsCmd(),
	)

	return c
}

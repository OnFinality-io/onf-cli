package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func AllCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "all",
		Short: "List cluster and node specifications",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).PrintWithTitle("Clusters", result.Clusters)
			printer.NewWithPrintFlag(printFlags).PrintWithTitle("NodeSpecs", result.NodeSpecs)
			// printer.NewWithPrintFlag(printFlags).PrintWithTitle("Protocols", result.Protocols)

		},
	}
	printFlags.AddFlags(c)
	return c
}

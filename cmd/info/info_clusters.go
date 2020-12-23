package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"

	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ClusterCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "cluster",
		Short: "List clusters",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.NewWithPrintFlag(printFlags).Print(result.Clusters)

		},
	}
	printFlags.AddFlags(c)
	return c
}

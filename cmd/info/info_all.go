package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

type OutputView struct {
	Clusters  []service.Clusters  `json:"Clusters"`
	NodeSpecs []service.NodeSpecs `json:"NodeSpecs"`
}

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
			if printFlags.OutputFormat != nil && *printFlags.OutputFormat != "" {
				outputView := &OutputView{Clusters: result.Clusters, NodeSpecs: result.NodeSpecs}
				printer.NewWithPrintFlag(printFlags).Print(outputView)
			} else {
				printer.NewWithPrintFlag(printFlags).PrintWithTitle("Clusters", result.Clusters)
				printer.NewWithPrintFlag(printFlags).PrintWithTitle("NodeSpecs", result.NodeSpecs)
			}

		},
	}
	printFlags.AddFlags(c)
	return c
}

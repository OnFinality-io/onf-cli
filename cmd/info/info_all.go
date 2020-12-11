package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func AllCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "all",
		Short: "List cluster and node specifications",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.New().PrintWithTitle("Clusters", result.Clusters)
			printer.New().PrintWithTitle("NodeSpecs", result.NodeSpecs)
			// printer.New().PrintWithTitle("Protocols", result.Protocols)

		},
	}
	return c
}

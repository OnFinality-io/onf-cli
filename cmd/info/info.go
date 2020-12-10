package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/utils/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "manage dedicated nodes in the OnFinality platform",
	}
	cmd.AddCommand(
		GetCmd(),
		ListImageVersionsCmd(),
		NodeRecommendsCmd(),
	)
	return cmd
}

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get supported Clusters, NodeSpecs and Protocols",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			printer.New().PrintWithTitle("Clusters", result.Clusters)
			printer.New().PrintWithTitle("NodeSpecs", result.NodeSpecs)
			printer.New().PrintWithTitle("Protocols", result.Protocols)
		},
	}
}

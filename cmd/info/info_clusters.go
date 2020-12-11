package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"

	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ClusterCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "cluster",
		Short: "Show cluster",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			printer.New().Print(result.Clusters)

		},
	}
	return c
}

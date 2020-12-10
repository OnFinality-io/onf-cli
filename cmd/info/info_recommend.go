package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/utils/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func NodeRecommendsCmd() *cobra.Command {
	var network string
	c := &cobra.Command{
		Use:   "network",
		Short: "Show recommended node",
		Run: func(cmd *cobra.Command, args []string) {
			if network != "" {
				result, err := service.NodeRecommends(network)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				printer.New().Print(result)
			}

		},
	}

	c.Flags().StringVarP(&network, "name", "n", "", "network name")
	_ = c.MarkFlagRequired("name")
	return c
}

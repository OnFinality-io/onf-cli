package info

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get supported Clusters and NodeSpecs",
	}

	cmd.AddCommand(
		// ListImageVersionsCmd(),
		NodeRecommendsCmd(),
		AllCmd(),
		ClusterCmd(),
		NodeSpecCmd(),
	)
	return cmd
}

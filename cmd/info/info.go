package info

import (
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get supported clusters and node specifications",
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

package node

import (
	"github.com/spf13/cobra"
)

func bootstrapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap a new network from configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// read config
			// create networkspec if exist in config
			// 	upload chainspec

			// loop to create validators from config
			// monitor validator node running status
			// update session key for each node

			// loop to create bootnode with validator p2p address
			// update netwrokspec with new bootnode p2p address
		},
	}
	cmd.PersistentFlags().Int64VarP(&wsID, "workspace", "w", 0, "Workspace ID")
	return cmd
}

package node

import (
	"github.com/spf13/cobra"
)

var wsID uint64
var nodeID uint64

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage dedicated nodes in the OnFinality platform",
	}
	cmd.PersistentFlags().Uint64VarP(&wsID, "workspace", "w", 0, "Workspace ID")

	cmd.AddCommand(
		listCmd(),
		showCmd(),
		createCmd(),
		updateCmd(),
		stopCmd(),
		resumeCmd(),
		restartCmd(),
		terminateCmd(),
		expandStorageCmd(),
	)
	return cmd
}

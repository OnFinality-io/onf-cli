package node

import "github.com/spf13/cobra"

var wsID int64
var nodeID int64

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "manage dedicated nodes in the OnFinality platform",
	}
	cmd.PersistentFlags().Int64VarP(&wsID, "workspace", "w", 0, "Workspace ID (required)")
	_ = cmd.MarkPersistentFlagRequired("workspace")

	cmd.AddCommand(
		listCmd(),
		showCmd(),
		createCmd(),
		updateCmd(),
		stopCmd(),
		resumeCmd(),
		terminateCmd(),
	)
	return cmd
}

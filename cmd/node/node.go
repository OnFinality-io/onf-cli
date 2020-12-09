package node

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "manage dedicated nodes in the OnFinality platform",
	}
	var workspaceId string
	cmd.PersistentFlags().StringVarP(&workspaceId, "workspace", "w", "", "Workspace ID (required)")
	//_ = cmd.MarkPersistentFlagRequired("workspace")
	cmd.AddCommand(
		listCmd(),
	)
	return cmd
}

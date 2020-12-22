package node

import (
	"github.com/spf13/cobra"
)

var wsID int64

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "fast network manage",
	}
	cmd.PersistentFlags().Int64VarP(&wsID, "workspace", "w", 0, "Workspace ID")

	cmd.AddCommand(
		bootstrapCmd(),
	)
	return cmd
}

package image

import (
	"github.com/spf13/cobra"
)

var wsID uint64

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Get image list of network specification",
	}
	cmd.PersistentFlags().Uint64VarP(&wsID, "workspace", "w", 0, "Workspace ID")
	cmd.AddCommand(
		listCmd(),
		addCmd(),
		deleteCmd(),
	)
	return cmd
}

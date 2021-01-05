package networkspec

import (
	"github.com/spf13/cobra"
)

var wsID uint64

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-spec",
		Short: "Manage network specs in the OnFinality platform",
	}
	cmd.PersistentFlags().Uint64VarP(&wsID, "workspace", "w", 0, "Workspace ID")

	cmd.AddCommand(
		listCmd(),
		CreateCmd(),
		DeleteCmd(),
		ShowCmd(),
		GenerateCmd(),
		BootstrapCmd(),
		UploadCmd(),
	)
	return cmd
}

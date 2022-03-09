package networkspec

import (
	"github.com/spf13/cobra"
)

var wsID uint64
var networkID string

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network-spec",
		Short: "Manage network specs in the OnFinality platform",
	}
	cmd.PersistentFlags().Uint64VarP(&wsID, "workspace", "w", 0, "Workspace ID")

	cmd.AddCommand(
		listCmd(),
		listBackupsCmd(),
		CreateCmd(),
		DeleteCmd(),
		ShowCmd(),
		UploadCmd(),
		UpdateCmd(),
	)
	return cmd
}

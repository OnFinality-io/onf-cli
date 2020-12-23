package helpers

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetWorkspaceID(cmd *cobra.Command) (uint64, error) {
	wsID, err := cmd.Flags().GetUint64("workspace")
	if err != nil || wsID != 0 {
		return wsID, err
	}
	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return wsID, err
	}
	wsID = viper.GetUint64(fmt.Sprintf("%s.default_workspace", profile))
	return wsID, nil
}

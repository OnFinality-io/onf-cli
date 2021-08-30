package helpers

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetWorkspaceID(cmd *cobra.Command) (uint64, error) {

	var err error
	// Get workspace id from ENV
	wsID := viper.GetUint64("workspace_id")
	if wsID > 0 {
		return wsID, nil
	}

	// Get workspace id from flags
	wsID, err = cmd.Flags().GetUint64("workspace")
	if err != nil || wsID != 0 {
		return wsID, err
	}

	// Get workspace id from configuration
	profile, err := cmd.Flags().GetString("profile")
	if err != nil {
		return wsID, err
	}
	wsID = viper.GetUint64(fmt.Sprintf("%s.default_workspace", profile))
	if wsID > 0 {
		return wsID, nil
	}
	list, err := service.GetWorkspaceList()
	if err != nil {
		fmt.Printf("Get workspace err %v\n", err)
		return 0, err
	}

	// If there is only one workspace
	if len(list) == 1 {
		return list[0].ID, err
	} else {
		return 0, fmt.Errorf("You have %d workspaces. Please specify one of the workspaces.", len(list))
	}
	return wsID, nil
}

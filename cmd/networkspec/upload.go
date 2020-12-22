package networkspec

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func UploadCmd() *cobra.Command {
	var networkID string
	var files []string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload chain spec",
		Run: func(cmd *cobra.Command, args []string) {
			if networkID != "" {
				wsID, err := helpers.GetWorkspaceID(cmd)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				ret, err := service.UploadChainSpec(wsID, networkID, files)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(ret))
			}
		},
	}
	cmd.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	_ = cmd.MarkFlagRequired("network")
	cmd.Flags().StringSliceVarP(&files, "file", "f", nil, "Chain spec file")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

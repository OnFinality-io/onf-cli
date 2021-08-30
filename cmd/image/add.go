package image

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
	"strings"
)

func addCmd() *cobra.Command {
	var networkID, version string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add new image version to network.",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			spec, err := service.GetNetworkSpec(wsID, networkID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			version = strings.TrimSpace(version)
			ImageCheckProcess(spec.ImageRepository, version, true, func() {
				err = service.UpsertImage(wsID, networkID, &service.ImagePayload{ImageRepository: spec.ImageRepository, Version: &version})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			})

		},
	}
	cmd.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	_ = cmd.MarkFlagRequired("network")
	cmd.Flags().StringVarP(&version, "version", "v", "", "Image version")
	_ = cmd.MarkFlagRequired("version")
	return cmd
}

package networkspec

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func GenerateCmd() *cobra.Command {
	var networkID, imageVersion string
	printFlags := printer.NewPrintFlags()
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate chanspec for the network",
		Run: func(cmd *cobra.Command, args []string) {
			if networkID != "" {
				wsID, err := helpers.GetWorkspaceID(cmd)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				payload := &service.GenerateChainSpecPayload{ImageVersion: imageVersion}
				// err = helpers.ApplyDefinitionFile(filePath, payload)
				// if err != nil {
				// 	fmt.Println(err.Error())
				// 	return
				// }
				specs, err := service.GenerateChainSpec(wsID, networkID, payload)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				printer.NewWithPrintFlag(printFlags).Print(specs)
			}
		},
	}
	printFlags.AddFlags(cmd)
	cmd.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	cmd.Flags().StringVarP(&imageVersion, "image-version", "i", "", "Image version")
	_ = cmd.MarkFlagRequired("network")
	_ = cmd.MarkFlagRequired("image-version")
	return cmd
}

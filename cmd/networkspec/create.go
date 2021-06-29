package networkspec

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/cmd/image"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)
var filePath string

func CreateCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "create (-f FILENAME)",
		Short: "Create your network",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := &service.CreateNetworkSpecPayload{}
			err = helpers.ApplyDefinitionFile(filePath, payload)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			image.ImageCheckProcess(payload.ImageRepository, *payload.Metadata.ImageVersion, true, func() {
				specs, err := service.CreateNetworkSpecs(wsID, payload)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				printer.NewWithPrintFlag(printFlags).Print(specs)
			})
		},
	}
	printFlags.AddFlags(c)
	c.Flags().StringVarP(&filePath, "file", "f", "", "definition file for create network, yaml or json")
	return c
}

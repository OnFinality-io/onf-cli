package node

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

var filePath string

func createCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "create (-f FILENAME)",
		Short: "create a new dedicate node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := &service.CreateNodePayload{}
			err = helpers.ApplyDefinitionFile(filePath, payload)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			// todo - filter out the minimum parameters can be put on the arguments instead of using a definition file
			node, err := service.CreateNode(wsID, payload)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Successfully created node, #ID:", node.ID)
		},
	}
	c.Flags().StringVarP(&filePath, "file", "f", "", "definition file for create node, yaml or json")
	return c
}

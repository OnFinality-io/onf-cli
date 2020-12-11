package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "update",
		Short: "update a new dedicate node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			payload := &service.UpdateNodePayload{}
			if filePath != "" {
				err = helpers.ApplyDefinitionFile(filePath, payload)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
			// todo - filter out the minimum parameters can be put on the arguments instead of using a definition file
			err = service.UpdateNode(wsID, nodeID, payload)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Successfully update node, #ID:", nodeID)
		},
	}
	c.Flags().StringVarP(&filePath, "file", "f", "", "definition file for update node, yaml or json")
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	return c
}

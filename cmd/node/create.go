package node

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

var filePath string

func createCmd() *cobra.Command {
	watcherFlags := watcher.NewWatcherFlags()
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "create (-f FILENAME)",
		Short: "Create a new dedicated node",
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
			if node.ID > 0 {
				watcherFlags.ToWatch(func(done chan bool) {
					node, _ := service.GetNodeStatus(wsID, node.ID)
					if printFlags.OutputFormat != nil && *printFlags.OutputFormat != "" {
						printer.NewWithPrintFlag(printFlags).Print(node)

					} else {
						fmt.Println("current status is", node.Status)
					}
					if node.Status == Running {
						done <- true
					}
				})
			}
		},
	}
	c.Flags().StringVarP(&filePath, "file", "f", "", "definition file for create node, yaml or json")
	watcherFlags.AddFlags(c, "Watch for creation status")
	printFlags.AddFlags(c)
	return c
}

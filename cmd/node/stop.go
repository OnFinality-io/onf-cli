package node

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

func stopCmd() *cobra.Command {
	watcherFlags := watcher.NewWatcherFlags()
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "stop",
		Short: "stop a running node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = service.StopNode(wsID, nodeID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("node", nodeID, "is stopped")
			if nodeID > 0 {
				watcherFlags.ToWatch(func(done chan bool) {
					node, _ := service.GetNodeStatus(wsID, int64(nodeID))
					fmt.Println("current status is ", node.Status)
					if printFlags.OutputFormat != nil && *printFlags.OutputFormat != "" {
						printer.NewWithPrintFlag(printFlags).Print(node)

					} else {
						fmt.Println("current status is", node.Status)
					}
					if node.Status == Stopped {
						done <- true
					}
				})
			}
		},
	}
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	watcherFlags.AddFlags(c, "Watch for stop status")
	printFlags.AddFlags(c)
	return c
}

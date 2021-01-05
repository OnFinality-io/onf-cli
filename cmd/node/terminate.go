package node

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

func terminateCmd() *cobra.Command {
	watcherFlags := watcher.NewWatcherFlags()
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "terminate",
		Short: "Terminate a node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = service.TerminateNode(wsID, nodeID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("node", nodeID, "is terminated")
			if nodeID > 0 {
				watcherFlags.ToWatch(func(done chan bool) {
					node, _ := service.GetNodeStatus(wsID, nodeID)
					if printFlags.OutputFormat != nil && *printFlags.OutputFormat != "" {
						printer.NewWithPrintFlag(printFlags).Print(node)

					} else {
						fmt.Println("current status is", node.Status)
					}
					if node.Status == Terminated {
						done <- true
					}
				})
			}
		},
	}
	c.Flags().Uint64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	watcherFlags.AddFlags(c, "Watch for terminate status")
	printFlags.AddFlags(c)
	return c
}

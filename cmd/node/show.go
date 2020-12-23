package node

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

type propertyView struct {
	Key   string      `header:"Property"`
	Value interface{} `header:"Value"`
}

type OutputView struct {
	Information []propertyView `json:"Information"`
	Endpoints   []propertyView `json:"Endpoints"`
}

func showCmd() *cobra.Command {
	watcherFlags := watcher.NewWatcherFlags()
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "show",
		Short: "Show the detailed information on a given node",
		Run: func(cmd *cobra.Command, args []string) {

			if watcherFlags.Watch {
				watcherFlags.ToWatch(func(done chan bool) {
					node, _ := service.GetNodeStatus(wsID, int64(nodeID))
					show(cmd, printFlags)
					if node.Status == Running {
						done <- true
					}
				})

			} else {
				show(cmd, printFlags)
			}
		},
	}
	printFlags.AddFlags(c)
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	watcherFlags.AddFlags(c, "Watch for node")
	return c
}

func show(cmd *cobra.Command, printFlags *printer.PrintFlags) {
	wsID, err := helpers.GetWorkspaceID(cmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	node, err := service.GetNodeDetail(wsID, nodeID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	informationView := []propertyView{
		{"ID:", node.ID},
		{"Name:", node.Name},
		{"Workspace:", node.WorkspaceID},
		{"Owner:", node.OwnerID},
		{"Cluster:", node.ClusterHash},
		{"Network:", node.NetworkSpec.DisplayName},
		{"Image:", node.Image},
		{"Node Size:", fmt.Sprintf("%s (CPU: %s, RAM: %s)", node.NodeSpec, node.CPU, node.Ram)},
		{"Storage Size:", node.Storage},
		{"Status:", node.Status},
	}
	endpointsView := []propertyView{
		{"HTTPs", node.Endpoints.RPC},
		{"Websocket", node.Endpoints.WS},
	}
	if printFlags.OutputFormat != nil && *printFlags.OutputFormat != "" {
		outputView := &OutputView{Information: informationView, Endpoints: endpointsView}
		printer.NewWithPrintFlag(printFlags).Print(outputView)
	} else {
		printer.New().PrintWithTitle("Information", informationView)
		fmt.Println("")
		printer.New().PrintWithTitle("Endpoints", endpointsView)
	}
}

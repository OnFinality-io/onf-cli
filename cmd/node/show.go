package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

type propertyView struct {
	Key   string      `header:"Property"`
	Value interface{} `header:"Value"`
}

func showCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "show the detail information on a given node",
		Run: func(cmd *cobra.Command, args []string) {
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
			view := []propertyView{
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
			printer.New().PrintWithTitle("Information", view)
			fmt.Println("")

			view = []propertyView{
				{"HTTPs", node.Endpoints.RPC},
				{"Websocket", node.Endpoints.WS},
			}
			printer.New().PrintWithTitle("Endpoints", view)
		},
	}
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	return c
}

package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "show the detail information on a given node",
		Run: func(cmd *cobra.Command, args []string) {
			node, err := service.GetNodeDetail(wsID, nodeID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			t := table.New("Node INFO", "")
			t.AddRow("ID:", node.ID)
			t.AddRow("Name:", node.Name)
			t.AddRow("Cluster:", node.ClusterHash)
			t.AddRow("Network:", node.NetworkSpec.DisplayName)
			t.AddRow("Image:", node.Image)
			t.AddRow("Node Size:", fmt.Sprintf("%s (CPU: %s, RAM: %s)", node.NodeSpec, node.CPU, node.Ram))
			t.AddRow("Storage Size:", node.Storage)
			t.AddRow("Status:", node.Status)
			t.Print()

			fmt.Println("")
			t = table.New("Endpoints", "")
			t.AddRow("HTTPs", node.Endpoints.RPC)
			t.AddRow("Websocket", node.Endpoints.WS)
			t.Print()
		},
	}
	c.Flags().Int64VarP(&nodeID, "node", "n", 0, "node id")
	_ = c.MarkFlagRequired("node")
	return c
}

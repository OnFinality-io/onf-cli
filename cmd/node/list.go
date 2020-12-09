package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list all the dedicated nodes in workspace",
		Run: func(cmd *cobra.Command, args []string) {
			nodes, err := service.GetNodeList(wsID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			t := table.New("ID", "Name", "Network", "Cluster", "Status")
			for _, n := range nodes {
				t.AddRow(n.ID, n.Name, n.NetworkSpecKey, n.ClusterHash, n.Status)
			}
			t.Print()
		},
	}
}

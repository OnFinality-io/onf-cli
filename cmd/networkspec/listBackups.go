package networkspec

import (
	"fmt"
	"strings"

	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func listBackupsCmd() *cobra.Command {
	printFlags := printer.NewPrintFlags()
	c := &cobra.Command{
		Use:   "list-backups",
		Short: "List all the networks that supports lightning restore/backups",
		Run: func(cmd *cobra.Command, args []string) {
			// get public network specs data
			specs, err := service.GetBackupNetworkSpecs()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// get backups supported public networks data
			backups, err := service.GetBackups()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			
			// get cluster data
			info, err := service.GetInfo()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			nsBackups := mapNetworkSpecBackups(specs, backups, info.Clusters)

			printer.NewWithPrintFlag(printFlags).Print(nsBackups)
		},
	}
	printFlags.AddFlags(c)
	return c
}

func mapNetworkSpecBackups(
	ns []service.NetworkSpecBackups,
	nsb []service.Backups,
	cl []service.Clusters,
) []service.NetworkSpecBackups {
	var (
		nsBackups []service.NetworkSpecBackups
		cloudRegion, nodeTypeStr string
		nodeTypes []string
	)
	for _, n := range ns {
		for _, b := range nsb {
			if n.Key == b.NetworkSpec { // filter out by networks with backups only
				for _, c := range cl {
					if b.ClusterHash == c.Hash && c.Active { // filter out by cluster key and active clusters
						nodeTypes = b.GetNodeTypeFromPruningModeAndProtocol()
						cloudRegion = strings.ToUpper(c.Cloud) + " - " + c.Region

						if len(nodeTypes) > 0 {
							nodeTypeStr = strings.Join(nodeTypes, " | ")
						} else {
							nodeTypeStr = "-"
						}

						nsBackups = append(nsBackups, service.NetworkSpecBackups{
							Key: n.Key,
							Name: n.Name,
							DisplayName: n.DisplayName,
							ProtocolKey: n.ProtocolKey,
							MinStorageSize: b.StorageSize,
							AvailNodeTypes: nodeTypeStr,
							AvailCloudRegion: cloudRegion,
							ClusterKey: c.Hash,
						})
					}
				}
			}
		}
	}
	return nsBackups
}


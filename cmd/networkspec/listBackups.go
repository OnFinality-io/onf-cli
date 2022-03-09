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
			backups, err2 := service.GetBackups()
			if err2 != nil {
				fmt.Println(err2.Error())
				return
			}
			
			// get cluster data
			info, err3 := service.GetInfo()
			if err3 != nil {
				fmt.Println(err3.Error())
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
		nodeTypes, cloudRegion string
	)
	for _, n := range ns {
		for _, b := range nsb {
			if n.Key == b.NetworkSpec { // filter out by networks with backups only
				for _, c := range cl {
					if b.ClusterHash == c.Hash && c.Active { // filter out by cluster key and active clusters
						nodeTypes = getNodeTypeFromPruningMode(b.PruningMode, b.Protocol)
						cloudRegion = strings.ToUpper(c.Cloud) + " - " + c.Region

						nsBackups = append(nsBackups, service.NetworkSpecBackups{
							Key: n.Key,
							Name: n.Name,
							DisplayName: n.DisplayName,
							ProtocolKey: n.ProtocolKey,
							MinStorageSize: b.StorageSize,
							AvailNodeTypes: nodeTypes,
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

func getNodeTypeFromPruningMode(pruningMode string, protocol string) string {
	// TODO: find a way to get available pruning modes and protocols from APIs
	switch strings.ToLower(pruningMode) {
		case "archive":
			if protocol == "polkadot-parachain" {
				return "archive | collator"
			} else if protocol == "substrate" {
				return "archive | validator"
			} else {
				return "archive"
			}
		case "none":
			return "full"
		default:
			return "-"
	}
}

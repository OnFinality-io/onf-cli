package node

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/utils"
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

			if payload.InitFromBackup { // skip this step if not initializing from backup
				isBelowMinStorage, minStorageSize := isBackupStorageSizeBelow(payload)
				if isBelowMinStorage {
					fmt.Printf("Warning: Minimum %dGi required to use Lightning Restore/Backups\n", minStorageSize)
				}
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

// to check if the provided storage size is valid for backups enabled networks
// based on the combination of nodeType, clusterKey & networkSpecKey values
func isBackupStorageSizeBelow(payload *service.CreateNodePayload) (bool, int) {
	nodeType := string(payload.NodeType)
	clusterKey := payload.ClusterHash
	networkSpecKey := payload.NetworkSpecKey
	storageSize, err := strconv.ParseUint(strings.TrimSuffix(*payload.Storage, "Gi"), 10, 32)
	if err != nil {
		fmt.Println(err.Error())
		return false, -1
	}

	// get backups supported public networks data
	backups, err := service.GetBackups()
	if err != nil {
		fmt.Println(err.Error())
		return false, -1
	}
	
	for _, b := range backups {
		convertedNodeType := b.GetNodeTypeFromPruningMode()
		// filter out by networks with backups only, by cluster key and by node type
		if networkSpecKey == b.NetworkSpec && clusterKey == b.ClusterHash && utils.Contains(convertedNodeType, nodeType) { 
			if uint(storageSize) < b.StorageSize {
				return true, int(b.StorageSize)
			}
		}
	}

	return false, -1
}

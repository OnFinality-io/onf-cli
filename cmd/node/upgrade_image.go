package node

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/cmd/image"
	"github.com/OnFinality-io/onf-cli/pkg/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/utils"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
	"math"
	"sync"
	"time"
)

var printFlags *printer.PrintFlags
var percent int

func updateImageCmd() *cobra.Command {
	printFlags = printer.NewPrintFlags()
	var version, networkID string
	c := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the image version of node",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			err = upgradeNode(wsID, networkID, version)

			fmt.Println("Successfully update nodes")
		},
	}
	c.Flags().StringVarP(&version, "version", "v", "", "Image version")
	c.Flags().StringVarP(&networkID, "network", "n", "", "Network id")
	c.Flags().IntVarP(&percent, "percent", "", 0, "percent ")

	_ = c.MarkFlagRequired("version")
	_ = c.MarkFlagRequired("network")

	printFlags.AddFlags(c)
	return c
}

func upgradeNode(wsID uint64, networkID string, version string) error {
	var err error
	network, err := service.GetNetworkSpec(wsID, networkID)
	if err != nil {
		fmt.Println(err.Error())
	}
	image.ImageCheckProcess(network.ImageRepository, version, true, func() {
		err := service.UpsertImage(wsID, networkID, &service.ImagePayload{ImageRepository: network.ImageRepository, Version: &version})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		batchUpgrade(wsID, networkID, err, version)
	})

	return err
}

func batchUpgrade(wsID uint64, networkID string, err error, version string) {
	nodeList, err := service.GetNodeList(wsID)
	nodeLen := len(nodeList)
	toBeUpgradeNodes := make(map[string]*[]service.NodeItem, nodeLen)
	for _, node := range nodeList {
		//if node.Status != Terminated && node.Status != Terminating && strings.Compare(nodeImageVersion, version) != 0 {
		if node.NetworkSpecKey == networkID {
			if node.Status == Running || node.Status == Error {
				if nodeList, ok := toBeUpgradeNodes[node.ClusterHash]; ok {
					*nodeList = append(*nodeList, node)
				} else {
					toBeUpgradeNodes[node.ClusterHash] = &[]service.NodeItem{node}
				}
			}
		}
	}

	var wg sync.WaitGroup

	for _, nodes := range toBeUpgradeNodes {
		wg.Add(1)
		go func(nodes *[]service.NodeItem) {
			doUpgrade(nodes, wsID, version)
			wg.Done()
		}(nodes)

	}
	wg.Wait()
}

func doUpgrade(item *[]service.NodeItem, wsID uint64, version string) {
	chunkSize := math.Floor(float64(len(*item)) * float64(utils.Min(100, percent)) / 100)
	if chunkSize == 0 {
		chunkSize = 1
	}
	chunkNodes := chunkSlice(*item, int(chunkSize))

	for _, chunkNodeArray := range chunkNodes {
		var tasksWG sync.WaitGroup
		for _, nodeItem := range chunkNodeArray {
			tasksWG.Add(1)
			go func(node service.NodeItem) {
				updateNode(wsID, node, version)
				tasksWG.Done()
			}(nodeItem)
		}
		tasksWG.Wait()
	}
}

func chunkSlice(slice []service.NodeItem, chunkSize int) [][]service.NodeItem {
	var chunks [][]service.NodeItem
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func updateNode(wsID uint64, node service.NodeItem, version string) {
	service.UpdateImage(wsID, node.ID, &service.UpdateNodeImagePayload{ImageVersion: &version})
	watch := &watcher.Watcher{Second: time.Duration(3)}
	watch.Run(func(done chan bool) {
		status, _ := service.GetNodeStatus(wsID, node.ID)
		if printFlags.OutputFormat != nil && *printFlags.OutputFormat != "" {
			printer.NewWithPrintFlag(printFlags).Print(status)
		} else {
			fmt.Printf("Node %s in %s current status is %s\n", node.Name, node.ClusterHash, status.Status)
		}
		if status.Status == Running || status.Status == Error {
			//if status.Status == Running  {
			done <- true
		}
	})
}

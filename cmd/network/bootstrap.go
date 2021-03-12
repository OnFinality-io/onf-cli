package network

import (
	"fmt"
	"strconv"

	"github.com/OnFinality-io/onf-cli/pkg/utils"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/cmd/node"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func bootstrapCmd() *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap a new network from configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// read config
			cfg, err := ReadConfig(filePath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Read config success")

			//TODO check
			// Is  the sessionKey and validator count equal?
			// check if chanspec file exist

			// setup workspace id
			wsID, err = helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// create network spec if exist in config
			spec, err := CreateNetworkSpec(&cfg.NetworkSpec.Config)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(fmt.Sprintf("Network spec %s created", spec.Key))

			// 	upload chainspec
			chainspecFile := cfg.NetworkSpec.ChainSpec
			ret, err := service.UploadChainSpec(wsID, spec.Key, []string{chainspecFile})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if !ret.Success {
				// Continue the follow-up process
				fmt.Println("failed to upload chain spec")
				return
			}
			fmt.Println("Chainspec uploaded")

			// Initial operation for validator and boot node
			cfg.Validator.Node.NetworkSpecKey = spec.Key
			cfg.Validator.Node.ImageVersion = spec.Metadata.ImageVersion

			cfg.BootNode.Node.NetworkSpecKey = spec.Key
			cfg.BootNode.Node.ImageVersion = spec.Metadata.ImageVersion

			// create networkspec if exist in config
			validatorRet := CreateValidator(cfg.Validator)
			isInterrupted := false
			for _, v := range validatorRet {
				if v.Error != nil {
					isInterrupted = true
					fmt.Println("Create validator fail ", v.Error)
				}
			}
			if isInterrupted {
				return
			}
			fmt.Println("Create validators success")

			// boot node startup parameter settings
			var extraArgs []string
			for _, validator := range validatorRet {
				extraArgs = append(extraArgs, fmt.Sprintf("--bootnodes=%s", validator.Node.Endpoints.P2p))
			}

			bootNodeRet := CreateBootNode(extraArgs, cfg.BootNode)
			isInterrupted = false
			for _, v := range bootNodeRet {
				if v.Error != nil {
					isInterrupted = true
					fmt.Println("Create boot node fail ", v.Error)
				}
			}
			if isInterrupted {
				return
			}
			fmt.Println("Create boot nodes success")

			// update network spec with new bootnode p2p address
			var bootNodeAddrs []interface{}
			for _, bootNode := range bootNodeRet {
				bootNodeAddrs = append(bootNodeAddrs, service.BootNode{
					NodeID: utils.String(strconv.FormatUint(bootNode.Node.ID, 10)),
				})
			}

			err = service.UpdateNetworkSpecMetadata(wsID, spec.Key, &service.NetworkSpecMetadata{
				BootNodes: bootNodeAddrs,
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Updated network spec with new bootnode list:")
			for _, a := range bootNodeAddrs {
				fmt.Println("\t", *a.(service.BootNode).NodeID)
			}
			fmt.Println("New network launched")
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "config file for bootstrap a new network, yaml or json")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func ReadConfig(filePath string) (*CfgBootstrap, error) {
	payload := &CfgBootstrap{}
	err := helpers.ApplyDefinitionFile(filePath, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func CreateNetworkSpec(payload *service.CreateNetworkSpecPayload) (*service.NetworkSpec, error) {
	specs, err := service.CreateNetworkSpecs(wsID, payload)
	if err != nil {
		return nil, err
	}
	return specs, nil
}

type CreateNodeResult struct {
	Node  *service.Node
	Error error
}

// loop to create validators from config
func CreateValidator(cfgValidator CfgValidator) []*CreateNodeResult {
	validatorCount := cfgValidator.Count

	nodeChan := make(chan *CreateNodeResult, validatorCount)
	for i := 0; i < validatorCount; i++ {
		conf := cfgValidator.Node
		conf.NodeName = fmt.Sprintf("%s-%d", conf.NodeName, i)
		createdNode, err := service.CreateNode(wsID, &conf)
		if err != nil {
			return []*CreateNodeResult{{Error: err}}
		}
		fmt.Println(fmt.Sprintf("Node %s (%d) created", conf.NodeName, createdNode.ID))
		if i == 0 {
			// update networkspec
			err = service.UpdateNetworkSpecMetadata(wsID, conf.NetworkSpecKey, &service.NetworkSpecMetadata{
				BootNodes: []interface{}{service.BootNode{
					NodeID: utils.String(strconv.FormatUint(createdNode.ID, 10)),
				}},
			})
			fmt.Println(fmt.Sprintf("Node %s (%d) is set as bootnode", conf.NodeName, createdNode.ID))
			if err != nil {
				nodeChan <- &CreateNodeResult{Error: err}
				break
			}
		}
		go func(idx int, nodeId uint64) {
			// monitor validator node running status
			w := watcher.Watcher{Second: 10}
			w.Run(func(done chan bool) {
				status, _ := service.GetNodeStatus(wsID, nodeId)
				switch status.Status {
				case node.Terminating:
					nodeChan <- &CreateNodeResult{Error: fmt.Errorf("create %s %s", conf.NodeName, node.Terminating)}
					done <- true
				case node.Terminated:
					nodeChan <- &CreateNodeResult{Error: fmt.Errorf("create %s %s", conf.NodeName, node.Terminated)}
					done <- true
				case node.Error:
					nodeChan <- &CreateNodeResult{Error: fmt.Errorf("create %s err", conf.NodeName)}
					done <- true
				case node.Running:
					fmt.Println(fmt.Sprintf("Node %s (%d) is %s", conf.NodeName, nodeId, status.Status))
					//TODO need to add retry
					nodeDetail, err := service.GetNodeDetail(wsID, nodeId)
					if err != nil {
						nodeChan <- &CreateNodeResult{Error: fmt.Errorf("Get %s of detail %s", conf.NodeName, err)}
						done <- true
						return
					}

					if cfgValidator.SessionsKey != nil && cfgValidator.SessionsKey[idx] != nil {
						// update session key for each node
						for _, key := range cfgValidator.SessionsKey[idx] {
							err = service.InsertSessionKey(nodeDetail.Endpoints.RPC, &key)
							if err != nil {
								fmt.Println("key err", err)
							}
						}
						fmt.Println(fmt.Sprintf("Node %s (%d)'s session keys are updated", conf.NodeName, nodeDetail.ID))

						err = service.RestartNode(createdNode.WorkspaceID, createdNode.ID)
						if err != nil {
							nodeChan <- &CreateNodeResult{Error: fmt.Errorf("failed to restart %s, err: %s", conf.NodeName, err)}
							done <- true
							return
						}
					}

					fmt.Println(fmt.Sprintf("Node %s (%d) is restarted", conf.NodeName, nodeDetail.ID))
					nodeChan <- &CreateNodeResult{Node: nodeDetail}
					done <- true
				default:
					fmt.Println(fmt.Sprintf("Node %s (%d): %s", conf.NodeName, nodeId, status.Status))
				}

			})
		}(i, createdNode.ID)
	}
	nodeRet := []*CreateNodeResult{}
	for i := 0; i < validatorCount; i++ {
		select {
		case ret := <-nodeChan:
			nodeRet = append(nodeRet, ret)
		}
	}

	return nodeRet
}

func CreateBootNode(extraArgs []string, bootNode CfgBootNode) []*CreateNodeResult {
	// loop to create bootnode with validator p2p address
	bootNodeCount := bootNode.Count
	nodeChan := make(chan *CreateNodeResult, bootNodeCount)
	for i := 0; i < bootNodeCount; i++ {
		go func(idx int) {
			conf := bootNode.Node
			conf.NodeName = fmt.Sprintf("%s-%d", conf.NodeName, idx)
			if conf.Metadata != nil {
				conf.Metadata.ExtraArgs = append(conf.Metadata.ExtraArgs, extraArgs...)
			} else {
				conf.Metadata = &service.NodeMetadata{
					ExtraArgs: extraArgs,
				}
			}
			node, err := service.CreateNode(wsID, &conf)
			if err != nil {
				nodeChan <- &CreateNodeResult{Error: fmt.Errorf("create %s err %s", conf.NodeName, err)}
				return
			}
			node, err = service.GetNodeDetail(wsID, node.ID)
			if err != nil {
				nodeChan <- &CreateNodeResult{Error: fmt.Errorf("Get %s of detail %s", conf.NodeName, err)}
				return
			}

			fmt.Println(fmt.Sprintf("Node %s (%d) is created", conf.NodeName, node.ID))
			nodeChan <- &CreateNodeResult{Node: node}
		}(i)
	}
	nodeRet := []*CreateNodeResult{}
	for i := 0; i < bootNodeCount; i++ {
		select {
		case ret := <-nodeChan:
			nodeRet = append(nodeRet, ret)
		}
	}

	return nodeRet
}

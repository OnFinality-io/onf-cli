package network

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/networkspec/payload"
	"github.com/OnFinality-io/onf-cli/pkg/models"
	"github.com/OnFinality-io/onf-cli/pkg/utils"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/cmd/image"
	"github.com/OnFinality-io/onf-cli/cmd/node"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

const (
	bootnodes            = "--bootnodes"
	ProtocolParachainKey = "polkadot-parachain"
	Substrate            = "substrate"
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

			// setup workspace id
			wsID, err = helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// create network spec if exist in config
			spec, err := CreateNetworkSpec(cfg.NetworkSpec)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(fmt.Sprintf("Network spec %s created", spec.Key))

			// Initial operation for validator and boot node
			cfg.Validator.Node.NetworkSpecKey = spec.Key
			cfg.Validator.Node.ImageVersion = spec.Metadata.ImageVersion

			cfg.BootNode.Node.NetworkSpecKey = spec.Key
			cfg.BootNode.Node.ImageVersion = spec.Metadata.ImageVersion

			// create networkspec if exist in config
			validatorRet := CreateValidator(cfg.Validator, spec)
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

			var validatorExArgs []string
			for _, validator := range validatorRet {
				if validator.Node.Endpoints.P2p != "" {
					validatorExArgs = append(validatorExArgs, fmt.Sprintf("%s=%s", bootnodes, validator.Node.Endpoints.P2p))
				}
			}

			bootNodeRet := CreateBootNode(validatorExArgs, cfg.BootNode, spec)
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
			var bootNodeAddrs []string
			for _, bootNode := range bootNodeRet {
				bootNodeAddrs = append(bootNodeAddrs, bootNode.Node.Endpoints.P2p)
			}

			var argPayloads []*payload.ArgPayload
			for _, addr := range bootNodeAddrs {
				argPayloads = append(argPayloads, &payload.ArgPayload{Key: utils.String(bootnodes), Value: utils.String(addr)})
			}


			nodeTypes,err := service.GetSupportedNodeTypes(spec.ProtocolKey)
			NodeTypes := map[models.NodeType]*payload.ConfigRule{}
			for _, nodeType := range nodeTypes {
				NodeTypes[nodeType] = &payload.ConfigRule{
					Args:argPayloads,
				}
			}

			err = service.UpdateNetworkSpec(wsID, spec.Key, &service.UpdateNetworkSpecPayload{
				Config: &payload.ConfigPayload{
					NodeTypes: NodeTypes,
				},
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Updated network spec with new bootnode list:")
			for _, a := range bootNodeAddrs {
				fmt.Println("\t", a)
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
	if payload.ImageVersion == nil {
		return nil, fmt.Errorf("ImageVersion is required.")
	}
	var specs *service.NetworkSpec
	var err error
	image.ImageCheckProcess(payload.ImageRepository, *payload.ImageVersion, true, func() {
		specs, err = service.CreateNetworkSpecs(wsID, payload)
	})
	if specs == nil {
		return nil, fmt.Errorf("Image check err.")
	}
	if err != nil {
		return nil, err
	}
	return specs, err
}

type CreateNodeResult struct {
	Node  *service.Node
	Error error
}

// loop to create validators from config
func CreateValidator(cfgValidator CfgValidator, spec *service.NetworkSpec) []*CreateNodeResult {
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
			fmt.Println(fmt.Sprintf("Node %s (%d) is set as bootnode", conf.NodeName, createdNode.ID))
			if err != nil {
				nodeChan <- &CreateNodeResult{Error: err}
				break
			}
			w := watcher.Watcher{Second: 10}
			w.Run(func(done chan bool) {
				status, _ := service.GetNodeStatus(wsID, createdNode.ID)
				switch status.Status {
				case node.Error:
					nodeChan <- &CreateNodeResult{Error: fmt.Errorf("create %s err", conf.NodeName)}
					done <- true
					return
				case node.Running:
					fmt.Println(fmt.Sprintf("Node %s (%d) is %s", conf.NodeName, createdNode.ID, status.Status))

					var nodeDetail *service.Node
					watcher := watcher.Watcher{Second: 5}
					watcher.Run(func(done chan bool) {
						nodeDetail, err = service.GetNodeDetail(wsID, createdNode.ID)
						if err != nil {
							done <- true
							//fmt.Printf("Get p2p address err. %v",err)
							return
						}
						if nodeDetail.Endpoints.P2p != "" {
							//fmt.Printf("nodeName:%s,p2p address:%s\n",conf.NodeName,nodeDetail.Endpoints.P2p)
							done <- true
							return
						}
					})

					nodeTypes,err := service.GetSupportedNodeTypes(spec.ProtocolKey)
					NodeTypes := map[models.NodeType]*payload.ConfigRule{}
					for _, nodeType := range nodeTypes {
						NodeTypes[nodeType] = &payload.ConfigRule{
							Args:[]*payload.ArgPayload{{Key: utils.String(bootnodes), Value: utils.String(nodeDetail.Endpoints.P2p)}},
						}
					}

					// update networkspec
					err = service.UpdateNetworkSpec(wsID, conf.NetworkSpecKey,
						&service.UpdateNetworkSpecPayload{
							Config: &payload.ConfigPayload{
								NodeTypes: NodeTypes,
							},
						},
					)
					if err != nil {
						nodeChan <- &CreateNodeResult{Error: fmt.Errorf("Update network spec %s, err:%s", conf.NetworkSpecKey, err)}
						done <- true
						return
					}

					if nodeDetail.Endpoints.P2p != "" || nodeDetail.Endpoints.P2pInternal != "" {
						fmt.Println(fmt.Sprintf("Node %s (%d): p2p endpoint is %s", conf.NodeName, createdNode.ID, nodeDetail.Endpoints.P2p))
						done <- true
					}
				default:
					fmt.Println(fmt.Sprintf("Waiting for %s (%d) running", conf.NodeName, createdNode.ID))
				}
			})
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
								done <- true
								return
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

func CreateBootNode(exArgs []string, bootNode CfgBootNode, spec *service.NetworkSpec) []*CreateNodeResult {
	// loop to create bootnode with validator p2p address
	bootNodeCount := bootNode.Count
	nodeChan := make(chan *CreateNodeResult, bootNodeCount)
	for i := 0; i < bootNodeCount; i++ {
		go func(idx int) {
			conf := bootNode.Node
			conf.NodeName = fmt.Sprintf("%s-%d", conf.NodeName, idx)

			extraArgs := make(service.ExtraArgs)
			if conf.Config == nil {
				if spec.ProtocolKey == Substrate {
					extraArgs["default"] = &exArgs
				} else if spec.ProtocolKey == ProtocolParachainKey {
					extraArgs["parachain"] = &exArgs
				}
				conf.Config = &service.NodeLaunchConfig{ExtraArgs: &extraArgs}
			} else {
				if spec.ProtocolKey == Substrate {
					argsPayload := extraArgs["default"]
					*argsPayload = append(*argsPayload, exArgs...)
				} else if spec.ProtocolKey == ProtocolParachainKey {
					argsPayload := extraArgs["parachain"]
					*argsPayload = append(*argsPayload, exArgs...)
				}

			}

			node, err := service.CreateNode(wsID, &conf)
			if err != nil {
				nodeChan <- &CreateNodeResult{Error: fmt.Errorf("create %s err %s", conf.NodeName, err)}
				return
			}

			watcher := watcher.Watcher{Second: 3}
			watcher.Run(func(done chan bool) {
				node, err = service.GetNodeDetail(wsID, node.ID)
				if err != nil {
					nodeChan <- &CreateNodeResult{Error: fmt.Errorf("Get %s of detail %s", conf.NodeName, err)}
					done <- true
					return
				}
				if node.Endpoints.P2p != "" {
					done <- true
					return
				}
			})
			if err != nil {
				nodeChan <- &CreateNodeResult{Error: fmt.Errorf("Get %s of detail %s", conf.NodeName, err)}
				return
			}

			fmt.Println(fmt.Sprintf("Node %s (%d) is created", conf.NodeName, node.ID))
			nodeChan <- &CreateNodeResult{Node: node}
		}(i)
	}
	var nodeRet []*CreateNodeResult
	for i := 0; i < bootNodeCount; i++ {
		select {
		case ret := <-nodeChan:
			nodeRet = append(nodeRet, ret)
		}
	}

	return nodeRet
}

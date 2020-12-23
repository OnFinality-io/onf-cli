package network

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/utils"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"sync"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
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
			fmt.Println("Chain spec uploaded")

			cfg.Validator.Node.NetworkSpecKey = spec.Key
			cfg.Validator.Node.ImageVersion = spec.Metadata.ImageVersion

			cfg.BootNode.Node.NetworkSpecKey = spec.Key
			cfg.BootNode.Node.ImageVersion = spec.Metadata.ImageVersion

			wg := sync.WaitGroup{}
			var validatorNodes []*service.Node
			var bootNodes []*service.Node
			var errs []error

			doError := func(err error) {
				errs = append(errs, err)
				wg.Done()
			}

			// loop to create validators from config
			for i := 0; i < cfg.Validator.Count; i++ {
				wg.Add(1)
				go func(idx int) {
					conf := cfg.Validator.Node
					conf.NodeName = fmt.Sprintf("%s-%d", conf.NodeName, idx)
					node, err := service.CreateNode(wsID, &conf)
					if err != nil {
						doError(err)
						return
					}
					fmt.Println(fmt.Sprintf("Node %s (%d) created", conf.NodeName, node.ID))

					// monitor validator node running status
					w := watcher.Watcher{Second: 10}
					w.Run(func(done chan bool) {
						status, _ := service.GetNodeStatus(node.WorkspaceID, node.ID)
						if status.Status == "running" {
							fmt.Println(fmt.Sprintf("Node %s (%d) is %s", conf.NodeName, node.ID, status.Status))
							node, err = service.GetNodeDetail(node.WorkspaceID, node.ID)
							if err != nil {
								doError(err)
								done <- true
								return
							}

							// update session key for each node
							for _, key := range cfg.Validator.SessionsKey[idx] {
								err = service.InsertSessionKey(node.Endpoints.RPC, &key)
								if err != nil {
									fmt.Println("key err", err)
								}
							}
							fmt.Println(fmt.Sprintf("Node %s (%d)'s session keys are updated", conf.NodeName, node.ID))

							validatorNodes = append(validatorNodes, node)
							wg.Done()

							done <- true
						} else {
							fmt.Println(fmt.Sprintf("Node %s (%d): %s", conf.NodeName, node.ID, status.Status))
						}
					})
				}(i)
			}
			wg.Wait()
			if len(errs) > 0 {
				fmt.Println(errs)
				return
			}
			fmt.Println(validatorNodes, errs)

			var extraArgs []string
			for _, n := range validatorNodes {
				extraArgs = append(extraArgs, "--bootnodes", n.Endpoints.P2p)
			}
			fmt.Println(extraArgs)

			// loop to create bootnode with validator p2p address
			errs = []error{}
			for i := 0; i < cfg.BootNode.Count; i++ {
				wg.Add(1)
				go func(idx int) {
					conf := cfg.BootNode.Node
					conf.NodeName = fmt.Sprintf("%s-%d", conf.NodeName, idx)
					conf.Metadata = &service.NodeMetadata{
						ExtraArgs: extraArgs,
					}
					node, err := service.CreateNode(wsID, &conf)
					if err != nil {
						doError(err)
						return
					}
					node, err = service.GetNodeDetail(wsID, node.ID)
					if err != nil {
						doError(err)
						return
					}
					bootNodes = append(bootNodes, node)
					fmt.Println(fmt.Sprintf("Node %s (%d) is created", conf.NodeName, node.ID))
					wg.Done()
				}(i)
			}
			wg.Wait()
			if len(errs) > 0 {
				fmt.Println(errs)
				return
			}

			// update network spec with new bootnode p2p address
			var bootNodeAddrs []service.BootNode
			for _, n := range bootNodes {
				bootNodeAddrs = append(bootNodeAddrs, service.BootNode{
					Address: utils.String(n.Endpoints.P2p),
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
	specs, err := service.CreateNetworkSpecs(wsID, payload)
	if err != nil {
		return nil, err
	}
	return specs, nil
}

package network

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/OnFinality-io/onf-cli/cmd/helpers"
	"github.com/OnFinality-io/onf-cli/cmd/node"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/spf13/cobra"
)

func bootstrapCmd() *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap a new network from configuration",
		Run: func(cmd *cobra.Command, args []string) {
			// read config
			bootstrap, err := ReadConfig(filePath)
			if err != nil {
				fmt.Println("Read config fail ", err.Error())
				return
			}
			// fmt.Println("read config bootstrap:", bootstrap.Validator.Node)

			//TODO check
			// Is  the sessionKey and validator count equal?
			// check if chanspec file exist
			fmt.Println("Read config success")

			// setup workspace id
			wsID, err = helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// create networkspec if exist in config
			networkSpecEntity, err := CreateNetworkSpec(&bootstrap.NetworkSpec.Config)
			if err != nil {
				fmt.Println("Create network specs fail ", err.Error())
				return
			}
			fmt.Println("Create network specs success")

			// 	upload chainspec
			chainspecFile := bootstrap.NetworkSpec.ChainSpec
			ret, err := UploadChanSpec(networkSpecEntity, []string{chainspecFile})
			if err != nil {
				fmt.Println("Upload chanspec fail ", err.Error())
				return
			}

			if ret.Success {
				fmt.Println("Upload chanspec success")

				// loop to create validators from config
				// monitor validator node running status
				// update session key for each node
				validatorRet := CreateValidator(bootstrap.Validator)
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

				// loop to create bootnode with validator p2p address
				// update netwrokspec with new bootnode p2p address
				bootNodeRet := CreateBootNode(validatorRet, bootstrap.BootNode, networkSpecEntity)
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
			}

		},
	}
	cmd.PersistentFlags().Int64VarP(&wsID, "workspace", "w", 0, "Workspace ID")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "definition file for create node, yaml or json")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

func ReadConfig(filePath string) (*Bootstrap, error) {
	payload := &Bootstrap{}
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

type UploadResult struct {
	Success bool `json:"success"`
}

func UploadChanSpec(networkSpec *service.NetworkSpec, files []string) (*UploadResult, error) {
	networkID := networkSpec.Key
	ret, err := service.UploadChainSpec(wsID, networkID, files)
	if err != nil {
		return nil, err
	}
	uploadRet := &UploadResult{}
	err = json.Unmarshal(ret, uploadRet)
	if err != nil {
		return nil, err
	}
	return uploadRet, nil
}

type CreateNodeResult struct {
	Node  *service.Node
	Error error
}

func CreateValidator(validator CfgValidator) []*CreateNodeResult {
	createNodePayload := validator.Node
	validatorCount := validator.Count

	// loop to create validators from config
	validRet := make(chan *CreateNodeResult, validatorCount)
	for i := 0; i < validatorCount; i++ {
		sessionKeySlice := validator.SessionsKey[i]
		go func(index string, sessionKeySlice []service.SessionKey) {
			cloneNodePayload := *createNodePayload
			if cloneNodePayload.NodeName == "" {
				cloneNodePayload.NodeName = "validator-" + index
			} else {
				cloneNodePayload.NodeName = cloneNodePayload.NodeName + "-" + index
			}
			cloneNodePayload.NodeType = "validator"
			createdNode, err := service.CreateNode(wsID, &cloneNodePayload)
			if err != nil {
				validRet <- &CreateNodeResult{Error: err}
				return
			}
			// monitor validator node running status
			isErr := false
			// monitor node running status
			if createdNode.ID > 0 {
				watch := &watcher.Watcher{Second: time.Duration(2)}
				watch.Run(func(done chan bool) {
					nodeStatus, _ := service.GetNodeStatus(wsID, int64(createdNode.ID))
					switch nodeStatus.Status {
					case node.Running:
						done <- true
					case node.Terminating:
						validRet <- &CreateNodeResult{Error: fmt.Errorf("create %s %s", cloneNodePayload.NodeName, node.Terminating)}
						done <- true
					case node.Terminated:
						validRet <- &CreateNodeResult{Error: fmt.Errorf("create %s %s", cloneNodePayload.NodeName, node.Terminated)}
						done <- true
					case node.Error:
						isErr = true
						validRet <- &CreateNodeResult{Error: fmt.Errorf("create %s err", cloneNodePayload.NodeName)}
						done <- true
					}
				})
			}
			if isErr {
				return
			}

			// Get node detail
			nodeDetail, err := service.GetNodeDetail(wsID, int64(createdNode.ID))
			if err != nil {
				// fmt.Println(err.Error())
				validRet <- &CreateNodeResult{Error: err}
				return
			}

			// update session key for each node
			rpcURL := nodeDetail.Endpoints.RPC
			for _, sessionKey := range sessionKeySlice {
				service.InsertSessionKey(rpcURL, &sessionKey)
			}

			validRet <- &CreateNodeResult{Node: nodeDetail}

		}(strconv.Itoa(i), sessionKeySlice)
	}

	nodeRet := []*CreateNodeResult{}
	for i := 0; i < validatorCount; i++ {
		select {
		case ret := <-validRet:
			nodeRet = append(nodeRet, ret)
		}
	}

	return nodeRet
}

func CreateBootNode(validatorNodeResult []*CreateNodeResult, bootNode CfgBootNode, networkSpec *service.NetworkSpec) []*CreateNodeResult {
	createNodePayload := bootNode.Node
	// validatorAddr := []string{}
	// for _, v := range validatorNodeResult {
	// 	val := "--sentry-nodes=" + v.Node.Endpoints.P2pInternal
	// 	validatorAddr = append(validatorAddr, val)
	// 	val = "--reserved-nodes=" + v.Node.Endpoints.P2pInternal
	// 	validatorAddr = append(validatorAddr, val)
	// 	val = "--reserved-only"
	// 	validatorAddr = append(validatorAddr, val)
	// }
	// createNodePayload.Metadata = &service.NodeMetadata{ExtraArgs: validatorAddr}

	// loop to create bootnode with validator p2p address
	bootNodeCount := bootNode.Count
	validRet := make(chan *CreateNodeResult, bootNodeCount)
	for i := 0; i < bootNodeCount; i++ {
		go func(index string) {
			cloneNodePayload := *createNodePayload
			if cloneNodePayload.NodeName == "" {
				cloneNodePayload.NodeName = "boot-node-" + index
			} else {
				cloneNodePayload.NodeName = cloneNodePayload.NodeName + "-" + index
			}
			createdNode, err := service.CreateNode(wsID, &cloneNodePayload)
			if err != nil {
				validRet <- &CreateNodeResult{Error: err}
				return
			}

			isErr := false
			// monitor node running status
			if createdNode.ID > 0 {
				watch := &watcher.Watcher{Second: time.Duration(2)}
				watch.Run(func(done chan bool) {
					nodeStatus, _ := service.GetNodeStatus(wsID, int64(createdNode.ID))
					switch nodeStatus.Status {
					case node.Running:
						done <- true
					case node.Terminating:
						validRet <- &CreateNodeResult{Error: fmt.Errorf("create %s %s", cloneNodePayload.NodeName, node.Terminating)}
						done <- true
					case node.Terminated:
						validRet <- &CreateNodeResult{Error: fmt.Errorf("create %s %s", cloneNodePayload.NodeName, node.Terminated)}
						done <- true
					case node.Error:
						isErr = true
						validRet <- &CreateNodeResult{Error: fmt.Errorf("create %s err", cloneNodePayload.NodeName)}
						done <- true
					}
				})
			}
			if isErr {
				return
			}

			// Get node detail
			nodeDetail, err := service.GetNodeDetail(wsID, int64(createdNode.ID))
			if err != nil {
				validRet <- &CreateNodeResult{Error: err}
				return
			}

			// update netwrokspec with new bootnode p2p address
			err = service.UpdateNetworkSpecMetadata(wsID, networkSpec.Key, &networkSpec.Metadata)
			if err != nil {
				validRet <- &CreateNodeResult{Error: err}
				return
			}

			validRet <- &CreateNodeResult{Node: nodeDetail}
		}(strconv.Itoa(i))
	}

	nodeRet := []*CreateNodeResult{}
	for i := 0; i < bootNodeCount; i++ {
		select {
		case ret := <-validRet:
			nodeRet = append(nodeRet, ret)
		}
	}

	return nodeRet
}

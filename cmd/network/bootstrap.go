package network

import (
	"fmt"

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
			bootstrap, err := ReadConfig(filePath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			// create networkspec if exist in config
			fmt.Println(&bootstrap.NetworkSpec.Config)
			networkSpecEntity, err := CreateNetworkSpec(&bootstrap.NetworkSpec.Config)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(networkSpecEntity)
			// 	upload chainspec

			// loop to create validators from config
			// monitor validator node running status
			// update session key for each node

			// loop to create bootnode with validator p2p address
			// update netwrokspec with new bootnode p2p address
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

func CreateNetworkSpec(payload *service.CreateNetworkSpecPayload) (*service.NetworkSpecEntity, error) {
	specs, err := service.CreateNetworkSpecs(wsID, payload)
	if err != nil {
		return nil, err
	}
	return specs, nil
}

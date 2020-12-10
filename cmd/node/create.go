package node

import (
	"encoding/json"
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/OnFinality-io/onf-cli/pkg/utils"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path/filepath"
)

var filePath string

func createCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "create",
		Short: "create a new dedicate node",
		Run: func(cmd *cobra.Command, args []string) {
			payload := &service.CreateNodePayload{}
			if filePath != "" {
				ext := filepath.Ext(filePath)
				if !utils.Contains([]string{".json", ".yaml", ".yml"}, ext) {
					fmt.Println("Error: definition file must be in JSON and YAML format")
					return
				}
				data, err := ioutil.ReadFile(filePath)
				if err != nil {
					fmt.Println("Error: Failed to read file", filePath)
					return
				}
				if ext == ".json" {
					err = json.Unmarshal(data, payload)
				} else {
					err = yaml.Unmarshal(data, payload)
				}
				if err != nil {
					fmt.Println("Error: Invalid definition file", filePath)
					return
				}
			}
			// todo - filter out the minimum parameters can be put on the arguments instead of using a definition file
			node, err := service.CreateNode(wsID, payload)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("Successfully created node, #ID:", node.ID)
		},
	}
	c.Flags().StringVarP(&filePath, "file", "f", "", "definition file for create node, yaml or json")
	return c
}

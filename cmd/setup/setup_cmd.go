package setup

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "initialize the configuration",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			Flow()
		},
	}
	return cmd
}

func Flow() {

	credential := &Credential{}

	// access key
	accessKeyPrompt := promptui.Prompt{
		Label: "Please input your access key",
	}
	result, err := accessKeyPrompt.Run()
	if err != nil {
		fmt.Printf("Fail to add access key %v\n", err)
		return
	}
	credential.AccessKey = result

	// secret key
	secretKeyPrompt := promptui.Prompt{
		Label: "Please input your secret key",
	}
	result, err = secretKeyPrompt.Run()
	if err != nil {
		fmt.Printf("Fail to add secret key %v\n", err)
		return
	}
	credential.SecretKey = result

	// workspace id key
	service.Init(credential.AccessKey, credential.SecretKey)
	list, err := service.GetWorkspaceList()
	if len(list) == 1 {
		credential.WorkspaceID = list[0].ID
	} else {
		var name []string
		for _, ws := range list {
			name = append(name, ws.Name)
		}
		workspaceIDPrompt := promptui.Select{
			Label: "Please select your workspace",
			Items: name,
		}
		index, _, err := workspaceIDPrompt.Run()
		if err != nil {
			fmt.Printf("Fail to add workspace id %v\n", err)
			return
		}
		credential.WorkspaceID = list[index].ID
	}

	config := &CredentialConfig{
		Credential: credential,
	}
	PersistentCredential(config)
}

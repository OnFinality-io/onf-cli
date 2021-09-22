package networkspec

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/cmd/helpers"

	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func UploadCmd() *cobra.Command {
	var files []string
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload chain specification",
		Run: func(cmd *cobra.Command, args []string) {
			wsID, err := helpers.GetWorkspaceID(cmd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			res, err := service.UploadPrivateFile(wsID, files)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Printf("chainspec upload successfully. Key:%s\n", *res.Key)

		},
	}
	cmd.Flags().StringSliceVarP(&files, "file", "f", nil, "Chain spec file")
	_ = cmd.MarkFlagRequired("file")
	return cmd
}

package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/cmd/utils/printer"
	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ListImageVersionsCmd() *cobra.Command {
	var image string
	c := &cobra.Command{
		Use:   "image",
		Short: "List the image version",
		Run: func(cmd *cobra.Command, args []string) {
			if image != "" {
				result, err := service.ListImageVersions(image)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				printer.New().Print(result)
			}

		},
	}

	c.Flags().StringVarP(&image, "name", "n", "", "image name")
	_ = c.MarkFlagRequired("name")
	return c
}

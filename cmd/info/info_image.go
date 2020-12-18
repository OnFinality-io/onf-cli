package info

import (
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/printer"

	"github.com/OnFinality-io/onf-cli/pkg/service"
	"github.com/spf13/cobra"
)

func ListImageVersionsCmd() *cobra.Command {
	var image string
	printFlags := printer.NewPrintFlags()
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

				printer.NewWithPrintFlag(printFlags).Print(result)
			}

		},
	}
	printFlags.AddFlags(c)
	c.Flags().StringVarP(&image, "name", "n", "", "image name")
	_ = c.MarkFlagRequired("name")
	return c
}

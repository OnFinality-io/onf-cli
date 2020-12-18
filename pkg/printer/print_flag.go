package printer

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type PrintFlags struct {
	JSONYamlPrintFlags *JSONYamlPrintFlags
	OutputFormat       *string
}

func (f *PrintFlags) AddFlags(cmd *cobra.Command) {
	if f.OutputFormat != nil {
		cmd.Flags().StringVarP(f.OutputFormat, "output", "o", *f.OutputFormat, fmt.Sprintf("Output format. One of: %s.", strings.Join(f.AllowedFormats(), "|")))
	}
}
func (f *PrintFlags) ToPrinter() (ResourcePrinter, error) {
	var printer ResourcePrinter
	var outputFormat string
	if f.OutputFormat != nil {
		outputFormat = strings.ToLower(*f.OutputFormat)
	}
	switch outputFormat {
	case "json":
		printer = &JSONPrinter{}
	case "yaml":
		printer = &YAMLPrinter{}
	default:
		printer = &TablePrinter{}
	}

	return printer, nil
}

func (f *PrintFlags) AllowedFormats() []string {
	ret := []string{}
	ret = append(ret, f.JSONYamlPrintFlags.AllowedFormats()...)
	return ret
}

func NewPrintFlags() *PrintFlags {
	outputFormat := ""
	return &PrintFlags{
		OutputFormat:       &outputFormat,
		JSONYamlPrintFlags: NewJSONYamlPrintFlags(),
	}
}

package printer

import (
	"io"
	"os"

	"github.com/lensesio/tableprinter"
)

type TablePrinter struct {
	BgBlackColor  int
	HeaderFgColor int
}

func (p *TablePrinter) PrintObj(obj interface{}, w io.Writer) error {
	printer := tableprinter.New(os.Stdout)
	// Optionally, customize the table, import of the underline 'tablewriter' package is required for that.
	// printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = ""
	printer.ColumnSeparator = ""
	printer.RowSeparator = ""
	printer.AutoFormatHeaders = false
	printer.HeaderLine = false
	printer.HeaderBgColor = p.BgBlackColor
	printer.HeaderFgColor = p.HeaderFgColor

	printer.Print(obj)
	return nil
}

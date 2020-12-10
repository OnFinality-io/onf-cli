package printer

import (
	"fmt"
	"os"

	"github.com/lensesio/tableprinter"
)

type Printer struct {
	BgBlackColor  int
	HeaderFgColor int
}

func New() *Printer {
	return &Printer{}
}

func (p Printer) Print(in interface{}, filters ...interface{}) {
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

	printer.Print(in, filters)
}

func (p Printer) PrintWithTitle(title string, in interface{}, filters ...interface{}) {
	fmt.Println(title)
	p.Print(in, filters)
}

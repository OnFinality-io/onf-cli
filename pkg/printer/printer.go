package printer

import (
	"fmt"
	"os"
)

type Printer struct {
	PrintFlags *PrintFlags
}

func New() *Printer {
	return &Printer{PrintFlags: NewPrintFlags()}
}
func NewWithPrintFlag(printFlags *PrintFlags) *Printer {
	return &Printer{PrintFlags: printFlags}
}

func (p Printer) Print(in interface{}) {
	out, _ := p.PrintFlags.ToPrinter()
	out.PrintObj(in, os.Stdout)
}

func (p Printer) PrintWithTitle(title string, in interface{}) {
	fmt.Println(title)
	p.Print(in)
}

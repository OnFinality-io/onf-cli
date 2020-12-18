package printer

import "io"

type ResourcePrinter interface {
	PrintObj(obj interface{}, w io.Writer) error
}

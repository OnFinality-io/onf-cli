package printer

import (
	"encoding/json"
	"io"
)

type JSONPrinter struct{}

func (p *JSONPrinter) PrintObj(obj interface{}, w io.Writer) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}

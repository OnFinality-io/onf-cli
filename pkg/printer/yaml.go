package printer

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v2"
)

type YAMLPrinter struct {
}

func (p *YAMLPrinter) PrintObj(obj interface{}, w io.Writer) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	data, err = JSONToYAML(data)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return err
}

// JSONToYAML Converts JSON to YAML.
func JSONToYAML(j []byte) ([]byte, error) {
	// Convert the JSON to an object.
	var jsonObj interface{}
	// We are using yaml.Unmarshal here (instead of json.Unmarshal) because the
	// Go JSON library doesn't try to pick the right number type (int, float,
	// etc.) when unmarshalling to interface{}, it just picks float64
	// universally. go-yaml does go through the effort of picking the right
	// number type, so we can preserve number type throughout this process.
	err := yaml.Unmarshal(j, &jsonObj)
	if err != nil {
		return nil, err
	}

	// Marshal this object into YAML.
	return yaml.Marshal(jsonObj)
}

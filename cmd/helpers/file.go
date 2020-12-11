package helpers

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"os"
)

func ApplyDefinitionFile(file string, payload interface{}) error {
	var data []byte
	var err error
	if file == "-" {
		data = dataFromStdin()
	} else {
		data, err = ioutil.ReadFile(file)
		if err != nil {
			return errors.New("failed to read file")
		}
	}

	err = json.Unmarshal(data, payload)
	if err != nil {
		err = yaml.Unmarshal(data, payload)
		if err != nil {
			return errors.New("invalid definition")
		}
	}
	return nil
}

func dataFromStdin() []byte {
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return []byte(string(output))
}

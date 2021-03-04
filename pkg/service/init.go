package service

import (
	"encoding/json"
	"fmt"

	"github.com/OnFinality-io/onf-cli/pkg/api"
	"github.com/parnurzeal/gorequest"
)

var instance *api.Api

func Init(accessKey, secretKey string, baseURL string) {
	instance = api.New(accessKey, secretKey, baseURL)
}

type errResponse struct {
	Message interface{} `json:"message"`
}

func checkError(resp gorequest.Response, data []byte, err []error) error {
	if resp != nil && resp.StatusCode >= 300 {
		r := errResponse{}
		err := json.Unmarshal(data, &r)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: %s", resp.Status, r.Message)
	}
	if len(err) > 0 {
		return err[0]
	}
	return nil
}

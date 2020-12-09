package service

import "github.com/OnFinality-io/onf-cli/pkg/api"

var instance *api.Api

func Init(accessKey, secretKey string) {
	instance = api.New(accessKey, secretKey)
}

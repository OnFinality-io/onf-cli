package base

import "github.com/OnFinality-io/onf-cli/pkg/onf_debug"

const (
	baseURL = "https://api.onfinality.io/api/v1"
)

func BaseUrl() string {
	debugBaseUrl := onf_debug.GetDebug().Baseurl
	if debugBaseUrl != "" {
		return debugBaseUrl
	}
	return baseURL
}

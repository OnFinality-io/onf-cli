package service

import (
	"github.com/OnFinality-io/onf-cli/pkg/api"
	"time"
)

type ImageCheckStatus string

const (
	Pending ImageCheckStatus = "pending"
	Success ImageCheckStatus = "success"
	Fail    ImageCheckStatus = "fail"
	Timeout ImageCheckStatus = "timeout"
)

type ImageCheckPayload struct {
	ImageRepository string  `json:"imageRepository" header:"imageRepository"`
	Version         *string `json:"version" header:"version"`
}

type ImageCheckResult struct {
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	ID              int              `json:"id"`
	ProtocolKey     string           `json:"protocolKey"`
	ImageRepository string           `json:"imageRepository"`
	Version         string           `json:"version"`
	Client          string           `json:"client"`
	Active          bool             `json:"active"`
	Config          Config           `json:"config"`
	Status          ImageCheckStatus `json:"status"`
	Image           string           `json:"image"`
}
type CheckStatus struct {
	Step   string `json:"step"`
	Reason string `json:"reason"`
}

type Config struct {
	Cmd            string          `json:"cmd"`
	CheckStatus    CheckStatus     `json:"checkStatus"`
	ClientFeatures map[string]bool `json:"clientFeatures"`
}

func CheckImage(payload *ImageCheckPayload) (*ImageCheckResult, error) {
	ret := &ImageCheckResult{}
	resp, d, errs := instance.Request(api.MethodPost, "/images/check", &api.RequestOptions{
		Body: payload,
	}).EndStruct(ret)
	return ret, checkError(resp, d, errs)
}

package service

import (
	"fmt"
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

type ImagePayload struct {
	ImageRepository string   `json:"imageRepository" header:"imageRepository"`
	Version         *string  `json:"version" header:"version"`
	Tags            []string `json:"tags"`
}

type ImageResult struct {
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	ID              uint64           `json:"id" header:"Id"`
	ProtocolKey     string           `json:"protocolKey"`
	ImageRepository string           `json:"imageRepository" header:"ImageRepository"`
	Version         string           `json:"version" header:"Version"`
	Client          string           `json:"client"`
	Active          bool             `json:"active" header:"Active"`
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

func CheckImage(payload *ImagePayload) (*ImageResult, error) {
	ret := &ImageResult{}
	resp, d, errs := instance.Request(api.MethodPost, "/images/check", &api.RequestOptions{
		Body: payload,
	}).EndStruct(ret)
	return ret, checkError(resp, d, errs)
}

func GetImage(wsID uint64, networkID string) ([]*ImageResult, error) {
	var ret []*ImageResult
	path := fmt.Sprintf("/partners/%d/networks/%s/images", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&ret)
	return ret, checkError(resp, d, errs)
}

func UpsertImage(wsID uint64, networkID string, payload *ImagePayload) error {
	path := fmt.Sprintf("/partners/%d/networks/%s/images", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: payload,
	}).End()
	return checkError(resp, []byte(d), errs)
}

func DeleteImage(wsID uint64, networkID string, imageId uint64) error {
	if imageId > 0 {
		path := fmt.Sprintf("/partners/%d/networks/%s/images/%d", wsID, networkID, imageId)
		resp, d, errs := instance.Request(api.MethodDelete, path, nil).End()
		return checkError(resp, []byte(d), errs)
	}
	return nil
}

func ActivateImage(wsID uint64, networkID string, imageId uint64) error {
	path := fmt.Sprintf("/partners/%d/networks/%s/images/%d/activate", wsID, networkID, imageId)
	resp, d, errs := instance.Request(api.MethodPut, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func DeactivateImage(wsID uint64, networkID string, imageId uint64) error {
	path := fmt.Sprintf("/partners/%d/networks/%s/images/%d/deactivate", wsID, networkID, imageId)
	resp, d, errs := instance.Request(api.MethodPut, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

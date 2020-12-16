package service

import (
	"fmt"
	"time"

	"github.com/OnFinality-io/onf-cli/pkg/api"
)

type NetworkSpecMetadata struct {
	ChainSpec string `json:"chainspec"`
}

type NetworkSpec struct {
	Key             string              `json:"key" header:"Key"`
	Name            string              `json:"name" `
	DisplayName     string              `json:"displayName" header:"Name"`
	ProtocolKey     string              `json:"protocolKey" header:"Protocol"`
	IsSystem        bool                `json:"isSystem" header:"System"`
	ImageRepository string              `json:"imageRepository" header:"Image"`
	WorkspaceID     uint64              `json:"workspaceId,string"`
	Status          string              `json:"status"`
	Metadata        NetworkSpecMetadata `json:"metadata"`
}

type CreateNetworkSpecPayload struct {
	Name            string         `json:"name"`
	DisplayName     string         `json:"displayName"`
	Protocol        string         `json:"protocol"`
	ImageRepository string         `json:"imageRepository"`
	Metadata        CreateMetadata `json:"metadata"`
}
type CreateMetadata struct {
	Chainspec    string `json:"chainspec"`
	ImageVersion string `json:"imageVersion"`
}

type NetworkSpecEntity struct {
	Key             string    `json:"key" header:"Key"`
	Name            string    `json:"name" header:"Name"`
	DisplayName     string    `json:"displayName" `
	ProtocolKey     string    `json:"protocolKey" `
	IsSystem        bool      `json:"isSystem" header:"IsSystem"`
	ImageRepository string    `json:"imageRepository" header:"Image"`
	WorkspaceID     string    `json:"workspaceId" `
	Metadata        Metadata  `json:"metadata" header:"Metadata"`
	Status          string    `json:"status" header:"Status"`
	CreatedAt       time.Time `json:"createdAt" `
	UpdatedAt       time.Time `json:"updatedAt" `
}
type Metadata struct {
	Chainspec string `json:"chainspec" header:"chainspec"`
}

func GetNetworkSpecs(wsID int64) ([]NetworkSpec, error) {
	var specs []NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

func CreateNetworkSpecs(wsID int64, payload *CreateNetworkSpecPayload) (*NetworkSpecEntity, error) {
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	node := &NetworkSpecEntity{}
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: payload,
	}).EndStruct(node)
	return node, checkError(resp, d, errs)
}
func DeleteNetworkSpecs(wsID int64, networkID string) error {
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodDelete, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func GetNetworkSpec(wsID int64, networkID string) (*NetworkSpec, error) {
	var specs *NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

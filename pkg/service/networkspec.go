package service

import (
	"fmt"
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

func GetNetworkSpecs(wsID int64) ([]NetworkSpec, error) {
	var specs []NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

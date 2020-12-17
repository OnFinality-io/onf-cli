package service

import (
	"encoding/json"
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

type GenerateChainSpecPayload struct {
	ImageVersion string   `json:"imageVersion" `
	CliArgs      []string `json:"cliArgs" `
}
type GenerateChainSpecResult struct {
	TaskId     string `json:"taskId"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

type BootstrapChainSpecPayload struct {
	Bootnode     BootstrapChainSpecNode `json:"bootnode"`
	Storage      string                 `json:"storage"`
	NodeSpecKey  string                 `json:"nodeSpecKey"`
	UseAPIKey    bool                   `json:"useApiKey"`
	ImageVersion string                 `json:"imageVersion"`
}
type BootstrapChainSpecMetadata struct {
}
type BootstrapChainSpecNode struct {
	Cluster  string                     `json:"cluster"`
	NodeName string                     `json:"nodeName"`
	Metadata BootstrapChainSpecMetadata `json:"metadata"`
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

func GenerateChainSpec(wsID int64, networkID string, payload *GenerateChainSpecPayload) (*GenerateChainSpecResult, error) {
	var result *GenerateChainSpecResult
	path := fmt.Sprintf("/workspaces/%d/private-chains/%s/chainSpec/generate", wsID, networkID)

	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: payload,
	}).EndBytes()
	err := json.Unmarshal(d, &result)
	if err != nil {
		if len(d) > 0 {
			result = &GenerateChainSpecResult{TaskId: string(d)}
		}
	}
	return result, checkError(resp, d, errs)
}
func BootstrapChainSpec(wsID int64, networkID string, payload *BootstrapChainSpecPayload) (*NetworkSpecEntity, error) {
	path := fmt.Sprintf("/workspaces/%d/private-chains/%s/bootstrap", wsID, networkID)
	node := &NetworkSpecEntity{}
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: payload,
	}).EndStruct(node)
	return node, checkError(resp, d, errs)
}

func UploadChainSpec(wsID int64, networkID string, files []string) error {
	path := fmt.Sprintf("/workspaces/%d/private-chains/%s/chainSpec/upload", wsID, networkID)
	// b, _ := ioutil.ReadFile(file)
	req := instance.Upload(path, nil)
	for _, file := range files {
		req.SendFile(file, "files")
	}

	resp, d, errs := req.EndBytes()
	fmt.Println(string(d))
	return checkError(resp, d, errs)
}

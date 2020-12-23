package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/OnFinality-io/onf-cli/pkg/api"
)

type BootNode struct {
	NodeID  *string `json:"nodeId,omitempty"`
	Address *string `json:"address,omitempty"`
}

type NetworkSpecMetadata struct {
	ChainSpec    *string    `json:"chainspec,omitempty"`
	ImageVersion *string    `json:"imageVersion,omitempty"`
	VersionList  []string   `json:"versionList,omitempty"`
	BootNodes    []BootNode `json:"bootnodes,omitempty"`
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
	Metadata        NetworkSpecMetadata `json:"metadata"  header:"Status"`
	CreatedAt       time.Time           `json:"createdAt" `
	UpdatedAt       time.Time           `json:"updatedAt" `
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

func GetNetworkSpecs(wsID uint64) ([]NetworkSpec, error) {
	var specs []NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

func CreateNetworkSpecs(wsID uint64, payload *CreateNetworkSpecPayload) (*NetworkSpec, error) {
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	spec := &NetworkSpec{}
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: payload,
	}).EndStruct(spec)
	return spec, checkError(resp, d, errs)
}

func DeleteNetworkSpecs(wsID uint64, networkID string) error {
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodDelete, path, nil).EndBytes()
	return checkError(resp, d, errs)
}

func GetNetworkSpec(wsID uint64, networkID string) (*NetworkSpec, error) {
	var specs *NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

func GenerateChainSpec(wsID uint64, networkID string, payload *GenerateChainSpecPayload) (*GenerateChainSpecResult, error) {
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
func BootstrapChainSpec(wsID int64, networkID string, payload *BootstrapChainSpecPayload) (*NetworkSpec, error) {
	path := fmt.Sprintf("/workspaces/%d/private-chains/%s/bootstrap", wsID, networkID)
	node := &NetworkSpec{}
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: payload,
	}).EndStruct(node)
	return node, checkError(resp, d, errs)
}

type UploadResult struct {
	Success bool `json:"success"`
}

func UploadChainSpec(wsID uint64, networkID string, files []string) (*UploadResult, error) {
	path := fmt.Sprintf("/workspaces/%d/private-chains/%s/chainSpec/upload", wsID, networkID)
	req := instance.Upload(path, &api.RequestOptions{Files: map[string]string{
		"chainspec.json": files[0],
	}})
	req.TargetType = req.ForceType

	r, err := req.MakeRequest()
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(req.Url)
	signature := instance.GetSign(r.Method, u.RequestURI(), r.Header)
	r.Header.Set("authorization", fmt.Sprintf("ONF %s:%s", instance.AccessKey, signature))

	if req.Debug {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(dump))
	}

	// Send request
	resp, err := req.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// Reset resp.Body so it can be use again
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := checkError(resp, body, []error{err}); err != nil {
		return nil, err
	}
	uploadRet := &UploadResult{}
	err = json.Unmarshal(body, uploadRet)
	if err != nil {
		return nil, err
	}
	return uploadRet, nil
}

func UpdateNetworkSpecMetadata(wsID uint64, networkID string, metadata *NetworkSpecMetadata) error {
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s/metadata", wsID, networkID)
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: metadata,
	}).EndBytes()
	fmt.Println(string(d))
	return checkError(resp, d, errs)
}

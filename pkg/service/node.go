package service

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/api"
)

type NodeMetadata struct {
	Labels               map[string]string `json:"labels"`
	LastStorageUpdatedAt int64             `json:"lastStorageUpdatedAt"`
	Client               string            `json:"client"`
}

type Endpoints struct {
	RPC         string `json:"rpc"`
	WS          string `json:"ws"`
	P2pInternal string `json:"p2p-internal"`
}

type Node struct {
	ID                 int64        `json:"id,string"`
	Name               string       `json:"name"`
	NetworkSpecKey     string       `json:"networkSpecKey"`
	WorkspaceID        int64        `json:"workspaceId,string"`
	OwnerID            int64        `json:"ownerId,string"`
	NodeType           string       `json:"nodeType"`
	NodeSpec           string       `json:"nodeSpec"`
	CPU                string       `json:"cpu"`
	Ram                string       `json:"ram"`
	NodeSpecMultiplier float32      `json:"nodeSpecMultiplier"`
	Storage            string       `json:"storage"`
	StorageType        string       `json:"storageType"`
	Image              string       `json:"image"`
	ClusterHash        string       `json:"clusterHash"`
	Status             string       `json:"status"`
	Metadata           NodeMetadata `json:"metadata"`
	NetworkSpec        *NetworkSpec `json:"networkSpec"`
	Endpoints          *Endpoints   `json:"endpoints"`
	AvailableVersions  []string     `json:"availableVersions"`
	HasUpgrade         bool         `json:"hasUpgrade"`
}

type CreateNodeMetadata struct {
	NodeKey    *string           `json:"nodeKey"`
	SkipName   bool              `json:"skipName"`
	ExtraArgs  []string          `json:"extraArgs"`
	Client     *string           `json:"client"`
	RpcMethods *string           `json:"rpcMethods"`
	Labels     map[string]string `json:"labels"`
}

type CreateNodePayload struct {
	NetworkSpecKey string              `json:"networkSpecKey"`
	NodeSpecKey    *string             `json:"nodeSpecKey"`
	NodeSpec       *string             `json:"nodeSpec"`
	NodeType       string              `json:"nodeType"`
	NodeName       string              `json:"nodeName"`
	ClusterHash    string              `json:"clusterKey"`
	Storage        *string             `json:"storage"`
	InitFromBackup bool                `json:"initFromBackup"`
	UseApiKey      bool                `json:"useApiKey"`
	ImageVersion   *string             `json:"imageVersion"`
	Client         *string             `json:"client"`
	PublicPort     bool                `json:"publicPort"`
	Metadata       *CreateNodeMetadata `json:"metadata"`
}

func GetNodeList(wsID int64) ([]Node, error) {
	var nodes []Node
	path := fmt.Sprintf("/workspaces/%d/nodes", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&nodes)
	return nodes, checkError(resp, d, errs)
}

func GetNodeDetail(wsID, nodeID int64) (*Node, error) {
	node := Node{}
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&node)
	return &node, checkError(resp, d, errs)
}

func StopNode(wsID, nodeID int64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/stop", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPost, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func ResumeNode(wsID, nodeID int64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/resume", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPost, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func TerminateNode(wsID, nodeID int64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodDelete, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func UpdateNode(wsID, nodeID int64, data interface{}) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPut, path, &api.RequestOptions{
		Body: data,
	}).End()
	return checkError(resp, []byte(d), errs)
}

func CreateNode(wsID int64, data *CreateNodePayload) (*Node, error) {
	path := fmt.Sprintf("/workspaces/%d/nodes", wsID)
	node := &Node{}
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: data,
	}).EndStruct(node)
	return node, checkError(resp, d, errs)
}

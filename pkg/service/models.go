package service

import (
	"github.com/OnFinality-io/onf-cli/pkg/models"
)

type Endpoints struct {
	RPC         string `json:"rpc"`
	WS          string `json:"ws"`
	P2pInternal string `json:"p2p-internal"`
	P2p         string `json:"p2p"`
	Metrics     string `json:"metrics"`
}

type NodeItem struct {
	ID             uint64 `json:"id,string" header:"ID"`
	Name           string `json:"name" header:"Name"`
	NetworkSpecKey string `json:"networkSpecKey" header:"Network"`
	ClusterHash    string `json:"clusterHash" header:"Cluster"`
	Status         string `json:"status" header:"Status"`
	Image          string `json:"image" header:"image"`
}

type Node struct {
	ID                 uint64       `json:"id,string" header:"ID"`
	Name               string       `json:"name" header:"Name"`
	NetworkSpecKey     string       `json:"networkSpecKey" header:"Network"`
	WorkspaceID        uint64       `json:"workspaceId,string"`
	OwnerID            uint64       `json:"ownerId,string"`
	NodeType           string       `json:"nodeType"`
	NodeSpec           string       `json:"nodeSpec"`
	CPU                string       `json:"cpu"`
	Ram                string       `json:"ram"`
	NodeSpecMultiplier float32      `json:"nodeSpecMultiplier"`
	Storage            string       `json:"storage"`
	StorageType        string       `json:"storageType"`
	Image              string       `json:"image"`
	ClusterHash        string       `json:"clusterHash" header:"Cluster"`
	Status             string       `json:"status" header:"Status"`
	Metadata           NodeMetadata `json:"metadata"`
	NetworkSpec        *NetworkSpec `json:"networkSpec"`
	Endpoints          *Endpoints   `json:"endpoints"`
	AvailableVersions  []string     `json:"availableVersions"`
	HasUpgrade         bool         `json:"hasUpgrade"`
}

type NodeMetadata struct {
	NodeKey    *string           `json:"nodeKey,omitempty"`
	SkipName   bool              `json:"skipName"`
	ExtraArgs  []string          `json:"extraArgs,omitempty"`
	Client     *string           `json:"client,omitempty"`
	RpcMethods *string           `json:"rpcMethods,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type ExtraArgs map[string]*[]string

type NodeLaunchConfig struct {
	Vars      []*Vars      `json:"vars"`
	ExtraArgs *ExtraArgs   `json:"extraArgs"`
	ExtraEnvs []*ExtraEnvs `json:"extraEnvs"`
}
type Value struct {
	ValueType models.VarValueType `json:"valueType"`
	Payload   interface{}         `json:"payload"`
}
type Vars struct {
	Key   *string `json:"key"`
	Value *Value  `json:"value"`
}

type ExtraEnvs struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

type CreateNodePayload struct {
	NetworkSpecKey string            `json:"networkSpecKey"`
	NodeSpecKey    *string           `json:"nodeSpecKey"`
	NodeSpec       *NodeSpec         `json:"nodeSpec"`
	NodeType       models.NodeType   `json:"nodeType"`
	NodeName       string            `json:"nodeName"`
	ClusterHash    string            `json:"clusterKey"`
	Storage        *string           `json:"storage"`
	InitFromBackup bool              `json:"initFromBackup"`
	UseApiKey      bool              `json:"useApiKey"`
	ImageVersion   *string           `json:"imageVersion"`
	Client         *string           `json:"client"`
	PublicPort     bool              `json:"publicPort"`
	Metadata       *NodeMetadata     `json:"metadata"`
	Config         *NodeLaunchConfig `json:"config"`
}

type UpdateNodePayload struct {
	NodeSpecKey  *string           `json:"nodeSpecKey"`
	NodeSpec     *NodeSpec         `json:"nodeSpec"`
	NodeType     *string           `json:"nodeType"`
	NodeName     *string           `json:"nodeName"`
	ImageVersion *string           `json:"imageVersion"`
	Metadata     *NodeMetadata     `json:"metadata"`
	Config       *NodeLaunchConfig `json:"config"`
}

type UpdateNodeImagePayload struct {
	ImageVersion *string `json:"imageVersion"`
}

type NodeSpec struct {
	Key        string `json:"key"`
	Multiplier int    `json:"multiplier"`
}

type NodeStatus struct {
	Status string `json:"status"`
}

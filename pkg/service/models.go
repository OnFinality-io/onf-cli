package service

type Endpoints struct {
	RPC         string `json:"rpc"`
	WS          string `json:"ws"`
	P2pInternal string `json:"p2p-internal"`
}
type NodeListItem struct {
	ID             uint64 `json:"id,string" header:"ID"`
	Name           string `json:"name" header:"Name"`
	NetworkSpecKey string `json:"networkSpecKey" header:"Network"`
	ClusterHash    string `json:"clusterHash" header:"Cluster"`
	Status         string `json:"status" header:"Status"`
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
	NodeKey    *string           `json:"nodeKey"`
	SkipName   bool              `json:"skipName"`
	ExtraArgs  []string          `json:"extraArgs"`
	Client     *string           `json:"client"`
	RpcMethods *string           `json:"rpcMethods"`
	Labels     map[string]string `json:"labels"`
}

type CreateNodePayload struct {
	NetworkSpecKey string        `json:"networkSpecKey"`
	NodeSpecKey    *string       `json:"nodeSpecKey"`
	NodeSpec       *string       `json:"nodeSpec"`
	NodeType       string        `json:"nodeType"`
	NodeName       string        `json:"nodeName"`
	ClusterHash    string        `json:"clusterKey"`
	Storage        *string       `json:"storage"`
	InitFromBackup bool          `json:"initFromBackup"`
	UseApiKey      bool          `json:"useApiKey"`
	ImageVersion   *string       `json:"imageVersion"`
	Client         *string       `json:"client"`
	PublicPort     bool          `json:"publicPort"`
	Metadata       *NodeMetadata `json:"metadata"`
}

type UpdateNodePayload struct {
	NodeSpecKey  *string       `json:"nodeSpecKey"`
	NodeSpec     *string       `json:"nodeSpec"`
	NodeType     *string       `json:"nodeType"`
	NodeName     *string       `json:"nodeName"`
	ImageVersion *string       `json:"imageVersion"`
	Metadata     *NodeMetadata `json:"metadata"`
}

type NodeStatus struct {
	Status string `json:"status"`
}

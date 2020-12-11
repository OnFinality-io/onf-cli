package service

type NetworkSpecMetadata struct {
	ChainSpec string `json:"chainspec"`
}

type NetworkSpec struct {
	Key             string              `json:"key"`
	Name            string              `json:"name"`
	DisplayName     string              `json:"displayName"`
	ProtocolKey     string              `json:"protocolKey"`
	IsSystem        bool                `json:"isSystem"`
	ImageRepository string              `json:"imageRepository"`
	WorkspaceID     uint64              `json:"workspaceId,string"`
	Status          string              `json:"status"`
	Metadata        NetworkSpecMetadata `json:"metadata"`
}

package network

import "github.com/OnFinality-io/onf-cli/pkg/service"

type Bootstrap struct {
	NetworkSpec NetworkSpec `json:"networkSpec"`
	Validator   Validator   `json:"validator"`
	BootNode    BootNode    `json:"bootNode"`
}

type NetworkSpec struct {
	Config    service.CreateNetworkSpecPayload `json:"config"`
	ChainSpec string                           `json:"chainSpec"`
}
type Node struct {
	NodeName       string `json:"nodeName"`
	ClusterKey     string `json:"clusterKey"`
	NodeSpecKey    string `json:"nodeSpecKey"`
	InitFromBackup bool   `json:"initFromBackup"`
	PublicPort     bool   `json:"publicPort"`
	UseAPIKey      bool   `json:"useApiKey"`
}
type Gran struct {
	Phrase    string `json:"phrase"`
	PublicKey string `json:"publicKey"`
}
type Arua struct {
	Phrase    string `json:"phrase"`
	PublicKey string `json:"publicKey"`
}

type SessionsKey struct {
	Gran Gran `json:"gran"`
	Arua Arua `json:"arua"`
}
type Validator struct {
	Count       int           `json:"count"`
	Node        Node          `json:"node"`
	SessionsKey []SessionsKey `json:"sessionsKey"`
}
type BootNode struct {
	Count int  `json:"count"`
	Node  Node `json:"node"`
}

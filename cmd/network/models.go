package network

import "github.com/OnFinality-io/onf-cli/pkg/service"

type Bootstrap struct {
	NetworkSpec CfgNetworkSpec `json:"networkSpec"`
	Validator   CfgValidator   `json:"validator"`
	BootNode    CfgBootNode    `json:"bootNode"`
}

type CfgNetworkSpec struct {
	Config    service.CreateNetworkSpecPayload `json:"config"`
	ChainSpec string                           `json:"chainSpec"`
}

type CfgValidator struct {
	Count       int                        `json:"count"`
	Node        *service.CreateNodePayload `json:"node"`
	SessionsKey [][]service.SessionKey     `json:"sessionsKey"`
}

type CfgBootNode struct {
	Count int                        `json:"count"`
	Node  *service.CreateNodePayload `json:"node"`
}

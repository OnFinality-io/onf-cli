package network

import (
	"github.com/OnFinality-io/onf-cli/pkg/service"
)

type CfgBootstrap struct {
	NetworkSpec *service.CreateNetworkSpecPayload `json:"networkSpec"`
	Validator   CfgValidator   `json:"validator"`
	BootNode    CfgBootNode    `json:"bootNode"`
}

type CfgValidator struct {
	Count       int                       `json:"count"`
	Node        service.CreateNodePayload `json:"node"`
	SessionsKey [][]service.SessionKey    `json:"sessionsKey"`
}

type CfgBootNode struct {
	Count int                       `json:"count"`
	Node  service.CreateNodePayload `json:"node"`
}
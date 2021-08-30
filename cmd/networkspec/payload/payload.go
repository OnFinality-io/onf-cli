package payload

import "github.com/OnFinality-io/onf-cli/pkg/models"

type ConfigPayload struct {
	NodeTypes map[models.NodeType]*ConfigRule `json:"nodeTypes"`
}

type ConfigRule struct {
	Args []*ArgPayload `json:"args"`
	Envs []*EnvPayload `json:"envs"`
}

type ArgPayload struct {
	Key    *string        `json:"key"`
	Value  *string        `json:"value"`
	File   *string        `json:"file"`
	Action *models.Action `json:"action"`
}

type EnvPayload struct {
	Key   *string `json:"key"`
	Value *string `json:"value"`
}

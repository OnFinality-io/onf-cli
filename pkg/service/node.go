package service

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/api"
)

func GetNodeList(wsID uint64) ([]NodeItem, error) {
	var nodes []NodeItem
	path := fmt.Sprintf("/workspaces/%d/nodes", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&nodes)
	return nodes, checkError(resp, d, errs)
}

func GetNodeDetail(wsID, nodeID uint64) (*Node, error) {
	node := Node{}
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&node)
	return &node, checkError(resp, d, errs)
}

func StopNode(wsID, nodeID uint64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/stop", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPost, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func ResumeNode(wsID, nodeID uint64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/resume", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPost, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func RestartNode(wsID, nodeID uint64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/restart", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPost, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func TerminateNode(wsID, nodeID uint64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodDelete, path, nil).End()
	return checkError(resp, []byte(d), errs)
}

func UpdateNode(wsID, nodeID uint64, data *UpdateNodePayload) error {
	//transformSection(data.Config)
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/update", wsID, nodeID)
	resp, d, errs := instance.Ver2().Request(api.MethodPost, path, &api.RequestOptions{
		Body: data,
	}).End()
	return checkError(resp, []byte(d), errs)
}

func UpdateImage(wsID, nodeID uint64, data *UpdateNodeImagePayload) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/update-image", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: data,
	}).End()
	return checkError(resp, []byte(d), errs)
}

func CreateNode(wsID uint64, data *CreateNodePayload) (*Node, error) {
	//transformSection(data.Config)
	path := fmt.Sprintf("/workspaces/%d/nodes", wsID)
	node := &Node{}
	resp, d, errs := instance.Ver2().Request(api.MethodPost, path, &api.RequestOptions{
		Body: data,
	}).EndStruct(node)
	return node, checkError(resp, d, errs)
}

func ExpandNodeStorage(wsID, nodeID, size uint64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPut, path, &api.RequestOptions{
		Body: map[string]string{
			"storage": fmt.Sprintf("%dGi", size),
		},
	}).End()
	return checkError(resp, []byte(d), errs)
}

func GetNodeStatus(wsID, nodeID uint64) (*NodeStatus, error) {
	node := NodeStatus{}
	path := fmt.Sprintf("/workspaces/%d/nodes/%d/status", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&node)
	return &node, checkError(resp, d, errs)
}

//func transformSection(config *NodeLaunchConfig) {
//	if config != nil {
//		sectionNum := 0
//		Section := sectionNum*200 + 100
//		var extraArgs []*ExtraArgs
//		for _, arg := range config.ExtraArgs {
//			if strings.Compare(*arg.Key, "--") == 0 {
//				sectionNum++
//				Section = sectionNum*200 + 100
//				continue
//			}
//			extraArgs = append(extraArgs, &ExtraArgs{Key: arg.Key, Value: arg.Value, Section: Section})
//		}
//		config.ExtraArgs = extraArgs
//	}
//}

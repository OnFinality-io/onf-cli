package service

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/api"
)

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

func UpdateNode(wsID, nodeID int64, data *UpdateNodePayload) error {
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

func ExpandNodeStorage(wsID, nodeID, size int64) error {
	path := fmt.Sprintf("/workspaces/%d/nodes/%d", wsID, nodeID)
	resp, d, errs := instance.Request(api.MethodPut, path, &api.RequestOptions{
		Body: map[string]string{
			"storage": fmt.Sprintf("%dGi", size),
		},
	}).End()
	return checkError(resp, []byte(d), errs)
}

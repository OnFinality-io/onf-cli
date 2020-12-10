package service

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/api"
	"time"
)

type Workspace struct {
	ID              int64      `json:"id,string"`
	Name            string     `json:"name"`
	Plan            string     `json:"plan"`
	OwnerID         int64      `json:"ownerId,string"`
	BillingType     string     `json:"billingType"`
	PaymentMethodID string     `json:"paymentMethodId"`
	Active          bool       `json:"active"`
	SuspendAt       *time.Time `json:"suspendAt"`
	CreateAt        *time.Time `json:"createdAt"`
	UpdateAt        *time.Time `json:"updateAt"`
}

type Member struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type InviteLog struct {
	ID     int64  `json:"id,string"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	IsDone bool   `json:"isDone"`
}

type InviteMemberPayload struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func GetWorkspaceList() ([]Workspace, error) {
	var list []Workspace
	resp, d, errs := instance.Request(api.MethodGet, "/workspaces", nil).EndStruct(&list)
	return list, checkError(resp, d, errs)
}

func GetMembers(wsID int64) ([]Member, error) {
	var members []Member
	path := fmt.Sprintf("/workspaces/%d/members", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&members)
	return members, checkError(resp, d, errs)
}

func GetInvitations(wsID int64) ([]InviteLog, error) {
	var logs []InviteLog
	path := fmt.Sprintf("/workspaces/%d/invitations", wsID)
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&logs)
	return logs, checkError(resp, d, errs)
}

func InviteMember(wsID int64, data *InviteMemberPayload) error {
	path := fmt.Sprintf("/workspaces/%d/invite", wsID)
	resp, d, errs := instance.Request(api.MethodPost, path, &api.RequestOptions{
		Body: data,
	}).End()
	return checkError(resp, []byte(d), errs)
}

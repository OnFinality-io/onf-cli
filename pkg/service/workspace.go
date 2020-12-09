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

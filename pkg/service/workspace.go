package service

import (
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

func GetWorkspaceList() ([]Workspace, []error) {
	var list []Workspace
	_, _, err := instance.Request(api.MethodGet, "/workspaces", nil).EndStruct(&list)
	return list, err
}

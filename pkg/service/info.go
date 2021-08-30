package service

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/models"
	"time"

	"github.com/OnFinality-io/onf-cli/pkg/api"
)

type Info struct {
	Clusters  []Clusters  `json:"clusters" header:"clusters"`
	NodeSpecs []NodeSpecs `json:"nodeSpecs" header:"nodeSpecs"`
	Protocols []Protocols `json:"protocols" header:"protocols"`
}
type Clusters struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name" header:"name"`
	Hash      string    `json:"hash" header:"hash"`
	Cloud     string    `json:"cloud" header:"cloud"`
	Region    string    `json:"region" header:"region"`
	Index     int       `json:"index" header:"index"`
	Active    bool      `json:"active" header:"active"`
}
type NodeSpecPrice struct {
	Available bool   `json:"available" header:"available"`
	Price     string `json:"price" header:"price"`
}
type NodeSpecs struct {
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Key           string        `json:"key" header:"key"`
	Name          string        `json:"name" header:"name"`
	CPU           string        `json:"cpu" header:"cpu"`
	Memory        string        `json:"memory" header:"memory"`
	Active        bool          `json:"active" header:"active"`
	Priority      int           `json:"priority" header:"priority"`
	Protocol      string        `json:"protocol" header:"protocol"`
	Network       string        `json:"network" header:"network"`
	Price         NodeSpecPrice `json:"price" header:"price"`
	MaxMultiplier int           `json:"maxMultiplier" header:"maxMultiplier"`
}
type Protocols struct {
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Key               string    `json:"key" header:"key"`
	Name              string    `json:"name" header:"name"`
	Derivable         bool      `json:"derivable" header:"derivable"`
	ImageRepositories []string  `json:"imageRepositories" header:"images"`
	Metadata          *Metadata `json:"metadata" header:"metadata"`
}
type ArgSectionItems []*ArgSectionItem

func (ast ArgSectionItems) GetItem(key string) *ArgSectionItem {
	for _, item := range ast {
		if item.ExtraArgKey != nil && *item.ExtraArgKey == key {
			return item
		}
	}
	return nil
}

type ArgumentSections struct {
	Sections  ArgSectionItems `json:"sections"`
	Delimiter *string         `json:"delimiter"`
}
type ArgSectionItem struct {
	Title       *string `json:"title"`
	ExtraArgKey *string `json:"extraArgKey"`
	Index       *int    `json:"index"`
}
type Metadata struct {
	Rules            map[models.NodeType]interface{} `json:"rules" header:"rules"`
	ArgumentSections *ArgumentSections               `json:"argumentSections" `
}

type NodeRecommendation struct {
	Network            string `json:"network" header:"network"`
	NodeSpec           string `json:"nodeSpec" header:"nodeSpec"`
	NodeSpecMultiplier int    `json:"nodeSpecMultiplier" header:"nodeSpecMultiplier"`
	StorageSize        int    `json:"storageSize" header:"storageSize"`
	ImageVersion       string `json:"imageVersion" header:"imageVersion"`
	Client             string `json:"client" header:"client"`
}

func GetInfo() (Info, error) {
	var result Info
	resp, d, errs := instance.Request(api.MethodGet, "/info", nil).EndStruct(&result)
	return result, checkError(resp, d, errs)
}

func ListImageVersions(image string) ([]string, error) {
	path := fmt.Sprintf("/info/images/%s/versions", image)
	var result []string
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&result)
	return result, checkError(resp, d, errs)
}

func NodeRecommends(network string) (NodeRecommendation, error) {
	path := fmt.Sprintf("/info/node-recommendation/%s", network)
	var result NodeRecommendation
	resp, d, errs := instance.Request(api.MethodGet, path, nil).EndStruct(&result)
	return result, checkError(resp, d, errs)
}

func GetArgumentSectionsByProtocol(protocolKey string) (*ArgumentSections, error) {
	info, err := GetInfo()
	if err != nil {
		return nil, err
	}
	for _, protocol := range info.Protocols {
		if protocol.Key == protocolKey {
			return protocol.Metadata.ArgumentSections, nil
		}
	}
	return nil, fmt.Errorf("%s is not supported.", protocolKey)
}

func GetSupportedNodeTypes(protocolKey string) ([]models.NodeType, error) {
	info, err := GetInfo()
	if err != nil {
		return nil, err
	}
	var nodeTypes []models.NodeType
	for _, protocol := range info.Protocols {
		if protocol.Key == protocolKey {
			for key, _ := range protocol.Metadata.Rules {
				nodeTypes = append(nodeTypes, key)
			}
			break
		}
	}
	return nodeTypes, nil
}

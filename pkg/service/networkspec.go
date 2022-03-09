package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/OnFinality-io/onf-cli/cmd/networkspec/payload"
	"github.com/OnFinality-io/onf-cli/pkg/models"
	"github.com/OnFinality-io/onf-cli/pkg/utils"

	"github.com/OnFinality-io/onf-cli/pkg/api"
)

type BootNode struct {
	NodeID  *string `json:"nodeId,omitempty"`
	Address *string `json:"address,omitempty"`
}

type NetworkSpecMetadata struct {
	Recommend        *Recommend `json:"recommend,omitempty"`
	ChainSpec        *string    `json:"chainspec,omitempty"`
	ImageVersion     *string    `json:"imageVersion,omitempty"`
	Command          *string    `json:"command,omitempty"`
	VersionList      []string   `json:"versionList,omitempty"`
	BootNodes        []string   `json:"bootnodes,omitempty"`
	FetchRcChainspec *string    `json:"fetchRcChainspec,omitempty"`
	RcChainspec      *string    `json:"rcChainspec,omitempty"`
	RcExtraArgs      []string   `json:"rcExtraArgs,omitempty"`
	ExtraArgs        []string   `json:"extraArgs,omitempty"`
	ParachainId      *int       `json:"parachainId,omitempty"`
	Cluster          *string    `json:"cluster,omitempty"`
}

type Recommend struct {
	ImageVersion string `json:"imageVersion"`
	NodeSpec     string `json:"nodeSpec"`
	StorageSize  uint   `json:"storageSize"`
}

type SpecNodeType struct {
	Key       string     `json:"key"`
	Recommend *Recommend `json:"recommend,omitempty"`
}

type NetworkSpec struct {
	Key             string              `json:"key" header:"Key"`
	Name            string              `json:"name" `
	DisplayName     string              `json:"displayName" header:"Name"`
	ProtocolKey     string              `json:"protocolKey" header:"Protocol"`
	IsPublic        bool                `json:"isPublic" header:"Public"`
	ImageRepository string              `json:"imageRepository" header:"Image"`
	WorkspaceID     uint64              `json:"workspaceId,string"`
	Status          string              `json:"status"`
	Metadata        NetworkSpecMetadata `json:"metadata"  header:"Status"`
	CreatedAt       time.Time           `json:"createdAt" `
	UpdatedAt       time.Time           `json:"updatedAt" `
	Recommend       *Recommend          `json:"recommend,omitempty"`
	NodeTypes       []SpecNodeType      `json:"nodeTypes,omitempty"`
	Config          *models.Config      `json:"config"`
}

type NetworkSpecBackups struct {
	Key             	string              `json:"key" header:"Network Spec Key"`
	Name            	string              `json:"name"`
	DisplayName     	string              `json:"displayName" header:"Name"`
	ProtocolKey     	string              `json:"protocolKey"`
	AvailNodeTypes  	string     			`json:"availNodeTypes" header:"Node Type"`
	MinStorageSize  	uint   				`json:"minStorageSize" header:"Min Storage Size (Gb)"`
	AvailCloudRegion  	string     			`json:"availCloudRegion" header:"Cloud & Region"`
	ClusterKey 			string     			`json:"cluserKey" header:"Cluster Key"`
}

type Backups struct {
	Id				string		`json:"id"`
	NetworkSpec		string		`json:"networkSpec"`
	Protocol	    string		`json:"protocol"`
	ClusterHash		string		`json:"clusterHash"`
	StorageSize     uint		`json:"storageSize"`
	PruningMode     string		`json:"pruningMode"`
}

func (c *NetworkSpec) MergeConfig(config *models.Config) {
	if c.Config == nil {
		c.Config = config
	} else {
		for nodeType, newOperation := range config.Operations {
			// Exist
			if originalOperation, ok := c.Config.Operations[nodeType]; ok {
				originalOperation.Var.Merge(newOperation.Var)
				originalOperation.Env.Merge(newOperation.Env)
				originalOperation.Arg.Merge(newOperation.Arg)
			} else {
				c.Config.Operations[nodeType] = newOperation
			}
		}
	}

}

type CreateNetworkSpecPayload struct {
	Name            string  `json:"name"`
	DisplayName     string  `json:"displayName"`
	Protocol        string  `json:"protocol"`
	ImageRepository string  `json:"imageRepository"`
	ImageVersion    *string `json:"imageVersion,omitempty"`
	//Metadata        *NetworkSpecMetadata     `json:"metadata"`
	Config *payload.ConfigPayload `json:"config"`
}

type UpdateNetworkSpecPayload struct {
	DisplayName *string                `json:"displayName"`
	Config      *payload.ConfigPayload `json:"config"`
}

type GenerateChainSpecPayload struct {
	ImageVersion string   `json:"imageVersion" `
	CliArgs      []string `json:"cliArgs" `
}
type GenerateChainSpecResult struct {
	TaskId     string `json:"taskId"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
}

type BootstrapChainSpecPayload struct {
	Bootnode     BootstrapChainSpecNode `json:"bootnode"`
	Storage      string                 `json:"storage"`
	NodeSpecKey  string                 `json:"nodeSpecKey"`
	UseAPIKey    bool                   `json:"useApiKey"`
	ImageVersion string                 `json:"imageVersion"`
}
type BootstrapChainSpecMetadata struct {
}
type BootstrapChainSpecNode struct {
	Cluster  string                     `json:"cluster"`
	NodeName string                     `json:"nodeName"`
	Metadata BootstrapChainSpecMetadata `json:"metadata"`
}
type CreateNetworkSpecModel struct {
	Name            string               `json:"name"`
	DisplayName     string               `json:"displayName"`
	Protocol        string               `json:"protocol"`
	ImageRepository string               `json:"imageRepository"`
	Metadata        *NetworkSpecMetadata `json:"metadata"`
	Config          *models.Config       `json:"config"`
}
type UpdateNetworkSpecModel struct {
	DisplayName *string        `json:"displayName"`
	Config      *models.Config `json:"config"`
}

func GetNetworkSpecs(wsID uint64) ([]NetworkSpec, error) {
	var specs []NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	resp, d, errs := instance.Ver2().Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

func GetBackupNetworkSpecs() ([]NetworkSpecBackups, error) {
	var specs []NetworkSpecBackups
	resp, d, errs := instance.Ver2().Request(api.MethodGet, "/network-specs", nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

func GetBackups() ([]Backups, error) {
	var backups []Backups
	resp, d, errs := instance.Ver2().Request(api.MethodGet, "/backup", nil).EndStruct(&backups)
	return backups, checkError(resp, d, errs)
}

func CreateNetworkSpecs(wsID uint64, payload *CreateNetworkSpecPayload) (*NetworkSpec, error) {
	argumentSections, err := GetArgumentSectionsByProtocol(payload.Protocol)
	if err != nil {
		return nil, err
	}
	config := TransformConfig(wsID, payload.Config, argumentSections)

	chainSpec := obtainChainSpec(config)
	p := &CreateNetworkSpecModel{
		DisplayName:     payload.DisplayName,
		Name:            payload.Name,
		ImageRepository: payload.ImageRepository,
		Protocol:        payload.Protocol,
		Metadata: &NetworkSpecMetadata{
			ImageVersion: payload.ImageVersion,
			ChainSpec:    chainSpec,
		},
		Config: config,
	}
	path := fmt.Sprintf("/workspaces/%d/network-specs", wsID)
	spec := &NetworkSpec{}
	resp, d, errs := instance.Ver2().Request(api.MethodPost, path, &api.RequestOptions{
		Body: p,
	}).EndStruct(spec)
	return spec, checkError(resp, d, errs)
}

func obtainChainSpec(config *models.Config) *string {
	if config == nil {
		return nil
	}
	var chainSpec *string
	for _, operations := range config.Operations {
		for _, variable := range operations.Var {
			if strings.HasPrefix(variable.Payload.Key, "--chain") {
				if variable.Payload.Value.ValueType == models.File {
					if file, ok := variable.Payload.Value.Payload.(models.FileTypeValue); ok {
						chainSpec = utils.String(fmt.Sprintf("/chain-data/%s", file.Destination))
					}
				} else if variable.Payload.Value.ValueType == models.String {
					chainSpec = variable.Payload.Value.Payload.(*string)
				}
			}
		}
	}
	return chainSpec
}

func DeleteNetworkSpecs(wsID uint64, networkID string) error {
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Ver2().Request(api.MethodDelete, path, nil).EndBytes()
	return checkError(resp, d, errs)
}

func GetNetworkSpec(wsID uint64, networkID string) (*NetworkSpec, error) {
	var specs *NetworkSpec
	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Ver2().Request(api.MethodGet, path, nil).EndStruct(&specs)
	return specs, checkError(resp, d, errs)
}

type UploadResult struct {
	Key *string `json:"key"`
}

func UploadPrivateFile(wsID uint64, files []string) (*UploadResult, error) {
	//path := fmt.Sprintf("/private-file/upload")
	path := fmt.Sprintf("/workspaces/%d/private-file/upload", wsID)
	f := files[0]
	_, file := filepath.Split(f)
	req := instance.Upload(path, &api.RequestOptions{Files: map[string]string{
		file: f,
	}})
	req.TargetType = req.ForceType

	r, err := req.MakeRequest()
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(req.Url)
	signature := instance.GetSign(r.Method, u.RequestURI(), r.Header)
	r.Header.Set("authorization", fmt.Sprintf("ONF %s:%s", instance.AccessKey, signature))

	if req.Debug {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(dump))
	}

	// Send request
	resp, err := req.Client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	// Reset resp.Body so it can be use again
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := checkError(resp, body, []error{err}); err != nil {
		return nil, err
	}
	uploadRet := &UploadResult{}
	err = json.Unmarshal(body, uploadRet)
	if err != nil {
		return nil, err
	}
	return uploadRet, nil
}

func UpdateNetworkSpec(wsID uint64, networkID string, payload *UpdateNetworkSpecPayload) error {

	spec, err := GetNetworkSpec(wsID, networkID)
	if err != nil {
		return err
	}
	argumentSections, err := GetArgumentSectionsByProtocol(spec.ProtocolKey)
	if err != nil {
		return err
	}
	transformConfig := TransformConfig(wsID, payload.Config, argumentSections)
	spec.MergeConfig(transformConfig)
	displayName := payload.DisplayName
	if payload.DisplayName == nil {
		displayName = &spec.DisplayName
	}
	p := &UpdateNetworkSpecModel{
		DisplayName: displayName,
		Config:      spec.Config,
	}

	path := fmt.Sprintf("/workspaces/%d/network-specs/%s", wsID, networkID)
	resp, d, errs := instance.Ver2().Request(api.MethodPut, path, &api.RequestOptions{
		Body: p,
	}).EndBytes()
	return checkError(resp, d, errs)
}

type ArgPayloadWrap struct {
	Key     *string
	Value   []*string
	File    *string
	Section *int
	Action  *models.Action
}

func TransformConfig(wsID uint64, payload *payload.ConfigPayload, argumentSections *ArgumentSections) *models.Config {
	if payload == nil {
		return nil
	}

	config := &models.Config{Operations: map[models.NodeType]*models.RuleOperationCollection{}}

	for nodeType, rule := range payload.NodeTypes {
		argPayloads := make([][]*ArgPayloadWrap, 5)
		sectionIndex := 0
		for _, argPayload := range rule.Args {
			if argumentSections.Delimiter != nil && strings.Compare(*argPayload.Key, *argumentSections.Delimiter) == 0 {
				sectionIndex++
				continue
			}
			argSections := argumentSections.Sections
			sort.SliceStable(argSections, func(i, j int) bool {
				return *argSections[i].Index < *argSections[j].Index
			})
			section := argumentSections.Sections[sectionIndex].Index
			skip := false
			for _, argWrap := range argPayloads[sectionIndex] {
				if strings.Compare(*argWrap.Key, *argPayload.Key) == 0 {
					argWrap.Value = append(argWrap.Value, argPayload.Value)
					skip = true
					break
				}
			}

			if skip == false {
				var value []*string
				if argPayload.Value != nil {
					value = []*string{argPayload.Value}
				}

				argPayloads[sectionIndex] = append(argPayloads[sectionIndex], &ArgPayloadWrap{
					Key:     argPayload.Key,
					File:    argPayload.File,
					Value:   value,
					Action:  argPayload.Action,
					Section: section,
				})
			}
		}
		vars, args, envs := transformConfigs(wsID, argPayloads)
		config.Operations[nodeType] = &models.RuleOperationCollection{Var: vars, Arg: args, Env: envs}
	}

	return config
}

func transformConfigs(wsID uint64, payload [][]*ArgPayloadWrap) ([]*models.Var, []*models.Arg, []*models.Env) {
	vars := make([]*models.Var, 0)
	args := make([]*models.Arg, 0)
	envs := make([]*models.Env, 0)
	for i, payload := range payload {
		for _, argPayload := range payload {
			if argPayload.Key == nil {
				continue
			}
			key := fmt.Sprintf("%s_%d", *argPayload.Key, i)
			file := argPayload.File
			value := argPayload.Value
			variable := &models.Var{
				RuleOperation: models.RuleOperation{Action: models.UPSERT},
				Payload: &models.VarModel{
					Key:      key,
					Options:  models.Options{Overwritable: true},
					Category: models.VAR,
					Value:    &models.VarValue{},
				}}
			if argPayload.Action != nil {
				variable.Action = *argPayload.Action
			}

			if file != nil {
				uploadRet, err := UploadPrivateFile(wsID, []string{*file})
				if err != nil {
					fmt.Println(err.Error())
				}
				_, fileName := filepath.Split(*file)
				fileBucket := uploadRet.Key
				variable.Payload.Value.ValueType = models.File
				variable.Payload.Value.Payload = models.FileTypeValue{Source: fileBucket, Destination: fileName}
			} else if value != nil && len(value) > 0 {
				if len(value) > 1 {
					variable.Payload.Value.ValueType = models.StringArray
					variable.Payload.Value.Payload = value
				} else {
					variable.Payload.Value.ValueType = models.String
					variable.Payload.Value.Payload = value[0]
				}
			} else {
				variable.Payload.Value.ValueType = models.Empty
			}
			vars = append(vars, variable)
			arg := &models.Arg{
				RuleOperation: models.RuleOperation{Action: models.UPSERT},
				Payload: &models.ArgModel{
					Key:      *argPayload.Key,
					Category: models.ARG,
					Options:  models.Options{Overwritable: true},
					Value: &models.ValueModel{
						InputType: models.Variable,
						Payload:   fmt.Sprintf("var.%s", key),
					},
					Section: argPayload.Section,
				}}
			args = append(args, arg)
		}
	}

	return vars, args, envs

}

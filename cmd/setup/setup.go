package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OnFinality-io/onf-cli/pkg/utils"

	"gopkg.in/ini.v1"
)

const (
	defaultOnfDir              = ".onf"
	defaultCredentialFileName  = "credentials"
	defaultCredentialsFileType = "ini"
)

type CredentialConfig struct {
	Section    string
	Credential *Credential
}

type Credential struct {
	AccessKey   string `ini:"onf_access_key" header:"onf_access_key"`
	SecretKey   string `ini:"onf_secret_key" header:"onf_secret_key"`
	WorkspaceID uint64 `ini:"default_workspace" header:"default_workspace"`
}

type CredentialFile struct {
	CredentialFileName  string
	CredentialsFileType string
}

func (cf *CredentialFile) GetCredentialPath(onfHomeDir string) (onfCredentialFile string) {
	if cf.CredentialFileName == "" {
		cf.CredentialFileName = defaultCredentialFileName
	}
	if cf.CredentialsFileType == "" {
		cf.CredentialsFileType = defaultCredentialsFileType
	}
	onfCredentialFile = filepath.Join(onfHomeDir, cf.CredentialFileName)
	return onfCredentialFile
}

func (cf *CredentialFile) CreateCredentialFileAt(onfHomeDir string) (onfCredentialFile string, err error) {
	onfCredentialFile = cf.GetCredentialPath(onfHomeDir)
	if exist, err := utils.Exists(onfCredentialFile); err == nil && exist {
	} else {
		err := utils.Touch(onfCredentialFile)
		if err != nil {
			return onfCredentialFile, err
		}
	}
	return onfCredentialFile, nil
}
func (cf *CredentialFile) IsExistAtOnfAtHomeDir() bool {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		d := filepath.Join(homeDir, defaultOnfDir)
		onfCredentialFile := cf.GetCredentialPath(d)
		if exist, err := utils.Exists(onfCredentialFile); err == nil && exist {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func CreateHomeDir(sysHomeDir, defaultOnfDir string) (onfHomeDir string) {
	onfHomeDir = filepath.Join(sysHomeDir, defaultOnfDir)
	if exist, err := utils.Exists(onfHomeDir); err == nil && exist {

	} else {
		err := utils.MkdirAll(onfHomeDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Create onf home " + onfHomeDir + " success")
	}
	return onfHomeDir
}
func PersistentCredential(credential *CredentialConfig) {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		onfHomeDir := CreateHomeDir(homeDir, defaultOnfDir)
		credentialFile := &CredentialFile{}
		onfCredentialFile, err := credentialFile.CreateCredentialFileAt(onfHomeDir)
		if err == nil {
			fmt.Println("create " + onfCredentialFile + " success")
			cfg, err := ini.Load(onfCredentialFile)
			if err != nil {
				fmt.Printf("Fail to read file: %v", err)
			} else {
				if credential.Section == "" {
					cfg.ReflectFrom(credential.Credential)
				} else {
					section, err := cfg.NewSection(credential.Section)
					if err != nil {
						fmt.Printf("Fail to save file: %v", err)
					}
					section.ReflectFrom(credential.Credential)
				}
				cfg.SaveTo(onfCredentialFile)
				if err != nil {
					fmt.Printf("Fail to save file: %v", err)
				}
			}
		} else {
			fmt.Println("Fail to create onf config file at " + onfHomeDir)
		}
	}
}

func (cf *CredentialFile) GetAllSetup() []*CredentialConfig {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		d := filepath.Join(homeDir, defaultOnfDir)
		onfCredentialFile := cf.GetCredentialPath(d)
		cfg, err := ini.Load(onfCredentialFile)
		if err != nil {
			fmt.Printf("Fail to read file: %v", err)
		} else {
			configArray := make([]*CredentialConfig, 0)
			for _, section := range cfg.Sections() {
				credential := &Credential{}
				section.MapTo(credential)
				config := &CredentialConfig{Credential: credential, Section: section.Name()}
				configArray = append(configArray, config)
			}
			return configArray
		}

	}
	return nil
}

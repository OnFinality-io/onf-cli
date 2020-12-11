package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/OnFinality-io/onf-cli/pkg/utils"

	"gopkg.in/ini.v1"
)

const (
	defaultDir          = ".onf"
	credentials         = "credentials"
	credentialsFileType = "ini"
)

type CredentialConfig struct {
	Credential *Credential
	Section    string
}

type Credential struct {
	AccessKey string `ini:"onf_access_key"`
	SecretKey string `ini:"onf_secret_key"`
}

func CreateHomeDir(homeDir, defaultDir string) (onfHomeDir string) {
	onfHomeDir = filepath.Join(homeDir, defaultDir)
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

func CreateCredentialFile(onfHomeDir string) (onfCredentialFile string) {
	onfCredentialFile = filepath.Join(onfHomeDir, credentials+"."+credentialsFileType)
	if exist, err := utils.Exists(onfCredentialFile); err == nil && exist {
	} else {
		err := utils.Touch(onfCredentialFile)
		if err != nil {
			fmt.Println("Fail to create onf config file at " + onfHomeDir)
			return onfCredentialFile
		}
		fmt.Println("create " + onfCredentialFile + " success")
	}
	return onfCredentialFile
}

func New(credential *CredentialConfig) {
	homeDir, err := os.UserHomeDir()
	if err == nil {
		onfHomeDir := CreateHomeDir(homeDir, defaultDir)
		onfCredentialFile := CreateCredentialFile(onfHomeDir)
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

	}
}

func IsCreated(onfCredentialFile string) bool {
	if exist, err := utils.Exists(onfCredentialFile); err == nil && exist {
		return true
	} else {
		return false
	}
}

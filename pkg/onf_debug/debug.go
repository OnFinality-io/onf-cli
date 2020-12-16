package onf_debug

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/OnFinality-io/onf-cli/pkg/constant"
	"github.com/OnFinality-io/onf-cli/pkg/utils"
)

var (
	_debugOnce sync.Once
	onf_debug  *ONFDebug
)

type ONFDebug struct {
	Baseurl string `json:"baseurl"`
}

func GetDebug() *ONFDebug {
	_debugOnce.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			d := filepath.Join(homeDir, constant.DefaultOnfDir)
			debugF := filepath.Join(d, "debug")
			data, err := utils.Read(debugF)
			if err == nil {
				err = json.Unmarshal(data, &onf_debug)
			}
			// fmt.Println("debugF", onf_debug)
			if onf_debug == nil {
				onf_debug = &ONFDebug{}
			}
		}
	})
	return onf_debug
}

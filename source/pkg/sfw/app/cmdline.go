package app

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"groot/comm/conf"
	"os"
	"runtime"
)

var (
	cmd_sp_config  = pflag.StringP("config", "C", "", "config file path")
	cmd_bp_version = pflag.BoolP("version", "V", false, "show cmd_bp_version info")
	cmd_bp_genconf = pflag.BoolP("genconf", "G", false, "generate default config to <config>.default")
	cmd_bp_syncdb  = pflag.BoolP("syncdb", "D", false, "auto migrate db tables")
)

//Info contains versioning information
type AppVerInfo struct {
	GitTag       string `json:"gitTag"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

//String returns info as a human-friendly cmd_bp_version string.
func (info AppVerInfo) String() string {
	return info.GitTag
}

func GetAppVerInfo() AppVerInfo {
	return AppVerInfo{
		GitTag:       gitTag,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func GetConfigPath() string {
	if cmd_sp_config == nil || *cmd_sp_config == "" {
		return "config.json"
	}
	return *cmd_sp_config
}

func parseCmdLine() {
	pflag.Parse()
	if *cmd_bp_version {
		v := GetAppVerInfo()
		//JSON显示
		marshalled, e := json.MarshalIndent(&v, "", " ")
		if e != nil {
			fmt.Printf("%v\n", e)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		os.Exit(0)
	}
	if *cmd_bp_genconf {
		filePath := "./config.json.default"
		if cmd_sp_config != nil && *cmd_sp_config != "" {
			filePath = *cmd_sp_config + ".default"
		}
		//fmt.Println("generating config file path:", filePath)
		conf.SaveConfig(filePath)
		os.Exit(0)
	}
}

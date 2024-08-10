package gonuts

import (
	"encoding/json"
	"os"
)

type VersionData struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	GitBranch string `json:"gitBranch"`
}

var vData VersionData = VersionData{Version: "0.0.0", GitCommit: "unknown", GitBranch: "unknown"}

func Init() {
	homedir, _ := os.Getwd()
	versionFile, err := os.ReadFile(homedir + "/" + "version.json")
	if err != nil {
		L.Errorf("[version.service] failed to load version file PANIC! \n%s", err)
		L.Panic(err)
	}
	err = json.Unmarshal([]byte(versionFile), &vData)
	if err != nil {
		L.Errorf("[version.service] failed to parse version file PANIC! \n%s", err)
		L.Panic(err)
	}
}

func GetVersion() string {
	return vData.Version
}

func GetGitCommit() string {
	return vData.GitCommit
}

func GetGitBranch() string {
	return vData.GitBranch
}

func GetVersionData() VersionData {
	return vData
}

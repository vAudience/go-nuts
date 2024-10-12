package gonuts

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type VersionData struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	GitBranch string `json:"gitBranch"`
}

var vData VersionData = VersionData{Version: "0.0.0", GitCommit: "unknown", GitBranch: "unknown"}

func InitVersion() {
	homedir, _ := os.Getwd()
	versionFilePath := filepath.Join(homedir, "version.json")

	// Try to read and parse the version file
	if err := readVersionFile(versionFilePath); err != nil {
		L.Warnf("[version.service] failed to read version file: %s", err)
		// If file doesn't exist or parsing fails, try to populate data
		populateVersionData(homedir)
	}

	// Write the version file if it doesn't exist or is incomplete
	if vData.Version == "0.0.0" || vData.GitCommit == "unknown" || vData.GitBranch == "unknown" {
		if err := writeVersionFile(versionFilePath); err != nil {
			L.Errorf("[version.service] failed to write version file: %s", err)
		}
	}
}

func readVersionFile(path string) error {
	versionFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(versionFile, &vData)
}

func writeVersionFile(path string) error {
	data, err := json.MarshalIndent(vData, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func populateVersionData(dir string) {
	// Try to get git information
	vData.GitCommit = getGitInfo(dir, "rev-parse", "HEAD")
	vData.GitBranch = getGitInfo(dir, "rev-parse", "--abbrev-ref", "HEAD")

	// If version is still default, try to get it from git tags
	if vData.Version == "0.0.0" {
		vData.Version = getGitInfo(dir, "describe", "--tags", "--abbrev=0")
		if vData.Version == "" {
			vData.Version = "0.0.0"
		}
	}
}

func getGitInfo(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		L.Warnf("[version.service] failed to get git info: %s", err)
		return "unknown"
	}
	return strings.TrimSpace(string(output))
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

// ForceUpdateVersionData re-reads git data and updates the version file
func ForceUpdateVersionData() error {
	homedir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Re-populate version data from git
	populateVersionData(homedir)

	// Write updated data to version file
	versionFilePath := filepath.Join(homedir, "version.json")
	if err := writeVersionFile(versionFilePath); err != nil {
		L.Errorf("[version.service] failed to write updated version file: %s", err)
		return err
	}

	L.Infof("[version.service] Version data forcefully updated: %+v", vData)
	return nil
}

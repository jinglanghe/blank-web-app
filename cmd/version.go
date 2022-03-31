package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"sigs.k8s.io/yaml"
	"strings"
)

// These are set during build time via -ldflags
var (
	gitCommit = "N/A"
	buildDate = "N/A"
)

// VersionInfo holds the version information of the driver
type VersionInfo struct {
	GitCommit string `json:"Git Commit"`
	BuildDate string `json:"Build Date"`
	GoVersion string `json:"Go Version"`
	Compiler  string `json:"Compiler"`
	Platform  string `json:"Platform"`
}

// GetVersion returns the version information of the driver
func GetVersion() VersionInfo {
	return VersionInfo{
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetVersionYAML returns the version information of the driver
// in YAML format
func GetVersionYAML() (string, error) {
	info := GetVersion()
	marshalled, err := yaml.Marshal(&info)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(marshalled)), nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Check service version",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := GetVersionYAML()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("%v\n", data)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package buildinfo

import (
	"fmt"
	"solopg/config"
	"strings"
)

type BuildInfo struct {
	version string
	commit  string
	date    string
}

type InfoTemplate struct {
	Version string
	Commit  string
	Date    string
}

func New(params InfoTemplate) BuildInfo {
	return BuildInfo{
		version: params.Version,
		commit:  params.Commit,
		date:    params.Date,
	}
}

func (b BuildInfo) String() string {
	sb := strings.Builder{}

	env := "dev"
	if !config.IsDev() {
		env = "prod"
	}

	fmt.Fprintf(&sb, "Version: %s/%s\n", b.version, env)
	fmt.Fprintf(&sb, "Commit: %s\n", b.commit)
	fmt.Fprintf(&sb, "Date: %s\n", b.date)
	return sb.String()
}

package goversion

import (
	"go/version"
	"strings"
)

const MinStdlibAutoMaxProcsGoVersion = "1.25.0"

func Resolve(goVersion string) string {
	if !version.IsValid("go" + goVersion) {
		return MinStdlibAutoMaxProcsGoVersion
	}
	if version.Compare("go"+goVersion, "go"+MinStdlibAutoMaxProcsGoVersion) < 0 {
		return MinStdlibAutoMaxProcsGoVersion
	}
	return goVersion
}

func Toolchain(goVersion string) string {
	return "go" + Resolve(goVersion)
}

func DockerTag(goVersion string) string {
	return strings.TrimPrefix(version.Lang(Toolchain(goVersion)), "go")
}

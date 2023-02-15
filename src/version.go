package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
	"strings"
)

const currentMajorVersion = 1
const currentMinorVersion = 0
const currentPatchVersion = 0

var versionPattern = regexp.MustCompile(`\d+\.\d+(\.\d+)?`)

type VersionReader struct {
}

func (_ VersionReader) read(node *yaml.Node, _ string, model *ArchitectureModel) []Issue {
	if node == nil {
		model.Version = fmt.Sprintf("%v.%v.%v", currentMajorVersion, currentMinorVersion, currentPatchVersion)
		return []Issue{}
	}
	version, issue := toString(node, "version")
	if issue != nil {
		return []Issue{*issue}
	}
	if !versionPattern.MatchString(version) {
		return []Issue{*NodeError("Version must be a semantic version as defined by https://semver.org", node)}
	}
	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	if major > currentMajorVersion {
		return []Issue{*NodeError(fmt.Sprintf("Undefined version: %s", version), node)}
	}
	if major == currentMajorVersion {
		minor, _ := strconv.Atoi(parts[1])
		if minor > currentMinorVersion {
			return []Issue{*NodeError(fmt.Sprintf("Undefined version: %s", version), node)}
		}
		if minor == currentMinorVersion {
			if len(parts) > 2 {
				patch, _ := strconv.Atoi(parts[2])
				if patch > currentPatchVersion {
					return []Issue{*NodeError(fmt.Sprintf("Undefined version: %s", version), node)}
				}
			}
		}
	}
	model.Version = version
	return []Issue{}
}

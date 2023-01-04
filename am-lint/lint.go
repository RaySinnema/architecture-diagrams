package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
	"strings"
)

func Lint(definition string) (*ArchitectureModel, []Issue) {
	m := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(definition), &m); err != nil {
		return nil, []Issue{*NewError("Invalid YAML")}
	}
	issues := make([]Issue, 0)
	model := ArchitectureModel{}
	version, issue := versionOf(m)
	if issue == nil {
		model.version = version
	} else {
		issues = append(issues, *issue)
	}

	return &model, issues
}

const defaultVersion = "1.0"
const maxMajorVersion = 1
const maxMinorVersion = 0
const maxPatchVersion = 0

var versionPattern = regexp.MustCompile(`\d+\.\d+(\.\d+)?`)

func versionOf(m map[string]interface{}) (string, *Issue) {
	raw, exists := m["version"]
	if !exists {
		return defaultVersion, nil
	}
	version, isString := raw.(string)
	if isString {
		if !versionPattern.MatchString(version) {
			return "", NewError("Version must be a semantic version as defined by https://semver.org")
		}
	} else {
		floatVersion, isFloat := raw.(float64)
		if !isFloat {
			return "", NewError("Version must be a semantic version as defined by https://semver.org")
		}
		version = fmt.Sprint(floatVersion)
		if !strings.Contains(version, ".") {
			version = fmt.Sprintf("%s.0", version)
		}
	}
	parts := strings.Split(version, ".")
	major, _ := strconv.Atoi(parts[0])
	if major > maxMajorVersion {
		return "", NewError(fmt.Sprintf("Undefined version: %s", version))
	} else if major == maxMajorVersion {
		minor, _ := strconv.Atoi(parts[1])
		if minor > maxMinorVersion {
			return "", NewError(fmt.Sprintf("Undefined version: %s", version))
		} else if minor == maxMinorVersion {
			if len(parts) > 2 {
				patch, _ := strconv.Atoi(parts[1])
				if patch > maxPatchVersion {
					return "", NewError(fmt.Sprintf("Undefined version: %s", version))
				}
			}
		}
	}
	return version, nil
}

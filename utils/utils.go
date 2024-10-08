package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetOwnerRepoPull(ref, repo string) (string, string, int, error) {

	if ref == "" || repo == "" {
		return "", "", 0, fmt.Errorf("ref is empty")
	}

	githubrefS := strings.Split(ref, "/")
	prnum := githubrefS[2]
	pull, _ := strconv.Atoi(prnum)

	repoS := strings.Split(repo, "/")
	owner := repoS[0]
	repoName := repoS[1]

	return owner, repoName, pull, nil
}

func ValidateModuleSource(source string, approvedModules map[string]bool) (bool, error) {
	for k, v := range approvedModules {
		if strings.Contains(source, k) {
			if v {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("module not approved")
}

func GetSourceLine(source string) string {
	sourceS := strings.Split(source, "\n")
	for _, line := range sourceS {
		if strings.Contains(line, "source") {
			sourceLine := strings.Split(line, "=")
			return sourceLine[1]
		}
	}
	return ""
}

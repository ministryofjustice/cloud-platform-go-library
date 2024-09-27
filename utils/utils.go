package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// GetOwnerRepoPull takes a ref and a repo and returns the owner, repo name, and pull request number.
// It returns the owner, repo name, pull request number, and an error if the ref or repo is empty.
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

// ValidateModuleSource checks if the provided module source is approved.
// It iterates through the approvedModules map and checks if the source contains
// any of the keys. If a key is found and its corresponding value is true, the function
// returns true and no error. If no approved module is found, it returns false and an error.
//
// Parameters:
//   - source: A string representing the module source to be validated.
//   - approvedModules: A map where keys are module identifiers and values are booleans
//     indicating whether the module is approved.
//
// Returns:
//   - bool: True if the module source is approved, otherwise false.
//   - error: An error if the module source is not approved, otherwise nil.
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

// GetSourceLine extracts the value of the "source" line from a given multi-line string.
// It splits the input string by newline characters, searches for a line containing the word "source",
// and then splits that line by the "=" character to return the value part.
//
// Parameters:
//   - source: A multi-line string containing various lines, one of which includes the "source" keyword.
//
// Returns:
//   - A string representing the value associated with the "source" keyword. If no such line is found, it returns an empty string.
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

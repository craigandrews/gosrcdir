package main

import (
	"errors"
	"fmt"
	"go/build"
	"net/url"
	"os"
	"path"
	"strings"
)

// parseStandardURL parses a standard URI formatted git URL
func parseStandardURL(repo string) (pathParts []string, err error) {
	parsedURL, err := url.Parse(repo)
	if err != nil {
		return
	}

	if parsedURL.Host == "" {
		err = errors.New("Missing host part")
		return
	}

	pathParts = append(pathParts, parsedURL.Host)

	for _, pathPart := range strings.Split(parsedURL.Path, "/") {
		if pathPart == "" {
			continue
		}
		pathParts = append(pathParts, pathPart)
	}

	return
}

// parseWeirdGitURL tries to make sense of a user@host:path format git URL
func parseWeirdGitURL(repo string) (pathParts []string, err error) {

	hostIndex := strings.Index(repo, "@") + 1
	pathIndex := strings.Index(repo, ":") + 1
	repoPath := strings.Split(string(repo[pathIndex:]), "/")

	// If there is no : then this is wrong
	if pathIndex == 0 || len(repoPath) == 0 {
		err = errors.New("Missing path part")
		return
	}

	// if index of @ is -1 then host is the first thing
	// if @ is after : then it's part of the path
	if hostIndex > pathIndex {
		hostIndex = 0
	}

	host := string(repo[hostIndex : pathIndex-1])
	pathParts = append(pathParts, host)

	for _, pathPart := range repoPath {
		if pathPart == "" {
			err = fmt.Errorf("Blank path segment")
			return
		}
		pathParts = append(pathParts, pathPart)
	}

	return
}

// calculateSourcePath works out the local filesystem path to directory above a given repo
func calculateSourcePath(goPath string, repo string) (repoPath string, err error) {
	pathParts, err := parseStandardURL(repo)
	if err != nil {
		pathParts, err = parseWeirdGitURL(repo)
		if err != nil {
			return
		}
	}

	if len(pathParts) < 2 {
		err = errors.New("Host and path required")
		return
	}

	pathParts = append([]string{goPath, "src"}, pathParts[:len(pathParts)-1]...)
	repoPath = path.Join(pathParts...)
	return
}

func getGoPath() string {
	goPath := os.Getenv("GOPATH")
	if goPath != "" {
		return goPath
	}
	return build.Default.GOPATH
}

func main() {
	args := os.Args[1:]

	// Require at least one repo
	if len(args) == 0 {
		os.Exit(1)
	}

	// Work out where GOPATH really is
	goPath := getGoPath()

	// Calculate path for all repos
	for _, repo := range args {
		repoPath, err := calculateSourcePath(goPath, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot parse repo URL %s: %s\n", repo, err)
			os.Exit(1)
		}
		fmt.Println(repoPath)
	}
}

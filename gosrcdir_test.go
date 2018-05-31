package main

import (
	"reflect"
	"testing"
)

type ParseTest struct {
	Name string
	Repo string
}

func TestParsesStandardURL(t *testing.T) {
	for _, test := range []ParseTest{
		{"HTTPS", "https://user@host.com/path/to/repo"},
		{"SSH", "ssh://user@host.com/path/to/repo"},
		{"Empty path elements", "https://user@host.com//path//to//repo"},
	} {
		repoPath, err := parseStandardURL(test.Repo)
		if err != nil {
			t.Errorf("%s: Unexpected error %q", test.Name, err)
		}

		expectedPath := []string{"host.com", "path", "to", "repo"}
		if !reflect.DeepEqual(repoPath, expectedPath) {
			t.Errorf("%s: Expected %q, got %q", test.Name, expectedPath, repoPath)
		}
	}
}

func TestParsesWeirdGitURL(t *testing.T) {
	for _, test := range []ParseTest{
		{"with user", "user@host.com:path/to/repo"},
		{"without user", "host.com:path/to/repo"},
	} {
		repoPath, err := parseWeirdGitURL(test.Repo)
		if err != nil {
			t.Errorf("%s: Unexpected error %q", test.Name, err)
		}

		expectedPath := []string{"host.com", "path", "to", "repo"}
		if !reflect.DeepEqual(repoPath, expectedPath) {
			t.Errorf("%s: Expected %q, got %q", test.Name, expectedPath, repoPath)
		}
	}
}

func TestIsIntolerantToEmptyPathElementsInWeirdURL(t *testing.T) {
	_, err := parseWeirdGitURL("user@host.com:path/to//repo")
	if err == nil {
		t.Errorf("Expected error, got %q", err)
	}
}

func TestSourcePath(t *testing.T) {
	for _, test := range []ParseTest{
		{"Standard URL", "https://user@host.com/path/to/repo"},
		{"Weird URL", "user@host.com:path/to/repo"},
	} {
		repoPath, err := calculateSourcePath("/home/user", test.Repo)
		if err != nil {
			t.Errorf("%s: Unexpected error %q", test.Name, err)
		}

		expectedPath := "/home/user/src/host.com/path/to"
		if repoPath != expectedPath {
			t.Errorf("%sExpected %q, got %q", test.Name, expectedPath, repoPath)
		}
	}
}

func TestSourcePathErrorsIfMissingHostOrPath(t *testing.T) {
	for _, test := range []ParseTest{
		{"Standard URL", "https://user@host.com/"},
		{"Weird URL", "user@host.com:"},
	} {
		_, err := calculateSourcePath("/home/user", "user@host.com:")
		if err == nil {
			t.Errorf("%s: Expected error, got %q", test.Name, err)
		}
	}
}

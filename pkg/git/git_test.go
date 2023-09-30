package git

import (
	"os"
	"testing"
)

const (
	testRepoURL           = "https://github.com/zcubbs/haproxy-test-repo-public.git"
	testHaproxyFilePath   = "/basic/haproxy.cfg"
	testHaproxyConfigPath = "/etc/haproxy/haproxy.cfg"
	testRepoBranch        = "main"
)

func TestNewGitHandler(t *testing.T) {
	handler := NewHandler(testRepoURL, testRepoBranch, testHaproxyFilePath, testHaproxyConfigPath)
	if handler == nil {
		t.Errorf("Failed to create a new GitHandler.")
	}
}

func TestCloneRepo(t *testing.T) {
	handler := NewHandler(testRepoURL, testRepoBranch, testHaproxyFilePath, testHaproxyConfigPath)
	err := handler.cloneRepo()
	if err != nil {
		t.Errorf("Failed to clone the repo: %v", err)
	}
	defer os.RemoveAll(handler.localRepoPath) // Cleanup after test
}

func TestPullRepo(t *testing.T) {
	handler := NewHandler(testRepoURL, testRepoBranch, testHaproxyFilePath, testHaproxyConfigPath)
	err := handler.cloneRepo()
	if err != nil {
		t.Errorf("Failed to clone the repo: %v", err)
	}
	_, err = handler.pullRepo()
	if err != nil {
		t.Errorf("Failed to pull the repo: %v", err)
	}
}

package git

import (
	"os"
	"testing"
)

const (
	testRepoURL           = "https://github.com/zcubbs/haproxy-test-repo-public.git"
	testHaproxyConfigPath = "/basic/haproxy.cfg"
	testRepoBranch        = "main"
)

func TestNewGitHandler(t *testing.T) {
	handler := NewHandler(testRepoURL, testRepoBranch, testHaproxyConfigPath)
	if handler == nil {
		t.Errorf("Failed to create a new GitHandler.")
	}
}

func TestCloneRepo(t *testing.T) {
	handler := NewHandler(testRepoURL, testRepoBranch, testHaproxyConfigPath)
	err := handler.cloneRepo()
	if err != nil {
		t.Errorf("Failed to clone the repo: %v", err)
	}
	defer os.RemoveAll(handler.localRepoPath) // Cleanup after test
}

func TestPullRepo(t *testing.T) {
	handler := NewHandler(testRepoURL, testRepoBranch, testHaproxyConfigPath)
	_, err := handler.pullRepo()
	if err != nil {
		t.Errorf("Failed to pull the repo: %v", err)
	}
}

package git

import (
	"fmt"
	"github.com/zcubbs/hpxd/pkg/cmd"
	"os"
	"path/filepath"
)

type Handler struct {
	repoURL           string
	branch            string
	localRepoPath     string
	path              string
	haproxyConfigPath string
}

// NewHandler initializes a new GitHandler
func NewHandler(repoURL, branch, path, haproxyConfigPath string) *Handler {
	return &Handler{
		repoURL:           repoURL,
		branch:            branch,
		localRepoPath:     filepath.Join(os.TempDir(), "hpxd-git-repo"),
		path:              path,
		haproxyConfigPath: haproxyConfigPath,
	}
}

// PullAndUpdate clones or pulls the git repo, and then returns the path to the new config if changes are detected
func (g *Handler) PullAndUpdate() (string, bool, error) {
	// Check if repo already exists locally
	if _, err := os.Stat(g.localRepoPath); os.IsNotExist(err) {
		// Clone repo if it doesn't exist
		return g.cloneRepo()
	}

	// Pull latest changes if repo exists
	return g.pullRepo()
}

func (g *Handler) getHAProxyConfigPath() string {
	return filepath.Join(g.localRepoPath, g.path)
}

func (g *Handler) cloneRepo() (string, bool, error) {
	if err := cmd.RunCmd("git", "clone", "-b", g.branch, g.repoURL, g.localRepoPath); err != nil {
		return "", false, err
	}
	return g.getHAProxyConfigPath(), true, nil
}

func (g *Handler) pullRepo() (string, bool, error) {
	output, err := cmd.RunCmdCombinedOutput("git", "-C", g.localRepoPath, "pull", "origin", g.branch)
	if err != nil {
		return "", false, fmt.Errorf("failed to pull repo: %v details: %s", err, string(output))
	}

	// Check if there were any updates from the pull
	if string(output) == "Already up to date.\n" {
		return "", false, nil
	}

	return g.getHAProxyConfigPath(), true, nil
}

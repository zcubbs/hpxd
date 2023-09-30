package git

import (
	"fmt"
	"github.com/zcubbs/hpxd/pkg/cmd"
	"log"
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

// PullAndUpdate clones or pulls the git repo, and then updates the HAProxy config if changes are detected
func (g *Handler) PullAndUpdate() (bool, error) {
	// Check if repo already exists locally
	if _, err := os.Stat(g.localRepoPath); os.IsNotExist(err) {
		// Clone repo if it doesn't exist
		err := g.cloneRepo()
		if err != nil {
			return false, fmt.Errorf("failed to clone repo: %v", err)
		}
		return true, nil // Since it's a new clone, we assume changes
	}

	// Pull latest changes if repo exists
	updated, err := g.pullRepo()
	if err != nil {
		return false, fmt.Errorf("failed to pull repo: %v", err)
	}

	return updated, nil
}

func (g *Handler) cloneRepo() error {
	if err := cmd.RunCmd("git",
		"clone", "-b", g.branch, g.repoURL, g.localRepoPath); err != nil {
		return err
	}
	return g.updateHAProxyConfig()
}

func (g *Handler) pullRepo() (bool, error) {
	output, err := cmd.RunCmdCombinedOutput("git", "-C", g.localRepoPath, "pull", "origin", g.branch)
	if err != nil {
		return false, fmt.Errorf("failed to pull repo: %v details: %s", err, string(output))
	}

	// Check if there were any updates from the pull
	if string(output) == "Already up to date.\n" {
		return false, nil
	}

	return true, g.updateHAProxyConfig()
}

func (g *Handler) updateHAProxyConfig() error {
	srcFile := filepath.Join(g.localRepoPath, g.path)
	dstFile := g.haproxyConfigPath

	input, err := os.ReadFile(filepath.Clean(srcFile))
	if err != nil {
		return err
	}

	err = os.WriteFile(dstFile, input, 0600)
	if err != nil {
		log.Println("Error creating", dstFile)
		return err
	}

	return nil
}

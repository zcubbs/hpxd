package git

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Handler struct {
	repoURL           string
	branch            string
	localRepoPath     string
	haproxyConfigPath string
}

// NewHandler initializes a new GitHandler
func NewHandler(repoURL, branch, haproxyConfigPath string) *Handler {
	return &Handler{
		repoURL:           repoURL,
		branch:            branch,
		localRepoPath:     filepath.Join(os.TempDir(), "hpxd-git-repo"),
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
			return false, err
		}
		return true, nil // Since it's a new clone, we assume changes
	}

	// Pull latest changes if repo exists
	updated, err := g.pullRepo()
	if err != nil {
		return false, err
	}

	return updated, nil
}

func (g *Handler) cloneRepo() error {
	cmd := exec.Command("git", "clone", "-b", g.branch, g.repoURL, g.localRepoPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return g.updateHAProxyConfig()
}

func (g *Handler) pullRepo() (bool, error) {
	cmd := exec.Command("git", "-C", g.localRepoPath, "pull", "origin", g.branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}

	// Check if there were any updates from the pull
	if string(output) == "Already up to date.\n" {
		return false, nil
	}

	return true, g.updateHAProxyConfig()
}

func (g *Handler) updateHAProxyConfig() error {
	srcFile := filepath.Join(g.localRepoPath, "path-to-config-inside-repo") // You need to specify the relative path of your HAProxy config file inside the Git repo
	dstFile := g.haproxyConfigPath

	input, err := os.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(dstFile, input, 0644)
	if err != nil {
		log.Println("Error creating", dstFile)
		return err
	}

	return nil
}

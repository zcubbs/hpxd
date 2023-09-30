// Package git provides utilities to interact with git repositories.
//
// This package is primarily designed to clone and pull updates from
// git repositories, specifically with support for optional credentials
// in the form of username and password from environment variables.
//
// Author: zakaria.elbouwab
package git

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zcubbs/hpxd/pkg/cmd"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Handler manages operations on a git repository.
//
// The Handler structure contains fields that represent details of the git repository,
// such as the repo URL, branch, and paths for both the local copy of the repository
// and the HAProxy configuration.
type Handler struct {
	repoURL           string
	branch            string
	username          string
	password          string
	localRepoPath     string
	path              string
	haproxyConfigPath string
}

// NewHandler initializes and returns a new Handler instance.
//
// This function constructs a Handler given details of the git repository and the
// path where the HAProxy configuration is located.
func NewHandler(repoURL, branch, username, password, path, haproxyConfigPath string) *Handler {
	return &Handler{
		repoURL:           repoURL,
		branch:            branch,
		username:          username,
		password:          password,
		localRepoPath:     filepath.Join(os.TempDir(), "hpxd-git-repo"),
		path:              path,
		haproxyConfigPath: haproxyConfigPath,
	}
}

// PullAndUpdate performs a git clone or pull operation.
//
// If the local copy of the repository doesn't exist, it clones the repo.
// If it does exist, it pulls the latest changes. This function then
// returns the path to the new configuration and a flag indicating if there
// were any updates.
func (g *Handler) PullAndUpdate() (string, bool, error) {
	// Check if repo already exists locally
	if _, err := os.Stat(g.localRepoPath); os.IsNotExist(err) {
		// Clone repo if it doesn't exist
		return g.cloneRepo()
	}

	// Pull latest changes if repo exists
	return g.pullRepo()
}

// getHAProxyConfigPath constructs and returns the complete file path
// for the HAProxy configuration within the local copy of the git repository.
func (g *Handler) getHAProxyConfigPath() string {
	return filepath.Join(g.localRepoPath, g.path)
}

// getRepoURLWithCredentials returns the repo URL with credentials embedded.
//
// If the username and password fields are empty, it returns the repo URL as is.
func (g *Handler) getRepoURLWithCredentials() string {
	if g.username != "" && g.password != "" {
		// Embed the credentials in the repo URL
		return strings.Replace(g.repoURL, "https://", fmt.Sprintf("https://%s:%s@", url.PathEscape(g.username), url.PathEscape(g.password)), 1)
	}
	return g.repoURL
}

// cloneRepo clones the git repository to the local machine.
//
// It returns the path to the HAProxy configuration and a flag indicating if
// the clone operation was successful.
func (g *Handler) cloneRepo() (string, bool, error) {
	if err := cmd.RunCmd("git", "clone", "-b", g.branch, g.getRepoURLWithCredentials(), g.localRepoPath); err != nil {
		return "", false, err
	}
	return g.getHAProxyConfigPath(), true, nil
}

// pullRepo fetches the latest changes from the git repository.
//
// It returns the path to the updated HAProxy configuration and a flag indicating
// if there were any changes during the pull operation.
func (g *Handler) pullRepo() (string, bool, error) {
	output, err := cmd.RunCmdCombinedOutput("git", "-C", g.localRepoPath, "pull", g.getRepoURLWithCredentials(), g.branch)
	if err != nil {
		logrus.Debugf("Failed to pull repo: %v, details: %s", err, string(output))
		return "", false, fmt.Errorf("failed to pull repo: %v, details: %s", err, string(output))
	}

	logrus.Debugf("Git pull output: %s", string(output))
	// Check if there were any updates from the pull
	if string(output) == "Already up to date.\n" {
		return "", false, nil
	}

	return g.getHAProxyConfigPath(), true, nil
}

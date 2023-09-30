// Package haproxy provides utilities to manage and interact with HAProxy.
//
// This package offers capabilities such as validating the HAProxy configuration
// and gracefully reloading HAProxy.
//
// Author: zakaria.elbouwab
package haproxy

import (
	"github.com/zcubbs/hpxd/pkg/cmd"
	"strings"
)

// Handler manages operations related to HAProxy.
//
// The Handler structure contains a field that represents the path
// to the HAProxy configuration.
type Handler struct {
	configPath string
}

// NewHandler initializes and returns a new Handler instance for HAProxy.
//
// This function constructs a Handler given the path to the HAProxy configuration.
func NewHandler(configPath string) *Handler {
	return &Handler{
		configPath: configPath,
	}
}

// ValidateConfig checks the validity of the current HAProxy configuration.
//
// This method runs the 'haproxy -c -f' command to validate the configuration.
// If the configuration is invalid, it returns an Error containing both the
// original error and the output from the validation command.
func (h *Handler) ValidateConfig() error {
	output, err := cmd.RunCmdCombinedOutput("haproxy", "-c", "-f", h.configPath)
	if err != nil {
		return &Error{OriginalError: err, Output: string(output)}
	}

	return nil
}

// Reload gracefully restarts HAProxy.
//
// This method runs the 'sudo systemctl reload haproxy' command to gracefully
// reload HAProxy. If there's an error during the reload, it returns an Error
// containing both the original error and the output from the reload command.
func (h *Handler) Reload() error {
	output, err := cmd.RunCmdCombinedOutput("sudo", "systemctl", "reload", "haproxy")
	if err != nil && len(output) > 0 {
		return &Error{OriginalError: err, Output: string(output)}
	}

	return nil
}

// Error is an error wrapper for capturing command output alongside the error message.
//
// This structure extends the built-in error type to provide more context about
// errors that occur when running commands related to HAProxy. It captures both
// the original error and the command's output.
type Error struct {
	OriginalError error
	Output        string
}

// Error returns a concatenated string of the original error message and the command output.
func (e *Error) Error() string {
	return strings.TrimSpace(e.OriginalError.Error() + ": " + e.Output)
}

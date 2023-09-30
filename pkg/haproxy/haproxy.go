package haproxy

import (
	"github.com/zcubbs/hpxd/pkg/cmd"
	"strings"
)

type Handler struct {
	configPath string
}

// NewHandler initializes a new Handler for HAProxy
func NewHandler(configPath string) *Handler {
	return &Handler{
		configPath: configPath,
	}
}

// ValidateConfig validates the current HAProxy configuration
func (h *Handler) ValidateConfig() error {
	output, err := cmd.RunCmdCombinedOutput("haproxy", "-c", "-f", h.configPath)
	if err != nil {
		return &Error{OriginalError: err, Output: string(output)}
	}

	return nil
}

// Reload gracefully reloads HAProxy
func (h *Handler) Reload() error {
	output, err := cmd.RunCmdCombinedOutput("systemctl", "reload", "haproxy")
	if err != nil && len(output) > 0 {
		return &Error{OriginalError: err, Output: string(output)}
	}

	return nil
}

// Error is an error wrapper that captures command output
type Error struct {
	OriginalError error
	Output        string
}

func (e *Error) Error() string {
	return strings.TrimSpace(e.OriginalError.Error() + ": " + e.Output)
}

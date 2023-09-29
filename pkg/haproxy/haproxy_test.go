package haproxy

import (
	"testing"
)

func TestValidateConfig_Valid(t *testing.T) {
	handler := NewHandler("./testdata/valid_haproxy.cfg")

	err := handler.ValidateConfig()
	if err != nil {
		t.Errorf("Expected configuration to be valid, but got error: %v", err)
	}
}

func TestValidateConfig_Invalid(t *testing.T) {
	handler := NewHandler("./testdata/invalid_haproxy.cfg")

	err := handler.ValidateConfig()
	if err == nil {
		t.Errorf("Expected configuration to be invalid, but validation passed.")
	}
}

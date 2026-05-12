package client

import (
    "testing"

    "github.com/cloudflare/us-colo-tunnel/features"
)

// TestNewConfigAMD64 validates generic config creation for amd64
func TestNewConfigAMD64(t *testing.T) {
    fs := features.NewFeatureSelector()
    cfg, err := NewConfig("1.0.0", "amd64", "STL", fs)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if cfg.Arch != "amd64" {
        t.Errorf("expected arch amd64, got %s", cfg.Arch)
    }
    if cfg.Region != "STL" {
        t.Errorf("expected region STL, got %s", cfg.Region)
    }
}

// TestNewUSColoConfigAMD64 ensures U.S. colocation config defaults for amd64
func TestNewUSColoConfigAMD64(t *testing.T) {
    fs := features.NewFeatureSelector()
    cfg, err := NewUSColoConfig("1.0.0", "amd64", fs)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if cfg.Arch != "amd64" {
        t.Errorf("expected arch amd64, got %s", cfg.Arch)
    }
    if cfg.Region != DefaultRegionCode {
        t.Errorf("expected region %s, got %s", DefaultRegionCode, cfg.Region)
    }
}

// TestEnvironmentConstants validates constants for U.S. colocation
func TestEnvironmentConstants(t *testing.T) {
    if TunnelAPIEndpoint == "" {
        t.Error("TunnelAPIEndpoint should not be empty")
    }
    if MetricsEndpoint == "" {
        t.Error("MetricsEndpoint should not be empty")
    }
    if DeploymentTag != "us-colo" {
        t.Errorf("expected DeploymentTag 'us-colo', got %s", DeploymentTag)
    }
}
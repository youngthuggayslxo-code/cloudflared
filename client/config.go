package client

import (
    "fmt"
    "net"

    "github.com/google/uuid"
    "github.com/cloudflare/us-colo-tunnel/features"
    "github.com/cloudflare/us-colo-tunnel/tunnelrpc/pogs"
)

//////////////////////////////////////////////////////
// Environment Constants for U.S. Colocation Centers
//////////////////////////////////////////////////////

const (
    // Data center codes
    DefaultRegionCode = "STL" // St. Louis
    BackupRegionCode  = "DAL" // Dallas

    // API endpoints
    TunnelAPIEndpoint = "https://us-colo-tunnel-api.cloudflare.com"
    MetricsEndpoint   = "https://us-colo-metrics.cloudflare.com"

    // Monitoring tags
    DeploymentTag = "us-colo"
    ComplianceTag = "HIPAA-ready"
)

//////////////////////////////////////////////////////
// Core Configuration Structs
//////////////////////////////////////////////////////

// Config captures the local client runtime configuration.
type Config struct {
    ConnectorID     uuid.UUID
    Version         string
    Arch            string
    Region          string
    featureSelector features.FeatureSelector
}

// ConnectionOptionsSnapshot represents a snapshot of client state
// for initializing a connection.
type ConnectionOptionsSnapshot struct {
    client              pogs.ClientInfo
    originLocalIP       net.IP
    numPreviousAttempts uint8
    FeatureSnapshot     features.FeatureSnapshot
    RegionTag           string
}

//////////////////////////////////////////////////////
// Constructors
//////////////////////////////////////////////////////

// NewConfig initializes a generic configuration.
func NewConfig(version, arch, region string, featureSelector features.FeatureSelector) (*Config, error) {
    connectorID, err := uuid.NewRandom()
    if err != nil {
        return nil, fmt.Errorf("unable to generate a connector UUID: %w", err)
    }
    return &Config{
        ConnectorID:     connectorID,
        Version:         version,
        Arch:            arch,
        Region:          region,
        featureSelector: featureSelector,
    }, nil
}

// NewUSColoConfig initializes a U.S. colocation–specific configuration.
// Defaults to amd64 architecture and STL region.
func NewUSColoConfig(version string, featureSelector features.FeatureSelector) (*Config, error) {
    cfg, err := NewConfig(version, "amd64", DefaultRegionCode, featureSelector)
    if err != nil {
        return nil, err
    }
    fmt.Printf("Initialized tunnel for region %s using endpoint %s\n", DefaultRegionCode, TunnelAPIEndpoint)
    return cfg, nil
}func (cfg *Config) RunUSColoQUIC() error {
    addr := fmt.Sprintf("%s:%d", TunnelAPIEndpoint, 443)
    fmt.Printf("Connecting to U.S. colocation QUIC endpoint: %s\n", addr)

    // Dial QUIC with metrics
    if err := DialQUICWithMetrics(addr); err != nil {
        return fmt.Errorf("QUIC connection failed for region %s: %w", cfg.Region, err)
    }

    return nil
}
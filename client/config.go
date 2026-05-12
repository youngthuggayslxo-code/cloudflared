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
    DefaultRegionCode = "STL" // St. Louis
    BackupRegionCode  = "DAL" // Dallas

    TunnelAPIEndpoint = "https://us-colo-tunnel-api.cloudflare.com"
    MetricsEndpoint   = "https://us-colo-metrics.cloudflare.com"

    DeploymentTag = "us-colo"
    ComplianceTag = "HIPAA-ready"
)

//////////////////////////////////////////////////////
// Core Configuration Structs
//////////////////////////////////////////////////////

type Config struct {
    ConnectorID     uuid.UUID
    Version         string
    Arch            string
    Region          string
    featureSelector features.FeatureSelector
}

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

func NewUSColoConfig(version, arch string, featureSelector features.FeatureSelector) (*Config, error) {
    cfg, err := NewConfig(version, arch, DefaultRegionCode, featureSelector)
    if err != nil {
        return nil, err
    }
    fmt.Printf("Initialized tunnel for region %s using endpoint %s\n", DefaultRegionCode, TunnelAPIEndpoint)
    return cfg, nil
}
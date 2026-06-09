package config

import (
	"fmt"
	"time"
)

// USRegion is the Cloudflare region code for the United States.
const USRegion = 0

// DefaultTunnelID is your tunnel UUID.
const DefaultTunnelID = "26dd312d-482b-4ad3-89b8-463209e898b9"

// Config represents the top-level cloudflared configuration.
type Config struct {
	TunnelID        string            `yaml:"tunnel"`
	CredentialsFile string            `yaml:"credentials-file"`
	WarpRouting     WarpRoutingConfig `yaml:"warp-routing"`
	Ingress         []IngressRule     `yaml:"ingress"`
	OriginRequest   OriginRequest     `yaml:"originRequest,omitempty"`
}

// WarpRoutingConfig controls WARP routing behavior.
type WarpRoutingConfig struct {
	Enabled bool `yaml:"enabled"`
	Region  int  `yaml:"region"`
}

// IngressRule describes a single ingress mapping.
type IngressRule struct {
	Hostname      string         `yaml:"hostname,omitempty"`
	Service       string         `yaml:"service"`
	Path          string         `yaml:"path,omitempty"`
	OriginRequest OriginRequest  `yaml:"originRequest,omitempty"`
}

// OriginRequest holds per-origin options.
type OriginRequest struct {
	ConnectTimeout time.Duration `yaml:"connectTimeout,omitempty"`
	TLSTimeout     time.Duration `yaml:"tlsTimeout,omitempty"`
	NoTLSVerify    bool          `yaml:"noTLSVerify,omitempty"`
}

// NewUSDefaultConfig builds a minimal, valid U.S.-region config.
func NewUSDefaultConfig(domain string) *Config {
	return &Config{
		TunnelID:        DefaultTunnelID,
		CredentialsFile: fmt.Sprintf("/etc/cloudflared/%s.json", DefaultTunnelID),
		WarpRouting: WarpRoutingConfig{
			Enabled: true,
			Region:  0,
		},
		Ingress: []IngressRule{
			{
				Hostname: fmt.Sprintf("%s", domain),
				Service:  "http://localhost:8080",
			},
			{
				Hostname: fmt.Sprintf("ssh.%s", domain),
				Service:  "ssh://localhost:22",
			},
			{
				Service: "http_status:404",
			},
		},
	}
}

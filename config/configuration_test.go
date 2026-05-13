package config

import "testing"

func TestNewUSDefaultConfig_BasicShape(t *testing.T) {
	domain := "example.com"
	cfg := NewUSDefaultConfig(domain)

	if cfg.TunnelID != DefaultTunnelID {
		t.Fatalf("expected tunnel ID %s, got %s", DefaultTunnelID, cfg.TunnelID)
	}

	if cfg.WarpRouting.Region != USRegion {
		t.Fatalf("expected region %d, got %d", USRegion, cfg.WarpRouting.Region)
	}

	if len(cfg.Ingress) != 3 {
		t.Fatalf("expected 3 ingress rules, got %d", len(cfg.Ingress))
	}

	if cfg.Ingress[0].Hostname != domain {
		t.Fatalf("expected first hostname %q, got %q", domain, cfg.Ingress[0].Hostname)
	}

	if cfg.Ingress[2].Service != "http_status:404" {
		t.Fatalf("expected last rule to be http_status:404, got %q", cfg.Ingress[2].Service)
	}
}
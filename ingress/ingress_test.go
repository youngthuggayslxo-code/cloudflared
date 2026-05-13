package ingress

import (
	"testing"

	"your/module/path/config"
)

func TestParse_ValidUSConfig(t *testing.T) {
	domain := "example.com"
	cfg := config.NewUSDefaultConfig(domain)

	ing, err := Parse(cfg)
	if err != nil {
		t.Fatalf("unexpected error parsing ingress: %v", err)
	}

	if len(ing.Rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(ing.Rules))
	}

	last := ing.Rules[len(ing.Rules)-1]
	if last.Service != "http_status:404" {
		t.Fatalf("expected last rule to be http_status:404, got %q", last.Service)
	}
}

func TestFindMatchingRule_HostMatch(t *testing.T) {
	domain := "example.com"
	cfg := config.NewUSDefaultConfig(domain)
	ing, _ := Parse(cfg)

	r := ing.FindMatchingRule("example.com", "/")
	if r == nil || r.Service != "http://localhost:8080" {
		t.Fatalf("expected match to http://localhost:8080, got %#v", r)
	}
}

func TestFindMatchingRule_FallbackToCatchAll(t *testing.T) {
	domain := "example.com"
	cfg := config.NewUSDefaultConfig(domain)
	ing, _ := Parse(cfg)

	r := ing.FindMatchingRule("unknown.com", "/whatever")
	if r == nil || r.Service != "http_status:404" {
		t.Fatalf("expected fallback to http_status:404, got %#v", r)
	}
}
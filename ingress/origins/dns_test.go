package dns

import "testing"

func TestDNSOriginValidation(t *testing.T) {
	origin := DNSOrigin{Address: "1.1.1.1", DNSSEC: true}
	if err := origin.Validate(); err != nil {
		t.Fatalf("expected valid origin, got %v", err)
	}
	if !origin.EnforceDNSSEC() {
		t.Fatalf("expected DNSSEC enforcement")
	}
}

func TestInvalidDNSOrigin(t *testing.T) {
	origin := DNSOrigin{Address: "not-an-ip"}
	if err := origin.Validate(); err == nil {
		t.Fatalf("expected error for invalid IP")
	}
}
package dns

import (
	"errors"
	"net"
)

// DNSOrigin represents a DNS resolver origin.
type DNSOrigin struct {
	Address string
	DNSSEC  bool
}

// Validate ensures the DNS origin is usable.
func (d *DNSOrigin) Validate() error {
	if d.Address == "" {
		return errors.New("DNS origin address cannot be empty")
	}
	if net.ParseIP(d.Address) == nil {
		return errors.New("invalid DNS origin address")
	}
	return nil
}

// EnforceDNSSEC checks if DNSSEC is required.
func (d *DNSOrigin) EnforceDNSSEC() bool {
	return d.DNSSEC
}
package ingress

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"your/module/path/config"
)

var (
	ErrNoIngressRules      = errors.New("no ingress rules defined")
	ErrLastRuleNotCatchAll = errors.New("last ingress rule must be a catch-all rule")
	ErrBadWildcard         = errors.New("wildcard may only appear once and only at the start of the hostname")
	ErrHostnameHasPort     = errors.New("hostname must not contain a port")
)

// Rule represents a compiled ingress rule.
type Rule struct {
	Hostname string
	Path     *regexp.Regexp
	Service  string
}

// Ingress is the compiled routing table.
type Ingress struct {
	Rules []Rule
}

// Parse builds an Ingress from configuration.
func Parse(conf *config.Config) (*Ingress, error) {
	if conf == nil || len(conf.Ingress) == 0 {
		return nil, ErrNoIngressRules
	}

	rules := make([]Rule, len(conf.Ingress))
	for i, r := range conf.Ingress {
		if err := validateHostname(r.Hostname, i, len(conf.Ingress), r.Path); err != nil {
			return nil, err
		}

		var pathRE *regexp.Regexp
		if r.Path != "" {
			re, err := regexp.Compile(r.Path)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid path regex in rule #%d", i+1)
			}
			pathRE = re
		}

		if err := validateService(r.Service); err != nil {
			return nil, errors.Wrapf(err, "invalid service in rule #%d", i+1)
		}

		rules[i] = Rule{
			Hostname: r.Hostname,
			Path:     pathRE,
			Service:  r.Service,
		}
	}

	return &Ingress{Rules: rules}, nil
}

// FindMatchingRule returns the first rule that matches hostname and path.
func (ing *Ingress) FindMatchingRule(hostname, path string) *Rule {
	host := stripPort(hostname)
	for i := range ing.Rules {
		r := &ing.Rules[i]
		if !matchHost(r.Hostname, host) {
			continue
		}
		if r.Path != nil && !r.Path.MatchString(path) {
			continue
		}
		return r
	}
	// Fallback to last rule (catch-all) if nothing matched.
	if len(ing.Rules) == 0 {
		return nil
	}
	return &ing.Rules[len(ing.Rules)-1]
}

func stripPort(hostport string) string {
	h, _, err := net.SplitHostPort(hostport)
	if err != nil {
		return hostport
	}
	return h
}

func matchHost(ruleHost, reqHost string) bool {
	if ruleHost == "" || ruleHost == "*" {
		return true
	}
	if ruleHost == reqHost {
		return true
	}
	if strings.HasPrefix(ruleHost, "*.") {
		suffix := strings.TrimPrefix(ruleHost, "*")
		return strings.HasSuffix(reqHost, suffix)
	}
	return false
}

func validateHostname(host string, index, total int, path string) error {
	if host == "" || host == "*" {
		// Only allowed as the last rule and only if no path.
		if index != total-1 || path != "" {
			return ErrLastRuleNotCatchAll
		}
		return nil
	}

	if _, _, err := net.SplitHostPort(host); err == nil {
		return ErrHostnameHasPort
	}

	if strings.LastIndex(host, "*") > 0 {
		return ErrBadWildcard
	}

	return nil
}

func validateService(s string) error {
	if strings.HasPrefix(s, "http_status:") {
		return nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("service %q must include scheme and host", s)
	}
	if u.Path != "" {
		return fmt.Errorf("service %q must not include a path", s)
	}
	return nil
}
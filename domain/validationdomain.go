package domain

import (
	"fmt"
	"net"
	"regexp"
	"unicode"
)

// ValidationDomain provides the domain logic for session validation
type ValidationDomain struct {
	stringBlacklist []regexp.Regexp
	ipBlacklist     []net.IP
}

// NewValidationDomain creates a new initalized Validation domain logic struct.
func NewValidationDomain(stringBlacklist []string, ipBlacklist []string) (*ValidationDomain, error) {
	ub := make([]regexp.Regexp, 0, len(stringBlacklist))
	for _, entry := range stringBlacklist {
		exp, err := regexp.Compile(entry)
		if err != nil {
			return nil, fmt.Errorf("Can't compile username blacklist regexp '%s': %w", entry, err)
		}
		ub = append(ub, *exp)
	}

	ib := make([]net.IP, 0, len(ipBlacklist))
	for _, entry := range ipBlacklist {
		ip := net.ParseIP(entry)
		if ip == nil {
			return nil, fmt.Errorf("Can't parse ip blacklist entry '%s'", entry)
		}
		ib = append(ib, ip)
	}

	return &ValidationDomain{ub, ib}, nil
}

// ValidateString validates a string against a regexp based blacklist and other rulesets. The validation has linear complexity.
func (d *ValidationDomain) ValidateString(s string) bool {
	if !d.isASCII(s) {
		return false
	}

	for _, entry := range d.stringBlacklist {
		if entry.MatchString(s) {
			return false
		}
	}

	return true
}

// ValdateIP validates an IP address against a IP blacklist. The validation has linear complexity.
func (d *ValidationDomain) ValdateIP(ip net.IP) bool {
	for _, entry := range d.ipBlacklist {
		if entry.Equal(ip) {
			return false
		}
	}

	return true
}

func (d *ValidationDomain) isASCII(s string) bool {
	for _, char := range s {
		if char > unicode.MaxASCII {
			return false
		}
	}

	return true
}

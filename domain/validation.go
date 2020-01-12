package domain

import (
	"fmt"
	"net"
	"regexp"
	"unicode"
)

// ValidationDomain provides the domain logic for session validation
type ValidationDomain struct {
	coreWhistelist map[string]bool
	stringBlacklist []regexp.Regexp
	ipBlacklist []net.IP
}

// NewValidationDomain creates a new initalized Validation domain logic struct.
func NewValidationDomain(coreWhielist []string, stringBlacklist []string, ipBlacklist []string) (*ValidationDomain, error) {
	cw := make(map[string]bool, len(coreWhielist))
	for _, entry := range coreWhielist {
		cw[entry] = true
	}	

	ub := make([]regexp.Regexp, 0, len(stringBlacklist))
	for _, entry := range stringBlacklist {
		exp, err := regexp.Compile(entry)
		if err!= nil {
			return nil, fmt.Errorf("Can't compile username blacklist regexp '%s': %w", entry, err)
		}
		ub = append(ub, *exp)
	}

	ib := make([]net.IP, 0, len(ipBlacklist))
	for _, entry := range ipBlacklist {
		ip :=  net.ParseIP(entry)
		if ip == nil {
			return nil, fmt.Errorf("Can't parse ip blacklist entry '%s'", entry)
		}
		ib = append(ib, ip)
	}

	return &ValidationDomain{cw, ub, ib}, nil
}

// ValidateCore validates a string against a core whitelist. Validation uses a map for more efficient comparison.
func (d* ValidationDomain) ValidateCore(corename string) bool {
	_, found := d.coreWhistelist[corename]
	return found
}

// ValidateString validates a string against a regexp based blacklist and other rulesets. The validation has linear complexity.
func (d* ValidationDomain) ValidateString(s string) bool {
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
func (d* ValidationDomain) ValdateIP(ip net.IP) bool {
	for _, entry := range d.ipBlacklist {
		if entry.Equal(ip) {
			return false
		}
	}

	return true
}

func (d* ValidationDomain) isASCII(s string) bool {
    for _, char := range s {
        if char > unicode.MaxASCII {
            return false
        }
    }
    return true
}

package domain

import (
	"fmt"
	"strings"
	"strconv"
)

// MitmInfo represents a relay server info.
type MitmInfo struct {
	Address string
	Port uint16
}

// MitmDomain abstracts the mitm logic for handling netplay relays.
type MitmDomain struct {
	server map[string]string
}

// NewMitmDomain creates a new MITM domain logic.
func NewMitmDomain(servers map[string]string) *MitmDomain {
	return &MitmDomain{servers}
}

// GetInfo translates a MITM server handle into an address/port pair.
func (d *MitmDomain) GetInfo(handle string) *MitmInfo {
	var server MitmInfo

	address, found := d.server[handle]
	if !found || address == "" {
		return nil
	}

	info := strings.Split(address, ":")
	if len(info) != 2 {
		return nil
	}

	addr := info[0]
	if addr == "" {
		return nil
	}
	port, err := strconv.ParseInt(info[1], 10, 32)
	if err != nil || port < 1 || port > 65535 {
		return nil
	}

	server.Address = addr
	server.Port = uint16(port)

	return &server
}

// PrintForRetroarch prints out the MITM information in a format that retroarch is expecting.
func (i *MitmInfo) PrintForRetroarch() string {
	return fmt.Sprintf("tunnel_addr=%s\ntunnel_port=%d\n", i.Address, i.Port)
}
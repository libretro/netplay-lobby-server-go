package domain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

// MitmSession represents a relay server session.
type MitmSession struct {
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

func (d *MitmDomain) getDefaultServer() string {
	var v string;
	for _,  v = range d.server {
		break 
	}
	return v
}

// OpenSession opens a new netplay session on the specified MITM server
func (d *MitmDomain) OpenSession(handle string) (*MitmSession, error) {
	var server MitmSession
	var port uint32 = 0
	data := make([]byte, 12)

	address, ok := d.server[handle]
	if !ok {
		address = d.getDefaultServer()
	}
	server.Address = strings.Split(address, ":")[0]

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Can't open connection to '%s': %w", address, err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Write([]byte{0x00,0x00,0x46,0x49,0x00,0x00,0x00,0x00})
	if err != nil {
		return nil, fmt.Errorf("Can't send open command to relay server '%s': %w", address, err)
	}
	
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Can't read data from relay server '%s': %w", address, err)
	}

	if res := bytes.Compare(data[0:8], []byte{0x00,0x00,0x46,0x4a,0x00,0x00,0x00,0x04}); res == 0 {
		if err := binary.Read(bytes.NewReader(data[8:12]), binary.BigEndian, &port); err != nil {
			return nil, fmt.Errorf("Can't convert data to port number: %w", err)
		}

		if port > math.MaxUint16 {
			return nil, fmt.Errorf("Recieved port is not in uint16 range: %d", port)
		}

		server.Port = uint16(port)
		return &server, nil
	}

	return nil, fmt.Errorf("Recieved invalid response by relay %s: %X", address, data)
}

// IsNewServerhandle compares if a given server handle is the same as the saved handle server address.
func (d *MitmDomain) IsNewServerhandle(handle string, server string) bool {
	address, ok := d.server[handle]
	if ok && strings.Split(address, ":")[0] == server {
		return false
	}
	return true
}

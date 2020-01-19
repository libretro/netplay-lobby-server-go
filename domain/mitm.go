package domain

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

// OpenSession opens a new netplay session on the specified MITM server
func (d *MitmDomain) OpenSession(handle string) (*MitmSession, error) {
	// TODO and implement a test using the Add function
	// TODO we need a MITM serve list and need to validate against it
	// TODO fallback is the first entry in the list
	return nil, nil
}

package domain

// MitmSession represents a relay server session.
type MitmSession struct {
	Path string
	Port uint16
}

// MitmServer holds the information of relay server.
type MitmServer struct {
	Handle string
	Name string
	Port uint16
}

// MitmDomain abstracts the mitm logic for handling netplay relays.
type MitmDomain struct {
	server map[string]MitmServer
}

// NewMitmDomain creates a new MITM domain logic.
func NewMitmDomain(servers []MitmServer) *MitmDomain {
	m := make(map[string]MitmServer, len(servers))
	for _, entry := range servers {
		m[entry.Handle] = entry
	}

	return &MitmDomain{m}
}

// OpenSession opens a new netplay session on the specified MITM server
func (d *MitmDomain) OpenSession(handle string) (*MitmSession, error)  {
	// TODO and implement a test using the Add function
	// TODO we need a MITM serve list and need to validate against it
	// TODO fallback is the first entry in the list
	return nil, nil
}

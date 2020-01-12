package domain

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/libretro/netplay-lobby-server-go/model/entity"
)

// SessionDeadline is lifespan of a session that hasn't recieved any updated in seconds.
const SessionDeadline = 60

// requestType enum
type requestType int
// SessionAddType enum value
const (
    SessionCreate requestType = iota
    SessionUpdate
    SessionTouch
)

// mitmSession represents a relay server session
type mitmSession struct {
	IP net.IP
	Port uint16
}

// ErrSessionRejected is thrown when a session got rejected by the domain logic. 
var ErrSessionRejected = errors.New("Session rejected")

// ErrRateLimited is thrown when the rate limit is reached for a particular session.
var ErrRateLimited = errors.New("Rate limit reached")

// SessionRepository interface to decouple the domain logic from the repository code.
type SessionRepository interface {
	Create(s *entity.Session) error
	GetByID(id string) (*entity.Session, error)
	GetAll(deadline time.Time) ([]entity.Session, error)
	Update(s *entity.Session) error
	Touch(id string) error
	PurgeOld(deadline time.Time) error
}

// SessionDomain abrsracts the domain logic for netplay session handling.
type SessionDomain struct {
	sessionRepo SessionRepository
	geopip2Domain *GeoIP2Domain
	validationDomain *ValidationDomain
}

// NewSessionDomain returns an initalized SessionDomain struct.
func NewSessionDomain(sessionRepo SessionRepository, geoIP2Domain *GeoIP2Domain, validationDomain *ValidationDomain) *SessionDomain {
	return &SessionDomain{sessionRepo, geoIP2Domain, validationDomain}
}

// Add adds or updates a session, based on the incomming session information.
// Returns ErrSessionRejected if session got rejected.
// Returns ErrRateLimited if rate limit for a session got reached.
func (d *SessionDomain) Add(session *entity.Session) (*entity.Session, error) {
	var err error
	var savedSession *entity.Session
	var requestType requestType = SessionCreate

	if (session.IP == nil || session.Port == 0) {
		return nil, errors.New("IP or port not set")
	}

	// Decide if this is an CREATE, UPDATE or TOUCH operation
	session.CalculateID()
	if savedSession, err = d.sessionRepo.GetByID(session.ID); err != nil {
		return nil, fmt.Errorf("Can't get saved session: %w", err)
	}
	session.CalculateContentHash()
	if savedSession != nil {
		requestType = SessionTouch
		if savedSession.ContentHash != session.ContentHash {
			requestType = SessionUpdate
		}
	}

	// Validate session on CREATE and UPDATE
	if (requestType == SessionCreate || requestType == SessionUpdate) {
		if !d.validateSession(session) {
			return nil, ErrSessionRejected
		}
	}

	// Open a game session on the selected MITM server if requested
	if (requestType == SessionCreate && session.HostMethod == entity.HostMethodMITM) {
		mitm, err := d.openMITMSession(session)
		if err != nil {
			return nil, fmt.Errorf("Can't open")
		}
		session.MitmIP = mitm.IP
		session.MitmPort = mitm.Port
	}

	// Persist session changes
	switch (requestType) {
	case SessionCreate:
		if session.Country, err = d.geopip2Domain.GetCountryCodeForIP(session.IP); err != nil {
			return nil, fmt.Errorf("Can't find country for given IP %s: %w", session.IP, err)
		}

		if err = d.sessionRepo.Create(session); err != nil {
			return nil, fmt.Errorf("Can't create new session: %w", err)
		}
	case SessionUpdate:
		session.Country = savedSession.Country

		if err = d.sessionRepo.Update(session); err != nil {
			return nil, fmt.Errorf("Can't update old session: %w", err)
		}
	case SessionTouch:
		if err = d.sessionRepo.Touch(session.ID); err != nil {
			return nil, fmt.Errorf("Can't touch old session: %w", err)
		}
	}

	return session, nil
}

// List returns a list of all sessions that are currently beeing hosted
func (d *SessionDomain) List() ([]entity.Session, error) {
	sessions, err := d.sessionRepo.GetAll(d.getDeadline())
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// PurgeOld removes all sessions that have not been updated for longer than 45 seconds.
func (d *SessionDomain) PurgeOld() error {
	if err := d.sessionRepo.PurgeOld(d.getDeadline()); err != nil {
		return err
	}
	return nil
}

// validateSession validaes an incomming session
func (d *SessionDomain) validateSession(s *entity.Session) bool {
	if len(s.Username) > 32 ||
		len(s.CoreName) > 255 ||
		len(s.GameName) > 255 ||
		len(s.GameCRC) != 8 ||
		len(s.RetroArchVersion) > 32 ||
		len(s.SubsystemName) > 255 ||
		len(s.Frontend) > 255 {
			return false;
	}

	if !d.validationDomain.ValidateString(s.Username) ||
		!d.validationDomain.ValidateString(s.CoreVersion) ||
		!d.validationDomain.ValidateString(s.Frontend) ||
		!d.validationDomain.ValidateString(s.SubsystemName) ||
		!d.validationDomain.ValidateString(s.RetroArchVersion) {
			return false;
	}

	return true
}

// openMITMSession opens a new netplay session on the specified MITM server
func (d *SessionDomain) openMITMSession(s *entity.Session) (*mitmSession, error)  {
	// TODO and implement a test using the Add function
	// TODO we need a MITM serve list and need to validate against it
	// TODO maybe this should be it's own domain logic?
	return nil, nil
}

func (d *SessionDomain) getDeadline() time.Time {
	return time.Now().Add(-SessionDeadline * time.Second)
}

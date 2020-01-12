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

// MitmSession represents a relay server session
type mitmSession struct {
	ip net.IP
	port uint16
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
// TODO controller has to set the IP
// TODO middleware for rate-limiting an IP on request basis (not more than 10 adds a minute?)
// TODO write test
func (d *SessionDomain) Add(session *entity.Session) (*entity.Session, error) {
	var err error
	var savedSession *entity.Session
	var requestType requestType = SessionCreate

	// Decide if this is an CREATE, UPDATE or TOUCH operation
	session.CalculateID()
	if savedSession, err = d.sessionRepo.GetByID(session.ID); err != nil {
		return nil, fmt.Errorf("Can't get saved session: %w", err)
	}
	if savedSession != nil {
		requestType = SessionTouch
		session.CalculateContentHash()
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

	if (requestType == SessionCreate) {
		// TODO
		d.openMITMSession(session)
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

//
func (d *SessionDomain) validateSession(s *entity.Session) bool {
	// TODO verify the string lengt of each field

	if d.validationDomain.ValidateString(s.Username) ||
		d.validationDomain.ValidateString(s.CoreVersion) ||
		d.validationDomain.ValidateString(s.Frontend) ||
		d.validationDomain.ValidateString(s.SubsystemName) ||
		d.validationDomain.ValidateString(s.RetroArchVersion) {
			return false;
	}

	return true
}

// TODO write test
func (d *SessionDomain) openMITMSession(s *entity.Session) (*mitmSession, error)  {
	// TODO
	return nil, nil
}

func (d *SessionDomain) getDeadline() time.Time {
	return time.Now().Add(-SessionDeadline * time.Second)
}
